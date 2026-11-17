'use strict';
// AILANG World's confirmation card (row 243 V3, design §2.5). It reads one stored verb request from
// the daemon, checks that its bytes hash to its digest, shows EVERY field of the signed bytes through
// one table (and refuses a request with a field it cannot show), and builds the 86-byte world-verb/v1
// message itself (never a hash from the server): Save the message and the ssh-keygen line. A state
// other than pending shows that state and no controls. Text is set with textContent only. The submit
// half (sprint 4, design §2.6): a passkey whose challenge is the message hash, or a pasted SSH
// signature, POSTed to /surface/confirm/<id>; the answer or the refusal is shown, then the card is
// read again (after a refusal or a 500 it is still pending and the controls return).
(() => {
  const NET_MS = 10000; // A3 §2.7 [r2]: every GET and POST, its body included
  const TIMEOUT_MS = 240000; // A2's ceremony bound (A2 §2.14 [r2])
  const PREFIX = 'world-verb/v1', NAMESPACE = 'world-verb@sunholo', LEN = 86; // approval.VerbPrefix, VerbNamespace, VerbMessageLen
  // Every key of a world/verb-request/v1 (approval.VerbRequest's json tags, in order), with its label
  // and JSON type. TestConfirmPageFieldTable compares this list with the Go type.
  const CONFIRM_FIELDS = /* CONFIRM_FIELDS */[
    {"key": "schema", "label": "Schema", "type": "string"},
    {"key": "verb", "label": "Verb", "type": "string"},
    {"key": "argv", "label": "Command words", "type": "strings"},
    {"key": "episode", "label": "Episode", "type": "string"},
    {"key": "prompt", "label": "Prompt", "type": "string"},
    {"key": "scope", "label": "Approval scope", "type": "string"},
    {"key": "title", "label": "Title", "type": "string"},
    {"key": "body", "label": "Body", "type": "string"},
    {"key": "committer", "label": "Committer", "type": "string"},
    {"key": "protected", "label": "Protected paths changed", "type": "strings"}
  ]/* /CONFIRM_FIELDS */;
  const CLOSED = { approved: 'Already confirmed', superseded: 'Superseded by a newer finish of this episode', expired: 'Expired' };

  const $ = (id) => document.getElementById(id);
  const setState = (s) => { document.body.dataset.state = s; };
  const refuse = (state, text) => { setState(state); $('what').textContent = text; $('what').toggleAttribute('data-error', true); };
  const unb64u = (s) => {
    if (typeof s !== 'string' || !/^[A-Za-z0-9_-]*$/.test(s)) throw new Error('not base64url');
    return Uint8Array.from(atob(s.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
  };
  const b64u = (buf) => {
    let s = '';
    for (const b of new Uint8Array(buf)) s += String.fromCharCode(b);
    return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
  };
  // A2 §2.7's assertion file, built here rather than by toJSON() (as decide.js does).
  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });
  const hex = (b) => Array.from(b, (x) => x.toString(16).padStart(2, '0')).join('');
  const unhex = (s) => (typeof s === 'string' && /^[0-9a-f]{64}$/.test(s) ? Uint8Array.from(s.match(/../g), (x) => parseInt(x, 16)) : null);
  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
  const time = (v) => `${new Date(v * 1000).toISOString()} (${v})`;
  // quote is verbs.Quote: a word outside safeWord's set is single-quoted, a ' inside as '"'"'.
  const SAFE = /^[A-Za-z0-9@%+=:,./_-]+$/;
  const quote = (argv) => argv.map((a) => (SAFE.test(a) ? a : "'" + a.split("'").join(`'"'"'`) + "'")).join(' ');

  async function net(url, init) {
    const c = new AbortController();
    const timer = setTimeout(() => c.abort(), NET_MS);
    try {
      const r = await fetch(url, Object.assign({ cache: 'no-store', credentials: 'omit', signal: c.signal, headers: { Accept: 'application/json' } }, init));
      const text = await r.text();
      let body = null;
      try { body = JSON.parse(text); } catch (_) { body = null; }
      return { status: r.status, body };
    } finally {
      clearTimeout(timer);
    }
  }
  const why = (e) => (e && e.name === 'AbortError' ? `no answer within ${NET_MS / 1000} s` : String(e && e.message || e));

  // unshown is the first key of request r the card cannot show (unlisted, of another type, or
  // missing), or ''. A value that is not a JSON object is itself unshowable.
  const typed = (f, v) => (f.type === 'string' ? typeof v === 'string' : Array.isArray(v) && v.every((x) => typeof x === 'string'));
  function unshown(r) {
    if (r === null || typeof r !== 'object' || Array.isArray(r)) return 'the request';
    for (const k of Object.keys(r)) {
      const f = CONFIRM_FIELDS.find((x) => x.key === k);
      if (!f || !typed(f, r[k])) return k;
    }
    for (const f of CONFIRM_FIELDS) if (!Object.prototype.hasOwnProperty.call(r, f.key)) return f.key;
    return '';
  }

  // message is the 86-byte verb message (design §2.2): PREFIX ‖ 0x00 ‖ digest ‖ challenge ‖ u64be(expiry).
  function message(digest, challenge, expiry) {
    const m = new Uint8Array(LEN);
    m.set(new TextEncoder().encode(PREFIX), 0);
    m.set(digest, 14);
    m.set(challenge, 46);
    new DataView(m.buffer).setBigUint64(78, BigInt(expiry));
    return m;
  }

  const row = (dl, id, label, text) => {
    const dt = document.createElement('dt');
    dt.textContent = label;
    const dd = document.createElement('dd');
    dd.id = id;
    dd.textContent = text;
    dl.append(dt, dd);
    return dd;
  };
  function render(r, s, h) {
    $('argv').textContent = quote(r.argv);
    $('prompt').textContent = r.prompt;
    const dl = $('signed');
    dl.textContent = '';
    row(dl, 'm-digest', 'Request digest (SHA-256)', s.digest);
    row(dl, 'm-challenge', 'Challenge', s.challenge);
    row(dl, 'm-expires', 'Expires', time(s.expires));
    row(dl, 'm-hash', 'Message hash (the passkey challenge)', hex(h));
    for (const f of CONFIRM_FIELDS) {
      const v = r[f.key];
      if (f.key === 'body') {
        const dd = row(dl, 'f-body-size', f.label, `${new TextEncoder().encode(v).length} bytes`);
        const details = document.createElement('details');
        const summary = document.createElement('summary');
        summary.textContent = 'the full body';
        const pre = document.createElement('pre');
        pre.id = 'f-body';
        pre.textContent = v;
        details.append(summary, pre);
        dd.append(details);
        continue;
      }
      row(dl, 'f-' + f.key, f.label, f.type === 'strings' ? v.join('\n') : v);
    }
    $('shown').hidden = false;
  }

  const id = (location.pathname.match(/^\/surface\/confirm\/(vr-[0-9a-f]{16})$/) || [])[1];
  let digestHex = '', hash = null, group = null; // the pending card's, for the submit half
  const buttons = (on) => { $('go').disabled = !on || !group; $('post-sig').disabled = !on; };
  async function load() {
    if (!id) return refuse('refused', 'This is not a confirmation card\'s address; nothing to confirm.');
    let res;
    try {
      res = await net('/surface/confirmations/' + id);
    } catch (e) {
      return refuse('failed', `The request could not be read (${why(e)}). Nothing to confirm here; reload, or answer at the terminal.`);
    }
    if (res.status !== 200 || !res.body) {
      return refuse('failed', `The daemon did not give the request (HTTP ${res.status}${res.body && res.body.reason ? ': ' + res.body.reason : ''}). Nothing to confirm here.`);
    }
    const s = res.body;
    let bytes;
    try {
      bytes = unb64u(s.request);
    } catch (_) {
      return refuse('refused', 'The daemon\'s request is not base64url; nothing to confirm.');
    }
    const challenge = unhex(s.challenge);
    if (s.id !== id || !unhex(s.digest) || !challenge || !Number.isSafeInteger(s.expires) || s.expires < 0 || typeof s.state !== 'string') {
      return refuse('refused', 'The daemon\'s answer is malformed; nothing to confirm.');
    }
    const digest = await sha256(bytes);
    if (hex(digest) !== s.digest) return refuse('refused', 'The request\'s bytes do not match its digest; nothing to confirm.');
    let r;
    try {
      r = JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
    } catch (e) {
      return refuse('refused', `The request is not UTF-8 JSON (${e.message}); nothing to confirm.`);
    }
    const bad = unshown(r);
    if (bad) return refuse('refused', `This request has a field the page cannot show (${bad}); nothing to confirm here; confirm at the terminal.`);
    const m = message(digest, challenge, s.expires), h = await sha256(m);
    render(r, s, h);
    $('what').textContent = `Confirm request ${id}: ${r.verb} ${r.episode}; expires ${time(s.expires)}; message hash ${hex(h)}`;
    if (s.state !== 'pending') {
      setState('closed');
      $('status').textContent = `${CLOSED[s.state] || 'This request is ' + s.state}${s.state === 'approved' ? ` by ${s.approvedBy}: decision ${s.decision}` : ''}; nothing to confirm here.`;
      return;
    }
    $('save').href = URL.createObjectURL(new Blob([m], { type: 'application/octet-stream' }));
    $('save').download = `world-verb-${id}.msg`;
    $('sshline').textContent = `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-verb-${id}.msg`;
    $('status').textContent = 'Waiting for your key: save the message and sign it with your enrolled SSH key.';
    digestHex = s.digest;
    hash = h;
    group = (s.passkeys || []).find((g) => g.origin === location.origin) || null;
    $('go').hidden = !group;
    $('controls').hidden = false;
    buttons(true);
    setState('ready');
  }

  // submit POSTs one key format; the answer or the refusal stays shown while the card is read again.
  async function submit(format) {
    buttons(false);
    setState('posting');
    let res;
    try {
      res = await net('/surface/confirm/' + id, { method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(Object.assign({ v: 1, digest: digestHex }, format)) });
    } catch (e) {
      $('answer').textContent = `Confirmation outcome unknown: ${why(e)}. Reading the request again.`;
      return reread();
    }
    const b = res.body || {};
    $('answer').textContent = res.status === 200 && b.confirmed === true
      ? `Confirmed by ${b.approvedBy}: decision ${b.decision}; ${b.next}.`
      : `Not confirmed (HTTP ${res.status} ${b.error || ''}): ${b.reason || 'no reason given'}`;
    return reread();
  }
  const reread = () => { $('controls').hidden = true; return load(); };

  async function confirmWithPasskey() {
    buttons(false);
    setState('waiting');
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), TIMEOUT_MS);
    let cred;
    try {
      cred = await navigator.credentials.get({ publicKey: { challenge: hash, rpId: group.rp,
        allowCredentials: group.allow.map((a) => ({ type: 'public-key', id: unb64u(a) })),
        userVerification: 'required', timeout: TIMEOUT_MS }, signal: controller.signal });
    } catch (e) {
      $('answer').textContent = `Not confirmed: the browser made no passkey response (${e.name}: ${e.message}). Nothing was sent.`;
      return reread();
    } finally {
      clearTimeout(timer);
    }
    return submit({ webauthn: assertionJSON(cred) });
  }
  $('go').addEventListener('click', () => { confirmWithPasskey(); });
  $('post-sig').addEventListener('click', () => { submit({ sshsig: $('sig').value }); });

  load(); // the daemon will serve this page only at /surface/confirm/vr-<16 hex>
})();
