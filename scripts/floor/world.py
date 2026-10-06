#!/usr/bin/env python3
"""Row 93 floor harness — the World side of one task-run (design §4.4, §4.5; AC2.2, AC2.4).

Talks to the daemon only over ``/v1/*`` and ``/mcp/`` (§4.1):
  * pre-flight: ``GET /v1/health`` is 200 and the session's ``tools/list`` is EXACTLY the 8 AILANG
    tools (``arms.world_allowed_tools``). Any failure is ``harness_setup`` (§4.6);
  * the log range: ``log_from`` = selected head's revision + 1 before the spawn, ``log_to`` = the
    head's revision after it (``GET /v1/head`` -> ``GET /v1/worlds/<ref>``);
  * solution provenance (AC2.4): the graded bytes must equal the file World's own log says the
    agent produced — the fold of every successful ``ailang-write``/``ailang-edit`` of that path in
    ``[log_from, log_to]`` for the episode, over the reset worktree's base content. Anything else
    is ``native_write_detected``.

The record chain read here (host code, cited so a drift is findable):
  log entry ``transitionRef`` -> ``world/invocation-record/v2`` {episodeId, skillId, effects[]}
  (host/coordinator/plan.go) -> each ``world/effect-record/v1`` {effect, allowed, failed,
  requestRef, resultRef} (host/broker/record.go) -> the request object, whose payload is
  ``"<len>:<effect><len>:<scope><cost>:<now>:" + payload`` (host/broker/broker.go requestBytes)
  where payload is the policy-tool request ``{op:"write", path, content}`` or
  ``{op:"edit", path, old_text, new_text}`` (packages/se-tools/se_tools/{write,edit}.ail) -> the
  result object, whose ``ok`` says whether policy-tool applied it.

The token is only ever an argument held in memory: it goes into an ``Authorization`` header and
nowhere else; nothing here logs or returns it.
"""
from __future__ import annotations

import base64
import hashlib
import json
import os
import posixpath
import re
import urllib.error
import urllib.parse
import urllib.request

import arms

RECORD_V2 = 'world/invocation-record/v2'
EFFECT_RECORD_V1 = 'world/effect-record/v1'
WORKSPACE_WRITE = 'Workspace.Write'
LOG_PAGE = 500
HTTP_TIMEOUT_S = 15.0
_OPENER = urllib.request.build_opener(urllib.request.ProxyHandler({}))  # loopback: never a proxy


class HarnessSetup(Exception):
    """A World pre-flight failure: the row is ``harness_setup`` with this as its cause_evidence."""


class WorldReadError(Exception):
    pass


class Client:
    def __init__(self, addr: str = arms.DEFAULT_ADDR, timeout: float = HTTP_TIMEOUT_S):
        self.base = f'http://{addr}'
        self.addr = addr
        self.timeout = timeout

    def _req(self, method: str, path: str, body: bytes | None = None, headers: dict | None = None):
        req = urllib.request.Request(self.base + path, data=body, method=method, headers=headers or {})
        try:
            with _OPENER.open(req, timeout=self.timeout) as r:
                return r.status, r.read(), r.headers.get('Content-Type', '')
        except urllib.error.HTTPError as e:
            return e.code, e.read(), e.headers.get('Content-Type', '') if e.headers else ''

    def get_json(self, path: str):
        status, body, _ = self._req('GET', path)
        if status == 404:
            return None
        if status != 200:
            raise WorldReadError(f'GET {path} -> HTTP {status}: {body[:200]!r}')
        return json.loads(body)

    def health(self) -> tuple[int, dict | None]:
        try:
            status, body, _ = self._req('GET', '/v1/health')
        except OSError as e:
            return 0, {'error': str(e)}
        try:
            return status, json.loads(body)
        except ValueError:
            return status, None

    def head_index(self) -> int:
        """The selected head's revision (= its log head's entry index); -1 when no head exists."""
        status, body, _ = self._req('GET', '/v1/head')
        if status == 404:
            return -1
        if status != 200:
            raise WorldReadError(f'GET /v1/head -> HTTP {status}')
        w = self.get_json('/v1/worlds/' + urllib.parse.quote(body.decode().strip(), safe=''))
        if w is None:
            raise WorldReadError('the selected head has no world row')
        return int(w['revision'])

    def log_page(self, start: int, limit: int = LOG_PAGE) -> list:
        return (self.get_json(f'/v1/log?from={int(start)}&limit={int(limit)}') or {}).get('items') or []

    def object(self, ref: str) -> dict | None:
        o = self.get_json('/v1/objects/' + urllib.parse.quote(ref, safe='') + '?payload=true')
        if o is not None and o.get('payload') is not None:
            o['payload_bytes'] = base64.b64decode(o['payload'])
        return o

    def mcp(self, token: str, method: str, params: dict | None = None, rpc_id: int = 1) -> dict:
        msg = {'jsonrpc': '2.0', 'id': rpc_id, 'method': method}
        if params is not None:
            msg['params'] = params
        status, body, _ = self._req('POST', '/mcp/', json.dumps(msg).encode(), {
            'Authorization': f'Bearer {token}', 'Content-Type': 'application/json',
            'Accept': 'application/json, text/event-stream'})
        if status != 200:
            raise WorldReadError(f'POST /mcp/ {method} -> HTTP {status}: {body[:200]!r}')
        return parse_sse_json(body)


