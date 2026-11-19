// The two projections share one document, one graph and a reversible camera.
(() => {
  if (typeof window === 'undefined') return;
  function create(stage, space, hooks) {
    const scene = WorldHeadScene.create(space, hooks);
    const input = WorldHeadTravelInput.create(scene.element, {step, shift:shiftLane, reset:() => move(-1, true)});
    let angle = 0, destination = 0, frame = 0, graph, base, layout, side;
    let selected = '', viewLane = '', keyboard = false;
    let departure = '', stops = [], flight, travelFrame = 0;
    const laneStates = new Map();
    const date = new Intl.DateTimeFormat(undefined, {month:'short', day:'numeric', hour:'2-digit', minute:'2-digit',timeZone:'UTC',timeZoneName:'short'});
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
      return layout?.heads.find(head => head.lane === layout.ground.centerLane)?.id;
    }
    function nextStop(id) {
      const node = graph?.N[id];
      if (!node || !layout?.points.has(id)) return;
      const index = graph.timeIndex.get(id);
      const time = state => graph.timeIndex.get(state.id);
      const earlier = state => state && !state.draft && time(state) < index && layout.points.has(state.id);
      const head = layout.heads.find(head => head.lane === node.lane);
      if (head?.id !== id && !node.draft && (!current() || index < graph.timeIndex.get(current()))) return id;
      const parents = (graph.parents[id] || []).map(key => graph.N[key]).filter(earlier)
        .sort((a,b) => time(b) - time(a) || b.rank - a.rank);
      // A depth step moves to an earlier instant. Simultaneous records
      // remain individually openable on the same plane.
      return parents.find(p => p.lane === node.lane)?.id ||
        laneStates.get(node.lane)?.find(earlier)?.id || parents[0]?.id ||
        (!node.draft && index < layout.time.referenceIndex ? id : undefined);
    }
    function updateControls() {
      const id = current(), node = graph?.N[id];
      const active = Boolean(layout && angle === 1 && destination === 1 && !space.hidden && !space.inert);
      const time = Date.parse(node?.at);
      const lane = layout?.ground?.centerLane || '';
      const title = lane ? graph.laneTitle(lane) : '';
      const when = Number.isFinite(time) ? date.format(time) : '';
      const lanes = layout?.laneOrder || [];
      const index = lanes.indexOf(lane);
      scene.element.dataset.travelState = id;
      scene.element.dataset.viewLane = lane;
      input.update({active, canIn:active && Boolean(nextStop(aim())), canOut:Boolean(id), label:title, when,
        left:index>0 ? graph.laneTitle(lanes[index-1]) : '', right:index>=0 && index+1<lanes.length ? graph.laneTitle(lanes[index+1]) : ''});
    }
    function slice() {
      return current() || (viewLane && viewLane !== base.ground.centerLane)
        ? WorldHeadCamera.layout(graph, {...layoutOptions(), base, viewLane, stop:current()}) : base;
    }
    function layoutOptions() {
      return {width:stage.clientWidth, height:stage.clientHeight, attention:hooks.attention(), timeXs:hooks.timeXs()};
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
      if (direction > 0) viewLane = graph.N[id].lane;
      const next = slice();
      const center = next.heads.find(head => head.lane === next.ground.centerLane);
      const focus = graph.N[id]?.lane === center?.lane ? id : center?.id || id;
      animateTo(next, focus, Boolean(document.activeElement?.closest('.hc-face,.hc-rear')));
    }
    function shiftLane(direction) {
      if (!layout || angle !== 1 || space.hidden || space.inert) return;
      const lanes = layout.laneOrder;
      const lane = lanes[lanes.indexOf(layout.ground.centerLane) + direction];
      if (!lane) return;
      viewLane = lane;
      const next = slice();
      animateTo(next, next.heads.find(head => head.lane === lane).id, true);
    }
    function animateTo(next, id, focusOnFace) {
      cancelTravel();
      if (!hooks.reduced()) flight = {...scene.startTravel(next), progress:0};
      layout = next;
      selected = id;
      hooks.select(id);
      paint();
      scene.center();
      updateControls();
      if (!flight) {
        if (focusOnFace) scene.focus(id);
        return;
      }
      const started = performance.now();
      function advance(now) {
        const t = Math.min(1, (now - started) / 440);
        if (t === 1 || hooks.reduced()) cancelTravel();
        else flight.progress = t * t * (3 - 2 * t);
        paint();
        if (flight) travelFrame = requestAnimationFrame(advance);
        else if (focusOnFace) scene.focus(id);
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
          laneStates.set(lane, [...(value.states || [])].sort((a,b) =>
            graph.timeIndex.get(b.id) - graph.timeIndex.get(a.id) || b.rank - a.rank));
        }
        stops = stops.filter(id => graph.N[id]);
      }
      cancelTravel();
      base = WorldHeadCamera.layout(graph, layoutOptions());
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
    return {turn, update, hide, present() { stops = []; departure = ''; viewLane = ''; cancelTravel(); },
      select(id) { if (id !== selected && angle === 1 && !space.hidden) { selected = id; paint(); updateControls(); } },
      get active() { return !space.hidden; },
      get heads() { return destination === 1; }
    };
  }
  window.WorldCameraController = {create};
})();
