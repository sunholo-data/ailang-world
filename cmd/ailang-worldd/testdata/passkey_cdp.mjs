// The passkey page's browser driver (row 201 A2 M7; XB1, XB2 and host/approval/testdata/webauthn/
// generate.sh). Test-only: Chrome headless with a temporary profile under argv[2], driven over the
// DevTools protocol, with a CDP VIRTUAL authenticator, a test key (D-WORLD-109): never a real
// profile, keychain, passkey or platform authenticator. Commands arrive as JSON lines on stdin;
// each answer is one JSON line on stdout. EOF on stdin kills Chrome and removes the profile.
import { spawn } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import readline from 'node:readline';

const dir = process.argv[2];
const profile = path.join(dir, 'profile');
fs.mkdirSync(profile);
const chrome = spawn(process.env.WORLD_CHROME, ['--headless=new', '--remote-debugging-port=0', '--user-data-dir=' + profile,
  '--no-first-run', '--no-default-browser-check', '--disable-background-networking', 'about:blank'], { stdio: 'ignore' });
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const until = async (ok, ms, what) => {
  for (const end = Date.now() + ms; Date.now() < end; await sleep(25)) if (await ok()) return;
  throw new Error('timed out waiting for ' + what);
};
let port;
await until(() => { // Chrome writes the file in two steps: wait for both lines
  try { port = fs.readFileSync(path.join(profile, 'DevToolsActivePort'), 'utf8'); } catch { return false; }
  return /^\d+\n\/devtools\/browser\/\S+/.test(port);
}, 30000, 'DevToolsActivePort');
const [dport, dpath] = port.trim().split('\n');
const ws = new WebSocket(`ws://127.0.0.1:${dport}${dpath}`);
await new Promise((r, j) => { ws.onopen = r; ws.onerror = j; });
let seq = 0, session = null, target = null, authenticatorId = null;
const replies = new Map(), events = [];
ws.onmessage = (m) => {
  const msg = JSON.parse(m.data);
  if (msg.id && replies.has(msg.id)) { replies.get(msg.id)(msg); replies.delete(msg.id); } else if (msg.method) events.push(msg);
};
const send = async (method, params = {}, sessionId = session) => {
  const msg = await new Promise((r) => { const id = ++seq; replies.set(id, r); ws.send(JSON.stringify({ id, method, params, sessionId: sessionId || undefined })); });
  if (msg.error) throw new Error(`${method}: ${msg.error.message}`);
  return msg.result;
};
const evaluate = async (expression) => {
  const r = await send('Runtime.evaluate', { expression: `JSON.stringify(${expression})`, returnByValue: true, awaitPromise: true });
  if (r.exceptionDetails) throw new Error(r.exceptionDetails.text);
  return JSON.parse(r.result.value);
};
const click = async (id) => {
  const box = await evaluate(`(() => { const e = document.getElementById('${id}'); e.scrollIntoView({ block: 'center' }); const r = e.getBoundingClientRect(); return { x: r.x + r.width / 2, y: r.y + r.height / 2, w: r.width }; })()`);
  if (box.w === 0) return false;
  for (const type of ['mousePressed', 'mouseReleased']) await send('Input.dispatchMouseEvent', { type, x: box.x, y: box.y, button: 'left', clickCount: 1 });
  return true;
};

// The injected credentials operations of XB2 (pending until aborted, or a late result) and of
// generate.sh (one algorithm). Each records what the page passed in window.__op.
const injections = {
  abort: `for (const m of ['create', 'get']) navigator.credentials[m] = (o) => new Promise((_, reject) => {
    const op = window.__op = { method: m, signal: !!(o && o.signal), t0: performance.now(), d0: Date.now() };
    if (op.signal) o.signal.addEventListener('abort', () => { op.abortedAfter = performance.now() - op.t0; op.abortedAfterDate = Date.now() - op.d0; op.aborted = o.signal.aborted; reject(o.signal.reason); });
  });`,
  late: `for (const m of ['create', 'get']) navigator.credentials[m] = (o) => new Promise((resolve) => {
    const op = window.__op = { method: m, signal: !!(o && o.signal), t0: performance.now(), d0: Date.now() };
    if (op.signal) o.signal.addEventListener('abort', () => { op.abortedAfter = performance.now() - op.t0; op.abortedAfterDate = Date.now() - op.d0; op.aborted = o.signal.aborted; });
    const b = (n) => new Uint8Array(n).buffer;
    setTimeout(() => { op.resolvedAfter = performance.now() - op.t0; op.resolvedAfterDate = Date.now() - op.d0; resolve({ id: 'AA', rawId: b(1), type: 'public-key',
      response: { clientDataJSON: b(2), attestationObject: b(3), authenticatorData: b(37), signature: b(8), userHandle: null } }); }, 245000);
  });`,
  eddsa: `{ const create = navigator.credentials.create.bind(navigator.credentials);
    navigator.credentials.create = (o) => { o.publicKey.pubKeyCredParams = o.publicKey.pubKeyCredParams.filter((p) => p.alg === -8); return create(o); }; }`,
};
const watch = `window.__cspv = []; document.addEventListener('securitypolicyviolation', (e) => window.__cspv.push(e.violatedDirective + ' ' + e.blockedURI));`;

