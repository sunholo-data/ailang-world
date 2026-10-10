#!/usr/bin/env python3
"""Offline census of immutable row-93 FINAL rows; never invokes a provider."""
import collections
import hashlib
import json
import pathlib
import statistics
import sys

ROOT = pathlib.Path(__file__).resolve().parents[3]
FINAL = ROOT / 'design_docs/verification/world-floor-final-2026-10-09'


def encoded(value):
    return json.dumps(value, ensure_ascii=False, separators=(',', ':')).encode()


def median(values):
    return statistics.median(values)


def profile():
    report = {'source': str(FINAL.relative_to(ROOT)), 'codex': {}, 'tasks': []}
    arms = {}
    for arm in ('shell', 'world'):
        rows = [json.loads(line) for f in sorted((FINAL/'runs/codex'/arm).glob('r*.jsonl'))
                for line in f.read_text().splitlines()]
        assert len(rows) == 69, (arm, len(rows))
        assert len({(r['task'], r['trial']) for r in rows}) == 69
        assert len({r['task'] for r in rows}) == 23
        arms[arm] = rows
        items = collections.Counter()
        tools = collections.Counter()
        result_bytes = collections.defaultdict(list)
        cached, uncached, turns = [], [], []
        duplicate_results = text_bytes = structured_bytes = 0
        for row in rows:
            raw = (FINAL/'transcripts'/row['transcript_file']).read_bytes()
            assert hashlib.sha256(raw).hexdigest() == row['transcript_sha256']
            events = [json.loads(line) for line in raw.splitlines()]
            usage = [e['usage'] for e in events if e['type'] == 'turn.completed']
            assert len(usage) == 1, row['transcript_file']
            assert usage[0]['input_tokens'] == row['input_tokens']
            turns.append(sum(e['type'] == 'turn.started' for e in events))
            cached.append(usage[0]['cached_input_tokens'])
            uncached.append(usage[0]['input_tokens'] - usage[0]['cached_input_tokens'])
            for event in events:
                if event['type'] != 'item.completed':
                    continue
                item = event['item']
                items[item['type']] += 1
                if item['type'] != 'mcp_tool_call':
                    continue
                name, result = item['tool'], item['result']
                tools[name] += 1
                result_bytes[name].append(len(encoded(result)))
                text = ''.join(block.get('text', '') for block in result['content'])
                structured = result.get('structured_content')
                assert json.loads(text) == structured, row['transcript_file']
                duplicate_results += 1
                text_bytes += len(text.encode())
                structured_bytes += len(encoded(structured))
        report['codex'][arm] = {
            'rows': len(rows), 'median_wall_ms_pooled_context_only': median([r['wall_ms'] for r in rows]),
            'median_input_tokens': median([r['input_tokens'] for r in rows]),
            'median_output_tokens': median([r['output_tokens'] for r in rows]),
            'median_cached_input_tokens': median(cached), 'median_uncached_input_tokens': median(uncached),
            'turn_started_counts': dict(collections.Counter(turns)), 'completed_item_types': dict(items),
            'completed_tools': dict(tools), 'duplicated_text_and_structured_results': duplicate_results,
            'result_text_utf8_bytes': text_bytes, 'structured_result_canonical_bytes': structured_bytes,
            'result_json_bytes': {name: {'count': len(v), 'total': sum(v), 'median': median(v), 'max': max(v)}
                                  for name, v in sorted(result_bytes.items())}}
    for task in sorted({r['task'] for r in arms['shell']}):
        s = [r for r in arms['shell'] if r['task'] == task]
        w = [r for r in arms['world'] if r['task'] == task]
        sm, wm = median([r['wall_ms'] for r in s]), median([r['wall_ms'] for r in w])
        report['tasks'].append({'task': task, 'shell_median_wall_ms': sm, 'world_median_wall_ms': wm,
                                'overhead': wm/sm - 1,
                                'world_median_completed_mcp_calls': median([r['tool_items']['mcp_tool_call'] for r in w]),
                                'shell_median_completed_command_and_file_items': median([
                                    sum(r['tool_items'].get(k, 0) for k in ('command_execution', 'file_change')) for r in s])})
    report['median_per_task_overhead'] = median([t['overhead'] for t in report['tasks']])
    manifest = ROOT/'packages/se-tools/transitions.json'
    descriptors = [{'name': t['id'], 'description': t['description'], 'inputSchema': t['inputSchema'],
                    'outputSchema': t['outputSchema']} for t in json.loads(manifest.read_text())
                   if t['id'] != 'workspace-exec']
    assert len(descriptors) == 8
    report['descriptor_reconstruction'] = {
        'qualification': 'Current manifest, canonical compact JSON; not a banked FINAL tools/list wire response.',
        'manifest_sha256': hashlib.sha256(manifest.read_bytes()).hexdigest(),
        'canonical_list_bytes': len(encoded(descriptors)),
        'description_utf8_bytes': sum(len(t['description'].encode()) for t in descriptors),
        'input_schema_canonical_bytes': sum(len(encoded(t['inputSchema'])) for t in descriptors),
        'output_schema_canonical_bytes': sum(len(encoded(t['outputSchema'])) for t in descriptors),
        'tools': {t['name']: len(encoded(t)) for t in descriptors}}
    return report


if __name__ == '__main__':
    json.dump(profile(), sys.stdout, indent=2, sort_keys=True)
    print()