def parse_sse_json(body: bytes) -> dict:
    """``/mcp/`` answers as SSE: the JSON-RPC object is the ``data:`` line (QUICKSTART §9)."""
    text = body.decode('utf-8', errors='replace')
    datas = [ln[len('data:'):].strip() for ln in text.splitlines() if ln.startswith('data:')]
    if not datas:
        return json.loads(text)
    return json.loads(datas[-1])


# ---------------------------------------------------------------- pre-flight (§4.4 step 2)

def preflight(client: Client, token: str) -> dict:
    """``/v1/health`` 200 and ``tools/list`` == exactly the 8 names, or HarnessSetup."""
    status, h = client.health()
    if status != 200:
        raise HarnessSetup(f'pre-flight: /v1/health -> HTTP {status}')
    try:
        resp = client.mcp(token, 'tools/list')
    except (OSError, ValueError, WorldReadError) as e:
        raise HarnessSetup(f'pre-flight: tools/list failed: {e}') from None
    names = sorted(t.get('name') for t in ((resp.get('result') or {}).get('tools') or []))
    want = sorted(n.split('__', 2)[2] for n in arms.world_allowed_tools())
    if names != want:
        raise HarnessSetup(f'pre-flight: tools/list {names} != the 8 AILANG tools {want}')
    return {'health': {k: (h or {}).get(k) for k in ('status', 'daemon_version', 'interpreter_ref')},
            'tools': names}


# ---------------------------------------------------------------- log range + provenance (AC2.4)

_REQ_HEAD = re.compile(rb'^(\d+):')


def parse_request_bytes(b: bytes) -> tuple[str, str, int, int, bytes]:
    """Inverse of broker.requestBytes: ``%d:%s%d:%s%d:%d:`` + payload."""
    def take_lp(buf, i):
        m = _REQ_HEAD.match(buf[i:])
        if not m:
            raise WorldReadError('effect request: malformed length prefix')
        n = int(m.group(1))
        j = i + m.end()
        return buf[j:j + n].decode(), j + n
    effect, i = take_lp(b, 0)
    scope, i = take_lp(b, i)
    m = re.match(rb'^(-?\d+):(-?\d+):', b[i:])
    if not m:
        raise WorldReadError('effect request: malformed cost/now')
    return effect, scope, int(m.group(1)), int(m.group(2)), b[i + m.end():]


def _norm(path: str) -> str:
    return posixpath.normpath(path.lstrip('/')) if path else path


def workspace_writes(client: Client, episode: str, log_from: int, log_to: int) -> list:
    """Every Workspace.Write effect of ``episode`` committed in [log_from, log_to], in log order:
    {index, op, path, content | old_text,new_text, ok}. ``ok`` = allowed, not failed, and the
    policy-tool result says ok:true."""
    out = []
    i = log_from
    while i <= log_to:
        items = client.log_page(i, min(LOG_PAGE, log_to - i + 1))
        if not items:
            break
        for e in items:
            idx = int(e['header']['entryIndex'])
            if idx <= log_to:
                out += _entry_writes(client, e, idx, episode)
        i = int(items[-1]['header']['entryIndex']) + 1
    return out


