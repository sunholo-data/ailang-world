"""Q1 probe (row 145 design): what a command under srt can do on loopback.
Run inside srt with settings that differ only in network.allowLocalBinding."""
import http.server
import socket
import threading
import urllib.request

try:
    urllib.request.urlopen('http://127.0.0.1:47330/v1/health', timeout=3)
    print('pre-bind: connect to the World daemon 127.0.0.1:47330: OK (reachable)')
except Exception as e:
    print('pre-bind: connect to the World daemon 127.0.0.1:47330: FAIL', type(e).__name__)
try:
    s = socket.socket()
    s.bind(('127.0.0.1', 0))
    s.listen()
    print('bind ok')
except Exception as e:
    print('bind FAIL', type(e).__name__, e)
    raise SystemExit(0)
srv = http.server.HTTPServer(('127.0.0.1', 0), http.server.SimpleHTTPRequestHandler)
threading.Thread(target=srv.handle_request, daemon=True).start()
try:
    r = urllib.request.urlopen(f'http://127.0.0.1:{srv.server_port}/', timeout=5)
    print('connect to own loopback server: ok', r.status)
except Exception as e:
    print('connect to own loopback server: FAIL', type(e).__name__)
try:
    urllib.request.urlopen('http://127.0.0.1:47330/v1/health', timeout=3)
    print('connect to the World daemon 127.0.0.1:47330: OK (reachable)')
except Exception as e:
    print('connect to the World daemon 127.0.0.1:47330: FAIL', type(e).__name__)
try:
    urllib.request.urlopen('https://example.com', timeout=5)
    print('external https: OK')
except Exception as e:
    print('external https: FAIL', type(e).__name__)
