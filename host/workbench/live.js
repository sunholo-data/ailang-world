/* Server HTML owns all markup; this consumer only follows cursors and swaps regions. */
(() => {
  "use strict";
  const STATE_SCHEMA = "world/agui-state/v1";
  const REGIONS = ["nav[aria-label=\"world browser\"]","section[aria-label=\"timeline\"]","section[aria-label=\"live\"]","section[aria-label=\"world graph\"]","section[aria-label=\"decisions\"]"];
  function splitFrames(tail, chunk) {
    const frames = (tail + chunk).split("\n\n");
    const rest = frames.pop();
    const events = frames.flatMap(frame => frame.split("\n")
      .filter(line => line.startsWith("data: ")).map(line => JSON.parse(line.slice(6))));
    return {events, tail: rest};
  }
  function applyDelta(state, delta) {
    for (const op of delta) {
      if (op.op !== "replace") throw new Error("unsupported state delta");
      if (op.path === "/lastIndex") state.lastIndex = op.value;
      else if (op.path === "/logHead") state.logHead = op.value;
      else throw new Error("unknown state path");
    }
    return state;
  }
  function nextBackoff(retryAfter, previous = 0) {
    const parsed = Number(retryAfter);
    const first = Number.isFinite(parsed) && parsed > 0 ? parsed : 1;
    return Math.min(30, previous > 0 ? previous * 2 : first);
  }
  function renderStatus(state) {
    if (state.layout) return "paused: page layout changed, reload required";
    if (state.mode === "paused") return "paused: " + state.reason + ", retrying in " + state.wait + "s";
    if (state.mode === "failed") return state.reason;
    if (state.count) return state.count + " new entries since " + state.checked;
    return "Live: checked " + state.checked + ", no new entries after entry " + state.cursor + "; next check within 18 s";
  }
  function applySwap(oldDocument, freshDocument, state) {
    for (const selector of REGIONS) {
      const old = oldDocument.querySelector(selector), fresh = freshDocument.querySelector(selector);
      if (!old || !fresh) { state.layout = true; continue; }
      if (selector === 'section[aria-label="live"]') {
        const prior = new Set(Array.from(old.querySelectorAll("li")).map(row => row.textContent));
        for (const row of fresh.querySelectorAll("li")) {
          if (!prior.has(row.textContent)) row.classList.add("fresh");
        }
      }
      old.replaceWith(oldDocument.importNode(fresh, true));
    }
    const span = oldDocument.querySelector("[data-live-status]");
    if (span) span.textContent = renderStatus(state);
    return !state.layout;
  }
  if (typeof module !== "undefined") module.exports = {splitFrames, applyDelta, nextBackoff, applySwap, renderStatus};
  if (typeof document === "undefined") return;
  const live = document.querySelector('section[aria-label="live"]');
  if (!live) return;
  const state = {mode: "live", cursor: Number(live.dataset.liveCursor), checked: new Date().toLocaleTimeString(), count: 0};
  let active = null, retryTimer = null, refreshTimer = null, refreshing = false, backoff = 0, run = 0;
  const status = () => {
    const span = document.querySelector("[data-live-status]");
    if (span) span.textContent = renderStatus(state);
  };
  const visible = () => document.visibilityState === "visible" && !state.layout;
  function pause(reason, retryAfter) {
    backoff = nextBackoff(retryAfter, backoff);
    Object.assign(state, {mode: "paused", reason, wait: backoff}); status();
    retryTimer = setTimeout(() => { retryTimer = null; follow(); }, backoff * 1000);
  }
  function scheduleRefresh() {
    if (refreshTimer || refreshing || state.layout) return;
    refreshTimer = setTimeout(async () => {
      refreshTimer = null; refreshing = true;
      try {
        const response = await fetch(location.href);
        if (!response.ok) throw new Error("refresh failed (" + response.status + " " + response.statusText + ")");
        const fresh = new DOMParser().parseFromString(await response.text(), "text/html");
        if (!applySwap(document, fresh, state)) {
          if (active) active.abort();
          if (retryTimer) clearTimeout(retryTimer);
          retryTimer = null;
        }
      } catch (error) { Object.assign(state, {mode: "failed", reason: error.message}); status(); }
      finally { refreshing = false; }
    }, 300);
  }
  async function follow() {
    if (!visible() || active || retryTimer) return;
    const controller = new AbortController(); active = controller;
    let terminal = false, failure = null, retryAfter;
    try {
      const response = await fetch("/agui/", {method: "POST", headers: {"Content-Type": "application/json"},
        body: JSON.stringify({threadId: "workbench", runId: "wb-" + run++, messages: [], state: {schema: STATE_SCHEMA, lastIndex: state.cursor}}), signal: controller.signal});
      if (!response.ok) {
        retryAfter = response.headers.get("Retry-After");
        throw new Error(response.status === 503 ? "stream limit" : "stream HTTP " + response.status);
      }
      state.mode = "live"; status();
      const reader = response.body.getReader(), decoder = new TextDecoder();
      let tail = "";
      while (!controller.signal.aborted && !state.layout) {
        const chunk = await reader.read();
        if (controller.signal.aborted || state.layout) break;
        const parsed = splitFrames(tail, decoder.decode(chunk.value || new Uint8Array(), {stream: !chunk.done}));
        tail = parsed.tail;
        for (const event of parsed.events) {
          if (event.type === "STATE_DELTA") {
            const next = applyDelta({lastIndex: state.cursor}, event.delta); state.cursor = next.lastIndex;
          } else if (event.type === "CUSTOM" && event.name === "world.entry.committed") {
            state.count++; scheduleRefresh(); status();
          } else if (event.type === "RUN_FINISHED") {
            state.cursor = event.result.lastIndex; terminal = true;
            Object.assign(state, {mode: "live", checked: new Date().toLocaleTimeString(), count: 0}); status();
          } else if (event.type === "RUN_ERROR") throw new Error(event.code || "stream error");
        }
        if (terminal || chunk.done) break;
      }
      if (!terminal && !controller.signal.aborted && !state.layout) throw new Error("connection ended");
    } catch (error) { if (!controller.signal.aborted) failure = error.message; }
    finally { controller.abort(); active = null; }
    if (!visible()) return;
    if (failure) pause(failure, retryAfter);
    else if (terminal) { backoff = 0; follow(); }
    else follow();
  }
  document.addEventListener("visibilitychange", () => {
    if (!visible()) {
      if (active) active.abort();
      if (retryTimer) clearTimeout(retryTimer);
      retryTimer = null;
      Object.assign(state, {mode: "paused", reason: "hidden tab", wait: 0}); status();
    } else follow();
  });
  status(); follow();
})();