def _entry_writes(client: Client, e: dict, idx: int, episode: str) -> list:
    rec_obj = client.object(e['transitionRef'])
    if not rec_obj or rec_obj.get('semanticId') != RECORD_V2:
        return []
    rec = json.loads(rec_obj['payload_bytes'])
    if rec.get('episodeId') != episode:
        return []
    out = []
    for ref in rec.get('effects') or []:
        eobj = client.object(ref)
        if not eobj or eobj.get('semanticId') != EFFECT_RECORD_V1:
            raise WorldReadError(f'entry {idx}: effect ref {ref} is not an effect record')
        er = json.loads(eobj['payload_bytes'])
        if er.get('effect') != WORKSPACE_WRITE:
            continue
        req_obj = client.object(er['requestRef'])
        if not req_obj:
            raise WorldReadError(f'entry {idx}: request {er["requestRef"]} does not resolve')
        eff, _scope, _cost, _now, payload = parse_request_bytes(req_obj['payload_bytes'])
        if eff != WORKSPACE_WRITE:
            raise WorldReadError(f'entry {idx}: request effect {eff!r} != record effect')
        body = json.loads(payload)
        ok = bool(er.get('allowed')) and not er.get('failed') and bool(er.get('resultRef'))
        if ok:
            res = client.object(er['resultRef'])
            ok = bool(res) and json.loads(res['payload_bytes']).get('ok') is True
        w = {'index': idx, 'op': body.get('op'), 'path': _norm(body.get('path') or ''), 'ok': ok}
        for k in ('content', 'old_text', 'new_text'):
            if k in body:
                w[k] = body[k]
        out.append(w)
    return out


def reconstruct(path: str, writes: list, base: bytes | None) -> dict:
    """Fold the successful writes/edits of ``path`` over ``base`` (the reset worktree's content,
    None when the file did not exist). ``content`` None means World's log cannot produce a file."""
    path = _norm(path)
    cur = base
    last = None
    for w in writes:
        if not w['ok'] or w['path'] != path:
            continue
        if w['op'] == 'write':
            cur = w['content'].encode('utf-8')
        elif w['op'] == 'edit':
            old, new = w['old_text'].encode('utf-8'), w['new_text'].encode('utf-8')
            if cur is None or cur.count(old) != 1:
                # policy-tool applied it, so the real file held old_text exactly once: the file
                # changed outside World's log between two World writes.
                return {'content': None, 'last_index': w['index'], 'reason': 'edit_base_mismatch'}
            cur = cur.replace(old, new, 1)
        else:
            continue
        last = w['index']
    return {'content': cur, 'last_index': last, 'reason': None if cur is not None else 'no_world_write'}


def solution_for_grading(worktree: str, rel: str = 'benchmark/solution.ail') -> bytes | None:
    """The bytes the grader reads (§4.2): the arm's WORKTREE file — never the agent's cwd/void."""
    p = posixpath.join(worktree, rel)
    try:
        with open(p, 'rb') as f:
            return f.read()
    except FileNotFoundError:
        return None


def check_task_run(client: Client, *, episode: str, log_from: int, log_to: int, worktree: str,
                   void_dir: str, base: bytes | None, rel: str = 'benchmark/solution.ail') -> dict:
    """AC2.4 for one World task-run: the graded bytes (read from the WORKTREE, as the grader reads
    them) against World's log, plus the void-dir listing (§4.4 step 4)."""
    graded = solution_for_grading(worktree, rel)
    writes = workspace_writes(client, episode, log_from, log_to)
    out = provenance(graded, rel, writes, base)
    out['world_writes'] = len(writes)
    out['void_listing'] = sorted(os.listdir(void_dir))
    return out


def provenance(graded: bytes | None, path: str, writes: list, base: bytes | None) -> dict:
    r = reconstruct(path, writes, base)
    sha = lambda b: hashlib.sha256(b).hexdigest() if b is not None else None  # noqa: E731
    ok = r['content'] is not None and graded is not None and graded == r['content']
    if r['content'] is None and graded is None:
        ok = True  # nothing written anywhere: a missing solution, graded as logic_error
    return {'solution_provenance': 'world' if ok else 'native_write_detected',
            'native_write_detected': not ok, 'graded_sha256': sha(graded),
            'world_sha256': sha(r['content']), 'last_world_write_index': r['last_index'],
            'reason': None if ok else (r['reason'] or 'bytes_differ')}
