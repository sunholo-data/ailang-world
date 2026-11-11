'use strict';
// AILANG World's passkey ceremony page (row 201 A2, design §2.14). It reads the link's fragment,
// runs one WebAuthn create or get in this browser, and offers the result as a file to save. It
// makes no network request and holds no data: the CLI judges the file. Text is set with
// textContent only.
(() => {
  const TIMEOUT_MS = 240000; // [r2] the application's own bound: WebAuthn's timeout is only a hint
  const $ = (id) => document.getElementById(id);
  const setState = (s) => { document.body.dataset.state = s; };
  const say = (text, error) => {
    $('status').textContent = text;
    $('status').toggleAttribute('data-error', error);
  };

  const b64u = (buf) => {
    let s = '';
    for (const b of new Uint8Array(buf)) s += String.fromCharCode(b);
    return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
  };
  const unb64u = (s) => {
    if (!/^[A-Za-z0-9_-]*$/.test(s)) throw new Error('not base64url');
    return Uint8Array.from(atob(s.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
  };

  // The enrolment code (design §2.8): base32 of the first 16 bytes of SHA-256("world-enrolment-code/v1"
  // ‖ 0x00 ‖ SHA-256(clientDataJSON) ‖ SHA-256(attestationObject)), shown in five groups.
  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
  async function enrolmentCode(clientDataJSON, attestationObject) {
    const label = new TextEncoder().encode('world-enrolment-code/v1\0');
    const msg = new Uint8Array(label.length + 64);
    msg.set(label);
    msg.set(await sha256(clientDataJSON), label.length);
    msg.set(await sha256(attestationObject), label.length + 32);
    const h = (await sha256(msg)).slice(0, 16);
    const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';
    let bits = 0, n = 0, out = '';
    for (const b of h) {
      bits = (bits << 8) | b;
      n += 8;
      while (n >= 5) { out += alphabet[(bits >>> (n - 5)) & 31]; n -= 5; }
    }
    out += alphabet[(bits << (5 - n)) & 31];
    return [out.slice(0, 5), out.slice(5, 10), out.slice(10, 15), out.slice(15, 20), out.slice(20)].join('-');
  }

  // The §2.7 files, built here rather than by toJSON().
  const registrationJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
    response: { clientDataJSON: b64u(c.response.clientDataJSON), attestationObject: b64u(c.response.attestationObject) } });
  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });

  function call(f, signal) {
    const challenge = unb64u(f.challenge);
    if (f.op === 'create') {
      return navigator.credentials.create({ publicKey: { challenge, rp: { id: f.rp, name: 'AILANG World' },
        user: { id: new TextEncoder().encode(f.name), name: f.name, displayName: f.name },
        pubKeyCredParams: [{ type: 'public-key', alg: -7 }, { type: 'public-key', alg: -8 }],
        authenticatorSelection: { userVerification: 'required', residentKey: 'preferred' },
        attestation: 'none', timeout: TIMEOUT_MS }, signal });
    }
    return navigator.credentials.get({ publicKey: { challenge, rpId: f.rp,
      allowCredentials: f.allow.map((id) => ({ type: 'public-key', id: unb64u(id) })),
      userVerification: 'required', timeout: TIMEOUT_MS }, signal });
  }

  async function ceremony(f) {
    const go = $('go');
    go.disabled = true;
    setState('pending');
    say('Waiting for your passkey…', false);
    const controller = new AbortController();
    let expired = false;
    const timer = setTimeout(() => {
      expired = true;
      controller.abort();
      setState('expired');
      say(`Timed out: no passkey answer within ${TIMEOUT_MS / 60000} minutes. Nothing can be saved from this attempt; try again.`, true);
    }, TIMEOUT_MS);
    let cred = null, failure = null;
    try {
      cred = await call(f, controller.signal);
    } catch (e) {
      failure = e;
    } finally {
      clearTimeout(timer);
    }
    if (expired) { // [r2] a result after the timer is ignored: the attempt is over, nothing is saved
      go.disabled = false;
      return;
    }
    if (failure !== null) {
      setState('failed');
      say(`The browser did not make a passkey response (${failure.name}: ${failure.message}). Nothing was saved; try again.`, true);
      go.disabled = false;
      return;
    }
    const create = f.op === 'create';
    const json = JSON.stringify(create ? registrationJSON(cred) : assertionJSON(cred));
    if (create) {
      $('code').value = await enrolmentCode(cred.response.clientDataJSON, cred.response.attestationObject);
      $('code-box').hidden = false;
    }
    $('result').value = json;
    $('result-box').hidden = false;
    const save = $('save');
    save.href = URL.createObjectURL(new Blob([json + '\n'], { type: 'application/json' }));
    save.download = create ? `world-registration-${f.name}.json` : `world-assertion-${f.request}.json`;
    save.hidden = false;
    setState('done');
    say(create ? 'Passkey created. Save the file, then type the enrolment code at the terminal.'
      : 'Signed by your passkey. Save the file, then run the approve command the terminal printed.', false);
  }

  let f;
  try {
    f = JSON.parse(new TextDecoder().decode(unb64u(location.hash.slice(1))));
    if (f.v !== 1 || (f.op !== 'create' && f.op !== 'get')) throw new Error('not a version 1 create or get');
  } catch (e) {
    setState('refused');
    $('what').textContent = `This is not a World passkey link (${e.message}). Open the whole link the terminal printed.`;
    return;
  }
  if (f.origin !== location.origin || f.rp !== location.hostname) {
    setState('refused');
    $('what').textContent = `This link is for ${f.origin} (passkey site ${f.rp}): open it there, not at ${location.origin}.`;
    return;
  }
  $('what').textContent = f.op === 'create'
    ? `Create a passkey for approver ${f.name} at ${f.rp}`
    : `Approve request ${f.request}, expires ${new Date(f.expires * 1000).toLocaleString()}, message hash ${f.hash16}`;
  const go = $('go');
  go.textContent = f.op === 'create' ? 'Create passkey' : 'Approve with passkey';
  go.addEventListener('click', () => { ceremony(f); });
  go.hidden = false;
  go.disabled = false;
  setState('ready');
})();
