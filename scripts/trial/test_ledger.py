"""Tests for ledger.py (row 187 evidence ledger). Run: cd scripts/trial && python3 -m unittest -v test_ledger"""
import json
import os
import tempfile
import unittest

import ledger


def _write(path, obj_or_lines):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, 'w', encoding='utf-8') as f:
        if isinstance(obj_or_lines, list):
            f.write('\n'.join(json.dumps(x) for x in obj_or_lines) + '\n')
        else:
            json.dump(obj_or_lines, f)


class LedgerTest(unittest.TestCase):
    def setUp(self):
        self.root = tempfile.mkdtemp()
        ev = os.path.join(self.root, 'evidence')
        claude_t = os.path.join(ev, 'proj', 'transcripts', 'c1.jsonl')
        _write(claude_t, [{'type': 'assistant'}, {'type': 'result', 'total_cost_usd': 1.5,
                          'usage': {'input_tokens': 10, 'cache_read_input_tokens': 2000, 'cache_creation_input_tokens': 300,
                                    'output_tokens': 40}}])
        _write(os.path.join(ev, 'proj', 'episodes', 'c1.json'), {
            'episode': 'c1', 'project': 'proj', 'agent': 'claude', 'model': 'm', 'task': 't',
            'started_at': '2026-10-16T10:00:00+00:00', 'finished_at': '2026-10-16T10:06:00+00:00',
            'agent_run': {'rc': 0, 'wall_ms': 300000, 'tool_items': {'mcp__world__ailang-read': 3, 'mcp__world__ailang-edit': 2,
                                                                   'mcp__world__workspace-exec': 1, 'mcp__world__builtins-search': 1}},
            'transcript': {'path': claude_t},
            'graders': [{'passed': True, 'outcome_record': {'entry': 9}}, {'passed': False}],
            'solution': {'diff_bytes': 77}})
        codex_t = os.path.join(ev, 'transcripts', 'x1.jsonl')  # the first trial's flat layout
        _write(codex_t, [{'type': 'item.completed', 'item': {'type': 'mcp_tool_call', 'tool': 'ailang-read'}},
                         {'type': 'item.completed', 'item': {'type': 'mcp_tool_call', 'tool': 'workspace-exec'}},
                         {'type': 'item.completed', 'item': {'type': 'agent_message'}},
                         {'type': 'turn.completed', 'usage': {'input_tokens': 500, 'cached_input_tokens': 400, 'output_tokens': 9}}])
        _write(os.path.join(ev, 'episodes', 'x1.json'), {
            'episode': 'x1', 'project': 'proj', 'agent': 'codex', 'model': 'g', 'task': 't2',
            'started_at': '2026-10-17T10:00:00+00:00', 'finished_at': '2026-10-17T10:02:00+00:00',
            'agent_run': {'rc': 0, 'tool_items': {'agent_message': 1, 'mcp_tool_call': 2}},
            'transcript': {'path': codex_t}, 'graders': []})

    def rows(self):
        out = []
        import contextlib, io
        buf = io.StringIO()
        with contextlib.redirect_stdout(buf):
            ledger.main(['--root', self.root, '--json'])
        return {r['episode']: r for r in json.loads(buf.getvalue())}

    def test_claude_row(self):
        r = self.rows()['c1']
        self.assertEqual(r['wall_min'], 5.0)  # wall_ms wins over the timestamps
        self.assertEqual(r['tools'], {'read': 3, 'edit': 2, 'exec': 1, 'other': 1})
        self.assertEqual(r['tokens'], {'input': 10, 'cache_read': 2000, 'cache_write': 300, 'output': 40, 'cost_usd_estimate': 1.5})
        self.assertEqual((r['graders_passed'], r['graders_total'], r['outcomes_recorded']), (1, 2, 1))

    def test_codex_counts_tools_from_the_transcript_and_reads_the_flat_layout(self):
        r = self.rows()['x1']
        self.assertEqual(r['tools'], {'read': 1, 'exec': 1})
        self.assertEqual(r['tokens']['cache_read'], 400)
        self.assertEqual(r['wall_min'], 2.0)  # from the timestamps when wall_ms is absent

    def test_markdown_has_one_line_per_episode_in_date_order(self):
        import contextlib, io
        buf = io.StringIO()
        with contextlib.redirect_stdout(buf):
            ledger.main(['--root', self.root])
        body = [l for l in buf.getvalue().splitlines() if l.startswith('| ') and not l.startswith('| episode')]
        self.assertEqual([l.split('|')[1].strip() for l in body], ['c1', 'x1'])

    def test_killed_shell_arm_sums_partial_usage_and_maps_native_tools(self):
        ev = os.path.join(self.root, 'evidence')
        t = os.path.join(ev, 'proj', 'transcripts', 's1.jsonl')
        _write(t, [{'type': 'assistant', 'message': {'id': 'm1', 'usage': {'input_tokens': 5, 'cache_read_input_tokens': 100, 'output_tokens': 7}}},
                   {'type': 'assistant', 'message': {'id': 'm1', 'usage': {'input_tokens': 5, 'cache_read_input_tokens': 100, 'output_tokens': 7}}},
                   {'type': 'assistant', 'message': {'id': 'm2', 'usage': {'input_tokens': 1, 'output_tokens': 3}}}])
        _write(os.path.join(ev, 'proj', 'episodes', 's1.json'), {
            'episode': 's1', 'arm': 'shell', 'pair_of': 'c1', 'project': 'proj', 'agent': 'claude', 'model': 'm', 'task': 't',
            'started_at': '2026-10-17T10:00:00+00:00', 'finished_at': '2026-10-17T11:00:00+00:00',
            'agent_run': {'rc': 143, 'error_category': 'timeout', 'wall_ms': 3600000,
                          'tool_items': {'Bash': 4, 'Edit': 2, 'Read': 3, 'Grep': 1}},
            'transcript': {'path': t}, 'graders': [{'passed': True}]})
        r = self.rows()['s1']
        self.assertEqual(r['tools'], {'exec': 4, 'edit': 2, 'read': 4})
        self.assertEqual(r['tokens'], {'input': 6, 'cache_read': 100, 'cache_write': 0, 'output': 10,
                                       'cost_usd_estimate': None, 'partial': True})  # m1 counted once
        self.assertEqual((r['arm'], r['pair_of'], r['finish']), ('shell', 'c1', 'timeout'))
        import contextlib, io
        buf = io.StringIO()
        with contextlib.redirect_stdout(buf):
            ledger.main(['--root', self.root])
        line = next(l for l in buf.getvalue().splitlines() if l.startswith('| s1 '))
        self.assertIn('shell (pair of c1)', line)
        self.assertIn('60.0 **timeout**', line)
        self.assertIn('(partial)', line)

    def _report(self, rel, episode, tools, plan=0, handler=0):
        _write(os.path.join(self.root, 'evidence', *rel), {
            'episode': episode, 'entries': [1, 9], 'total_calls': sum(tools.values()), 'calls_by_tool': tools,
            'plan_refusals': {'count': plan, 'by_tool': {}, 'first': []},
            'handler_refusals': {'count': handler, 'by_tool': {}, 'first': []}})

    def test_world_row_counts_calls_and_refusals_from_the_banked_report(self):
        # World's record disagrees with the transcript's 7 calls: the report wins for a World arm.
        self._report(('proj', 'reports', 'c1.json'), 'c1', {'ailang-read': 5, 'ailang-cli': 1, 'workspace-exec': 2}, plan=1, handler=2)
        r = self.rows()['c1']
        self.assertEqual(r['tools'], {'read': 5, 'cli': 1, 'exec': 2})
        self.assertEqual((r['calls_source'], r['refusals']), ('world-report', {'plan': 1, 'handler': 2}))
        self.assertEqual(r['tokens']['output'], 40)  # tokens still from the transcript
        import contextlib, io
        buf = io.StringIO()
        with contextlib.redirect_stdout(buf):
            ledger.main(['--root', self.root])
        line = next(l for l in buf.getvalue().splitlines() if l.startswith('| c1 '))
        self.assertIn('| 8 (5/0/0/2/1; refused 1 plan + 2 handler) |', line)

    def test_codex_world_row_needs_no_special_case_with_a_report(self):
        # x1's runner tool_items are Codex item TYPES; with a report they are not consulted at all.
        self._report(('reports', 'x1.json'), 'x1', {'ailang-edit': 3})
        r = self.rows()['x1']
        self.assertEqual((r['tools'], r['calls_source']), ({'edit': 3}, 'world-report'))

    def test_without_a_report_the_transcript_is_the_fallback(self):
        rows = self.rows()
        self.assertEqual((rows['c1']['calls_source'], rows['c1']['refusals']), ('transcript', None))
        self.assertEqual(rows['x1']['tools'], {'read': 1, 'exec': 1})  # the Codex fallback still applies

    def test_a_report_for_another_episode_or_a_shell_arm_is_ignored(self):
        self._report(('proj', 'reports', 'c1.json'), 'someone-else', {'ailang-read': 99})
        self.assertEqual(self.rows()['c1']['calls_source'], 'transcript')
        ev = os.path.join(self.root, 'evidence')
        t = os.path.join(ev, 'proj', 'transcripts', 's2.jsonl')
        _write(t, [{'type': 'result', 'usage': {'output_tokens': 1}}])
        _write(os.path.join(ev, 'proj', 'episodes', 's2.json'), {
            'episode': 's2', 'arm': 'shell', 'project': 'proj', 'agent': 'claude', 'model': 'm', 'task': 't',
            'started_at': '2026-10-18T10:00:00+00:00', 'agent_run': {'tool_items': {'Bash': 2}},
            'transcript': {'path': t}, 'graders': []})
        self._report(('proj', 'reports', 's2.json'), 's2', {'ailang-read': 99})
        r = self.rows()['s2']
        self.assertEqual((r['tools'], r['calls_source']), ({'exec': 2}, 'transcript'))

    def test_recorded_report_path_is_preferred(self):
        self._report(('elsewhere', 'c1-report.json'), 'c1', {'ailang-write': 4})
        p = os.path.join(self.root, 'evidence', 'proj', 'episodes', 'c1.json')
        with open(p, encoding='utf-8') as f:
            d = json.load(f)
        d['report'] = {'path': os.path.join(self.root, 'evidence', 'elsewhere', 'c1-report.json')}
        _write(p, d)
        self.assertEqual(self.rows()['c1']['tools'], {'write': 4})

    def test_missing_transcript_leaves_tokens_blank(self):
        self.assertEqual(ledger.transcript_usage(os.path.join(self.root, 'nope.jsonl')), {})


if __name__ == '__main__':
    unittest.main()
