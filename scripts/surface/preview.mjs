#!/usr/bin/env node
// Live read-only UI preview server for /surface.
//
// Reads host/surfacepage/app/* directly from disk on every request (no `go build`,
// no `ailang-worldd service upgrade`, Cache-Control: no-store) while forwarding
// read-only JSON and /agui/ SSE requests to the local ailang-worldd daemon.
//
// Usage on Studio:
//   node scripts/surface/preview.mjs [--port 7645] [--daemon http://127.0.0.1:7644] [--worktree <dir>]
//   tailscale serve --bg --https=8454 http://127.0.0.1:7645

import fs from 'node:fs';
import http from 'node:http';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const SCRIPT_PARTS = [
  'head-drafts.js',
  'head-camera.js',
  'head-travel-input.js',
  'head-ground-input.js',
  'head-lanes.js',
  'head-node-preview.js',
  'head-lineage.js',
  'head-scene.js',
  'camera-controller.js',
  'node-trajectory.js',
  'app.js',
  'composer.js',
];

const STYLE_PARTS = [
  'app.css',
  'head-scene.css',
  'head-inspect.css',
  'head-drafts.css',
  'head-travel.css',
  'head-lineage.css',
  'head-lanes.css',
  'head-status.css',
  'time-ribbon.css',
  'node-trajectory.css',
];

const STYLESHEET_LINK = '<link rel="stylesheet" href="/surface/app.css">';
const SURFACE_CSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; img-src 'self' data:; base-uri 'none'; form-action 'self'; frame-ancestors 'none'";

const READ_ONLY_GET_PATHS = new Set([
  '/surface/index.json',
  '/surface/follow-status.json',
  '/surface/home-summary.json',
  '/surface/model.json',
  '/surface/stack.json',
  '/surface/file.json',
  '/surface/trajectory.json',
  '/surface/episode-files.json',
  '/surface/verb.json',
  '/v1/health',
]);

const EPISODE_RE = /^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$/;
const TOKEN_SCRUB_RE = /\b(wld_[A-Za-z0-9_-]{8,}|gh[pousr]_[A-Za-z0-9_]{12,}|github_pat_[A-Za-z0-9_]{12,}|sk-[A-Za-z0-9_-]{16,})\b/g;

function scrubText(raw) {
  if (!raw) return '';
  let s = String(raw).replace(TOKEN_SCRUB_RE, '[redacted]');
  s = s.replace(/<CONTEXT_SUMMARY>[\s\S]*?<\/CONTEXT_SUMMARY>/g, '');
  s = s.replace(/<SYSTEM_MESSAGE>[\s\S]*?<\/SYSTEM_MESSAGE>/g, '');
  s = s.replace(/<ADDITIONAL_METADATA>[\s\S]*?<\/ADDITIONAL_METADATA>/g, '');
  const reqMatch = s.match(/<USER_REQUEST>([\s\S]*?)<\/USER_REQUEST>/);
  if (reqMatch) s = reqMatch[1];
  s = s.trim();
  if (s.length > 12000) s = s.slice(0, 12000) + '\n… [truncated]';
  return s;
}

function extractUserMedia(raw, episode) {
  if (!raw || !EPISODE_RE.test(episode || '')) return undefined;
  const seen = new Set(), media = [];
  const imgRe = /\.user_uploaded\/([A-Za-z0-9._-]+\.(?:png|jpe?g|gif|webp))/gi;
  let m;
  while ((m = imgRe.exec(String(raw))) !== null) {
    const file = m[1];
    if (seen.has(file)) continue;
    seen.add(file);
    media.push({
      name: file,
      url: '/surface/conversation-media?ep=' + encodeURIComponent(episode) + '&file=' + encodeURIComponent(file),
    });
  }
  return media.length ? media : undefined;
}

function shortPathLabel(p) {
  const clean = String(p || '').trim();
  const parts = clean.split('/').filter(Boolean);
  return parts.length > 2 ? parts.slice(-2).join('/') : clean;
}

function describeToolCall(tc) {
  if (!tc || typeof tc !== 'object') return null;
  let name = String(tc.name || tc.tool_name || tc.toolName || '').trim();
  let args = tc.arguments ?? tc.args ?? tc.input ?? tc.argumentsJson ?? tc.arguments_json ?? tc.argsJson ?? tc.args_json ?? tc.parameters;
  if (typeof args === 'string') {
    try { args = JSON.parse(args); } catch (_) { args = null; }
  }
  if ((name === 'call_mcp_tool' || /^mcp__/i.test(name)) && args && typeof args === 'object') {
    const innerName = args.ToolName || args.toolName || args.tool_name || '';
    if (innerName) name = String(innerName).trim();
    let inner = args.Arguments ?? args.arguments ?? args.args ?? args;
    if (typeof inner === 'string') {
      try { inner = JSON.parse(inner); } catch (_) {}
    }
    if (inner && typeof inner === 'object') args = inner;
  }
  name = name.replace(/^mcp__[^_]+__/, '');
  if (!name) return null;
  let target = '', file = '';
  if (args && typeof args === 'object') {
    const p = args.path || args.Path || '';
    const cmd = args.command || args.CommandLine || '';
    const pat = args.pattern || args.query || '';
    if (typeof p === 'string' && p.trim() && !p.startsWith('/Users/')) {
      file = p.trim();
      target = shortPathLabel(file);
    } else if (typeof p === 'string' && p.trim()) {
      target = path.basename(p.trim());
    } else if (typeof cmd === 'string' && cmd.trim()) {
      target = cmd.trim().slice(0, 32);
    } else if (typeof pat === 'string' && pat.trim()) {
      target = '"' + pat.trim().slice(0, 22) + '"';
    }
  }
  return { name, ...(target ? { target } : {}), ...(file ? { file } : {}) };
}

