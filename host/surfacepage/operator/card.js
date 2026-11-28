(() => {
  'use strict';
  const cardPath = /^\/surface\/operator-actions\/(oa-[0-9a-f]{32})(?![\s\S])/;
  const episodeID = /^[a-z0-9][a-z0-9-]{0,63}(?![\s\S])/;
  const hex64 = /^[0-9a-f]{64}(?![\s\S])/;
  const hexBytes = /^(?:[0-9a-f]{2})+(?![\s\S])/;
  const pathID = path => typeof path === 'string' ? cardPath.exec(path)?.[1] : undefined;
  const states = new Set(['requested', 'confirmed', 'claimed', 'running', 'awaiting-pr-confirmation',
    'succeeded', 'failed', 'expired', 'cancelled', 'uncertain', 'launched', 'stopped', 'capped', 'spawn-refused']);
  function statusFor(s, path) {
    if (!s || s.v !== 1 || s.id !== pathID(path) || s.href !== path || s.statusHref !== path + '/status.json' ||
        typeof s.episode !== 'string' || !episodeID.test(s.episode) || !['start', 'finish'].includes(s.kind) ||
        !states.has(s.state) || ![s.now, s.createdAt, s.expiresAt].every(Number.isSafeInteger) ||
        ![s.terminalAt, s.runnerSeenAt].every(n => n === undefined || Number.isSafeInteger(n)) ||
        (s.reason !== undefined && typeof s.reason !== 'string')) return null;
    const report = '/v1/episodes/' + s.episode + '/report';
    return { ...s, reportHref: s.reportHref === report ? report : '' };
  }
  function ceremonyFor(c, path) {
    return c && c.post === path + '/confirm' && typeof c.digest === 'string' && hex64.test(c.digest) &&
      typeof c.messageHex === 'string' && hexBytes.test(c.messageHex) ? c : null;
  }
  async function passkeyProof(c, origin, hostname, crypto, credentials, signal) {
    const group = Array.isArray(c.passkeys) && c.passkeys.find(k => k.origin === origin && k.rp === hostname);
    if (!group || !Array.isArray(group.allow) || group.allow.length === 0) throw new Error('No passkey for this origin');
    const challenge = await crypto.subtle.digest('SHA-256', Uint8Array.from(c.messageHex.match(/../g), x => parseInt(x, 16)));
    const allowCredentials = group.allow.map(id => ({ type: 'public-key', id: Uint8Array.from(atob(id.replace(/-/g, '+').replace(/_/g, '/')), ch => ch.charCodeAt(0)) }));
    const key = await credentials.get({ publicKey: { challenge, rpId: group.rp, allowCredentials, userVerification: 'required', timeout: 15000 }, signal });
    return { webauthn: key.toJSON() };
  }
  if (typeof module !== 'undefined') module.exports = { pathID, statusFor, ceremonyFor, passkeyProof };
  if (typeof document === 'undefined') return;

  const $ = id => document.getElementById(id);
  const path = location.pathname;
  if (!pathID(path)) return;
  let ceremony = null;
  try { ceremony = ceremonyFor(JSON.parse($('operator-ceremony')?.textContent), path); } catch (_) { /* status stays readable */ }
  let pending = false, polling, last = $('operator-progress').textContent, state = $('operator-state').dataset.state;
  const controls = $('operator-controls'), progress = $('operator-progress'), answer = $('operator-answer');
  function buttons() {
    controls.hidden = !ceremony || state !== 'requested';
    $('operator-key').hidden = !Array.isArray(ceremony?.passkeys) ||
      !ceremony.passkeys.some(k => k.origin === location.origin && k.rp === location.hostname);
    for (const id of ['operator-key', 'operator-sign', 'operator-refresh']) $(id).disabled = pending;
  }
  async function read(url, options) {
    const res = await fetch(url, { ...options, redirect: 'error', credentials: 'omit', cache: 'no-store',
      signal: AbortSignal.timeout(15000) });
    return { code: res.status, body: await res.json() };
  }
  async function refresh() {
    if (pending) return;
    pending = true;
    clearTimeout(polling);
    buttons();
    try {
      const { code, body } = await read(path + '/status.json');
      const s = code === 200 && statusFor(body, path);
      if (!s) throw new Error('Invalid status');
      state = s.state;
      $('operator-state').textContent = state;
      last = progress.textContent = [s.reason, 'Status at ' + s.now, 'Requested ' + s.createdAt, 'Expires ' + s.expiresAt,
        s.runnerSeenAt && 'Runner seen ' + s.runnerSeenAt, s.terminalAt && 'Completed ' + s.terminalAt].filter(Boolean).join(' · ');
      $('operator-report').hidden = !s.reportHref;
      if (s.reportHref) $('operator-report').href = s.reportHref;
      else $('operator-report').removeAttribute('href');
    } catch (_) {
      progress.textContent = 'Disconnected. ' + last + '. Use Refresh status to reconnect.';
    } finally {
      pending = false;
      buttons();
      polling = setTimeout(refresh, 5000);
    }
  }
  async function confirm(proof) {
    if (pending || !ceremony || state !== 'requested') return;
    pending = true;
    buttons();
    answer.textContent = 'Waiting for confirmation…';
    answer.focus();
    try {
      const artifact = await proof();
      const { code, body } = await read(path + '/confirm', { method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ v: 1, digest: ceremony.digest, ...artifact }) });
      answer.textContent = code === 200 && statusFor(body, path) && ['confirmed', 'already-confirmed'].includes(body.result) && typeof body.approvedBy === 'string' && body.approvedBy !== ''
        ? 'Confirmed by ' + body.approvedBy
        : 'Not confirmed (HTTP ' + code + '): ' + (typeof body?.error === 'string' ? body.error + ' ' : '') +
          (typeof body?.reason === 'string' ? body.reason : 'unreadable confirmation response');
    } catch (_) {
      answer.textContent = 'Confirmation unavailable. Read the action status before trying again.';
    } finally {
      pending = false;
      await refresh();
    }
  }
  $('operator-key').addEventListener('click', () => confirm(() => passkeyProof(ceremony, location.origin, location.hostname,
    crypto, navigator.credentials, AbortSignal.timeout(15000))));
  $('operator-sign').addEventListener('click', () => confirm(() => ({ sshsig: $('operator-signature').value })));
  $('operator-refresh').addEventListener('click', refresh);
  buttons();
  $('operator-title').focus();
  refresh();
})();
