// The two projections share one document, one graph and a reversible camera.
(() => {
  if (typeof window === 'undefined') return;
  function create(stage, space, hooks) {
    const scene = WorldHeadScene.create(space, {...hooks, pan:updateControls});
    const input = WorldHeadTravelInput.create(scene.element, {step, reset:() => move(-1, true)});
    let angle = 0, destination = 0, frame = 0, graph, base, layout, side;
    let selected = '', keyboard = false;
    let departure = '', stops = [], flight, travelFrame = 0;
    const laneStates = new Map();
    const date = new Intl.DateTimeFormat(undefined, {month:'short', day:'numeric', hour:'2-digit', minute:'2-digit'});
    const paint = () => {
      stage.style.setProperty('--side-visibility', String(Math.max(0, 1 - angle * 2)));
      scene.render(graph, layout, side, angle, selected, flight);
    };
    const current = () => stops.at(-1) || '';
    function targetID(target) {
      const el = target?.closest?.('[data-head-state],[data-rear-state]');
      return el?.dataset.headState || el?.dataset.rearState || '';
    }
    function aim() {
      return scene.aim(selected) || layout?.heads[0]?.id;
    }
    function nextStop(id) {
      const node = graph?.N[id];
      if (!node || !layout?.points.has(id)) return;
      const face = layout.heads.find(head => head.lane === node.lane);
      if (face?.id !== id && !node.draft) return id;
      const parents = (graph.parents[id] || []).map(p => graph.N[p]).filter(p =>
        p && layout.points.has(p.id) && !p.draft && p.rank <= node.rank && !stops.includes(p.id));
      parents.sort((a,b) => b.rank - a.rank);
      const onPath = parents.find(p => p.lane === node.lane);
      if (onPath) return onPath.id;
      // Stay on this brane when it has an earlier pane. At its beginning
      // only a recorded parent can carry the camera onto another brane.
      const lane = laneStates.get(node.lane);
      const index = lane?.indexOf(id) ?? -1;
      if (index >= 0 && index + 1 < lane.length && layout.points.has(lane[index + 1])) return lane[index + 1];
      return parents[0]?.id;
    }
    function updateControls() {
      const id = current(), node = graph?.N[id];
      const active = Boolean(layout && angle === 1 && destination === 1 && !space.hidden && !space.inert);
      const time = Date.parse(node?.at);
      const label = node ? [Number.isFinite(time) && date.format(time), graph.laneTitle(node.lane)].filter(Boolean).join(' · ') : '';
      scene.element.dataset.travelState = id;
      input.update({active, canIn:active && Boolean(nextStop(aim())), canOut:Boolean(id), label});
    }
    function slice() {
      const stop = graph?.N[current()];
      return stop ? WorldHeadTravel.project(graph, base, stop.lane, stop.id) : base;
    }
    function cancelTravel() {
      cancelAnimationFrame(travelFrame);
      flight = null;
      travelFrame = 0;
    }
    function move(direction, all, target) {
      if (!layout || angle !== 1 || space.hidden || space.inert) return;
      let id;
      if (direction > 0) {
        const anchor = layout.points.has(target) ? target : aim();
        id = nextStop(anchor);
        if (!id || stops.includes(id)) return;
        if (!stops.length) departure = anchor;
        stops.push(id);
      } else {
        if (!stops.length) return;
        if (all) stops.length = 0; else stops.pop();
        id = current() || departure;
      }
      const focusOnFace = Boolean(document.activeElement?.closest('.hc-face,.hc-rear'));
      cancelTravel();
      const next = slice();
      if (!hooks.reduced()) flight = {...scene.startTravel(next), progress:0};
      layout = next;
      selected = id;
      hooks.select(id);
      paint();
      if (focusOnFace) scene.focus(id);
      scene.reveal(id, !hooks.reduced());
      updateControls();
      if (!flight) return;
      const started = performance.now();
      function advance(now) {
        const t = Math.min(1, (now - started) / 440);
        if (t === 1 || hooks.reduced()) cancelTravel();
        else flight.progress = t * t * (3 - 2 * t);
        paint();
        if (flight) travelFrame = requestAnimationFrame(advance);
      }
      travelFrame = requestAnimationFrame(advance);
    }
    function step(direction, target) { move(direction, false, targetID(target)); }
    function update() {
      if (space.hidden) return;
      const fresh = hooks.model();
      if (!fresh) return;
      if (graph !== fresh) {
        graph = fresh;
        laneStates.clear();
        for (const [lane, value] of Object.entries(graph.LANES)) {
          laneStates.set(lane, (value.states || []).slice().reverse().sort((a,b) => b.rank - a.rank).map(n => n.id));
        }
        stops = stops.filter(id => graph.N[id]);
      }
      cancelTravel();
      base = WorldHeadCamera.layout(graph, {width:stage.clientWidth, height:stage.clientHeight,
        selected:current() ? departure : selected, attention:hooks.attention()});
      layout = slice();
      side = hooks.sidePoints();
      paint();
      updateControls();
    }
    function settle() {
      frame = 0;
      angle = destination;
      if (angle === 1) {
        paint();
        space.inert = false;
        stage.dataset.view = 'heads';
        if (keyboard) scene.focus(selected);
      } else {
        space.hidden = true;
        delete stage.dataset.camera;
        stage.dataset.view = graph ? 'timeline' : stage.dataset.view;
        hooks.sideVisible(keyboard);
      }
      updateControls();
    }
    function turn(heads, id, instant, focus) {
      if (!hooks.model()) return;
      cancelAnimationFrame(frame);
      destination = heads ? 1 : 0;
      selected = id || selected;
      keyboard = focus;
      if (!heads && space.hidden) return;
      space.hidden = false;
      space.inert = true;
      stage.dataset.camera = 'turn';
      stage.dataset.view = 'turning';
      update();
      if (instant || hooks.reduced() || destination === angle) return settle();
      const start = angle, time = performance.now();
      const span = 1250 * Math.abs(destination - start);
      const step = now => {
        const t = Math.min(1, (now - time) / span);
        if (t === 1 || hooks.reduced()) return settle();
        angle = start + (destination - start) * (t*t*(3 - 2*t));
        paint();
        frame = requestAnimationFrame(step);
      };
      frame = requestAnimationFrame(step);
    }
    function hide() {
      cancelAnimationFrame(frame);
      frame = 0;
      space.hidden = space.inert = true;
      delete stage.dataset.camera;
      cancelTravel();
      updateControls();
    }
    return {turn, update, hide, present() { stops = []; departure = ''; cancelTravel(); },
      select(id) { if (id !== selected && angle === 1 && !space.hidden) { selected = id; paint(); updateControls(); } },
      get active() { return !space.hidden; },
      get heads() { return destination === 1; }
    };
  }
  window.WorldCameraController = {create};
})();