function parseTranscriptJSONL(content, agentHint, episode = '') {
  const lines = String(content || '').split('\n');
  let detected = agentHint || '';
  const rawTurns = [];
  for (const rawLine of lines) {
    const line = rawLine.trim();
    if (!line || line.startsWith('{"world_launch":')) continue;
    let obj;
    try { obj = JSON.parse(line); } catch (_) { continue; }
    if (obj.type === 'USER_INPUT' || obj.type === 'PLANNER_RESPONSE') {
      detected ||= 'agy';
      const ts = obj.timestamp || obj.created_at || '';
      if (obj.type === 'USER_INPUT') {
        const media = extractUserMedia(obj.content || '', episode);
        const text = scrubText(obj.content || '');
        if (text || media) rawTurns.push({ role: 'user', at: ts, text: text || '(Uploaded image)', ...(media ? { media } : {}) });
      } else {
        const text = scrubText(obj.content || obj.response || '');
        const tools = (obj.tool_calls || obj.toolCalls || []).map(describeToolCall).filter(Boolean);
        if (text || tools.length) rawTurns.push({ role: 'assistant', at: ts, text, tools });
      }
      continue;
    }
    if (Array.isArray(obj.steps)) {
      detected ||= 'agy';
      for (const step of obj.steps) {
        const st = step?.type || '';
        const ts = step?.metadata?.createdAt || '';
        if (st === 'CORTEX_STEP_TYPE_USER_INPUT') {
          const items = Array.isArray(step.userInput?.items) ? step.userInput.items : [];
          const parts = items.map(it => it?.text || '').filter(Boolean);
          if (!parts.length && step.userInput?.userResponse) parts.push(step.userInput.userResponse);
          const text = scrubText(parts.join('\n\n'));
          if (text) rawTurns.push({ role: 'user', at: ts, text });
        } else if (st === 'CORTEX_STEP_TYPE_PLANNER_RESPONSE') {
          const pr = step.plannerResponse || {};
          const text = scrubText(pr.modifiedResponse || pr.response || '');
          const tools = [];
          for (const tc of pr.toolCalls || []) {
            let name = tc?.name || '';
            if (name === 'call_mcp_tool' && tc?.argumentsJson) {
              try {
                const parsed = JSON.parse(tc.argumentsJson);
                if (parsed?.ToolName) name = parsed.ToolName;
              } catch (_) {}
            }
            if (name) tools.push(name);
          }
          if (text || tools.length) rawTurns.push({ role: 'assistant', at: ts, text, tools });
        }
      }
      continue;
    }
    if (obj.type === 'user' || obj.type === 'assistant') {
      detected ||= 'claude';
      const role = obj.type;
      const ts = obj.timestamp || '';
      const msg = obj.message || {};
      const parts = [];
      const tools = [];
      if (typeof msg.content === 'string') {
        parts.push(msg.content);
      } else if (Array.isArray(msg.content)) {
        for (const block of msg.content) {
          if (block?.type === 'text' && block.text) parts.push(block.text);
          else if (block?.type === 'tool_use' && block.name) {
            const d = describeToolCall(block);
            if (d) tools.push(d);
          }
        }
      }
      const text = scrubText(parts.join('\n\n'));
      if (text || tools.length) rawTurns.push({ role, at: ts, text, tools });
      continue;
    }
    if (obj.type === 'event_msg' || obj.type === 'response_item') {
      detected ||= 'codex';
      const ts = obj.timestamp || '';
      const p = obj.payload || {};
      if (obj.type === 'event_msg' && (p.type === 'user_message' || p.type === 'agent_message')) {
        const text = scrubText(p.message || '');
        if (text) rawTurns.push({ role: p.type === 'user_message' ? 'user' : 'assistant', at: ts, text });
      } else if (obj.type === 'response_item' && (p.role === 'user' || p.role === 'assistant') && Array.isArray(p.content)) {
        const text = scrubText(p.content.map(c => c?.text || '').filter(Boolean).join('\n\n'));
        if (text) rawTurns.push({ role: p.role, at: ts, text });
      }
    }
  }
  const merged = [];
  for (const t of rawTurns) {
    const prev = merged[merged.length - 1];
    if (t.role === 'assistant' && prev && prev.role === 'assistant' && (!prev.text || !t.text)) {
      if (!prev.text && t.text) {
        prev.tools = [...(prev.tools || []), ...(t.tools || [])];
        prev.text = t.text;
        if (t.at) prev.at = t.at;
        continue;
      }
      if (!prev.text && !t.text) {
        prev.tools = [...(prev.tools || []), ...(t.tools || [])];
        if (t.at) prev.at = t.at;
        continue;
      }
    }
    merged.push({ ...t, tools: t.tools?.length ? t.tools : undefined });
  }
  const total = merged.length;
  const maxTurns = 120;
  const messages = total > maxTurns ? merged.slice(total - maxTurns) : merged;
  return { agent: detected, total, truncated: total > maxTurns, messages };
}

