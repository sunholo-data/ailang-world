'use strict';
// AILANG World's request card (row 201 A3, design §2.7). It reads one mint request from the daemon,
// checks that the stored bytes hash to the request's digest, shows EVERY signed field of those bytes
// (and refuses a request with a field it cannot show), computes the 90-byte message and its hash
// itself (never a hash from the server), and approves with a passkey or a pasted SSH signature by a
// POST the daemon judges with the A2 law. Text is set with textContent only.
(() => {
  const TIMEOUT_MS = 240000; // A2's ceremony bound (A2 §2.14 [r2])
  const NET_MS = 10000; // [r2] every GET and POST, its body included
  const RETRIES = 3, RETRY_MS = 1000; // [r2] only after an explicit 503 Busy
  const NAMESPACE = 'world-approval@sunholo';
  // Every JSON key path of approval.MintRequest ([] = each element), with its label, group and type.
  // TestApprovePageRendersEveryCanonicalField compares this list with the Go struct.
  const FIELDS = /* FIELDS */[
    {"path": "episode", "label": "Episode", "group": "what", "type": "string"},
    {"path": "agent", "label": "Agent", "group": "what", "type": "string"},
    {"path": "model", "label": "Model", "group": "what", "type": "string"},
    {"path": "task", "label": "Task file", "group": "what", "type": "string"},
    {"path": "grants[].effect", "label": "Grant effect", "group": "what", "type": "string"},
    {"path": "grants[].scope", "label": "Grant scope", "group": "what", "type": "string"},
    {"path": "grants[].budget", "label": "Grant budget", "group": "what", "type": "int"},
    {"path": "ttl_s", "label": "Session lifetime (s)", "group": "what", "type": "int"},
    {"path": "worktree", "label": "Worktree", "group": "where", "type": "string"},
    {"path": "reuse", "label": "Worktree use", "group": "where", "type": "bool", "no": "creates a new worktree", "yes": "uses an existing worktree"},
    {"path": "degraded", "label": "Mode", "group": "where", "type": "bool", "no": "git worktree", "yes": "NON-GIT, DEGRADED"},
    {"path": "repo", "label": "Repository", "group": "where", "type": "string"},
    {"path": "branch", "label": "Branch", "group": "where", "type": "string"},
    {"path": "git", "label": "Git binary", "group": "where", "type": "string"},
    {"path": "workspace_root", "label": "Workspace root", "group": "where", "type": "string"},
    {"path": "profile.project", "label": "Exec profile", "group": "bound", "type": "string"},
    {"path": "profile.digest", "label": "Profile digest", "group": "bound", "type": "string"},
    {"path": "profiles_dir", "label": "Profiles directory", "group": "bound", "type": "string"},
    {"path": "world_db", "label": "World store", "group": "stores", "type": "string"},
    {"path": "credentials_db", "label": "Credentials store", "group": "stores", "type": "string"},
    {"path": "session_file", "label": "Token written to", "group": "stores", "type": "string"},
    {"path": "mcp_config", "label": "MCP config", "group": "stores", "type": "string"},
    {"path": "daemon.addr", "label": "Daemon address", "group": "daemon", "type": "string"},
    {"path": "daemon.vcs_revision", "label": "Daemon revision", "group": "daemon", "type": "string"},
    {"path": "daemon.vcs_modified", "label": "Daemon build", "group": "daemon", "type": "bool", "no": "clean build", "yes": "modified build"},
    {"path": "daemon.executable_sha256", "label": "Daemon binary SHA-256", "group": "daemon", "type": "string"},
    {"path": "daemon.go_version", "label": "Daemon Go version", "group": "daemon", "type": "string"},
    {"path": "schema", "label": "Request schema", "group": "daemon", "type": "string"}
  ]/* /FIELDS */;
  const GROUPS = [['what', 'What runs'], ['where', 'Where'], ['bound', 'Bound to'], ['stores', 'Stores and files'], ['daemon', 'Requesting daemon']];
  const TYPES = { string: 'string', bool: 'boolean', int: 'number' };

  const $ = (id) => document.getElementById(id);
  const setState = (s) => { document.body.dataset.state = s; };
  const say = (text, error) => {
    $('status').textContent = text;
    $('status').toggleAttribute('data-error', !!error);
  };
  const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
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
  const unhex = (s) => Uint8Array.from(s.match(/../g), (x) => parseInt(x, 16));
  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
  // A2 §2.7's assertion file, built here rather than by toJSON() (as passkey.js does).
  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });

  // [r2] Each GET and POST, response body included, has its own 10 s deadline; leaving the page
  // cancels every outstanding one.
  const outstanding = new Set();
  addEventListener('pagehide', () => { for (const c of outstanding) c.abort(); });
  async function net(url, init) {
    const c = new AbortController();
    outstanding.add(c);
    const timer = setTimeout(() => c.abort(), NET_MS);
    try {
      const r = await fetch(url, Object.assign({ cache: 'no-store', signal: c.signal }, init));
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

  // leaves is every leaf of v: its key path pattern ([] for an element), its concrete path (the
  // element's index) and its value. An empty object or array has none, so its fields are missing.
  function leaves(v, pattern, concrete, out, index) {
    if (Array.isArray(v)) {
      v.forEach((e, i) => leaves(e, pattern + '[]', concrete + '.' + i, out, i));
    } else if (v !== null && typeof v === 'object') {
      for (const k of Object.keys(v)) leaves(v[k], pattern ? pattern + '.' + k : k, concrete ? concrete + '.' + k : k, out, index);
    } else {
      out.push({ pattern, concrete, value: v, index });
    }
    return out;
  }
  // unshown is the first key path the card cannot show (unlisted, of another type, or missing), or ''.
  // A request that is not an object has a leaf at '' or under '[]', which no FIELDS entry names.
  function unshown(ls) {
    const byPath = new Map(FIELDS.map((f) => [f.path, f]));
    const seen = new Set();
    for (const l of ls) {
      const f = byPath.get(l.pattern);
      if (!f || typeof l.value !== TYPES[f.type] || (f.type === 'int' && !Number.isSafeInteger(l.value))) return l.pattern;
      seen.add(l.pattern);
    }
    for (const f of FIELDS) if (!seen.has(f.path)) return f.path;
    return '';
  }
  const shown = (f, v) => (f.type === 'bool' ? (v ? f.yes : f.no) : String(v));

  function render(ls) {
    const card = $('card');
    for (const [group, title] of GROUPS) {
      const section = document.createElement('section');
      section.dataset.group = group;
      const h = document.createElement('h2');
      h.textContent = title;
      const dl = document.createElement('dl');
      for (const f of FIELDS.filter((x) => x.group === group)) {
        for (const l of ls.filter((x) => x.pattern === f.path)) {
          const dt = document.createElement('dt');
          dt.textContent = f.label + (l.index === undefined ? '' : ' ' + (l.index + 1)) + (f.path === 'ttl_s' ? ` (${(l.value / 3600).toFixed(1)} h)` : '');
          const dd = document.createElement('dd');
          dd.id = 'f-' + l.concrete;
          dd.textContent = shown(f, l.value);
          dl.append(dt, dd);
        }
      }
      section.append(h, dl);
      card.append(section);
    }
    card.hidden = false;
  }

  const id = (location.pathname.match(/^\/surface\/approve\/(mr-[0-9a-f]{16})$/) || [])[1];
  let card = null, digest = '', message = null, challenge = null, group = null, again = null;
  const buttons = (on) => { $('go').disabled = !on; $('post-sig').disabled = !on; };
  const offer = (label, action) => { again = action; $('again').textContent = label; $('again').hidden = false; $('again').disabled = false; };
  const refuse = (text) => { setState('refused'); $('what').textContent = text; };

  async function load() {
    $('again').hidden = true;
    setState('loading');
    let res;
    try {
      res = await net('/surface/requests/' + id, { headers: { Accept: 'application/json' } });
    } catch (e) {
      setState('failed');
      say(`The request could not be read (${why(e)}). Nothing to approve yet.`, true);
      return offer('Read it again', load);
    }
    if (res.status !== 200 || !res.body) {
      setState('failed');
      say(`The daemon did not give the request (HTTP ${res.status}${res.body && res.body.reason ? ': ' + res.body.reason : ''}).`, true);
      return offer('Read it again', load);
    }
    card = res.body;
    let bytes, obj;
    try {
      bytes = unb64u(card.request);
    } catch (e) {
      return refuse('The daemon\'s request is not base64url; nothing to approve.');
    }
    if (card.id !== id || hex(await sha256(bytes)) !== card.digest) {
      return refuse('The daemon\'s request does not match its digest; nothing to approve.');
    }
    try {
      obj = JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
    } catch (e) {
      return refuse(`The daemon's request is not UTF-8 JSON (${e.message}); nothing to approve.`);
    }
    const ls = leaves(obj, '', '', []);
    const bad = unshown(ls);
    if (bad) {
      return refuse(`This request has a field the page cannot show (${bad}); nothing to approve; approve it at a terminal.`);
    }
    if (!/^[0-9a-f]{64}$/.test(card.challenge) || !Number.isSafeInteger(card.expires) || card.expires < 0) {
      return refuse('The daemon\'s challenge or expiry is malformed; nothing to approve.');
    }
    // The 90-byte message (design §2.7 step 3): "world-approval/v1" ‖ 0x00 ‖ digest ‖ nonce ‖ u64be(expires).
    digest = card.digest;
    message = new Uint8Array(90);
    message.set(new TextEncoder().encode('world-approval/v1'), 0);
    message.set(unhex(card.digest), 18);
    message.set(unhex(card.challenge), 50);
    new DataView(message.buffer).setBigUint64(82, BigInt(card.expires));
    challenge = await sha256(message);
    render(ls);
    $('what').textContent = `Approve request ${id} for episode ${obj.episode}, expires ${new Date(card.expires * 1000).toLocaleString()}, message hash ${hex(challenge.slice(0, 8))}`;
    if (card.consumed) {
      setState('used');
      say(`Already used: approved by ${card.approvedBy}.`, false);
      return;
    }
    group = (card.passkeys || []).find((p) => p.origin === location.origin) || null;
    if (group) {
      $('go').hidden = false;
    } else {
      const others = (card.passkeys || []).map((p) => p.origin);
      $('origin-note').textContent = others.length ? `No passkey is enrolled for ${location.origin}: open this card at ${others.join(' or ')}.` : 'No passkey is enrolled; approve with an SSH key below.';
      $('origin-note').hidden = false;
    }
    const save = $('save');
    save.href = URL.createObjectURL(new Blob([message], { type: 'application/octet-stream' }));
    save.download = `world-approval-${id}.msg`;
    save.hidden = false;
    $('sshline').textContent = `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-approval-${id}.msg`;
    $('ssh').hidden = false;
    buttons(true);
    setState('ready');
  }

  // status reads the request's state from the decision list: one GET, under its deadline.
  async function status() {
    $('again').disabled = true;
    let res;
    try {
      res = await net('/surface/requests.json', { headers: { Accept: 'application/json' } });
    } catch (e) {
      say(`The status could not be read (${why(e)}).`, true);
      return offer('Check again', status);
    }
    const b = res.body || {};
    const p = (b.pending || []).find((x) => x.id === id);
    $('again').hidden = true;
    if (p && p.recorded) {
      setState('recorded');
      say(`Recorded: approved by ${p.recorded}, waiting for the minter (approve ${id} --await).`, false);
    } else if (p) {
      setState('ready');
      say('Not recorded: the request is still pending with no approval. You may approve it again.', true);
      buttons(true);
    } else {
      say(res.status === 200 ? 'The request is no longer pending (used or expired): read the card again.' : `The status could not be read (HTTP ${res.status}).`, true);
      offer('Check again', status);
    }
  }

  async function post(body) {
    buttons(false);
    setState('posting');
    say('Recording the approval…', false);
    for (let attempt = 0; ; attempt++) {
      let res;
      try {
        res = await net('/surface/approve/' + id, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
      } catch (e) {
        // [r2] No answer: the daemon may or may not have recorded it. Never retried, never another format.
        setState('unknown');
        say(`Approval recording outcome unknown: ${why(e)}. It may or may not have been recorded; nothing was minted by this page. Check its status before approving again.`, true);
        return offer('Check status', status);
      }
      const b = res.body || {};
      if (res.status === 503 && b.error === 'Busy' && attempt < RETRIES) {
        say(`The credentials store is busy; trying again (${attempt + 1} of ${RETRIES})…`, false);
        await sleep(RETRY_MS);
        continue;
      }
      if (res.status === 200) {
        setState('recorded');
        say(`Approved by ${b.approvedBy}: ${b.next}.`, false);
        return;
      }
      setState('rejected');
      say(`Not recorded (HTTP ${res.status} ${b.error || ''}): ${b.reason || 'no reason given'}`, true);
      buttons(true);
      return;
    }
  }

  async function approvePasskey() {
    buttons(false);
    setState('pending');
    say('Waiting for your passkey…', false);
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
      cred = await navigator.credentials.get({ publicKey: { challenge, rpId: location.hostname,
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
    await post({ v: 1, digest, webauthn: assertionJSON(cred) });
  }

  $('go').addEventListener('click', () => { approvePasskey(); });
  $('post-sig').addEventListener('click', () => { post({ v: 1, digest, sshsig: $('sig').value }); });
  $('again').addEventListener('click', () => { if (again) again(); });
  load(); // the daemon serves this page only at /surface/approve/mr-<16 hex>
})();
