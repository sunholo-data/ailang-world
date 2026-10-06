#!/usr/bin/env python3
"""A fake World daemon for the floor harness tests (no Go, no store): the read surface
``/v1/health``, ``/v1/head``, ``/v1/worlds/<ref>``, ``/v1/log``, ``/v1/objects/<ref>`` and a
``/mcp/`` that answers ``tools/list`` as SSE for one bearer token.

Objects are built with the host's codecs, byte for byte where the floor reads bytes:
``broker.requestBytes`` (host/broker/broker.go), ``recordWire`` field order (host/broker/record.go)
and ``recordV2`` (host/coordinator/plan.go).
"""
from __future__ import annotations

import base64
import hashlib
import http.server
import json
import threading
import urllib.parse

import arms


def sha(b: bytes) -> str:
    return 'sha256:' + hashlib.sha256(b).hexdigest()


def request_bytes(effect: str, scope: str, cost: int, now: int, payload: bytes) -> bytes:
    return f'{len(effect)}:{effect}{len(scope)}:{scope}{cost}:{now}:'.encode() + payload


class FakeWorld:
    def __init__(self, token: str, tools=None):
        self.token = token
        self.tools = list(tools) if tools is not None else [n.split('__', 2)[2] for n in arms.world_allowed_tools()]
        self.objects: dict = {}
        self.entries: list = []
        self.health_status = 200
        self._srv = None

    # -- building the log
    def put(self, semantic_id: str, payload: bytes) -> str:
        h = sha(payload)
        self.objects[h] = {'hash': h, 'interfaceHash': sha(semantic_id.encode()), 'semanticId': semantic_id,
                           'provenance': 'fake', 'payload': base64.b64encode(payload).decode()}
        return h

    def _entry(self, record_ref: str):
        i = len(self.entries)
        self.entries.append({'header': {'entryIndex': i, 'semanticsEpoch': 1, 'transitionFn': sha(b'fn'),
                                        'interpreter': sha(b'pin'), 'prevEntryHash': sha(b'p%d' % i),
                                        'writtenBy': 'coordinator:a2a'},
                             'entryHash': sha(b'e%d' % i), 'transitionRef': record_ref})
        return i

    def genesis(self):
        return self._entry(self.put('world/demo/genesis-goal', b'{"goal":"floor"}'))

    def invoke(self, episode: str, skill: str, effects: list) -> int:
        """effects: [(effect, payload_dict, outcome)] with outcome 'ok' | 'refused' | 'denied' | 'failed'."""
        refs = []
        for k, (effect, payload, outcome) in enumerate(effects):
            req = request_bytes(effect, 'worktree', 1, 1700000000 + len(self.entries) * 10 + k,
                                json.dumps(payload).encode())
            req_ref = self.put('world/effect-request/v1', req)
            result_ref = ''
            if outcome in ('ok', 'refused'):
                res = {'ok': outcome == 'ok', 'tool': sha(b'tool'), 'policy_digest': 'x'}
                if outcome == 'refused':
                    res['refused'] = 'policy'
                result_ref = self.put('world/effect-result/v1', json.dumps(res).encode())
            allowed = outcome != 'denied'
            rec = {'effect': effect, 'scope': 'worktree', 'cost': 1, 'budgetBefore': 10,
                   'budgetAfter': 9 if allowed else 10, 'allowed': allowed, 'failed': outcome == 'failed',
                   'denial': '' if allowed else 'denied:budget', 'requestRef': req_ref, 'resultRef': result_ref}
            refs.append(self.put('world/effect-record/v1', json.dumps(rec).encode()))
        record = {'invocationId': f'a2a:{len(self.entries)}', 'episodeId': episode, 'skillId': skill,
                  'transitionFn': sha(b'fn'), 'interpreter': sha(b'pin'), 'semanticsEpoch': 1,
                  'input': sha(b'in'), 'output': sha(b'out'), 'plan': sha(b'plan'), 'effects': refs}
        return self._entry(self.put('world/invocation-record/v2', json.dumps(record).encode()))

    def write(self, episode: str, path: str, content: str, outcome: str = 'ok') -> int:
        return self.invoke(episode, 'ailang-write',
                           [('Workspace.Write', {'op': 'write', 'path': path, 'content': content}, outcome)])

    def edit(self, episode: str, path: str, old: str, new: str, outcome: str = 'ok') -> int:
        return self.invoke(episode, 'ailang-edit', [('Workspace.Write', {'op': 'edit', 'path': path,
                                                                         'old_text': old, 'new_text': new}, outcome)])

    def read(self, episode: str, path: str) -> int:
        return self.invoke(episode, 'ailang-read', [('Workspace.Read', {'op': 'read', 'path': path}, 'ok')])

    # -- serving
    def start(self) -> str:
        fw = self

        class H(http.server.BaseHTTPRequestHandler):
            def log_message(self, *a):
                pass

            def _send(self, code, body, ctype='application/json'):
                b = body if isinstance(body, bytes) else json.dumps(body).encode()
                self.send_response(code)
                self.send_header('Content-Type', ctype)
                self.send_header('Content-Length', str(len(b)))
                self.end_headers()
                self.wfile.write(b)

            def do_GET(self):
                u = urllib.parse.urlparse(self.path)
                q = urllib.parse.parse_qs(u.query)
                p = urllib.parse.unquote(u.path)
                if p == '/v1/health':
                    return self._send(fw.health_status, {'status': 'ok', 'daemon_version': 'fake'})
                if p == '/v1/head':
                    if not fw.entries:
                        return self._send(404, {'class': 'NotFound'})
                    return self._send(200, sha(b'world%d' % (len(fw.entries) - 1)).encode(), 'text/plain')
                if p.startswith('/v1/worlds/'):
                    ref = p[len('/v1/worlds/'):]
                    for i in range(len(fw.entries)):
                        if sha(b'world%d' % i) == ref:
                            return self._send(200, {'ref': ref, 'revision': i})
                    return self._send(404, {'class': 'NotFound'})
                if p == '/v1/log':
                    start = int(q.get('from', ['0'])[0])
                    limit = min(int(q.get('limit', ['100'])[0] or 100), 500)
                    return self._send(200, {'items': fw.entries[start:start + limit]})
                if p.startswith('/v1/objects/'):
                    o = fw.objects.get(p[len('/v1/objects/'):])
                    if o is None:
                        return self._send(404, {'class': 'NotFound'})
                    o = dict(o)
                    if q.get('payload', [''])[0] != 'true':
                        o.pop('payload')
                    return self._send(200, o)
                return self._send(404, {'class': 'NotFound'})

            def do_POST(self):
                n = int(self.headers.get('Content-Length') or 0)
                msg = json.loads(self.rfile.read(n) or b'{}')
                if self.path != '/mcp/' or self.headers.get('Authorization') != f'Bearer {fw.token}':
                    return self._send(401, {'class': 'SessionAbsent'})
                result = {'tools': [{'name': t, 'inputSchema': {'type': 'object'}} for t in fw.tools]}
                body = {'jsonrpc': '2.0', 'id': msg.get('id'), 'result': result}
                return self._send(200, ('event: message\ndata: ' + json.dumps(body) + '\n\n').encode(),
                                  'text/event-stream')

        self._srv = http.server.ThreadingHTTPServer(('127.0.0.1', 0), H)
        threading.Thread(target=self._srv.serve_forever, daemon=True).start()
        return f'127.0.0.1:{self._srv.server_address[1]}'

    def stop(self):
        if self._srv:
            self._srv.shutdown()
            self._srv.server_close()