function readAgyFull(txPath) {
  if (path.basename(txPath) === 'transcript.jsonl') {
    const full = path.join(path.dirname(txPath), 'transcript_full.jsonl');
    try { if (fs.existsSync(full) && fs.statSync(full).size <= 16 * 1024 * 1024) return fs.readFileSync(full, 'utf8'); } catch (_) {}
  }
  return fs.readFileSync(txPath, 'utf8');
}

function findEpisodeTranscript(episode, agentHint) {
  if (!EPISODE_RE.test(episode || '')) return null;
  const home = os.homedir();
  const sealed = path.join(home, '.ailang', 'world', 'sessions', `${episode}.transcript.jsonl`);
  if (fs.existsSync(sealed)) {
    return { status: 'verified', content: fs.readFileSync(sealed, 'utf8'), agent: agentHint };
  }
  const brainRoot = path.join(home, '.gemini', 'antigravity-cli', 'brain');
  for (const dir of [
    path.join(home, '.ailang', 'world', 'launch', episode),
    path.join(home, '.ailang', 'world-ws', episode),
  ]) {
    const metaFile = path.join(dir, '.agents', '.world-session-meta.json');
    if (fs.existsSync(metaFile)) {
      try {
        const meta = JSON.parse(fs.readFileSync(metaFile, 'utf8'));
        const txCandidate = meta?.transcriptPath || meta?.transcript_path || '';
        if (txCandidate && fs.existsSync(txCandidate)) {
          const realBrain = fs.realpathSync(brainRoot);
          const realTx = fs.realpathSync(txCandidate);
          if (realTx.startsWith(realBrain + path.sep)) {
            const convDir = path.dirname(path.dirname(path.dirname(realTx)));
            return { status: 'live', content: readAgyFull(realTx), agent: agentHint || 'agy', convDir };
          }
        }
      } catch (_) {}
    }
  }
  if ((!agentHint || agentHint === 'agy') && fs.existsSync(brainRoot)) {
    try {
      const candidates = [];
      for (const entry of fs.readdirSync(brainRoot, { withFileTypes: true })) {
        if (!entry.isDirectory()) continue;
        const txFile = path.join(brainRoot, entry.name, '.system_generated', 'logs', 'transcript.jsonl');
        if (fs.existsSync(txFile)) {
          candidates.push({ txFile, convDir: path.join(brainRoot, entry.name), mtime: fs.statSync(txFile).mtimeMs });
        }
      }
      candidates.sort((a, b) => b.mtime - a.mtime);
      const needle = '/' + episode;
      for (const c of candidates.slice(0, 20)) {
        const content = fs.readFileSync(c.txFile, 'utf8');
        if (content.includes(needle)) {
          return { status: 'live', content: readAgyFull(c.txFile), agent: 'agy', convDir: c.convDir };
        }
      }
    } catch (_) {}
  }
  const claudeProjects = path.join(home, '.claude', 'projects');
  for (const dir of [
    path.join(home, '.ailang', 'world', 'launch', episode),
    path.join(home, '.ailang', 'world-ws', episode),
  ]) {
    let realDir = dir;
    try { realDir = fs.realpathSync(dir); } catch (_) {}
    for (const cand of [dir, realDir]) {
      const slug = cand.replace(/[^A-Za-z0-9]/g, '-');
      const projDir = path.join(claudeProjects, slug);
      if (fs.existsSync(projDir)) {
        const files = fs.readdirSync(projDir)
          .filter(f => f.endsWith('.jsonl'))
          .map(f => ({ full: path.join(projDir, f), mtime: fs.statSync(path.join(projDir, f)).mtimeMs }))
          .sort((a, b) => b.mtime - a.mtime);
        if (files.length) {
          return { status: 'live', content: fs.readFileSync(files[0].full, 'utf8'), agent: agentHint || 'claude' };
        }
      }
    }
  }
  return null;
}

