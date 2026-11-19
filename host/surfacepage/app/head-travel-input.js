// Translate scene-local gestures into discrete moves along the recorded graph.
(() => {
  'use strict';
  if (typeof window === 'undefined') return;

  function create(scene, hooks) {
    const events = new AbortController();
    const controls = document.createElement('div');
    controls.className = 'hc-travel';
    controls.setAttribute('role', 'group');
    controls.setAttribute('aria-label', 'History depth');
    controls.innerHTML = `
      <span class="hc-travel-label" role="status"></span>
      <button type="button" class="hc-travel-reset" aria-label="Return to present">Present</button>
      <button type="button" class="hc-travel-out" aria-label="Move toward present" title="Move toward present (−)">−</button>
      <button type="button" class="hc-travel-in" aria-label="Move into history" title="Move into history (+)">+</button>`;
    controls.hidden = true;
    scene.append(controls);

    const label = controls.querySelector('.hc-travel-label');
    const older = controls.querySelector('.hc-travel-in');
    const newer = controls.querySelector('.hc-travel-out');
    const reset = controls.querySelector('.hc-travel-reset');
    let state = {active:false, canIn:false, canOut:false};
    let wheelAt = 0, wheelDirection = 0, wheelAmount = 0, lastWheelStep = -Infinity;
    let pinch = null, suppressClickUntil = 0;
    const active = () => state.active && !scene.closest('[hidden], [inert]');

    function listen(target, type, handler, options) {
      target.addEventListener(type, handler, {...options, signal:events.signal});
    }

    function sceneTarget(target) {
      return target instanceof Element && scene.contains(target) && !controls.contains(target) ? target : null;
    }

    function pointerTarget(x, y, original) {
      return sceneTarget(document.elementFromPoint(x, y)) || sceneTarget(original);
    }

    function available(direction, pointed) {
      return active() && (direction > 0 ? state.canIn || Boolean(pointed?.closest('[data-head-state], [data-rear-state]')) : state.canOut);
    }

    function step(direction, target) {
      if (available(direction)) hooks.step(direction, target);
    }

    listen(older, 'click', () => step(1, null));
    listen(newer, 'click', () => step(-1, null));
    listen(reset, 'click', () => {
      if (available(-1)) hooks.reset();
    });
    listen(scene, 'wheel', event => {
      if (!active() || event.defaultPrevented || event.shiftKey || event.metaKey || event.ctrlKey ||
          Math.abs(event.deltaX) >= Math.abs(event.deltaY)) return;
      event.preventDefault();
      const now = performance.now();
      const direction = Math.sign(event.deltaY);
      if (now - wheelAt > 160 || direction !== wheelDirection) {
        wheelAmount = 0;
        lastWheelStep = -Infinity;
      }
      wheelAt = now;
      wheelDirection = direction;
      const unit = event.deltaMode === 1 ? 16 : event.deltaMode === 2 ? scene.clientHeight : 1;
      const delta = Math.abs(event.deltaY) * unit;
      if (delta < 4 || now - lastWheelStep < 220) return;
      wheelAmount += delta;
      if (wheelAmount < (Number.isFinite(lastWheelStep) ? 100 : 42)) return;
      const target = pointerTarget(event.clientX, event.clientY, event.target);
      if (!available(direction, target)) return;
      lastWheelStep = now;
      wheelAmount = 0;
      hooks.step(direction, target);
    }, {passive:false});

    function measure(touches) {
      const [a, b] = touches;
      return {distance:Math.hypot(a.clientX - b.clientX, a.clientY - b.clientY),
        x:(a.clientX + b.clientX) / 2, y:(a.clientY + b.clientY) / 2};
    }
    listen(scene, 'touchstart', event => {
      if (!active() || event.touches.length !== 2 || controls.contains(event.target)) return;
      const start = measure(event.touches);
      suppressClickUntil = performance.now() + 450;
      pinch = {distance:Math.max(1, start.distance),
        target:pointerTarget(start.x, start.y, event.target), stepped:false};
      if (event.cancelable) event.preventDefault();
    }, {passive:false});
    listen(scene, 'touchmove', event => {
      if (!active() || !pinch || event.touches.length !== 2) return;
      if (event.cancelable) event.preventDefault();
      suppressClickUntil = performance.now() + 450;
      if (pinch.stepped) return;
      const ratio = measure(event.touches).distance / pinch.distance;
      const direction = ratio >= 1.18 ? 1 : ratio <= .85 ? -1 : 0;
      const target = sceneTarget(pinch.target);
      if (!direction || !available(direction, target)) return;
      pinch.stepped = true;
      hooks.step(direction, target);
    }, {passive:false});
    function finishTouch(event) {
      if (!pinch) return;
      if (event.cancelable) event.preventDefault();
      suppressClickUntil = performance.now() + 450;
      if (event.touches.length < 2) pinch = null;
    }
    listen(scene, 'touchend', finishTouch, {passive:false});
    listen(scene, 'touchcancel', finishTouch, {passive:false});
    listen(scene, 'click', event => {
      if (event.detail === 0 || performance.now() >= suppressClickUntil) return;
      event.preventDefault();
      event.stopImmediatePropagation();
    }, {capture:true});

    listen(window, 'keydown', event => {
      if (!active() || event.defaultPrevented || event.isComposing || event.repeat || event.altKey ||
          event.metaKey || event.ctrlKey || event.target?.closest?.('dialog, [role="dialog"], input, textarea, select, [contenteditable]:not([contenteditable="false"]), [role="textbox"]')) return;
      const direction = event.key === '+' || event.key === '=' ? 1 : event.key === '-' ? -1 : 0;
      if (!direction || !available(direction)) return;
      event.preventDefault();
      step(direction, sceneTarget(document.activeElement));
    });

    return {
      update(next) {
        const becameActive = next.active && !state.active;
        state = next;
        controls.hidden = !state.active;
        scene.classList.toggle('hc-travel-active', state.active);
        if (becameActive) {
          wheelAmount = 0;
          lastWheelStep = -Infinity;
          wheelAt = 0;
          pinch = null;
        }
        older.setAttribute('aria-disabled', String(!state.canIn));
        newer.setAttribute('aria-disabled', String(!state.canOut));
        if (reset === document.activeElement && !state.canOut && state.active) older.focus({preventScroll:true});
        reset.hidden = !state.canOut;
        const description = state.label || '';
        if (label.textContent !== description) label.textContent = description;
        label.title = description;
        label.hidden = !description;
      },
      destroy() {
        events.abort();
        controls.remove();
        scene.classList.remove('hc-travel-active');
      },
    };
  }

  window.WorldHeadTravelInput = {create};
})();
