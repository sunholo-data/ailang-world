#!/usr/bin/env python3
"""AC1.1 — the corpus loader: exactly 23 gate-tier ids with per-YAML sha256; refuses 22 and 24.

Fixture corpora are generated in a temp dir. Rig-only arms (skipped unless FLOOR_AILANG_REPO
points at an ailang checkout holding the pinned commit) load the REAL corpus at 76a5aef and
compare it with the checked-in manifest, cross-check every parsed spec against Ruby's YAML (psych)
when ruby is present, and check test_grade's inline V8 specs against the pinned YAMLs.
"""
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import corpus  # noqa: E402

PINNED_COMMIT = '76a5aef65c81ac2fb4657fafaeff4a89e9faaf75'
MANIFEST = os.path.join(HERE, 'testdata', 'corpus_manifest_core_76a5aef.json')
REPO = os.environ.get('FLOOR_AILANG_REPO')


def spec_text(tid, tier_line):
    body = (f'id: {tid}\n'
            f'description: "fixture {tid}"\n'
            '# a comment line\n'
            'languages: ["python", "ailang"]\n'
            'entrypoint: "main"\n'
            'caps: ["IO"]\n'
            'task_prompt: |\n'
            '  Write a program in <LANG>.\n'
            '  tier: core\n'
            '\n'
            '  It prints ok.\n'
            'expected_stdout: |\n'
            '  ok\n')
    if tier_line is not None:
        body += tier_line + '\n'
    return body + 'tags: [functional]\n'


def make_dir(n_core, n_comment=4, n_untiered=0, extra=()):
    d = tempfile.mkdtemp(prefix='floor-corpus-')
    i = 0
    for k in range(n_core - n_untiered):
        line = 'tier: core # promoted stretch->core, see CURATION.md' if k < n_comment else 'tier: core'
        with open(os.path.join(d, f'core_{i:02d}.yml'), 'w') as f:
            f.write(spec_text(f'core_{i:02d}', line))
        i += 1
    for k in range(n_untiered):  # spec.go: a missing tier defaults to core
        with open(os.path.join(d, f'core_{i:02d}.yml'), 'w') as f:
            f.write(spec_text(f'core_{i:02d}', None))
        i += 1
    for k in range(5):
        with open(os.path.join(d, f'smoke_{k}.yml'), 'w') as f:
            f.write(spec_text(f'smoke_{k}', 'tier: smoke'))
    with open(os.path.join(d, 'stretch_0.yml'), 'w') as f:
        f.write(spec_text('stretch_0', 'tier: stretch # demoted core->stretch'))
    with open(os.path.join(d, 'events.yml'), 'w') as f:  # the Go loader's meta-file, never a spec
        f.write('# suite events\n- version: v1\n  label: x\n')
    for name, text in extra:
        with open(os.path.join(d, name), 'w') as f:
            f.write(text)
    return d


class LoaderCount(unittest.TestCase):
    def tearDown(self):
        for d in getattr(self, 'dirs', []):
            shutil.rmtree(d, ignore_errors=True)

    def mk(self, *a, **k):
        d = make_dir(*a, **k)
        self.dirs = getattr(self, 'dirs', []) + [d]
        return d

    def test_exactly_23_with_trailing_comments(self):
        c = corpus.load_corpus(bench_dir=self.mk(23, n_comment=4))
        self.assertEqual(len(c.tasks), 23)
        self.assertEqual(c.ids(), sorted(c.ids()))
        for t in c.tasks:
            self.assertRegex(t.sha256, r'^[0-9a-f]{64}$')
            self.assertEqual(t.spec['expected_stdout'], 'ok\n')
        self.assertEqual(set(c.manifest()['yaml_sha256']), set(c.ids()))

    def test_refuses_22(self):
        with self.assertRaisesRegex(corpus.CorpusError, 'found 22 specs, expected exactly 23'):
            corpus.load_corpus(bench_dir=self.mk(22))

    def test_refuses_24(self):
        with self.assertRaisesRegex(corpus.CorpusError, 'found 24 specs, expected exactly 23'):
            corpus.load_corpus(bench_dir=self.mk(24))

    def test_untiered_counts_as_core_like_go(self):
        self.assertEqual(len(corpus.load_corpus(bench_dir=self.mk(23, n_untiered=1)).tasks), 23)

    def test_tier_inside_a_block_scalar_is_not_the_tier(self):
        # every fixture's task_prompt carries an indented "tier: core" line; smoke must stay smoke
        c = corpus.load_corpus(bench_dir=self.mk(23), tier='smoke', expected=5)
        self.assertEqual(len(c.tasks), 5)

    def test_id_must_match_file(self):
        d = self.mk(22, extra=[('core_99.yml', spec_text('other', 'tier: core'))])
        with self.assertRaisesRegex(corpus.CorpusError, 'does not match the file name'):
            corpus.load_corpus(bench_dir=d)

    def test_digest_moves_with_any_byte(self):
        d = self.mk(23)
        a = corpus.load_corpus(bench_dir=d).digest()
        with open(os.path.join(d, 'core_03.yml'), 'a') as f:
            f.write('# touched\n')
        self.assertNotEqual(a, corpus.load_corpus(bench_dir=d).digest())