async function readEpisodeLedger(episode) {
  if (!EPISODE_RE.test(episode || '')) return null;
  const dbPath = path.join(os.homedir(), '.ailang', 'world', 'world.db');
  if (!fs.existsSync(dbPath)) return null;
  try {
    const { DatabaseSync } = await import('node:sqlite');
    const db = new DatabaseSync(dbPath, { readOnly: true });
    try {
      const getObj = db.prepare('SELECT payload FROM objects WHERE hash_ref = ?');
      const getWorld = db.prepare('SELECT revision FROM worlds WHERE world_ref = ?');
      const listIntents = db.prepare(
        "SELECT o.payload AS ip, oc.payload AS op FROM journal j " +
        "JOIN objects o ON o.hash_ref = j.object_ref AND o.semantic_id = 'world/journal-intent/v1' " +
        "JOIN journal jo ON jo.invocation_id = j.invocation_id AND jo.kind = 'outcome' " +
        "JOIN objects oc ON oc.hash_ref = jo.object_ref AND oc.semantic_id = 'world/journal-outcome/v1' " +
        "WHERE j.kind = 'intent' AND j.invocation_id >= ? AND j.invocation_id < ? ORDER BY j.seq"
      );
      const decode = b => (typeof b === 'string' ? b : Buffer.from(b).toString('utf8'));
      const calls = [];
      const writesMap = new Map(), readsMap = new Map();
      for (const surface of ['mcp', 'cli', 'rest', 'a2a']) {
        const rows = listIntents.all(`${surface}:${episode}:`, `${surface}:${episode};`);
        for (const r of rows) {
          const it = JSON.parse(decode(r.ip)), out = JSON.parse(decode(r.op));
          if (out.status !== 'committed') continue;
          const recRow = getObj.get(it.transitionRef);
          const worldRow = getWorld.get(out.resultRef);
          if (!recRow || !worldRow) continue;
          const rec = JSON.parse(decode(recRow.payload));
          if (rec.episodeId !== episode) continue;
          const tool = String(rec.skillId || 'tool');
          let target = '', file = '', status = 'ok';
          if (rec.input) {
            const inRow = getObj.get(rec.input);
            if (inRow) {
              try {
                const args = JSON.parse(decode(inRow.payload));
                if (typeof args.path === 'string' && args.path.trim() && !args.path.startsWith('/')) {
                  file = args.path.trim().replace(/^\.\//, '');
                  target = file;
                } else if (typeof args.command === 'string' && args.command.trim()) {
                  target = args.command.trim();
                } else if (typeof args.pattern === 'string' && args.pattern.trim()) {
                  target = '"' + args.pattern.trim().slice(0, 24) + '"';
                }
              } catch (_) {}
            }
          }
          const entry = Number(worldRow.revision);
          const atMs = Number(it.logicalTime || 0) * 1000;
          const at = atMs > 0 ? new Date(atMs).toISOString().replace(/\.\d{3}Z$/, 'Z') : '';
          calls.push({ entry, at, atMs, tool, surface, target, file, status });
          if (file && file !== '.') {
            const isWrite = /write|edit|delete/i.test(tool);
            const map = isWrite ? writesMap : (/read/i.test(tool) ? readsMap : null);
            if (map) {
              const prev = map.get(file);
              if (prev) { prev.count++; prev.op = tool; prev.lastEntry = entry; }
              else map.set(file, { path: file, op: tool, count: 1, firstEntry: entry });
            }
          }
        }
      }
      calls.sort((a, b) => a.entry - b.entry);
      const writes = [...writesMap.values()].sort((a, b) => b.count - a.count || a.path.localeCompare(b.path));
      const reads = [...readsMap.values()].sort((a, b) => b.count - a.count || a.path.localeCompare(b.path));
      return { calls, writes, reads };
    } finally { db.close(); }
  } catch (_) { return null; }
}

function hydrateTurnsWithLedger(messages, calls) {
  if (!Array.isArray(messages) || !Array.isArray(calls) || !calls.length) return messages;
  const times = messages.map(m => Date.parse(m.at || '') || 0);
  return messages.map((m, idx) => {
    if (m.role !== 'assistant' || !Array.isArray(m.tools)) return m;
    const hasGeneric = m.tools.some(t => (typeof t === 'string' ? t : t?.name) === 'call_mcp_tool');
    if (!hasGeneric) return m;
    let startMs = 0, endMs = Infinity;
    for (let j = idx - 1; j >= 0; j--) if (times[j] > 0) { startMs = times[j] - 2000; break; }
    for (let j = idx + 1; j < messages.length; j++) if (messages[j].role === 'user' && times[j] > 0) { endMs = times[j] + 1000; break; }
    const matched = calls.filter(c => c.atMs >= startMs && c.atMs <= endMs).map(c => ({
      name: c.tool,
      ...(c.file ? { target: shortPathLabel(c.file), file: c.file } : c.target ? { target: c.target } : {}),
    }));
    if (!matched.length) return m;
    const nonGeneric = m.tools.filter(t => (typeof t === 'string' ? t : t?.name) !== 'call_mcp_tool');
    return { ...m, tools: [...nonGeneric, ...matched] };
  });
}

function fetchDaemonJSON(daemonBase, urlPath) {
  return new Promise((resolve, reject) => {
    const target = new URL(urlPath, daemonBase);
    const req = http.request(
      { hostname: target.hostname, port: target.port, path: target.pathname + target.search, method: 'GET', timeout: 4000 },
      res => {
        const chunks = [];
        res.on('data', c => chunks.push(c));
        res.on('end', () => {
          const raw = Buffer.concat(chunks).toString('utf8');
          let body = null;
          try { body = JSON.parse(raw); } catch (_) {}
          resolve({ status: res.statusCode || 502, raw, body });
        });
      }
    );
    req.on('error', reject);
    req.on('timeout', () => { req.destroy(new Error('daemon timeout')); });
    req.end();
  });
}

const PREVIEW_FALLBACK_CATALOG = JSON.stringify([
  {
    name: 'ailang-world',
    repo: 'sunholo-data/ailang-world',
    base: 'dev',
    commit: 'HEAD',
    host: 'agy',
    model: '',
    seconds: 3600,
    profile: 'ailang-world',
    grants: [
      { Effect: 'Workspace.Read', Scope: 'worktree', Budget: 400 },
      { Effect: 'Workspace.Write', Scope: 'worktree', Budget: 120 },
      { Effect: 'Workspace.Exec', Scope: 'ailang-world', Budget: 80 },
    ],
    checks: [
      { name: 'go-test', argv: ['go', 'test', './...'], timeout_s: 60, sandbox: 'confined' },
      { name: 'go-vet', argv: ['go', 'vet', './...'], timeout_s: 30, sandbox: 'confined' },
    ],
  },
]);

function escapeAttr(s) {
  return String(s).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

function parseArgs(argv) {
  const opts = {
    bind: process.env.SURFACE_PREVIEW_BIND || '127.0.0.1',
    port: Number(process.env.SURFACE_PREVIEW_PORT || 7645),
    daemon: process.env.SURFACE_PREVIEW_DAEMON || 'http://127.0.0.1:7644',
    worktree: process.env.SURFACE_PREVIEW_WORKTREE || '',
    wsRoot: process.env.SURFACE_PREVIEW_WS_ROOT || path.join(os.homedir(), '.ailang', 'world-ws'),
  };
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === '--port' && argv[i + 1]) opts.port = Number(argv[++i]);
    else if (arg === '--bind' && argv[i + 1]) opts.bind = argv[++i];
    else if (arg === '--daemon' && argv[i + 1]) opts.daemon = argv[++i];
    else if (arg === '--worktree' && argv[i + 1]) opts.worktree = path.resolve(argv[++i]);
    else if (arg === '--ws-root' && argv[i + 1]) opts.wsRoot = path.resolve(argv[++i]);
    else if (arg === '--help' || arg === '-h') {
      process.stdout.write(
        'Usage: node scripts/surface/preview.mjs [--port 7645] [--bind 127.0.0.1] [--daemon http://127.0.0.1:7644] [--worktree DIR]\n'
      );
      process.exit(0);
    }
  }
  return opts;
}

