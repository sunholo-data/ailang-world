#!/usr/bin/env python3
"""Row 93 floor harness — AC2.7: no bearer token in the evidence.

Scans every file under an evidence dir (``.gz`` files decompressed) and fails on:
  * any KNOWN token — read from the session ``--out`` files given with ``--token-file`` (a World
    session token is 64 lowercase hex, host/authority/mint.go; a bare 64-hex pattern would match
    every sha256 in the evidence, so known tokens are matched literally);
  * any ``Authorization: Bearer <value>`` (any quoting, ``:`` or ``=``), and any
    ``Bearer <value>`` whose value looks like a credential (>= 16 token chars).
A value that is a placeholder — starting ``$``, ``<``, ``{`` or ``%`` (``Bearer $(cat f)``,
``Bearer <token>``) — is not a token.

Findings name the file, line and the matched value's LENGTH only; the value is never printed.

    python3 scripts/floor/tokens.py <evidence-dir> [--token-file <session-file> ...]
exits 0 when clean, 1 on any finding, 2 on usage.
"""
from __future__ import annotations

import argparse
import gzip
import os
import re
import sys

AUTH_BEARER = re.compile(rb'(?i)authorization["\']?\s*[:=]\s*["\']?\s*bearer\s+([^\s"\'\\,}]+)')
BARE_BEARER = re.compile(rb'(?i)\bbearer\s+([A-Za-z0-9._~+/=-]{16,})')
PLACEHOLDER = (b'$', b'<', b'{', b'%')


def _read(path: str) -> bytes:
    with open(path, 'rb') as f:
        data = f.read()
    if path.endswith('.gz'):
        try:
            data = gzip.decompress(data)
        except OSError:
            pass
    return data


def scan_bytes(data: bytes, tokens: tuple[bytes, ...] = ()) -> list:
    found = []
    for lineno, line in enumerate(data.splitlines(), 1):
        for t in tokens:
            if t and t in line:
                found.append((lineno, 'known-token', len(t)))
        for kind, rx in (('authorization-bearer', AUTH_BEARER), ('bearer-value', BARE_BEARER)):
            for m in rx.finditer(line):
                v = m.group(1)
                if not v.startswith(PLACEHOLDER):
                    found.append((lineno, kind, len(v)))
    return found


def scan_dir(root: str, tokens: tuple[bytes, ...] = ()) -> list:
    out = []
    for dirpath, _dirs, files in os.walk(root):
        for name in sorted(files):
            p = os.path.join(dirpath, name)
            if os.path.islink(p) or not os.path.isfile(p):
                continue
            for lineno, kind, n in scan_bytes(_read(p), tokens):
                out.append({'file': os.path.relpath(p, root), 'line': lineno, 'kind': kind,
                            'value_len': n})
    return out


def load_tokens(paths) -> tuple[bytes, ...]:
    toks = []
    for p in paths or ():
        with open(p, 'rb') as f:
            t = f.read().strip()
        if t:
            toks.append(t)
    return tuple(toks)


def main(argv=None) -> int:
    ap = argparse.ArgumentParser(description='AC2.7 bearer-token grep over a floor evidence dir')
    ap.add_argument('evidence_dir')
    ap.add_argument('--token-file', action='append', default=[])
    a = ap.parse_args(argv)
    if not os.path.isdir(a.evidence_dir):
        print(f'not a directory: {a.evidence_dir}', file=sys.stderr)
        return 2
    findings = scan_dir(a.evidence_dir, load_tokens(a.token_file))
    for f in findings:
        print(f"TOKEN {f['kind']} {f['file']}:{f['line']} (value length {f['value_len']})")
    print(f'AC2.7 {"FAIL" if findings else "PASS"}: {len(findings)} finding(s) under {a.evidence_dir}'
          f' ({len(a.token_file)} known token(s) checked)')
    return 1 if findings else 0


if __name__ == '__main__':
    sys.exit(main())
