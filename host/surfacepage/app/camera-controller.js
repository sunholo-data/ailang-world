// The two projections share one document, one graph and a reversible camera.
(() => {
  if (typeof window === 'undefined') return;
  function create(stage, space, hooks) {
    const scene = WorldHeadScene.create(space, hooks);
    const input = WorldHeadTravelInput.create(scene.element, {step, shift:shiftLane, reset:() => move(-1, true)});
    let angle = 0, destination = 0, frame = 0, graph, base, scope, layout, side;
    let selected = '', viewLane = '', keyboard = false;
    let stop = '', followHead = false, history = [], flight, travelFrame = 0;
    const laneStates = new Map();
    const date = new Intl.DateTimeFormat(undefined, {month:'short', day:'numeric', hour:'2-digit', minute:'2-digit',timeZone:'UTC',timeZoneName:'short'});
    const paint = () => {
      stage.style.setProperty('--side-visibility', String(Math.max(0, 1 - angle * 2)));
      scene.render(graph, layout, side, angle, selected, flight);
    };
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
      if (head?.id !== id && !node.draft && index < layout.time.referenceIndex) return id;
      const parents = (graph.parents[id] || []).map(key => graph.N[key]).filter(earlier)
        .sort((a,b) => time(b) - time(a) || b.rank - a.rank);
      // A depth step moves to an earlier instant. Simultaneous records
      // remain individually openable on the same plane.
      return parents.find(p => p.lane === node.lane)?.id ||
        laneStates.get(node.lane)?.find(earlier)?.id ||
        (!node.draft && index < layout.time.referenceIndex ? id : undefined);
    }
    function neighbor(direction) {
      if (!scope) return;
      const lanes = scope.laneOrder, center = lanes.indexOf(layout.ground.centerLane);
      if (center < 0) return;
      for (let i = center + direction; i >= 0 && i < lanes.length; i += direction) {
        if (scope.heads.some(head => head.lane === lanes[i])) return lanes[i];
      }
    }
    function updateControls() {
      const id = layout?.time.stop || '', node = graph?.N[id];
      const active = Boolean(layout && angle === 1 && destination === 1 && !space.hidden && !space.inert);
      const time = Date.parse(node?.at);
      const lane = layout?.ground?.centerLane || '';
      const title = lane ? graph.laneTitle(lane) : '';
      const when = Number.isFinite(time) ? date.format(time) : '';
      const left = neighbor(-1), right = neighbor(1);
      scene.element.dataset.travelState = id;
      scene.element.dataset.viewLane = lane;
      input.update({active, canIn:active && Boolean(nextStop(aim())), canOut:Boolean(id), label:title, when,
        left:left && graph.laneTitle(left), right:right && graph.laneTitle(right)});
    }
    function slice() {
      const options = {...layoutOptions(), base};
      // Browse lanes at the user's chosen moment. Framing an older tip moves
      // the camera to its timestamp without losing the way back to newer lanes.
      scope = stop ? WorldHeadCamera.layout(graph, {...options, stop}) : base;
      const tip = followHead && scope.heads.find(head => head.lane === viewLane);
      const landing = tip && tip.timeIndex < scope.time.referenceIndex ? tip.id : stop;
      if (landing === stop && (!viewLane || viewLane === scope.ground.centerLane)) return scope;
      return WorldHeadCamera.layout(graph, {...options, viewLane, stop:landing});
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
      if (direction > 0) {
        const anchor = layout.points.has(target) ? target : aim();
        const id = nextStop(anchor);
        if (!id) return;
        history.push({stop, followHead});
        stop = id;
        viewLane = graph.N[anchor].lane;
        followHead = false;
      } else if (all) {
        history = [];
        stop = '';
        followHead = false;
      } else if (followHead && layout.time.referenceIndex < scope.time.referenceIndex) {
        followHead = false;
      } else {
        if (!history.length) return;
        ({stop, followHead} = history.pop());
      }
      const next = slice();
      const center = next.heads.find(head => head.lane === next.ground.centerLane);
      animateTo(next, center?.id, Boolean(document.activeElement?.closest('.hc-face,.hc-rear')));
    }
    function shiftLane(direction) {
      if (!layout || angle !== 1 || space.hidden || space.inert) return;
      const lane = neighbor(direction);
      if (!lane) return;
      viewLane = lane;
      followHead = true;
      const next = slice();
      animateTo(next, next.heads.find(head => head.lane === lane).id, true);
    }
    function animateTo(next, id, focusOnFace) {
      cancelTravel();
      if (!hooks.reduced()) flight = {...scene.startTravel(next), progress:0, focusOnFace};
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
      const trackingTip = followHead && selected === aim();
      const restoreFocus = flight?.focusOnFace || (trackingTip && targetID(document.activeElement) === selected);
      if (graph !== fresh) {
        graph = fresh;
        laneStates.clear();
        for (const [lane, value] of Object.entries(graph.LANES)) {
          laneStates.set(lane, [...(value.states || [])].sort((a,b) =>
            graph.timeIndex.get(b.id) - graph.timeIndex.get(a.id) || b.rank - a.rank));
        }
        history = history.filter(place => !place.stop || graph.N[place.stop]);
        if (stop && !graph.N[stop]) ({stop, followHead} = history.pop() || {stop:'', followHead:false});
        if (viewLane && !graph.LANES[viewLane]) viewLane = '';
      }
      cancelTravel();
      base = WorldHeadCamera.layout(graph, layoutOptions());
      layout = slice();
      const tip = trackingTip && aim();
      if (tip && selected !== tip) {
        selected = tip;
        hooks.select(selected);
      }
      side = hooks.sidePoints();
      paint();
      if (restoreFocus && angle === 1 && !space.inert) scene.focus(selected);
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
    return {turn, update, hide, present() { history = []; stop = ''; followHead = false; viewLane = ''; cancelTravel(); },
      select(id) { if (id !== selected && angle === 1 && !space.hidden) { selected = id; paint(); updateControls(); } },
      get active() { return !space.hidden; },
      get heads() { return destination === 1; }
    };
  }
  window.WorldCameraController = {create};
})();