const repoFallback = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..', '..');

function hasSurfaceTree(dir) {
  try {
    return fs.statSync(path.join(dir, 'host', 'surfacepage', 'app', 'index.html')).isFile();
  } catch {
    return false;
  }
}

function newestMtimeMs(appDir) {
  let latest = 0;
  for (const name of ['index.html', ...SCRIPT_PARTS, ...STYLE_PARTS]) {
    try {
      const st = fs.statSync(path.join(appDir, name));
      if (st.mtimeMs > latest) latest = st.mtimeMs;
    } catch { /* ignore missing optional file */ }
  }
  return latest;
}

export function resolveWorktree(opts) {
  if (opts.worktree && hasSurfaceTree(opts.worktree)) return opts.worktree;
  let bestDir = hasSurfaceTree(repoFallback) ? repoFallback : '';
  let bestTime = bestDir ? newestMtimeMs(path.join(bestDir, 'host', 'surfacepage', 'app')) : -1;
  try {
    for (const entry of fs.readdirSync(opts.wsRoot, { withFileTypes: true })) {
      if (!entry.isDirectory()) continue;
      const candidate = path.join(opts.wsRoot, entry.name);
      if (!hasSurfaceTree(candidate)) continue;
      const mtime = newestMtimeMs(path.join(candidate, 'host', 'surfacepage', 'app'));
      if (mtime > bestTime) {
        bestTime = mtime;
        bestDir = candidate;
      }
    }
  } catch { /* wsRoot may not exist */ }
  return bestDir || repoFallback;
}

export function assembleBundle(worktree) {
  const base = path.join(worktree, 'host', 'surfacepage');
  const appDir = path.join(base, 'app');
  const script = SCRIPT_PARTS.map(f => fs.readFileSync(path.join(appDir, f), 'utf8')).join('\n');
  const style = STYLE_PARTS.map(f => fs.readFileSync(path.join(appDir, f), 'utf8')).join('\n');
  const index = fs.readFileSync(path.join(appDir, 'index.html'), 'utf8');
  if (!index.includes(STYLESHEET_LINK) || style.toLowerCase().includes('</style')) {
    throw new Error('surfacepage: index.html must link /surface/app.css once and CSS must not close </style>');
  }
  const shell = index.replace(STYLESHEET_LINK, () => `<style>\n${style}</style>`);
  return { base, appDir, script, style, shell };
}

function sendAsset(res, status, contentType, body, csp = "default-src 'none'", extraHeaders = {}) {
  res.writeHead(status, {
    'Content-Type': contentType,
    'Cache-Control': 'no-store',
    'X-Content-Type-Options': 'nosniff',
    'Referrer-Policy': 'no-referrer',
    'Content-Security-Policy': csp,
    ...extraHeaders,
  });
  res.end(body);
}

