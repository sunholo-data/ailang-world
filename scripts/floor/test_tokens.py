#!/usr/bin/env python3
"""AC2.7 — the evidence token grep finds a planted token (plain and gzipped), any
``Authorization: Bearer <value>``, never prints the value, and passes a clean dir whose text
merely talks about bearer tokens and env-var names."""
import contextlib
import gzip
import io
import os
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import tokens  # noqa: E402

TOKEN = '9f' * 32  # the shape authority.Mint produces: 64 lowercase hex

CLEAN = '\n'.join([
    '{"transport":{"type":"streamable_http","bearer_token_env_var":"WORLD_SESSION"},"auth_status":"bearer_token"}',
    'curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" http://127.0.0.1:7644/mcp/',
    'the .mcp.json holds the raw bearer token (V32); Authorization: Bearer <token>',
    '"hash":"sha256:' + 'ab' * 32 + '"',
    'headers: {"Authorization": "Bearer ${WORLD_SESSION}"}',
])


class TokenGrep(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.d = self.tmp.name
        os.makedirs(os.path.join(self.d, 'transcripts'))
        self.write('summary.md', CLEAN)
        self.tokfile = os.path.join(self.d, '..', os.path.basename(self.d) + '.session')
        with open(self.tokfile, 'w') as f:
            f.write(TOKEN + '\n')

    def tearDown(self):
        os.unlink(self.tokfile)
        self.tmp.cleanup()

    def write(self, rel, text, gz=False):
        p = os.path.join(self.d, rel)
        data = text.encode()
        with open(p, 'wb') as f:
            f.write(gzip.compress(data) if gz else data)

    def run_main(self):
        buf = io.StringIO()
        with contextlib.redirect_stdout(buf):
            rc = tokens.main([self.d, '--token-file', self.tokfile])
        return rc, buf.getvalue()

    def test_clean_dir_passes(self):
        rc, out = self.run_main()
        self.assertEqual(rc, 0, out)
        self.assertIn('AC2.7 PASS', out)

    def test_planted_known_token_fails_without_printing_it(self):
        self.write('runs.jsonl', '{"argv":["x","%s"]}' % TOKEN)
        rc, out = self.run_main()
        self.assertEqual(rc, 1)
        self.assertIn('known-token runs.jsonl:1', out)
        self.assertNotIn(TOKEN, out)

    def test_planted_token_in_gzipped_transcript_fails(self):
        self.write('transcripts/cc.jsonl.gz', 'line1\n{"cfg":"%s"}\n' % TOKEN, gz=True)
        rc, out = self.run_main()
        self.assertEqual(rc, 1)
        self.assertIn('transcripts/cc.jsonl.gz:2', out)

    def test_authorization_bearer_with_any_value_fails(self):
        for i, line in enumerate(['Authorization: Bearer s3cr3t-value',
                                  '{"headers":{"Authorization":"Bearer abcdefgh"}}',
                                  "authorization='bearer xyz123'"]):
            with self.subTest(line=line):
                self.write(f'leak{i}.txt', line)
                rc, out = self.run_main()
                self.assertEqual(rc, 1, line)
                self.assertNotIn(line.split()[-1].strip('"}\''), out)
                os.unlink(os.path.join(self.d, f'leak{i}.txt'))

    def test_bare_bearer_credential_fails(self):
        self.write('x.log', 'sent Bearer ' + 'Zz9' * 8)
        rc, _ = self.run_main()
        self.assertEqual(rc, 1)


if __name__ == '__main__':
    unittest.main()
