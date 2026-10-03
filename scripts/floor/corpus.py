#!/usr/bin/env python3
"""Row 93 floor harness — corpus loader (design §4.2, V5–V8; AC1.1).

Reads the benchmark YAMLs of the ailang repo AT A PINNED COMMIT (``git show
<commit>:benchmarks/<f>``, so a dirty checkout can never leak in) or, for fixtures, from a
plain directory. It selects the tier exactly as the Go loader does
(``internal/eval_harness/spec.go``: a missing ``tier`` defaults to ``core``;
``cmd/ailang/eval_helpers.go``: ``events.yml`` is a meta-file, never a spec) and refuses any
count other than the expected one (23 for the gate tier).

The ``tier`` line is matched with its trailing comment allowed (V5: 4 of the 23 core files
carry ``tier: core # ...``); a column-0 match only, so an indented ``tier:`` inside a block
scalar is never read as the spec's tier.

There is no PyYAML on the rig or guaranteed in CI, so the spec files are parsed with a STRICT
subset parser (``parse_yaml_subset``): every construct it does not implement is a refusal, not
a best effort. The subset is what the benchmark specs use: top-level mappings; plain, single-
and double-quoted scalars; flow sequences of scalars; ``|``/``|-``/``|+`` literal blocks; one
nested level of block mappings and block sequences of scalars. It is cross-checked against a
real YAML implementation on the rig (see ``FLOOR_YAML_CROSSCHECK`` in test_corpus.py) and
against the banked rows' own ``expected_stdout``/``caps`` in agreement.py.
"""
from __future__ import annotations

import dataclasses
import hashlib
import json
import os
import re
import subprocess
from typing import Any

GATE_TIER = 'core'
EXPECTED_COUNT = 23
META_FILES = frozenset({'events.yml'})
DEFAULT_TIER = 'core'  # spec.go: `Missing defaults to "core"`

_TIER_RE = re.compile(r'^tier:[ \t]*([^\s#]*)[ \t]*(?:#.*)?$')


class CorpusError(Exception):
    """The corpus cannot be loaded as pinned; never a partial result."""


class YAMLSubsetError(CorpusError):
    """A construct outside the implemented YAML subset (refused, never guessed)."""


# ---------------------------------------------------------------- YAML subset parser

def _strip_comment(s: str) -> str:
    """Remove a trailing ` #...` comment from a plain (unquoted) scalar or flow value."""
    out = []
    quote = None
    i = 0
    while i < len(s):
        c = s[i]
        if quote:
            out.append(c)
            if quote == '"' and c == '\\' and i + 1 < len(s):
                out.append(s[i + 1])
                i += 2
                continue
            if c == quote:
                quote = None
        else:
            if c in '"\'' and (not out or out[-1] in ' \t[,'):
                quote = c
            elif c == '#' and (not out or out[-1] in ' \t'):
                break
            out.append(c)
        i += 1
    if quote:
        raise YAMLSubsetError(f'unterminated quoted scalar: {s!r}')
    return ''.join(out).rstrip(' \t')


_DQ_ESCAPES = {'0': '\0', 'a': '\a', 'b': '\b', 't': '\t', '\t': '\t', 'n': '\n', 'v': '\v',
               'f': '\f', 'r': '\r', 'e': '\x1b', ' ': ' ', '"': '"', '/': '/', '\\': '\\',
               'N': '\x85', '_': '\xa0', 'L': ' ', 'P': ' '}


def _double_quoted(s: str) -> str:
    if len(s) < 2 or s[0] != '"' or s[-1] != '"':
        raise YAMLSubsetError(f'not a complete double-quoted scalar: {s!r}')
    body = s[1:-1]
    out = []
    i = 0
    while i < len(body):
        c = body[i]
        if c == '"':
            raise YAMLSubsetError(f'unescaped quote inside double-quoted scalar: {s!r}')
        if c != '\\':
            out.append(c)
            i += 1
            continue
        if i + 1 >= len(body):
            raise YAMLSubsetError(f'dangling escape: {s!r}')
        e = body[i + 1]
        if e in _DQ_ESCAPES:
            out.append(_DQ_ESCAPES[e])
            i += 2
        elif e in 'xuU':
            n = {'x': 2, 'u': 4, 'U': 8}[e]
            hexs = body[i + 2:i + 2 + n]
            if len(hexs) != n or not re.fullmatch(r'[0-9A-Fa-f]+', hexs):
                raise YAMLSubsetError(f'bad \\{e} escape: {s!r}')
            out.append(chr(int(hexs, 16)))
            i += 2 + n
        else:
            raise YAMLSubsetError(f'unknown escape \\{e}: {s!r}')
    return ''.join(out)