function proxyToDaemon(req, res, daemonOrigin) {
  const target = new URL(req.url, daemonOrigin);
  const upstream = http.request(
    {
      hostname: target.hostname,
      port: target.port || 80,
      path: target.pathname + target.search,
      method: req.method,
      agent: false,
      headers: {
        accept: req.headers.accept || 'application/json',
        connection: 'close',
        ...(req.headers['content-type'] ? { 'content-type': req.headers['content-type'] } : {}),
        ...(req.headers['content-length'] ? { 'content-length': req.headers['content-length'] } : {}),
        ...(req.headers['last-event-id'] ? { 'last-event-id': req.headers['last-event-id'] } : {}),
      },
    },
    upRes => {
      const headers = {
        'Cache-Control': 'no-store',
        'X-Content-Type-Options': 'nosniff',
      };
      if (upRes.headers['content-type']) headers['Content-Type'] = upRes.headers['content-type'];
      if (upRes.headers['retry-after']) headers['Retry-After'] = upRes.headers['retry-after'];
      res.writeHead(upRes.statusCode || 502, headers);
      upRes.pipe(res);
    }
  );
  upstream.on('error', err => {
    if (res.headersSent || res.writableEnded) {
      res.destroy();
      return;
    }
    sendAsset(
      res,
      502,
      'application/json; charset=utf-8',
      JSON.stringify({ refused: true, reason: 'daemon unreachable: ' + err.message })
    );
  });
  res.on('close', () => {
    if (!res.writableEnded) upstream.destroy();
  });
  if (req.method === 'GET' || req.method === 'HEAD') {
    upstream.end();
  } else {
    req.pipe(upstream);
  }
}

