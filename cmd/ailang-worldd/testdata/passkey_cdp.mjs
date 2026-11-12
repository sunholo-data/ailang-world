// The passkey page's browser driver (row 201 A2 M7; XB1, XB2 and host/approval/testdata/webauthn/
// generate.sh), extended for the request card (row 201 A3 sprint 3; XA1-XA3, XA5). Test-only: Chrome headless with a temporary profile under argv[2], driven over the
// DevTools protocol, with a CDP VIRTUAL authenticator, a test key (D-WORLD-109): never a real
// profile, keychain, passkey or platform authenticator. Commands arrive as JSON lines on stdin;
// each answer is one JSON line on stdout. EOF on stdin kills Chrome and removes the profile.
import { spawn } from 'node:child_process';
import { createHash } from 'node:crypto';
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
  if (msg.method === 'Fetch.requestPaused') rewritePaused(msg.params).catch((e) => { rewriteErrors.push(String(e && e.message || e)); });
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
  const box = await evaluate(`(() => { const e = document.getElementById('${id}'); e.scrollIntoView({ block: 'center' }); const r = e.getClientRects()[0] || e.getBoundingClientRect(); return { x: r.x + r.width / 2, y: r.y + r.height / 2, w: r.width }; })()`);
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
  // getRecorder (XA1, XA2): every navigator.credentials.get call's options, then the real get.
  getRecorder: `window.__gets = []; if (navigator.credentials) { const get = navigator.credentials.get.bind(navigator.credentials);
    const b = (x) => btoa(String.fromCharCode(...new Uint8Array(x.buffer || x))).replace(/\\+/g, '-').replace(/\\//g, '_').replace(/=+$/, '');
    navigator.credentials.get = (o) => { const k = o.publicKey; window.__gets.push({ challenge: b(k.challenge), rpId: k.rpId, allow: (k.allowCredentials || []).map((c) => b(c.id)),
      userVerification: k.userVerification, timeout: k.timeout, signal: !!o.signal }); return get(o); }; }`,
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

// fetchRewrite (XA2) rewrites the card's GET /surface/requests/<id> answer in the browser: 'bytes'
// changes the stored bytes and keeps the digest; 'forge' changes both; 'extra' adds a key, 'drop:<key>'
// deletes one, 'type' makes reuse a string, 'float' makes ttl_s 3600.5 and 'utf8' puts a 0xff byte in the
// model, each re-hashed so the digest matches; 'id' names another request and 'nonce' sends a
// non-hex challenge, 'expiry' a negative expiry and 'b64' bytes that are not base64url. Every paused response is counted.
let rewriteMode = null, rewritten = 0;
const rewriteErrors = [];
async function rewritePaused(p) {
  const { body, base64Encoded } = await send('Fetch.getResponseBody', { requestId: p.requestId });
  const card = JSON.parse(base64Encoded ? Buffer.from(body, 'base64').toString('utf8') : body);
  const req = JSON.parse(Buffer.from(card.request, 'base64url').toString('utf8'));
  if (rewriteMode === 'extra') req.extra = 'x';
  else if (rewriteMode.startsWith('drop:')) delete req[rewriteMode.slice(5)];
  else if (rewriteMode === 'type') req.reuse = String(req.reuse);
  else if (rewriteMode === 'float') req.ttl_s += 0.5;
  else if (rewriteMode === 'bytes' || rewriteMode === 'forge') req.ttl_s += 3600;
  let bytes = Buffer.from(JSON.stringify(req), 'utf8');
  if (rewriteMode === 'utf8') { const k = bytes.indexOf('"model":"') + 9; bytes = Buffer.concat([bytes.subarray(0, k), Buffer.from([0xff]), bytes.subarray(k)]); }
  if (rewriteMode === 'id') card.id = 'mr-0000000000000000';
  else if (rewriteMode === 'nonce') card.challenge = 'zz' + card.challenge.slice(2);
  else if (rewriteMode === 'expiry') card.expires = -1;
  else if (rewriteMode === 'b64') card.request = '*' + card.request;
  else card.request = bytes.toString('base64url');
  if (!['bytes', 'id', 'nonce', 'expiry', 'b64'].includes(rewriteMode)) card.digest = createHash('sha256').update(bytes).digest('hex');
  rewritten++;
  await send('Fetch.fulfillRequest', { requestId: p.requestId, responseCode: p.responseStatusCode, responseHeaders: p.responseHeaders,
    body: Buffer.from(JSON.stringify(card)).toString('base64') });
}
async function fetchRewrite({ mode }) {
  rewriteMode = mode;
  if (mode) await send('Fetch.enable', { patterns: [{ urlPattern: '*/surface/requests/mr-*', requestStage: 'Response' }] });
  else await send('Fetch.disable');
  return { rewritten, errors: rewriteErrors };
}

