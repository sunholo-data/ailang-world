// The two projections share one document, one graph and a reversible camera.
(() => {
  if (typeof window === 'undefined') return;
  function create(stage, space, hooks) {
    const scene = WorldHeadScene.create(space, hooks);
    const scroll = scene.element.querySelector('.hc-scroll');
    const toolbar = space.querySelector('#head-toolbar');
    const lanes = WorldHeadLanes.create(toolbar, lane => selectLane(lane));
    const input = WorldHeadTravelInput.create(scene.element, {step, shift:(count, fast) => shiftLane(count, {fast}), reset:() => move(-1, true)});
    let angle = 0, destination = 0, frame = 0, graph, base, scope, layout, side;
    let selected = '', viewLane = '', keyboard = false;
    // The scope bounds keyboard navigation. The camera can stay on a
    // particular cut inside it, or follow the selected lane's tip (null).
    let scopeStop = '', pinnedStop = '', history = [], flight, travelFrame = 0;
    const laneStates = new Map();
    WorldHeadGroundInput.create(scene.element, {
      view:() => angle === 1 && destination === 1 && !space.hidden && !space.inert && layout ?
        {layout, stops:hooks.stops()} : null,
      seek:(id, lane) => seek(id, lane, true),
      shift:count => shiftLane(count, {dragging:true}),
      grabPan:scene.grabPan, pan:scene.pan,
      settle(changedLane) {
        const head = changedLane ? aim() : undefined;
        scene.center(head, true);
        // A depth seek can interrupt a lane flight before it transfers focus.
        if (head && flight) flight.focusOnFace = true;
        else if (head && targetID(document.activeElement) !== head) scene.focus(head, false);
      },
      hop(stop, worldX) {
        const lanes = layout.laneOrder.filter(lane => laneStates.get(lane)?.some(node =>
          (!stop.id || !node.draft) && graph.timeIndex.get(node.id) <= stop.index));
        const distance = lane => Math.abs(layout.plan.lanes.get(lane).origin * layout.plan.pitch - worldX);
        if (lanes.length) seek(stop.id, lanes.reduce((a,b) => distance(a) <= distance(b) ? a : b));
      }
    });
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
    function neighbor(count, at = scope) {
      if (!at) return;
      const lanes = at.laneOrder, center = lanes.indexOf(layout.ground.centerLane);
      if (center < 0) return;
      const direction = Math.sign(count);
      let remaining = Math.abs(count), target;
      for (let i = center + direction; remaining && i >= 0 && i < lanes.length; i += direction) {
        if (at.heads.some(head => head.lane === lanes[i])) {
          target = lanes[i];
          remaining--;
        }
      }
      return target;
    }
    function updateControls() {
      const id = layout?.time.stop || '', node = graph?.N[id];
      const active = Boolean(layout && angle === 1 && destination === 1 && !space.hidden && !space.inert);
      const time = Date.parse(node?.at);
      const lane = layout?.ground?.centerLane || '';
      const title = lane ? graph.laneTitle(lane) : '';
      const when = Number.isFinite(time) ? date.format(time) : '';
      hooks.timeChanged(layout?.time.referenceIndex, active);
      lanes.update({active, graph, current:lane, heads:scope?.heads || [],
        historical:Boolean(id), order:scope?.laneOrder || []});
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
      scope = scopeStop ? WorldHeadCamera.layout(graph, {...options, stop:scopeStop}) : base;
      const tip = pinnedStop === null && scope.heads.find(head => head.lane === viewLane);
      const landing = pinnedStop ?? (tip && tip.timeIndex < scope.time.referenceIndex ? tip.id : scopeStop);
      if (landing === scopeStop && (!viewLane || viewLane === scope.ground.centerLane)) return scope;
      return WorldHeadCamera.layout(graph, {...options, viewLane, stop:landing});
    }
    function layoutOptions() {
      // Frame the street from its scroll box, including the reserved travel controls.
      const height = scroll.offsetHeight || stage.clientHeight - toolbar.offsetHeight - (parseFloat(getComputedStyle(scroll).bottom) || 0);
      return {width:stage.clientWidth, height:Math.max(1, height), attention:hooks.attention(), timeXs:hooks.timeXs()};
    }
    function cancelTravel() {
      cancelAnimationFrame(travelFrame);
      flight = null;
      travelFrame = 0;
    }
    function move(direction, all, target, fast) {
      if (!layout || angle !== 1 || space.hidden || space.inert) return;
      if (direction > 0) {
        const anchor = layout.points.has(target) ? target : aim();
        const id = nextStop(anchor);
        if (!id) return;
        history.push({scopeStop, pinnedStop});
        scopeStop = pinnedStop = id;
        viewLane = graph.N[anchor].lane;
      } else if (all) {
        history = [];
        scopeStop = pinnedStop = '';
      } else if (layout.time.referenceIndex < scope.time.referenceIndex) {
        pinnedStop = scopeStop;
      } else {
        if (!history.length) return;
        ({scopeStop, pinnedStop} = history.pop());
      }
      const next = slice();
      const center = next.heads.find(head => head.lane === next.ground.centerLane);
      const focus = document.activeElement?.closest('.hc-face,.hc-rear') ||
        (document.activeElement === document.body && flight?.focusOnFace);
      animateTo(next, center?.id, {focusOnFace:Boolean(focus), fast});
    }
    function shiftLane(count, motion) {
      const target = neighbor(count, motion.dragging ? layout : scope);
      if (target) selectLane(target, {focusOnFace:true, ...motion});
      return target;
    }
    function selectLane(lane, motion) {
      if (!layout || angle !== 1 || space.hidden || space.inert || !scope.heads.some(head => head.lane === lane)) return;
      viewLane = lane;
      pinnedStop = motion?.dragging ? layout.time.stop : null;
      const next = slice();
      animateTo(next, next.heads.find(head => head.lane === lane).id, motion);
    }
    function seek(id, lane = layout?.ground.centerLane, dragging = false) {
      if (!layout || angle !== 1 || space.hidden || space.inert) return;
      const record = graph.N[id], index = graph.timeIndex.get(id);
      if (id && (!record || record.draft || !Number.isInteger(index) || index > graph.presentIndex)) return;
      if (id === scopeStop && pinnedStop === scopeStop && lane === layout.ground.centerLane) {
        if (!dragging) scene.center();
        return;
      }
      // A ribbon drag sets one absolute cut, regardless of how many pixels it crosses.
      history = id ? [{scopeStop:'', pinnedStop:''}] : [];
      scopeStop = pinnedStop = id;
      const visible = !id || laneStates.get(lane)?.some(state => !state.draft && graph.timeIndex.get(state.id) <= index);
      viewLane = visible ? lane : record.lane;
      const next = slice();
      const head = next.heads.find(head => head.lane === next.ground.centerLane)?.id || id;
      if (dragging) return animateTo(next, head, {dragging:true});
      cancelTravel();
      layout = next;
      selected = head;
      hooks.select(selected);
      paint();
      scene.center();
      updateControls();
    }
    function animateTo(next, id, {focusOnFace = false, fast = false, dragging = false} = {}) {
      const duration = dragging ? 300 : (fast || flight ? 110 : 220);
      // Retarget from the last drawn frame. Closely spaced pointer events
      // must not keep resetting the camera clock before it can move.
      const previousFrame = dragging ? flight?.drawnAt : undefined;
      cancelTravel();
      if (!hooks.reduced()) flight = {...scene.startTravel(next, dragging, id), progress:0, focusOnFace};
      layout = next;
      selected = id;
      hooks.select(id);
      paint();
      if (!dragging) scene.center();
      updateControls();
      if (!flight) {
        if (focusOnFace) scene.focus(id, !dragging);
        return;
      }
      const now = performance.now();
      const started = Math.max(previousFrame ?? now, now - 32);
      flight.drawnAt = started;
      function advance() {
        const now = performance.now();
        const t = Math.max(0, Math.min(1, (now - started) / duration));
        flight.drawnAt = now;
        const focusAtArrival = flight.focusOnFace;
        if (t === 1 || hooks.reduced()) cancelTravel();
        else if (dragging) {
          // Ease the selected face in screen space. Perspective otherwise
          // compresses most of a long approach into the last few frames.
          const q = 1 - (1 - t) ** 3;
          const weighted = q * flight.depthRatio;
          flight.progress = weighted / (1 - q + weighted);
        } else flight.progress = t * t * (3 - 2 * t);
        paint();
        if (flight) travelFrame = requestAnimationFrame(advance);
        else if (focusAtArrival) scene.focus(id, !dragging);
      }
      travelFrame = requestAnimationFrame(advance);
    }
    function step(direction, target, fast) {
      // While a flight is settling DOM focus can still be on the previous
      // face. Continue from the current stop on the next keypress.
      move(direction, false, flight ? undefined : targetID(target), fast);
    }
    function update(anchor) {
      if (space.hidden) return;
      const fresh = hooks.model();
      if (!fresh) return;
      const trackingTip = !anchor && Boolean(viewLane) && selected === aim();
      const restoreFocus = flight?.focusOnFace || (trackingTip && targetID(document.activeElement) === selected);
      if (graph !== fresh) {
        graph = fresh;
        laneStates.clear();
        for (const [lane, value] of Object.entries(graph.LANES)) {
          laneStates.set(lane, [...(value.states || [])].sort((a,b) =>
            graph.timeIndex.get(b.id) - graph.timeIndex.get(a.id) || b.rank - a.rank));
        }
        const exists = id => !id || graph.N[id];
        history = history.filter(place => exists(place.scopeStop) && exists(place.pinnedStop));
        if (!exists(scopeStop)) ({scopeStop, pinnedStop} = history.pop() || {scopeStop:'', pinnedStop:''});
        if (!exists(pinnedStop)) pinnedStop = scopeStop;
        if (viewLane && !graph.LANES[viewLane]) viewLane = '';
      }
      cancelTravel();
      base = WorldHeadCamera.layout(graph, layoutOptions());
      layout = slice();
      if (anchor) scene.frameNode(layout, anchor);
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
        hooks.placeStatus(false);
        space.hidden = true;
        delete stage.dataset.camera;
        stage.dataset.view = graph ? 'timeline' : stage.dataset.view;
        hooks.sideVisible(keyboard);
      }
      updateControls();
    }
    function turn(heads, id, {instant = false, keyboard:focus = false, frameNode = false} = {}) {
      const model = hooks.model();
      if (!model) return;
      const node = frameNode && model.timeIndex.get(id) <= model.presentIndex && model.N[id];
      if (node) {
        viewLane = node.lane;
        const cut = node.draft ? '' : id;
        // Returning to the same stop keeps the path back through history.
        // Entering at another node establishes its exact cut, including ties.
        if (!layout || layout.time.stop !== cut) {
          history = cut ? [{scopeStop:'', pinnedStop:''}] : [];
          scopeStop = pinnedStop = cut;
        } else pinnedStop = cut;
      }
      cancelAnimationFrame(frame);
      destination = heads ? 1 : 0;
      selected = id || selected;
      keyboard = focus || Boolean(node);
      if (!heads && space.hidden) return;
      space.hidden = false;
      space.inert = true;
      if (heads) hooks.placeStatus(true);
      stage.dataset.camera = 'turn';
      stage.dataset.view = 'turning';
      update(node && id);
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
      hooks.placeStatus(false);
      delete stage.dataset.camera;
      cancelTravel();
      updateControls();
    }
    return {turn, update, hide, seek,
      followLane(lane) {
        if (!graph?.LANES[lane]) return;
        history = []; scopeStop = pinnedStop = ''; viewLane = lane;
        update();
        const id = aim();
        if (id) { selected = id; hooks.select(id); paint(); scene.center(); scene.focus(id, false); }
      },
      present() { history = []; scopeStop = pinnedStop = ''; viewLane = ''; cancelTravel(); },
      select(id) { if (id !== selected && angle === 1 && !space.hidden) { selected = id; paint(); updateControls(); } },
      get active() { return !space.hidden; },
      get heads() { return destination === 1; }
    };
  }
  window.WorldCameraController = {create};
})();
