// Empty ground is a grab surface. Cards and the scene controls retain their own gestures.
(() => {
  if (typeof window === 'undefined') return;
  function create(scene, hooks) {
    const scroll = scene.querySelector('.hc-scroll');
    const floor = document.createElement('div');
    floor.className = 'hc-floor-hit';
    floor.setAttribute('aria-hidden', 'true');
    scene.querySelector('.hc-space').prepend(floor);
    const canvas = scene.querySelector('.hc-graph');
    const isGround = target => target === floor || target === canvas;
    let press, tap, dragged = false, lastTouchHop = -Infinity;

    // Find the closest recorded stop on a monotone axis without scanning the graph on every move.
    function nearest(stops, target, value) {
      let lo = 0, hi = stops.length - 1;
      while (lo < hi) {
        const mid = (lo + hi) >>> 1;
        if (value(stops[mid]) < target) lo = mid + 1;
        else hi = mid;
      }
      const a = stops[Math.max(0, lo - 1)], b = stops[lo];
      return Math.abs(value(a) - target) < Math.abs(value(b) - target) ? a : b;
    }
    function groundPoint(event) {
      const view = hooks.view();
      if (!view || !view.stops.length) return;
      const rect = scroll.getBoundingClientRect();
      const x = event.clientX - rect.left + scroll.scrollLeft, y = event.clientY - rect.top + scroll.scrollTop;
      const g = view.layout.ground;
      if (y <= g.horizon) return;
      return {...view, x, y, depth:(g.streetY - g.horizon) / (y - g.horizon)};
    }
    function hop(event) {
      const hit = groundPoint(event);
      if (!hit) return;
      const {ground, time} = hit.layout;
      const visible = hit.stops.filter(stop => stop.index <= time.referenceIndex);
      if (!visible.length) return;
      const stop = nearest(visible, hit.y, item => time.groundYAt(item.index));
      const worldX = ground.cameraX + (hit.x - ground.vx) * time.depthAt(stop.index);
      hooks.hop(stop, worldX);
    }
    function cancel() {
      scene.classList.remove('is-floor-dragging');
      press = null;
    }
    scene.addEventListener('pointerdown', event => {
      if (!event.isPrimary) { cancel(); tap = null; return; }
      if (!isGround(event.target) || event.button !== 0 || event.altKey || event.metaKey || event.ctrlKey || event.shiftKey) {
        tap = null;
        return;
      }
      const hit = groundPoint(event);
      if (!hit) return;
      const {ground, time} = hit.layout;
      // The local slope of the ground projection makes a near grab precise
      // and a distant grab faster. Limit the slope at the vanishing point.
      const scale = time.unitsPerZ * Math.min(3.5, hit.depth) ** 2 / (ground.streetY - ground.horizon);
      press = {id:event.pointerId, x:event.clientX, y:event.clientY, offset:scroll.scrollLeft - ground.vx, lane:ground.centerLane,
        stops:hit.stops, time, scale, index:time.referenceIndex, dragged:false};
    });
    window.addEventListener('pointermove', event => {
      if (event.pointerId !== press?.id) return;
      if (event.pointerType === 'mouse' && !(event.buttons & 1)) { cancel(); return; }
      const dx = event.clientX - press.x, dy = event.clientY - press.y;
      if (!press.dragged && Math.hypot(dx, dy) < 6) return;
      if (!press.dragged) {
        press.dragged = true;
        scene.classList.add('is-floor-dragging');
      }
      tap = null;
      const {time} = press;
      const stop = nearest(press.stops, time.referenceX - dy * press.scale, s => time.xAt(s.index));
      if (stop.index !== press.index) {
        press.index = stop.index;
        hooks.seek(stop.id, press.lane);
      }
      const view = hooks.view();
      if (view) scroll.scrollLeft = view.layout.ground.vx + press.offset - dx;
      event.preventDefault();
      event.stopPropagation();
    }, {capture:true});
    window.addEventListener('pointerup', event => {
      if (event.pointerId !== press?.id) return;
      dragged = press.dragged;
      if (!press.dragged && event.pointerType === 'touch') {
        const now = performance.now();
        if (tap && now - tap.at <= 360 && Math.hypot(tap.x - event.clientX, tap.y - event.clientY) < 20) {
          hop(event);
          lastTouchHop = now;
          tap = null;
        } else tap = {at:now, x:event.clientX, y:event.clientY};
      }
      cancel();
    }, {capture:true});
    scene.addEventListener('dblclick', event => {
      if (!isGround(event.target) || dragged || event.button !== 0 || event.altKey || event.metaKey || event.ctrlKey || event.shiftKey ||
          performance.now() - lastTouchHop < 400) return;
      event.preventDefault();
      hop(event);
    });
    window.addEventListener('pointercancel', event => { if (event.pointerId === press?.id) cancel(); }, {capture:true});
    window.addEventListener('blur', () => { cancel(); tap = null; });
  }
  window.WorldHeadGroundInput = {create};
})();
