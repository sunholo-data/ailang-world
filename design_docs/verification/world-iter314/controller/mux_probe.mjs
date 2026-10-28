const S = process.argv[2];
const { createMuxProxyServer } = await import(S + '/dist/sandbox/mux-proxy.js');
const { createHttpProxyServer } = await import(S + '/dist/sandbox/http-proxy.js');
const net = await import('node:net');
const http = createHttpProxyServer({ filter: async () => ({ deniedReason: 'probe' }) });
const mux = createMuxProxyServer({ httpServer: http, handleSocksConnection: s => s.destroy() });
await mux.listenHttpBackend();
await new Promise(r => mux.server.listen(0, '127.0.0.1', r));
const port = mux.getPort();
function trial(send, waitMs) { return new Promise(res => {
  const t0 = Date.now(); let buf = ''; const c = net.connect(port, '127.0.0.1');
  c.on('connect', () => { if (send !== null) c.write(send); });
  c.on('data', d => buf += d);
  const done = why => { c.destroy(); res(`${why} after ${Date.now()-t0}ms bytes=${buf.length} first=${JSON.stringify(buf.split('\r\n')[0])}`); };
  c.on('close', () => done('close')); c.on('error', e => done('error ' + e.code));
  setTimeout(() => done('deadline'), waitMs);
});}
for (let i = 0; i < 3; i++) console.log('HELLO  :', await trial('WORLD-PROBE-HELLO\n', 1500));
for (let i = 0; i < 2; i++) console.log('SILENT :', await trial(null, 1500));
await mux.close?.(); process.exit(0);