def _single_quoted(s: str) -> str:
    if len(s) < 2 or s[0] != "'" or s[-1] != "'":
        raise YAMLSubsetError(f'not a complete single-quoted scalar: {s!r}')
    body = s[1:-1]
    if re.search(r"(?<!')'(?!')", body.replace("''", '')):
        raise YAMLSubsetError(f'stray quote in single-quoted scalar: {s!r}')
    return body.replace("''", "'")


_INT_RE = re.compile(r'[-+]?[0-9]+')


def _scalar(raw: str) -> Any:
    s = _strip_comment(raw).strip(' \t')
    if s.startswith('"'):
        return _double_quoted(s)
    if s.startswith("'"):
        return _single_quoted(s)
    if s.startswith('[') or s.startswith('{'):
        raise YAMLSubsetError(f'flow collection where a scalar is required: {raw!r}')
    if s[:1] in ('&', '*', '!', '|', '>', '%', '@', '`'):
        raise YAMLSubsetError(f'unsupported scalar indicator: {raw!r}')
    if s in ('true', 'True', 'TRUE'):
        return True
    if s in ('false', 'False', 'FALSE'):
        return False
    if s in ('', '~', 'null', 'Null', 'NULL'):
        return None
    if _INT_RE.fullmatch(s):
        return int(s)
    if ': ' in s or s.endswith(':'):
        raise YAMLSubsetError(f'plain scalar containing a mapping indicator: {raw!r}')
    return s


def _flow_seq(raw: str) -> list:
    s = _strip_comment(raw).strip(' \t')
    if not (s.startswith('[') and s.endswith(']')):
        raise YAMLSubsetError(f'unterminated/multi-line flow sequence: {raw!r}')
    inner = s[1:-1]
    if '[' in inner.replace('"[', '').replace("'[", '') and re.search(r'(^|,)\s*\[', inner):
        raise YAMLSubsetError(f'nested flow collection: {raw!r}')
    items, cur, quote = [], [], None
    i = 0
    while i < len(inner):
        c = inner[i]
        if quote:
            cur.append(c)
            if quote == '"' and c == '\\' and i + 1 < len(inner):
                cur.append(inner[i + 1])
                i += 2
                continue
            if c == quote:
                quote = None
        elif c in '"\'':
            quote = c
            cur.append(c)
        elif c == ',':
            items.append(''.join(cur))
            cur = []
        elif c in '{}':
            raise YAMLSubsetError(f'flow mapping inside a flow sequence: {raw!r}')
        else:
            cur.append(c)
        i += 1
    tail = ''.join(cur)
    if tail.strip() or items:
        items.append(tail)
    out = []
    for it in items:
        if not it.strip():
            raise YAMLSubsetError(f'empty flow sequence entry: {raw!r}')
        out.append(_scalar(it))
    return out


def _indent(line: str) -> int:
    n = len(line) - len(line.lstrip(' '))
    if line[n:n + 1] == '\t':
        raise YAMLSubsetError(f'tab indentation: {line!r}')
    return n


def _is_blank_or_comment(line: str) -> bool:
    st = line.strip(' \t')
    return st == '' or st.startswith('#')


def _block_literal(lines: list[str], i: int, header: str, parent_indent: int) -> tuple[str, int]:
    """Parse a `|`, `|-` or `|+` literal block whose content starts at lines[i]."""
    m = re.fullmatch(r'\|([-+]?)', header)
    if not m:
        raise YAMLSubsetError(f'unsupported block scalar header {header!r} (only |, |-, |+)')
    chomp = m.group(1)
    content_indent = None
    body: list[str] = []
    j = i
    while j < len(lines):
        line = lines[j]
        if line.strip(' ') == '':
            body.append('')
            j += 1
            continue
        ind = _indent(line)
        if content_indent is None:
            if ind <= parent_indent:
                break
            content_indent = ind
        if ind < content_indent:
            break
        body.append(line[content_indent:])
        j += 1
    # trailing blank lines belong to the block only for chomping purposes; they were consumed
    text_lines = body
    trailing = 0
    while text_lines and text_lines[-1] == '':
        text_lines = text_lines[:-1]
        trailing += 1
    if not text_lines:
        text = ''
    else:
        text = '\n'.join(text_lines)
        if chomp == '':
            text += '\n'
        elif chomp == '+':
            text += '\n' + '\n' * trailing
    return text, j


