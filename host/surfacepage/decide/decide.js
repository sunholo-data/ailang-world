'use strict';
// AILANG World's decision card (row 201 A4, design §2.8). It reads one decision subject from the
// daemon, checks that the request's and the packet's bytes hash to their refs and that the packet
// carries this request, shows EVERY field of both objects (and refuses an object with a field it
// cannot show), builds the 91-byte decision message itself (never a hash from the server), and
// approves or denies with a passkey or a pasted SSH signature by a POST the daemon judges with the
// decision law. Text is set with textContent only.
(() => {
  const TIMEOUT_MS = 240000; // A2's ceremony bound (A2 §2.14 [r2])
  const NET_MS = 10000; // A3 §2.7 [r2]: every GET and POST, its body included
  const FRESH_S = 60; // re-read the subject before signing when less than this remains of its window
  const PREFIX = 'world-decision/v1', NAMESPACE = 'world-decision@sunholo';
  const VERBS = { approve: 1, deny: 2 }; // approval.VerbApprove, approval.VerbDeny
  // Every key of the request (world/approval-request/v1) and of the packet (world/decision-packet/v1),
  // with its label and JSON type. TestDecidePageFieldTable compares this list with the stored wires.
  const DECISION_FIELDS = /* DECISION_FIELDS */[
    {"object": "request", "path": "effect", "label": "Effect", "type": "string"},
    {"object": "request", "path": "scope", "label": "Scope", "type": "string"},
    {"object": "request", "path": "cost", "label": "Cost", "type": "int"},
    {"object": "request", "path": "requester", "label": "Requested by", "type": "string"},
    {"object": "request", "path": "now", "label": "Requested at", "type": "time"},
    {"object": "packet", "path": "packetHash", "label": "Packet hash", "type": "string"},
    {"object": "packet", "path": "proposalHash", "label": "Proposal", "type": "string"},
    {"object": "packet", "path": "requestRef", "label": "Request", "type": "string"},
    {"object": "packet", "path": "createdAt", "label": "Created at", "type": "time"},
    {"object": "packet", "path": "deadlineAt", "label": "Deadline", "type": "time"},
    {"object": "packet", "path": "escalationsRemaining", "label": "Escalations remaining", "type": "int"},
    {"object": "packet", "path": "policy", "label": "Policy", "type": "string"}
  ]/* /DECISION_FIELDS */;
  const GROUPS = [['request', 'The request'], ['packet', 'The decision packet']];

  const $ = (id) => document.getElementById(id);
  const setState = (s) => { document.body.dataset.state = s; };
  const say = (text, error) => {
    $('status').textContent = text;
    $('status').toggleAttribute('data-error', !!error);
  };
  const b64u = (buf) => {
    let s = '';
    for (const b of new Uint8Array(buf)) s += String.fromCharCode(b);
    return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
  };
  const unb64u = (s) => {
    if (typeof s !== 'string' || !/^[A-Za-z0-9_-]*$/.test(s)) throw new Error('not base64url');
    return Uint8Array.from(atob(s.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
  };
  const hex = (b) => Array.from(b, (x) => x.toString(16).padStart(2, '0')).join('');
  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
  const Verb = (verb) => (verb === 'deny' ? 'Deny' : 'Approve');
  // A2 §2.7's assertion file, built here rather than by toJSON() (as approve.js does).
  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });

  const outstanding = new Set();
  addEventListener('pagehide', () => { for (const c of outstanding) c.abort(); });
  async function net(url, init) {
    const c = new AbortController();
    outstanding.add(c);
    const timer = setTimeout(() => c.abort(), NET_MS);
    try {
      const r = await fetch(url, Object.assign({ cache: 'no-store', credentials: 'omit', signal: c.signal }, init));
      const text = await r.text();
      let body = null;
      try { body = JSON.parse(text); } catch (_) { body = null; }
      return { status: r.status, body };
    } finally {
      clearTimeout(timer);
      outstanding.delete(c);
    }
  }
  const why = (e) => (e && e.name === 'AbortError' ? `no answer within ${NET_MS / 1000} s` : String(e && e.message || e));

  // unshown is the first key of object name the card cannot show (unlisted, of another type, or
  // missing), or ''. A value that is not a JSON object is itself unshowable.
  function unshown(name, obj) {
    if (obj === null || typeof obj !== 'object' || Array.isArray(obj)) return name;
    const fields = DECISION_FIELDS.filter((f) => f.object === name);
    for (const k of Object.keys(obj)) {
      const f = fields.find((x) => x.path === k);
      if (!f || typeof obj[k] !== (f.type === 'string' ? 'string' : 'number') || (f.type !== 'string' && !Number.isSafeInteger(obj[k]))) return name + '.' + k;
    }
    for (const f of fields) if (!Object.prototype.hasOwnProperty.call(obj, f.path)) return name + '.' + f.path;
    return '';
  }
  const shown = (f, v) => (f.type === 'time' ? `${new Date(v * 1000).toISOString()} (${v})` : String(v));
  const parse = (bytes) => JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));

  // message is the 91-byte decision message (design §2.3): PREFIX ‖ 0x00 ‖ verb ‖ P ‖ Q ‖ u64be(expiry).
  function message(verb, P, Q, expiry) {
    const m = new Uint8Array(91);
    m.set(new TextEncoder().encode(PREFIX), 0);
    m[18] = VERBS[verb];
    m.set(P, 19);
    m.set(Q, 51);
    new DataView(m.buffer).setBigUint64(83, BigInt(expiry));
    return m;
  }

  function render(objects, entry) {
    const card = $('card');
    card.textContent = '';
    for (const [name, title] of GROUPS) {
      const section = document.createElement('section');
      section.dataset.group = name;
      const h = document.createElement('h2');
      h.textContent = title;
      const dl = document.createElement('dl');
      for (const f of DECISION_FIELDS.filter((x) => x.object === name)) {
        const dt = document.createElement('dt');
        dt.textContent = f.label;
        const dd = document.createElement('dd');
        dd.id = 'f-' + name + '.' + f.path;
        dd.textContent = shown(f, objects[name][f.path]);
        dl.append(dt, dd);
      }
      if (name === 'packet') {
        const dt = document.createElement('dt');
        dt.textContent = 'Committed at entry';
        const dd = document.createElement('dd');
        dd.id = 'f-entry';
        const a = document.createElement('a');
        a.id = 'proposal';
        a.href = '/workbench?object=' + encodeURIComponent(objects.packet.proposalHash);
        a.textContent = `${entry} · the proposal`;
        dd.append(a);
        dl.append(dt, dd);
      }
      section.append(h, dl);
      card.append(section);
    }
    card.hidden = false;
  }

  const ref = (location.pathname.match(/^\/surface\/decide\/(sha256:[0-9a-f]{64})$/) || [])[1];
  const ref16 = ref ? ref.slice(7, 23) : '';
  let subject = null, P = null, Q = null, expiry = 0, messages = {}, group = null, renderedAt = 0, again = null;
  const buttons = (on) => { for (const b of ['go', 'deny', 'post-sig', 'post-sig-deny']) $(b).disabled = !on; };
  const offer = (label, action) => { again = action; $('again').textContent = label; $('again').hidden = false; $('again').disabled = false; };
  const refuse = (text) => { setState('refused'); $('what').textContent = text; };

  // load reads the subject; keep leaves the status line as it is (after a decision).
  async function load(keep) {
    $('again').hidden = true;
    for (const id of ['go', 'deny', 'ssh', 'save', 'save-deny', 'origin-note']) $(id).hidden = true;
    buttons(false);
    setState('loading');
    let res;
    try {
      res = await net('/surface/decisions/' + ref, { headers: { Accept: 'application/json' } });
    } catch (e) {
      setState('failed');
      say(`The request could not be read (${why(e)}). Nothing to decide yet.`, true);
      return offer('Read it again', load);
    }
    if (res.status !== 200 || !res.body) {
      setState('failed');
      say(`The daemon did not give the request (HTTP ${res.status}${res.body && res.body.reason ? ': ' + res.body.reason : ''}).`, true);
      return offer('Read it again', load);
    }
    const s = res.body;
    let rb, pb, objects;
    try {
      rb = unb64u(s.requestBytes);
      pb = unb64u(s.packetBytes);
    } catch (e) {
      return refuse('The daemon\'s objects are not base64url; nothing to decide.');
    }
    const q = await sha256(rb), p = await sha256(pb);
    if (s.request !== ref || 'sha256:' + hex(q) !== ref || 'sha256:' + hex(p) !== s.packet) {
      return refuse('The daemon\'s objects do not match their hashes; nothing to decide.');
    }
    try {
      objects = { request: parse(rb), packet: parse(pb) };
    } catch (e) {
      return refuse(`The daemon's objects are not UTF-8 JSON (${e.message}); nothing to decide.`);
    }
    const bad = unshown('request', objects.request) || unshown('packet', objects.packet);
    if (bad) {
      return refuse(`This request has a field the page cannot show (${bad}); nothing to decide; decide it at a terminal: ailang-worldd approve ${ref}.`);
    }
    if (objects.packet.requestRef !== ref) {
      return refuse('The packet does not carry this request; nothing to decide.');
    }
    if (!Number.isSafeInteger(s.now) || !s.window || !Number.isSafeInteger(s.window.max) || s.window.max < 0) {
      return refuse('The daemon\'s decision window is malformed; nothing to decide.');
    }
    subject = s; P = p; Q = q; expiry = s.window.max; renderedAt = performance.now();
    messages = { approve: message('approve', P, Q, expiry), deny: message('deny', P, Q, expiry) };
    const h = { approve: await sha256(messages.approve), deny: await sha256(messages.deny) };
    render(objects, s.entry);
    const r = objects.request;
    $('what').textContent = `Decide request ${ref}: ${r.effect} for ${r.scope}, requested by ${r.requester}; a decision made here expires ${new Date(expiry * 1000).toLocaleString()}; message hash approve ${hex(h.approve.slice(0, 8))}, deny ${hex(h.deny.slice(0, 8))}`;
    if (s.decision) {
      setState('decided');
      if (!keep) say(`Already decided: ${s.decision} by ${s.decidedBy}.`, false);
      return;
    }
    if (s.reserved) {
      setState('reserved');
      say(`A ${s.reserved} approval: use ${s.reserved === 'publish' ? 'world-publish approve' : 'ailang-worldd episode finish'} at a terminal. ${s.refusal}`, true);
      return;
    }
    if (s.refusal) {
      setState('refused');
      say(`This request cannot be decided here: ${s.refusal}`, true);
      return;
    }
    group = (s.passkeys || []).find((g) => g.origin === location.origin) || null;
    if (group) {
      $('go').hidden = false;
      $('deny').hidden = false;
    } else {
      const others = (s.passkeys || []).map((g) => g.origin);
      $('origin-note').textContent = others.length ? `No passkey is enrolled for ${location.origin}: open this card at ${others.join(' or ')}.` : 'No passkey is enrolled; decide with an SSH key below.';
      $('origin-note').hidden = false;
    }
    for (const [id, verb] of [['save', 'approve'], ['save-deny', 'deny']]) {
      $(id).href = URL.createObjectURL(new Blob([messages[verb]], { type: 'application/octet-stream' }));
      $(id).download = `world-decision-${verb}-${ref16}.msg`;
      $(id).hidden = false;
    }
    $('sshline').textContent = ['approve', 'deny'].map((v) => `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-decision-${v}-${ref16}.msg`).join('\n');
    $('ssh').hidden = false;
    buttons(true);
    setState('ready');
  }

  // fresh re-reads the subject once when less than FRESH_S seconds remain of its window.
  async function fresh() {
    if (subject.window.max - subject.now - (performance.now() - renderedAt) / 1000 >= FRESH_S) return true;
    await load(false);
    return document.body.dataset.state === 'ready';
  }

  async function post(verb, format) {
    buttons(false);
    setState('posting');
    say(`${Verb(verb)} request ${ref}: sending the decision…`, false);
    let res;
    try {
      res = await net('/surface/decide/' + ref, { method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(Object.assign({ v: 1, verb, packet: hex(P), request: hex(Q), expiry }, format)) });
    } catch (e) {
      setState('unknown');
      say(`Decision outcome unknown: ${why(e)}. It may or may not have been decided; read the request again before deciding again.`, true);
      return offer('Read it again', load);
    }
    const b = res.body || {};
    if (res.status === 200) {
      say(`${b.decided === 'deny' ? 'Denied' : 'Approved'} by ${b.decidedBy}: decision ${b.decision}.`, false);
      return load(true);
    }
    setState('rejected');
    say(`Not decided (HTTP ${res.status} ${b.error || ''}): ${b.reason || 'no reason given'}`, true);
    buttons(true);
  }

  async function decideWithPasskey(verb) {
    buttons(false);
    if (!(await fresh())) return;
    setState('pending');
    say(`${Verb(verb)} request ${ref}: waiting for your passkey…`, false);
    const challenge = await sha256(messages[verb]);
    const controller = new AbortController();
    let expired = false;
    const timer = setTimeout(() => {
      expired = true;
      controller.abort();
      setState('expired');
      say(`Timed out: no passkey answer within ${TIMEOUT_MS / 60000} minutes. Nothing was sent; try again.`, true);
      buttons(true);
    }, TIMEOUT_MS);
    let cred = null, failure = null;
    try {
      cred = await navigator.credentials.get({ publicKey: { challenge, rpId: group.rp,
        allowCredentials: group.allow.map((a) => ({ type: 'public-key', id: unb64u(a) })),
        userVerification: 'required', timeout: TIMEOUT_MS }, signal: controller.signal });
    } catch (e) {
      failure = e;
    } finally {
      clearTimeout(timer);
    }
    if (expired) return; // a result after the timer is ignored: nothing is sent
    if (failure !== null) {
      setState('failed');
      say(`The browser did not make a passkey response (${failure.name}: ${failure.message}). Nothing was sent; try again.`, true);
      buttons(true);
      return;
    }
    await post(verb, { webauthn: assertionJSON(cred) });
  }

  $('go').addEventListener('click', () => { decideWithPasskey('approve'); });
  $('deny').addEventListener('click', () => { decideWithPasskey('deny'); });
  $('post-sig').addEventListener('click', () => { post('approve', { sshsig: $('sig').value }); });
  $('post-sig-deny').addEventListener('click', () => { post('deny', { sshsig: $('sig').value }); });
  $('again').addEventListener('click', () => { if (again) again(); });
  load(false); // the daemon serves this page only at /surface/decide/sha256:<64 hex>
})();