class YAMLSubset(unittest.TestCase):
    def test_constructs(self):
        y = ('id: x\n'
             'description: "a \\"q\\" \\u00e9 \\n end"\n'
             "single: 'it''s'\n"
             'timeout: 120  # comment\n'
             'flag: true\n'
             'plain: hard\n'
             'tags: [a, "b c", \'d\']\n'
             'clip: |\n  one\n\n  two\n\n\n'
             'strip: |-\n  s\n'
             'keep: |+\n  k\n\n'
             'files:\n  numbers.txt: |\n    1\n    2\n  other: "z"\n'
             'args:\n  - "numbers.txt"\n  - plain\n'
             'empty:\n'
             'last: 1\n')
        d = corpus.parse_yaml_subset(y)
        self.assertEqual(d['description'], 'a "q" \u00e9 \n end')
        self.assertEqual(d['single'], "it's")
        self.assertEqual(d['timeout'], 120)
        self.assertIs(d['flag'], True)
        self.assertEqual(d['plain'], 'hard')
        self.assertEqual(d['tags'], ['a', 'b c', 'd'])
        self.assertEqual(d['clip'], 'one\n\ntwo\n')
        self.assertEqual(d['strip'], 's')
        self.assertEqual(d['keep'], 'k\n\n')
        self.assertEqual(d['files'], {'numbers.txt': '1\n2\n', 'other': 'z'})
        self.assertEqual(d['args'], ['numbers.txt', 'plain'])
        self.assertIsNone(d['empty'])
        self.assertEqual(d['last'], 1)

    def test_refusals(self):
        for bad in ('a: &x 1\n', 'a: >\n  folded\n', 'a: {b: 1}\n', 'a: [1, [2]]\n', 'a: 1\na: 2\n',
                    '---\na: 1\n', 'a:\n\t- x\n', 'a: "unterminated\n', 'a:\n  - k: v\n', 'a: b: c\n'):
            with self.subTest(bad):
                with self.assertRaises(corpus.YAMLSubsetError):
                    corpus.parse_yaml_subset(bad)


@unittest.skipUnless(REPO, 'rig-only: set FLOOR_AILANG_REPO to an ailang checkout holding the pinned commit')
class PinnedCorpus(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.c = corpus.load_corpus(repo=REPO, commit=PINNED_COMMIT)

    def test_matches_checked_in_manifest(self):
        with open(MANIFEST) as f:
            want = json.load(f)
        got = self.c.manifest()
        got['digest'] = self.c.digest()
        self.assertEqual(got, want)

    def test_smoke_tier_is_23_too(self):
        self.assertEqual(len(corpus.load_corpus(repo=REPO, commit=PINNED_COMMIT, tier='smoke').tasks), 23)

    def test_v8_specs_used_by_test_grade(self):
        import test_grade
        by = self.c.by_id()
        for s in (test_grade.SPEC_PIPELINE, test_grade.SPEC_CLI_ARGS, test_grade.SPEC_API):
            for k, v in s.items():
                self.assertEqual(by[s['id']].spec.get(k), v, (s['id'], k))

    @unittest.skipUnless(shutil.which('ruby'), 'no ruby for the psych cross-check')
    def test_parser_agrees_with_psych(self):
        for t in self.c.tasks:
            raw = subprocess.run(['git', '-C', REPO, 'show', f'{PINNED_COMMIT}:benchmarks/{t.file}'],
                                 capture_output=True, check=True).stdout
            rb = subprocess.run(['ruby', '-ryaml', '-rjson', '-e', 'puts JSON.generate(YAML.safe_load(STDIN.read))'],
                                input=raw, capture_output=True, check=True)
            self.assertEqual(json.loads(rb.stdout), t.spec, t.id)


if __name__ == '__main__':
    unittest.main()