export function createPreviewServer(opts) {
  return http.createServer((req, res) => {
    const parsed = new URL(req.url || '/', 'http://localhost');
    const pathname = parsed.pathname;
    const isRead = req.method === 'GET' || req.method === 'HEAD';
    const worktree = resolveWorktree(opts);
    const label = path.basename(worktree);
    const meta = { 'X-World-Surface-Worktree': label };
    const reply = (status, contentType, body, csp, extra) =>
      sendAsset(res, status, contentType, req.method === 'HEAD' ? '' : body, csp, extra);

    try {
      if (isRead && pathname === '/health') {
        reply(200, 'application/json; charset=utf-8', JSON.stringify({ ok: true, worktree: label }), "default-src 'none'", meta);
        return;
      }
      if (isRead && (pathname === '/' || pathname === '/surface')) {
        if (pathname === '/') {
          res.writeHead(302, { Location: '/surface', 'Cache-Control': 'no-store' });
          res.end();
          return;
        }
        const { shell } = assembleBundle(worktree);
        const badged = shell
          .replace('id="episode-draft" data-catalog=""', `id="episode-draft" data-catalog="${escapeAttr(PREVIEW_FALLBACK_CATALOG)}"`)
          .replace('<title>World</title>', `<title>World · Preview (${label})</title>`)
          .replace(
            '<span class="brand-name">World</span>',
            `<span class="brand-name">World</span><span style="margin-left:6px;padding:2px 6px;border-radius:4px;background:#e73c17;color:#fff;font:700 10px/1.2 var(--mono);letter-spacing:.04em" title="${label}">PREVIEW</span>`
          );
        reply(200, 'text/html; charset=utf-8', badged, SURFACE_CSP, meta);
        return;
      }
      if (isRead && pathname === '/surface/app.js') {
        const { script } = assembleBundle(worktree);
        reply(200, 'text/javascript; charset=utf-8', script, "default-src 'none'", meta);
        return;
      }
      if (isRead && pathname === '/surface/app.css') {
        const { style } = assembleBundle(worktree);
        reply(200, 'text/css; charset=utf-8', style, "default-src 'none'", meta);
        return;
      }
      if (isRead && pathname === '/surface/surface.js') {
        const file = path.join(worktree, 'host', 'surfacepage', 'surface.js');
        reply(200, 'text/javascript; charset=utf-8', fs.readFileSync(file), "default-src 'none'", meta);
        return;
      }
      if (isRead && pathname === '/surface/ailang-logo.svg') {
        const file = path.join(worktree, 'host', 'surfacepage', 'app', 'ailang-logo.svg');
        reply(200, 'image/svg+xml', fs.readFileSync(file), "default-src 'none'", meta);
        return;
      }
      if (isRead && pathname.startsWith('/surface/fonts/')) {
        const name = path.basename(pathname);
        const fontFile = path.join(worktree, 'host', 'surfacepage', 'fonts', name);
        if (name.endsWith('.ttf') && fs.existsSync(fontFile)) {
          reply(200, 'font/ttf', fs.readFileSync(fontFile), "default-src 'none'", meta);
          return;
        }
        reply(404, 'text/plain; charset=utf-8', 'font not found');
        return;
      }
      if (isRead && pathname === '/surface/requests.json') {
        reply(
          200,
          'application/json; charset=utf-8',
          JSON.stringify({ pending: [], decisions: [], confirmations: [], unavailable: 'Read-only preview (approve on main surface)' }),
          "default-src 'none'",
          meta
        );
        return;
      }
      if (isRead && pathname === '/surface/conversation.json') {
        (async () => {
          try {
            const up = await fetchDaemonJSON(opts.daemon, '/surface/conversation.json' + parsed.search);
            if (up.status === 200) {
              reply(200, 'application/json; charset=utf-8', up.raw, "default-src 'none'", meta);
              return;
            }
            const stRes = await fetchDaemonJSON(opts.daemon, '/surface/stack.json' + parsed.search);
            if (stRes.status !== 200 || !stRes.body?.front) {
              reply(
                stRes.status || 404,
                'application/json; charset=utf-8',
                stRes.raw || JSON.stringify({ schema: 'world/surface-conversation/v1', status: 'unavailable', reason: 'state not found' }),
                "default-src 'none'",
                meta
              );
              return;
            }
            const front = stRes.body.front;
            const episode = front.trajectoryEpisode || front.run?.episode || front.episode || '';
            const agentHint = front.agent || front.run?.harness || '';
            if (!episode) {
              reply(200, 'application/json; charset=utf-8', JSON.stringify({ schema: 'world/surface-conversation/v1', sourceState: front.id, status: 'unavailable', reason: 'state has no recorded episode' }), "default-src 'none'", meta);
              return;
            }
            const found = findEpisodeTranscript(episode, agentHint);
            if (!found) {
              reply(200, 'application/json; charset=utf-8', JSON.stringify({ schema: 'world/surface-conversation/v1', episode, sourceState: front.id, status: 'unavailable', reason: 'no live or sealed transcript found on host' }), "default-src 'none'", meta);
              return;
            }
            const parsedTx = parseTranscriptJSONL(found.content, found.agent, episode);
            const ledger = await readEpisodeLedger(episode);
            const hydrated = hydrateTurnsWithLedger(parsedTx.messages, ledger?.calls);
            reply(
              200,
              'application/json; charset=utf-8',
              JSON.stringify({
                schema: 'world/surface-conversation/v1',
                episode,
                sourceState: front.id,
                agent: parsedTx.agent,
                status: found.status,
                total: parsedTx.total,
                truncated: parsedTx.truncated,
                messages: hydrated,
              }),
              "default-src 'none'",
              meta
            );
          } catch (err) {
            reply(502, 'application/json; charset=utf-8', JSON.stringify({ schema: 'world/surface-conversation/v1', status: 'unavailable', reason: 'preview conversation read failed: ' + err.message }), "default-src 'none'", meta);
          }
        })();
        return;
      }
      if (isRead && pathname === '/surface/stack.json') {
        (async () => {
          try {
            const stRes = await fetchDaemonJSON(opts.daemon, '/surface/stack.json' + parsed.search);
            const perms = stRes.body?.front?.run?.permissions;
            const ep = stRes.body?.front?.run?.episode || stRes.body?.front?.trajectoryEpisode || stRes.body?.front?.episode || '';
            const dbPath = path.join(os.homedir(), '.ailang', 'world', 'world.db');
            if (stRes.status === 200 && perms?.sessions?.some(s => (s.grants || []).some(g => g.used == null)) && EPISODE_RE.test(ep) && fs.existsSync(dbPath)) {
              try {
                const { DatabaseSync } = await import('node:sqlite');
                const db = new DatabaseSync(dbPath, { readOnly: true });
                try {
                  const rows = db.prepare("SELECT o.payload FROM journal j JOIN objects o ON o.hash_ref = j.object_ref WHERE j.kind = 'intent' AND j.invocation_id >= ? AND j.invocation_id < ? AND o.semantic_id = 'world/effect-intent/v1'").all(`effect:${ep}:`, `effect:${ep};`);
                  const byEx = new Map(), byEf = new Map();
                  for (const r of rows) {
                    const it = JSON.parse(typeof r.payload === 'string' ? r.payload : Buffer.from(r.payload).toString('utf8'));
                    if (it.episodeId !== ep) continue;
                    const c = Number(it.cost || 0), sc = /[\/\\;:@ \t\r\n]/.test(it.scope || '') ? '[private scope]' : (it.scope || '');
                    byEx.set(`${it.effect}\0${sc}`, (byEx.get(`${it.effect}\0${sc}`) || 0) + c);
                    byEf.set(it.effect, (byEf.get(it.effect) || 0) + c);
                  }
                  for (const s of perms.sessions) {
                    const cnt = new Map();
                    for (const g of s.grants || []) cnt.set(g.effect, (cnt.get(g.effect) || 0) + 1);
                    for (const g of s.grants || []) {
                      const k = `${g.effect}\0${g.scope}`;
                      const used = byEx.has(k) ? byEx.get(k) : (cnt.get(g.effect) === 1 ? (byEf.get(g.effect) || 0) : 0);
                      g.used = used; g.reserved = 0;
                      if (g.budget != null) {
                        const left = Math.max(0, g.budget - used);
                        g.ledgerRemaining = left;
                        const st = g.state || s.state;
                        if (st === 'active') g.remaining = left; else if (st === 'expired') g.remaining = 0;
                      }
                    }
                  }
                  perms.truncated = false;
                } finally { db.close(); }
              } catch (_) {}
            }
            reply(stRes.status || 502, 'application/json; charset=utf-8', stRes.body ? JSON.stringify(stRes.body) : (stRes.raw || '{}'), "default-src 'none'", meta);
          } catch (_) { proxyToDaemon(req, res, opts.daemon); }
        })();
        return;
      }
      if (isRead && pathname === '/surface/conversation-media') {
        const ep = parsed.searchParams.get('ep') || '', file = parsed.searchParams.get('file') || '';
        if (!EPISODE_RE.test(ep) || !/^[A-Za-z0-9._-]+\.(png|jpe?g|gif|webp)$/i.test(file)) {
          reply(400, 'text/plain; charset=utf-8', 'invalid media request'); return;
        }
        const found = findEpisodeTranscript(ep, 'agy');
        const imgPath = found?.convDir ? path.join(found.convDir, '.user_uploaded', file) : '';
        if (!imgPath || !fs.existsSync(imgPath)) { reply(404, 'text/plain; charset=utf-8', 'media not found'); return; }
        const ext = path.extname(file).toLowerCase();
        const mime = ext === '.png' ? 'image/png' : ext === '.gif' ? 'image/gif' : ext === '.webp' ? 'image/webp' : 'image/jpeg';
        reply(200, mime, fs.readFileSync(imgPath), "default-src 'none'", meta);
        return;
      }
      if (isRead && pathname === '/surface/episode-files.json') {
        (async () => {
          try {
            const up = await fetchDaemonJSON(opts.daemon, '/surface/episode-files.json' + parsed.search);
            if (up.status === 200 && up.body?.status === 'partial' && EPISODE_RE.test(up.body.episode || '')) {
              const ledger = await readEpisodeLedger(up.body.episode);
              if (ledger) {
                const enriched = { schema: 'world/surface-episode-files/v1', episode: up.body.episode, sourceState: up.body.sourceState, status: 'complete', reads: ledger.reads, writes: ledger.writes, produced: up.body.produced || [] };
                reply(200, 'application/json; charset=utf-8', JSON.stringify(enriched), "default-src 'none'", meta);
                return;
              }
            }
            reply(up.status || 502, 'application/json; charset=utf-8', up.raw || '{}', "default-src 'none'", meta);
          } catch (_) { proxyToDaemon(req, res, opts.daemon); }
        })();
        return;
      }
      if (isRead && pathname === '/surface/trajectory.json') {
        (async () => {
          try {
            const up = await fetchDaemonJSON(opts.daemon, '/surface/trajectory.json' + parsed.search);
            if (up.status === 200 && up.body?.world?.truncated && EPISODE_RE.test(up.body.episode || '')) {
              const ledger = await readEpisodeLedger(up.body.episode);
              if (ledger?.calls?.length) {
                const total = ledger.calls.length, calls = total > 120 ? ledger.calls.slice(total - 120) : ledger.calls;
                up.body.world = { available: true, total, truncated: total > 120, calls: calls.map(({ entry, at, tool, surface, target, status }) => ({ entry, at, tool, surface, ...(target ? { target } : {}), ...(status ? { status } : {}) })) };
                reply(200, 'application/json; charset=utf-8', JSON.stringify(up.body), "default-src 'none'", meta);
                return;
              }
            }
            reply(up.status || 502, 'application/json; charset=utf-8', up.raw || '{}', "default-src 'none'", meta);
          } catch (_) { proxyToDaemon(req, res, opts.daemon); }
        })();
        return;
      }
      if (isRead && pathname === '/surface/file.json') {
        (async () => {
          try {
            const up = await fetchDaemonJSON(opts.daemon, '/surface/file.json' + parsed.search);
            if (up.status === 200) { reply(200, 'application/json; charset=utf-8', up.raw, "default-src 'none'", meta); return; }
            const at = parsed.searchParams.get('at') || '', relPath = parsed.searchParams.get('path') || '';
            const epMatch = at.match(/:episode:([A-Za-z0-9._-]+)$/);
            if (epMatch && relPath && !relPath.startsWith('/') && !relPath.includes('..')) {
              const wsDir = path.join(os.homedir(), '.ailang', 'world-ws', epMatch[1]);
              const full = path.join(wsDir, relPath);
              if (fs.existsSync(wsDir) && fs.existsSync(full) && fs.realpathSync(full).startsWith(fs.realpathSync(wsDir) + path.sep)) {
                const stat = fs.statSync(full);
                if (stat.isFile() && stat.size <= 512 * 1024) {
                  const text = fs.readFileSync(full, 'utf8');
                  reply(200, 'application/json; charset=utf-8', JSON.stringify({ path: relPath, kind: 'text', status: 'M', text, truncated: false }), "default-src 'none'", meta);
                  return;
                }
              }
            }
            reply(up.status || 404, 'application/json; charset=utf-8', up.raw || '{}', "default-src 'none'", meta);
          } catch (_) { proxyToDaemon(req, res, opts.daemon); }
        })();
        return;
      }
      if (isRead && READ_ONLY_GET_PATHS.has(pathname)) {
        proxyToDaemon(req, res, opts.daemon);
        return;
      }
      if (req.method === 'POST' && pathname === '/agui/') {
        proxyToDaemon(req, res, opts.daemon);
        return;
      }
      sendAsset(res, 405, 'application/json; charset=utf-8', JSON.stringify({ refused: true, reason: 'read-only surface preview' }));
    } catch (err) {
      sendAsset(res, 500, 'text/plain; charset=utf-8', 'Preview error: ' + err.message);
    }
  });
}

if (import.meta.main) {
  const opts = parseArgs(process.argv.slice(2));
  const server = createPreviewServer(opts);
  server.listen(opts.port, opts.bind, () => {
    const active = resolveWorktree(opts);
    process.stdout.write(
      `World /surface live preview listening on http://${opts.bind}:${opts.port}/surface\n` +
      `  daemon:   ${opts.daemon} (read-only JSON + /agui/)\n` +
      `  worktree: ${active}${opts.worktree ? ' (pinned)' : ' (auto-tracks newest in ' + opts.wsRoot + ')'}\n`
    );
  });
}