// A virtual authenticator that can verify the user (uv) or has no user verification at all.
const addAuthenticator = async (uv) => {
  ({ authenticatorId } = await send('WebAuthn.addVirtualAuthenticator', { options: { protocol: 'ctap2', ctap2Version: 'ctap2_1', transport: 'internal',
    hasResidentKey: true, hasUserVerification: uv, isUserVerified: uv, automaticPresenceSimulation: true } }));
};

async function tab({ authenticator = null, inject = null }) {
  if (target) await send('Target.closeTarget', { targetId: target }, null);
  const { targetId } = await send('Target.createTarget', { url: 'about:blank' }, null);
  const { sessionId } = await send('Target.attachToTarget', { targetId, flatten: true }, null);
  session = sessionId;
  target = targetId;
  for (const d of ['Page', 'Runtime', 'Log']) await send(d + '.enable');
  await send('Page.addScriptToEvaluateOnNewDocument', { source: watch + (inject ? injections[inject] : '') });
  if (authenticator) {
    await send('WebAuthn.enable', { enableUI: false });
    await addAuthenticator(authenticator === 'uv');
  }
  return {};
}

const pageState = () => evaluate(`(() => { const $ = (id) => document.getElementById(id); return {
  state: document.body.dataset.state, what: $('what').textContent, status: $('status').textContent, statusError: $('status').hasAttribute('data-error'),
  button: $('go').textContent, buttonHidden: $('go').hidden, buttonDisabled: $('go').disabled, code: $('code').value, result: $('result').value,
  saveHidden: $('save').hidden, saveHref: $('save').getAttribute('href') || '', saveName: $('save').download,
  violations: window.__cspv || [], op: window.__op || null }; })()`);

async function open({ url }) {
  for (const u of ['about:blank', url]) { // a new document each time, even when only the fragment differs
    const from = events.length;
    await send('Page.navigate', { url: u });
    await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the page load');
  }
  await until(async () => (await evaluate(`document.body && document.body.dataset.state`)) !== 'loading', 5000, 'the page script');
  return pageState();
}

async function ceremony({ virtualBudget = 0, download = null, settleMs = 30000 }) {
  const from = events.length;
  if (download) await send('Browser.setDownloadBehavior', { behavior: 'allow', downloadPath: download, eventsEnabled: true }, null);
  // In virtual time the page's clock is paused from before the click, so the page's timer and the
  // operation start at one instant, however loaded the machine is.
  if (virtualBudget > 0) await send('Emulation.setVirtualTimePolicy', { policy: 'pause' });
  if (!(await click('go'))) throw new Error('the ceremony button is not visible');
  const pending = await pageState();
  const t0 = Date.now();
  if (virtualBudget > 0) {
    await send('Emulation.setVirtualTimePolicy', { policy: 'advance', budget: virtualBudget });
    await until(() => events.slice(from).some((e) => e.method === 'Emulation.virtualTimeBudgetExpired'), settleMs, 'the virtual time budget');
  } else {
    await until(async () => (await evaluate(`document.body.dataset.state`)) !== 'pending', settleMs, 'the ceremony');
  }
  const settled = await pageState();
  const realMs = Date.now() - t0;
  let saveClicked = false;
  if (download) {
    saveClicked = await click('save');
    const done = () => events.slice(from).some((e) => e.method === 'Browser.downloadProgress' && e.params.state !== 'inProgress');
    if (saveClicked) await until(done, settleMs, 'the download');
    else await sleep(500);
  }
  const downloads = events.slice(from).filter((e) => e.method.startsWith('Browser.download'))
    .map((e) => ({ event: e.method, state: e.params.state || '', name: e.params.suggestedFilename || '', url: e.params.url || '' }));
  const errors = events.slice(from).filter((e) => e.method === 'Runtime.exceptionThrown' || (e.method === 'Log.entryAdded' && e.params.entry.level === 'error'))
    .map((e) => e.params.entry ? e.params.entry.text : e.params.exceptionDetails.text);
  return { pending, settled, after: await pageState(), saveClicked, downloads, errors, realMs };
}

// uv moves the authenticator's credentials (keys and counters) to a new virtual authenticator that
// can verify the user, or has no user verification at all.
async function uv({ verified }) {
  const { credentials } = await send('WebAuthn.getCredentials', { authenticatorId });
  await send('WebAuthn.removeVirtualAuthenticator', { authenticatorId });
  await addAuthenticator(verified);
  for (const credential of credentials) await send('WebAuthn.addCredential', { authenticatorId, credential });
  return { moved: credentials.length };
}

// advance runs budget ms of page time at once (virtual time) and returns the page.
async function advance({ budget }) {
  const from = events.length;
  await send('Emulation.setVirtualTimePolicy', { policy: 'advance', budget });
  await until(() => events.slice(from).some((e) => e.method === 'Emulation.virtualTimeBudgetExpired'), 30000, 'the virtual time budget');
  return pageState();
}

const product = (await send('Browser.getVersion', {}, null)).product;
console.log(JSON.stringify({ ready: true, chrome: product, pid: chrome.pid }));
const commands = { tab, open, ceremony, uv, advance };
for await (const line of readline.createInterface({ input: process.stdin })) {
  const { cmd, ...args } = JSON.parse(line);
  try {
    console.log(JSON.stringify({ ok: await commands[cmd](args) }));
  } catch (e) {
    console.log(JSON.stringify({ error: String(e && e.message || e) }));
  }
}
ws.close();
chrome.kill('SIGKILL');
await sleep(300);
fs.rmSync(profile, { recursive: true, force: true });
process.exit(0);