def _split_key(line: str, indent: int) -> tuple[str, str]:
    body = line[indent:]
    m = re.match(r'([A-Za-z0-9_][A-Za-z0-9_.\-/ ]*?|"[^"]*"|\'[^\']*\'):(?:[ \t]+(.*)|[ \t]*)$', body)
    if not m:
        raise YAMLSubsetError(f'not a `key: value` line: {line!r}')
    key = m.group(1)
    if key[:1] in '"\'':
        key = _scalar(key)
    return key, (m.group(2) or '')


def _value(lines: list[str], i: int, rest: str, indent: int) -> tuple[Any, int]:
    """Value for a key at `indent` whose inline remainder is `rest`; lines[i] is the next line."""
    r = _strip_comment(rest).strip(' \t')
    if r.startswith('|'):
        return _block_literal(lines, i, r, indent)
    if r.startswith('>'):
        raise YAMLSubsetError('folded block scalars (>) are not supported')
    if r.startswith('['):
        return _flow_seq(rest), i
    if r.startswith('{'):
        raise YAMLSubsetError('flow mappings are not supported')
    if r != '':
        return _scalar(rest), i
    # empty inline value: a nested block (mapping or sequence) or null
    j = i
    while j < len(lines) and _is_blank_or_comment(lines[j]):
        j += 1
    if j >= len(lines):
        return None, i
    child_indent = _indent(lines[j])
    child = lines[j][child_indent:]
    if (child.startswith('- ') or child == '-') and child_indent >= indent:
        return _block_seq(lines, j, child_indent)
    if child_indent > indent:
        return _mapping(lines, j, child_indent)
    return None, i


def _block_seq(lines: list[str], i: int, indent: int) -> tuple[list, int]:
    out = []
    while i < len(lines):
        line = lines[i]
        if _is_blank_or_comment(line):
            i += 1
            continue
        ind = _indent(line)
        if ind < indent:
            break
        if ind > indent:
            raise YAMLSubsetError(f'unexpected indentation in block sequence: {line!r}')
        body = line[ind:]
        if not (body.startswith('- ') or body == '-'):
            break
        item = body[2:] if body.startswith('- ') else ''
        st = item.strip(' \t')
        if st == '' or st.startswith(('|', '>', '[', '{', '- ')) or re.match(r'[A-Za-z0-9_"\'][^:]*:(\s|$)', st):
            raise YAMLSubsetError(f'block sequence entries must be inline scalars: {line!r}')
        out.append(_scalar(item))
        i += 1
    return out, i


def _mapping(lines: list[str], i: int, indent: int) -> tuple[dict, int]:
    out: dict = {}
    while i < len(lines):
        line = lines[i]
        if _is_blank_or_comment(line):
            i += 1
            continue
        ind = _indent(line)
        if ind < indent:
            break
        if ind > indent:
            raise YAMLSubsetError(f'unexpected indentation: {line!r}')
        if line[ind:].startswith('- '):
            if indent == 0:
                raise YAMLSubsetError(f'top-level sequence: {line!r}')
            break
        key, rest = _split_key(line, ind)
        if key in out:
            raise YAMLSubsetError(f'duplicate key {key!r}')
        val, i = _value(lines, i + 1, rest, ind)
        out[key] = val
    return out, i


def parse_yaml_subset(text: str) -> dict:
    """Parse a benchmark spec. Refuses (YAMLSubsetError) anything outside the subset."""
    if text.startswith('﻿'):
        text = text[1:]
    if '\t' in ''.join(l[:len(l) - len(l.lstrip(' \t'))] for l in text.split('\n')):
        raise YAMLSubsetError('tab in indentation')
    lines = text.split('\n')
    for l in lines:
        if l.startswith(('---', '...', '%')):
            raise YAMLSubsetError(f'document markers/directives are not supported: {l!r}')
    doc, i = _mapping(lines, 0, 0)
    while i < len(lines):
        if not _is_blank_or_comment(lines[i]):
            raise YAMLSubsetError(f'unparsed trailing content at line {i + 1}: {lines[i]!r}')
        i += 1
    return doc


# ---------------------------------------------------------------- corpus

@dataclasses.dataclass(frozen=True)
class Task:
    id: str
    file: str
    sha256: str
    spec: dict

    @property
    def caps(self) -> list:
        return list(self.spec.get('caps') or [])