// The card page's view: its state and texts, every rendered field (f-<path> → text), the buttons,
// the message link, the recorded get calls, CSP violations and the page's errors since the tab.
const pageErrors = () => events.filter((e) => e.method === 'Runtime.exceptionThrown' || (e.method === 'Log.entryAdded' && e.params.entry.level === 'error'))
  .map((e) => e.params.entry ? e.params.entry.text : e.params.exceptionDetails.text);
const cardState = async () => Object.assign(await evaluate(`(() => { const $ = (id) => document.getElementById(id); return {
  state: document.body.dataset.state, what: $('what').textContent, status: $('status').textContent, statusError: $('status').hasAttribute('data-error'),
  buttonHidden: $('go').hidden, buttonDisabled: $('go').disabled, againHidden: $('again').hidden, again: $('again').textContent,
  saveHidden: $('save').hidden, saveHref: $('save').getAttribute('href') || '', saveName: $('save').download, sshLine: $('sshline').textContent,
  fields: Object.fromEntries(Array.from(document.querySelectorAll('[id^="f-"]'), (e) => [e.id.slice(2), e.textContent])),
  note: $('origin-note').hidden ? '' : $('origin-note').textContent, ttlLabel: $('f-ttl_s') ? $('f-ttl_s').previousElementSibling.textContent : '', op: window.__op || null,
  gets: window.__gets || null, violations: window.__cspv || [] }; })()`), { errors: pageErrors() });

// openCard opens a card page and waits until it leaves 'loading' (at most settleMs); raw opens any
// page and returns its location.
async function openCard({ url, settleMs = 15000, raw = false }) {
  for (const u of ['about:blank', url]) { // a new document each time
    const from = events.length;
    await send('Page.navigate', { url: u });
    await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the page load');
  }
  if (raw) return { location: await evaluate('location.href') };
  await until(async () => (await evaluate(`document.body && document.body.dataset.state`)) !== 'loading', settleMs, 'the card script');
  return cardState();
}

// press clicks element id; with page, it waits for the new document; else, until the card's state
// is one of states (at most settleMs), or with virtualBudget, after that much page time from a clock
// paused before the click. It returns the card's view, or with page only the location.
async function press({ id, page = false, states = null, settleMs = 30000, virtualBudget = 0 }) {
  const from = events.length;
  if (virtualBudget > 0) await send('Emulation.setVirtualTimePolicy', { policy: 'pause' });
  if (!(await click(id))) throw new Error(id + ' is not visible');
  if (virtualBudget > 0) {
    await send('Emulation.setVirtualTimePolicy', { policy: 'advance', budget: virtualBudget });
    await until(() => events.slice(from).some((e) => e.method === 'Emulation.virtualTimeBudgetExpired'), settleMs, 'the virtual time budget');
  }
  if (page) {
    await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the next page');
    return { location: await evaluate('location.href') };
  }
  if (states) await until(async () => states.includes(await evaluate(`document.body.dataset.state`)), settleMs, 'a state in ' + states);
  return cardState();
}

// credential gives the virtual authenticator a test key's credential (base64 id and PKCS#8 key).
async function credential({ id, rpId, privateKey }) {
  await send('WebAuthn.addCredential', { authenticatorId, credential: { credentialId: id, isResidentCredential: false, rpId, privateKey, signCount: 1 } });
  return {};
}

// fill types value into element id (as a paste: one input event).
async function fill({ id, value }) {
  await evaluate(`(() => { const e = document.getElementById(${JSON.stringify(id)}); e.value = ${JSON.stringify(value)}; e.dispatchEvent(new Event('input')); return 1; })()`);
  return cardState();
}

// save clicks element id and waits for its download into dir; it returns the download's events.
async function save({ id, dir }) {
  const from = events.length;
  await send('Browser.setDownloadBehavior', { behavior: 'allow', downloadPath: dir, eventsEnabled: true }, null);
  if (!(await click(id))) throw new Error(id + ' is not visible');
  await until(() => events.slice(from).some((e) => e.method === 'Browser.downloadProgress' && e.params.state !== 'inProgress'), 15000, 'the download');
  return events.slice(from).filter((e) => e.method.startsWith('Browser.download'))
    .map((e) => ({ event: e.method, state: e.params.state || '', name: e.params.suggestedFilename || '', url: e.params.url || '' }));
}

// waitFor evaluates expression until it is truthy (at most ms) and returns its value.
async function waitFor({ expression, ms = 15000 }) {
  let v;
  await until(async () => (v = await evaluate(expression)), ms, expression);
  return v;
}

const product = (await send('Browser.getVersion', {}, null)).product;
console.log(JSON.stringify({ ready: true, chrome: product, pid: chrome.pid }));
const commands = { tab, open, ceremony, uv, advance, fetchRewrite, openCard, press, waitFor, credential, fill, save, card: cardState };
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
