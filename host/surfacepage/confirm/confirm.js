'use strict';
// AILANG World's confirmation card (row 243 V3, design §2.5). It reads one stored verb request from
// the daemon, checks that its bytes hash to its digest, shows every field of the signed bytes through
// the review and signed details (and refuses a request with a field it cannot show), and builds the 86-byte world-verb/v1
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
  const CLOSED = new Map([
    ['approved', 'Approval recorded'],
    ['superseded', 'Superseded by a newer finish'],
    ['expired', 'Request expired']
  ]);
  // Match the broker's canonical finish scope. Review fields come from the scope,
  // since an argv flag may be omitted and the prompt includes free-form text.
  const FINISH_SCOPE = /^github-pr:([A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+)#remote=[^&#= \t\n\v\f\r]+&head=([^&#= \t\n\v\f\r]+)&sha=([0-9a-f]{40})&base=([^&#= \t\n\v\f\r]+)&title-sha256=[0-9a-f]{64}&body-sha256=[0-9a-f]{64}&label=ailang-world&draft=(true|false)&create-label=(true|false)&expires=[1-9][0-9]*(?![\s\S])/;

  const $ = (id) => document.getElementById(id);
  const setState = (s) => { document.body.dataset.state = s; };
  const status = (title, detail = '', error = false) => {
    $('status').textContent = title;
    $('what').textContent = detail;
    $('what').toggleAttribute('data-error', error);
  };
  const refuse = (state, detail) => {
    $('shown').hidden = true;
    $('overview').hidden = true;
    $('controls').hidden = true;
    $('receipt-actions').hidden = true;
    delete document.body.dataset.requestState;
    $('review-title').textContent = 'Review pull request';
    status(state === 'failed' ? 'Request unavailable' : 'Request refused', detail, true);
    setState(state);
  };
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
  };
  function render(r, s, h, scope) {
    const [, repository, head, sha, base, draft, createLabel] = scope;
    const mode = draft === 'true' ? 'Draft' : 'Ready for review';
    $('action').textContent = draft === 'true' ? 'Open a draft pull request' : 'Open a pull request';
    $('repository').textContent = repository;
    $('episode').textContent = r.episode;
    $('overview').hidden = false;
    $('review-title').textContent = r.title;
    $('head').textContent = head;
    $('head-commit').textContent = sha;
    $('base').textContent = base;
    $('pr-state').textContent = mode;
    $('label').textContent = createLabel === 'true' ? 'ailang-world (will be created)' : 'ailang-world';
    const paths = $('review-protected');
    paths.textContent = '';
    for (const path of r.protected) {
      const item = document.createElement('li');
      item.textContent = path;
      item.dir = 'auto';
      paths.append(item);
    }
    paths.hidden = r.protected.length === 0;
    $('no-protected').hidden = !paths.hidden;
    $('f-body').textContent = r.body;
    $('f-body-size').textContent = `(${new TextEncoder().encode(r.body).length} bytes)`;
    $('argv').textContent = quote(r.argv);
    $('prompt').textContent = r.prompt;
    const dl = $('signed');
    dl.textContent = '';
    row(dl, 'm-id', 'Request ID', s.id);
    row(dl, 'm-digest', 'Request digest (SHA-256)', s.digest);
    row(dl, 'm-challenge', 'Challenge', s.challenge);
    row(dl, 'm-expires', 'Expires', time(s.expires));
    row(dl, 'm-hash', 'Message hash (the passkey challenge)', hex(h));
    for (const f of CONFIRM_FIELDS) {
      const v = r[f.key];
      if (f.key === 'body') continue; // Full PR body is its own disclosure beside the review.
      row(dl, 'f-' + f.key, f.label, f.type === 'strings' ? v.join('\n') : v);
    }
    $('shown').hidden = false;
    document.body.dataset.requestState = s.state;
    $('review-facts').open = s.state === 'pending';
  }

  function verifiedEpisodeHref(href, episode) {
    if (typeof href !== 'string' || !href.startsWith('/surface?') || href.includes('#') || href.length > 512) return '';
    const q = new URLSearchParams(href.slice('/surface?'.length));
    if ([...q].length !== 2 || q.getAll('brane').length !== 1 || q.getAll('at').length !== 1) return '';
    const brane = q.get('brane'), at = q.get('at');
    const lane = /^([a-z][a-z0-9-]{0,31}):brane:([^\u0000-\u001f]+)(?![\s\S])/u.exec(brane);
    if (!lane || at !== lane[1] + ':episode:' + episode) return '';
    return '/surface?' + new URLSearchParams({brane, at});
  }

  function renderReceipt(s, episode, repository) {
    const episodeHref = verifiedEpisodeHref(s.episodeHref, episode);
    $('receipt-back').href = episodeHref || '/surface';
    $('receipt-back').textContent = episodeHref ? 'Return to episode' : 'Return to World';
    const pr = typeof s.prHref === 'string' && /^https:\/\/github\.com\/([A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+)\/pull\/([1-9][0-9]*)(?![\s\S])/.exec(s.prHref);
    const prHref = pr && pr[1].toLowerCase() === repository.toLowerCase() ? pr[0] : '';
    $('receipt-pr').hidden = !prHref;
    if (prHref) {
      $('receipt-pr').href = prHref;
      $('receipt-pr').textContent = `View PR #${pr[2]} on GitHub`;
    }
    $('receipt-actions').hidden = false;
    return prHref;
  }

  function approverMessage(id) {
    if (typeof id !== 'string' || !id) return '';
    const separator = id.indexOf(':');
    const method = id.slice(0, separator), name = id.slice(separator + 1);
    if (name && method === 'webauthn') return `${name} approved with a passkey.`;
    if (name && method === 'ssh') return `${name} approved with an SSH signature.`;
    return `Approved by ${id}.`;
  }

  const id = (location.pathname.match(/^\/surface\/confirm\/(vr-[0-9a-f]{16})$/) || [])[1];
  try {
    const theme = JSON.parse(localStorage.getItem('world.appearance') || '{}')?.theme;
    document.documentElement.dataset.theme = theme === 'night' || theme === 'slate' ||
      (theme === 'system' && matchMedia('(prefers-color-scheme: dark)').matches) ? 'slate' : 'paper';
  } catch (_) { /* Paper is the World default when preferences are unavailable. */ }
  if (location.hash === '#activity') {
    $('back').href = '/surface#activity';
    $('back').textContent = '← Back to Activity';
  } else if (location.hash) {
    const place = new URLSearchParams(location.hash.slice(1));
    const brane = place.get('brane'), at = place.get('at'), view = place.get('view') || '';
    if (brane && at && brane.length < 256 && at.length < 256) {
      const fragment = /^(heads(?:\/(?:files|run|trajectory))?|files|run|trajectory)(?![\s\S])/.test(view) ? '#' + view : '';
      $('back').href = '/surface?' + new URLSearchParams({brane, at}) + fragment;
      $('back').textContent = '← Return to selected state';
    }
  }
  let submittedReceipt = null;
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
    if (s.id !== id || !unhex(s.digest) || !challenge || !Number.isSafeInteger(s.expires) || s.expires < 0 || s.expires > 8640000000000 || typeof s.state !== 'string') {
      return refuse('refused', 'The daemon\'s answer is malformed; nothing to confirm.');
    }
    let digest;
    try {
      digest = await sha256(bytes);
    } catch (_) {
      return refuse('failed', 'This browser could not verify the request. Confirm at the terminal.');
    }
    if (hex(digest) !== s.digest) return refuse('refused', 'The request\'s bytes do not match its digest; nothing to confirm.');
    let r;
    try {
      r = JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
    } catch (e) {
      return refuse('refused', `The request is not UTF-8 JSON (${e.message}); nothing to confirm.`);
    }
    const bad = unshown(r);
    if (bad) return refuse('refused', `This request has a field the page cannot show (${bad}); nothing to confirm here; confirm at the terminal.`);
    const scope = FINISH_SCOPE.exec(r.scope);
    if (!scope || r.verb !== 'episode finish') {
      return refuse('refused', 'This request does not describe a recognizable pull request. Confirm at the terminal.');
    }
    const m = message(digest, challenge, s.expires), h = await sha256(m);
    render(r, s, h, scope);
    if (s.state !== 'pending') {
      let pullRequest = '';
      if (s.state === 'approved') {
        pullRequest = renderReceipt(s, r.episode, scope[1]);
        const draft = scope[5] === 'true';
        $('action').textContent = pullRequest ? (draft ? 'Draft pull request' : 'Pull request')
          : (draft ? 'Requested a draft pull request' : 'Requested a pull request');
        row($('signed'), 'approval-signer', 'Recorded approver', s.approvedBy);
        row($('signed'), 'approval-decision', 'Recorded decision', s.decision);
        if (submittedReceipt && submittedReceipt.approvedBy === s.approvedBy && submittedReceipt.decision === s.decision) $('answer').textContent = '';
      }
      const detail = s.state === 'approved'
        ? [approverMessage(s.approvedBy), pullRequest ? '' : 'The finish command reports the pull request outcome.'].filter(Boolean).join(' ')
        : s.state === 'superseded' ? 'Return to World for the latest request.'
        : s.state === 'expired' ? 'Return to the terminal to start a new finish request.'
        : `Request state: ${s.state}. Check the episode finish at the terminal.`;
      status(pullRequest ? 'Pull request opened' : CLOSED.get(s.state) || 'Request closed', detail);
      setState('closed');
      return;
    }
    $('save').href = URL.createObjectURL(new Blob([m], { type: 'application/octet-stream' }));
    $('save').download = `world-verb-${id}.msg`;
    $('sshline').textContent = `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-verb-${id}.msg`;
    digestHex = s.digest;
    hash = h;
    group = (s.passkeys || []).find((g) => g.origin === location.origin) || null;
    $('go').hidden = !group;
    if (!group) $('ssh-fallback').open = true;
    $('controls').hidden = false;
    const expiry = new Date(s.expires * 1000).toISOString().slice(11, 16);
    const fallback = group ? '' : ' Use an SSH signature, or rerun finish at the terminal without --confirm-on-surface.';
    status('Awaiting approval', 'Expires ' + expiry + ' UTC.' + fallback);
    buttons(true);
    setState('ready');
  }

  // submit POSTs one key format; the answer or the refusal stays shown while the card is read again.
  async function submit(format) {
    buttons(false);
    setState('posting');
    status('Submitting approval…');
    let res;
    try {
      res = await net('/surface/confirm/' + id, { method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(Object.assign({ v: 1, digest: digestHex }, format)) });
    } catch (e) {
      $('answer').textContent = `Confirmation outcome unknown: ${why(e)}. Reading the request again.`;
      return reread();
    }
    const b = res.body || {};
    submittedReceipt = res.status === 200 && b.confirmed === true ? b : null;
    $('answer').textContent = submittedReceipt
      ? `Approval recorded by ${b.approvedBy}. Decision ${b.decision}. Check the episode finish at the terminal for the PR outcome.`
      : `Not confirmed (HTTP ${res.status} ${b.error || ''}): ${b.reason || 'no reason given'}`;
    return reread();
  }
  const reread = () => {
    $('controls').hidden = true;
    status('Checking request status…');
    return load();
  };

  async function confirmWithPasskey() {
    buttons(false);
    setState('waiting');
    status('Waiting for your passkey…');
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
  $('go').addEventListener('click', () => {
    $('answer').textContent = '';
    confirmWithPasskey();
    $('status').focus();
  });
  $('post-sig').addEventListener('click', () => {
    $('answer').textContent = '';
    submit({ sshsig: $('sig').value });
    $('status').focus();
  });

  load(); // the daemon will serve this page only at /surface/confirm/vr-<16 hex>
})();