@dataclasses.dataclass(frozen=True)
class Corpus:
    source: str        # 'git:<repo>@<commit>' or 'dir:<path>'
    commit: str | None
    tier: str
    tasks: tuple

    def ids(self) -> list:
        return [t.id for t in self.tasks]

    def by_id(self) -> dict:
        return {t.id: t for t in self.tasks}

    def manifest(self) -> dict:
        """The pinned identity of the corpus: commit plus per-YAML sha256 (sorted)."""
        return {'commit': self.commit, 'tier': self.tier, 'count': len(self.tasks),
                'yaml_sha256': {t.id: t.sha256 for t in self.tasks}}

    def digest(self) -> str:
        blob = json.dumps(self.manifest(), sort_keys=True, separators=(',', ':')).encode()
        return hashlib.sha256(blob).hexdigest()


def tier_of(text: str) -> str:
    """The spec's tier from its column-0 `tier:` line (trailing comment allowed, V5)."""
    found = None
    for line in text.split('\n'):
        if line.startswith('tier:'):
            m = _TIER_RE.match(line)
            if not m:
                raise CorpusError(f'unreadable tier line: {line!r}')
            if found is not None:
                raise CorpusError('duplicate tier line')
            found = m.group(1)
    if found is None or found == '':
        return DEFAULT_TIER
    if found[:1] in '"\'':
        found = _scalar(found)
    return found


def _git(repo: str, *args: str) -> bytes:
    p = subprocess.run(['git', '-C', repo, *args], capture_output=True)
    if p.returncode != 0:
        raise CorpusError(f'git {" ".join(args)} failed in {repo}: {p.stderr.decode(errors="replace").strip()}')
    return p.stdout


def _files_from_git(repo: str, commit: str) -> dict:
    full = _git(repo, 'rev-parse', '--verify', commit + '^{commit}').decode().strip()
    names = _git(repo, 'ls-tree', '--name-only', full, 'benchmarks/').decode().split('\n')
    out = {}
    for n in names:
        base = os.path.basename(n)
        if n and (base.endswith('.yml') or base.endswith('.yaml')):
            out[base] = _git(repo, 'show', f'{full}:{n}')
    return full, out


def _files_from_dir(path: str) -> dict:
    out = {}
    for base in sorted(os.listdir(path)):
        if base.endswith('.yml') or base.endswith('.yaml'):
            with open(os.path.join(path, base), 'rb') as f:
                out[base] = f.read()
    return out


def load_corpus(*, repo: str | None = None, commit: str | None = None, bench_dir: str | None = None,
                tier: str = GATE_TIER, expected: int = EXPECTED_COUNT) -> Corpus:
    """Load exactly `expected` specs of `tier`. Refuses any other count (AC1.1)."""
    if (repo is None) == (bench_dir is None):
        raise CorpusError('give exactly one of repo+commit or bench_dir')
    if repo is not None:
        if not commit:
            raise CorpusError('a repo source must be pinned by commit')
        full, files = _files_from_git(repo, commit)
        source = f'git:{repo}@{full}'
    else:
        full, files = None, _files_from_dir(bench_dir)
        source = f'dir:{bench_dir}'
    if not files:
        raise CorpusError(f'no benchmark YAMLs found in {source}')
    tasks = []
    for base in sorted(files):
        if base in META_FILES:
            continue
        raw = files[base]
        text = raw.decode('utf-8')
        if tier_of(text) != tier:
            continue
        spec = parse_yaml_subset(text)
        stem = base.rsplit('.', 1)[0]
        if spec.get('id') != stem:
            raise CorpusError(f'{base}: id {spec.get("id")!r} does not match the file name')
        if not spec.get('languages') or 'ailang' not in spec['languages']:
            raise CorpusError(f'{base}: not an ailang benchmark')
        tasks.append(Task(id=stem, file=base, sha256=hashlib.sha256(raw).hexdigest(), spec=spec))
    if len(tasks) != expected:
        raise CorpusError(f'tier {tier!r}: found {len(tasks)} specs, expected exactly {expected}: '
                          f'{[t.id for t in tasks]}')
    return Corpus(source=source, commit=full, tier=tier, tasks=tuple(tasks))


def main(argv=None) -> int:
    import argparse
    ap = argparse.ArgumentParser(description='Print the pinned corpus manifest (ids + sha256).')
    ap.add_argument('--repo')
    ap.add_argument('--commit')
    ap.add_argument('--bench-dir')
    ap.add_argument('--tier', default=GATE_TIER)
    ap.add_argument('--expected', type=int, default=EXPECTED_COUNT)
    a = ap.parse_args(argv)
    c = load_corpus(repo=a.repo, commit=a.commit, bench_dir=a.bench_dir, tier=a.tier, expected=a.expected)
    m = c.manifest()
    m['digest'] = c.digest()
    print(json.dumps(m, indent=2, sort_keys=True))
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
