// Carry one selected brane across the two documents. URLs carry navigation;
// the short-lived geometry is only used when both pages can animate it.
(() => {
  'use strict';
  if (typeof window === 'undefined' || typeof document === 'undefined') return;
  const key = 'world.projection-motion';
  const duration = 1080;
  const sideAngle = 68;
  const maxGlyphSize = 22;
  const motionPreference = window.matchMedia('(prefers-reduced-motion: reduce)');
  const reduced = () => motionPreference.matches ||
    document.documentElement.dataset.motion === 'reduced' ||
    document.getElementById('root')?.classList.contains('reduced');
  const plainClick = event => !event.defaultPrevented && event.button === 0 &&
    !event.altKey && !event.ctrlKey && !event.metaKey && !event.shiftKey;
  const box = (x, y, w, h) => ({x, y, w, h});
  const rect = el => {
    const bounds = el?.getBoundingClientRect();
    if (!bounds?.width || !bounds?.height || bounds.right <= 0 || bounds.bottom <= 0 ||
        bounds.left >= innerWidth || bounds.top >= innerHeight) return null;
    return box(bounds.left, bounds.top, bounds.width, bounds.height);
  };
  const midpoint = bounds => ({x: bounds.x + bounds.w / 2, y: bounds.y + bounds.h / 2});
  const limit = (n, lo, hi) => Math.max(lo, Math.min(n, hi));
  const here = url => {
    const address = new URL(url, location.href);
    address.searchParams.sort();
    return address.href;
  };

  function findGlyph(anchor) {
    if (anchor?.matches('.g, svg')) return anchor;
    return anchor?.querySelector('.g, svg') || null;
  }

  function svgInk(element) {
    if (element?.localName !== 'svg') return null;
    try {
      const shape = element.getBBox(), matrix = element.getScreenCTM();
      if (!matrix || !shape.width || !shape.height) return null;
      const pad = Math.min(shape.width, shape.height) * .1;
      const view = new DOMRect(shape.x - pad, shape.y - pad, shape.width + pad * 2, shape.height + pad * 2);
      const corners = [[view.x, view.y], [view.right, view.y],
        [view.x, view.bottom], [view.right, view.bottom]]
        .map(([x, y]) => new DOMPoint(x, y).matrixTransform(matrix));
      const xs = corners.map(point => point.x), ys = corners.map(point => point.y);
      const x = Math.min(...xs), y = Math.min(...ys);
      return {view, bounds: box(x, y, Math.max(...xs) - x, Math.max(...ys) - y)};
    } catch (_) { return null; }
  }

  function measure(anchor, face) {
    const panel = rect(face);
    if (!panel) return null;
    const glyph = findGlyph(anchor);
    const glyphBounds = rect(glyph);
    // Small SVGs already have a deliberate viewport. Large wrappers can
    // include a tiny, offset glyph, so measure their vector ink instead.
    const ink = glyphBounds && Math.max(glyphBounds.w, glyphBounds.h) > maxGlyphSize && svgInk(glyph);
    const bounds = ink?.bounds || glyphBounds;
    const point = midpoint(bounds || rect(anchor) || panel);
    return {
      panel, point: { x: limit(point.x, 0, innerWidth), y: limit(point.y, 0, innerHeight) },
      glyph, bounds, ink
    };
  }

  function save(href, identity, glyph, face, colour, side) {
    try {
      sessionStorage.removeItem(key);
      if (reduced() || !Element.prototype.animate) return;
      const start = measure(glyph, face);
      if (start) sessionStorage.setItem(key, JSON.stringify({
        url: here(href), identity, point: start.point, panel: start.panel, side,
        viewport: { w: innerWidth, h: innerHeight },
        colour: /^#[a-f\d]{3}(?:[a-f\d]{3})?$/i.test(colour) ? colour : '#8496a5',
        expires: Date.now() + 4500
      }));
    } catch (_) { /* Links work without storage or animation support. */ }
  }

  function take() {
    try {
      const from = JSON.parse(sessionStorage.getItem(key));
      sessionStorage.removeItem(key);
      if (from?.url === here(location.href) && from.expires > Date.now() && !reduced()) return from;
    } catch (_) { /* A fresh visit can render without motion. */ }
    return null;
  }

  const position = (element, bounds) => Object.assign(element.style, {
    left: bounds.x + 'px', top: bounds.y + 'px', width: bounds.w + 'px', height: bounds.h + 'px'
  });

  function paneAt(endpoint, width, height) {
    // Keep the glyph on the pane's hinge. A front face opens away from its own
    // edge, and the timeline's row intersects the hinge across its centre.
    const {point, panel, side} = endpoint;
    const x = side ? point.x - width / 2 : Math.min(point.x - 5, panel.x);
    const y = side ? point.y - height / 2 : Math.min(point.y - 5, panel.y);
    return box(limit(x, 4, innerWidth - width - 4) - point.x,
      limit(y, 4, innerHeight - height - 4) - point.y, width, height);
  }

  function footprint(endpoint, pane) {
    const element = document.createElement('div');
    element.className = 'projection-footprint' + (endpoint.side ? ' is-side' : '');
    const {point, panel, side} = endpoint;
    if (side) {
      const width = Math.min(innerWidth - 8, pane.w * 1.3);
      position(element, box(limit(point.x - width / 2, 4, innerWidth - width - 4),
        limit(panel.y + panel.h / 2, point.y - 22, point.y + 22), width, 0));
    } else {
      position(element, box(point.x + pane.x, point.y + pane.y, pane.w, pane.h));
    }
    return element;
  }

  function markerFor(target) {
    const {glyph, bounds, ink} = target;
    if (!glyph || !bounds?.w || !bounds?.h) return null;
    const clone = glyph.cloneNode(true);
    clone.removeAttribute('id');
    clone.removeAttribute('role');
    const paint = getComputedStyle(glyph);
    for (const property of ['background', 'border', 'border-radius', 'clip-path', 'color', 'fill', 'stroke', 'stroke-width']) {
      const value = paint.getPropertyValue(property);
      if (value) clone.style.setProperty(property, value);
    }
    if (ink) {
      const {x, y, width, height} = ink.view;
      clone.setAttribute('viewBox', [x, y, width, height].join(' '));
    }
    const factor = Math.min(1, maxGlyphSize / Math.max(bounds.w, bounds.h));
    const width = bounds.w * factor, height = bounds.h * factor;
    const marker = document.createElement('div');
    marker.className = 'projection-marker';
    position(marker, box(-width / 2, -height / 2, width, height));
    marker.append(clone);
    return {marker, mask: factor === 1 ? glyph : null};
  }

  let stop = null;
  function animate(from, identity, glyph, face, side) {
    if (!from || from.identity !== identity || from.expires <= Date.now() ||
        reduced() || !Element.prototype.animate) return;
    const target = measure(glyph, face);
    if (!target) return;
    target.side = Boolean(side);
    stop?.();

    const sx = innerWidth / (from.viewport?.w || innerWidth);
    const sy = innerHeight / (from.viewport?.h || innerHeight);
    const start = {
      point: {x: from.point.x * sx, y: from.point.y * sy},
      panel: box(from.panel.x * sx, from.panel.y * sy, from.panel.w * sx, from.panel.h * sy),
      side: Boolean(from.side)
    };
    const turning = start.side !== target.side;
    const width = Math.min(innerWidth - 8, 220);
    const height = Math.min(innerHeight - 8, 96);
    const origin = paneAt(start, width, height), destination = paneAt(target, width, height);
    const carrier = document.createElement('div');
    carrier.className = 'projection-carrier';
    carrier.setAttribute('aria-hidden', 'true');
    carrier.inert = true;
    carrier.style.setProperty('--projection-ink', from.colour);
    const originCue = footprint(start, origin);
    const destinationCue = footprint(target, destination);
    const flight = document.createElement('div');
    flight.className = 'projection-flight';
    const plane = document.createElement('div');
    plane.className = 'projection-edge';
    const movingGlyph = markerFor(target);
    flight.append(plane);
    if (movingGlyph) flight.append(movingGlyph.marker);
    carrier.append(originCue, destinationCue, flight);
    document.body.append(carrier);

    const travel = (ratio, lift = 0) => {
      const x = start.point.x + (target.point.x - start.point.x) * ratio;
      const y = start.point.y + (target.point.y - start.point.y) * ratio - lift;
      return `translate3d(${x}px, ${y}px, 0)`;
    };
    const frame = (pane, yaw, pitch, opacity) => ({
      left: pane.x + 'px', top: pane.y + 'px', width: pane.w + 'px', height: pane.h + 'px',
      transformOrigin: `${-pane.x}px ${-pane.y}px`,
      transform: `perspective(700px) rotateY(${yaw}deg) rotateX(${pitch}deg)`,
      opacity
    });
    const begin = start.side ? sideAngle : 0, end = target.side ? sideAngle : 0;
    const between = box((origin.x + destination.x) / 2, (origin.y + destination.y) / 2, width, height);
    const lift = turning ? 0 : Math.min(36, innerHeight * .07);
    const ease = 'cubic-bezier(.37,.03,.23,1)';
    const soft = 'cubic-bezier(.22,.68,.34,1)';
    const settle = .83;
    const animations = [];
    const run = (element, frames) => {
      const animation = element.animate(frames, {duration, fill: 'both'});
      animations.push(animation);
      return animation;
    };
    const preferenceChanged = () => { if (reduced()) cleanup(); };
    const observer = new MutationObserver(preferenceChanged);
    function cleanup() {
      if (stop !== cleanup) return;
      stop = null;
      observer.disconnect();
      motionPreference.removeEventListener('change', preferenceChanged);
      window.removeEventListener('pagehide', cleanup);
      for (const animation of animations) animation.cancel();
      carrier.remove();
    }
    stop = cleanup;
    try {
      motionPreference.addEventListener('change', preferenceChanged);
      window.addEventListener('pagehide', cleanup, {once: true});
      observer.observe(document.documentElement, {attributes:true, attributeFilter:['data-motion']});
      const root = document.getElementById('root');
      if (root) observer.observe(root, {attributes:true, attributeFilter:['class']});
      const motion = run(flight, [
        {transform: travel(0)},
        {transform: travel(0), offset: .13, easing: ease},
        {transform: travel(.57, lift), offset: .52, easing: soft},
        {transform: travel(1), offset: settle},
        {transform: travel(1)}
      ]);
      run(plane, [
        frame(origin, begin, 0, .82),
        {...frame(origin, begin, 0, .92), offset: .13, easing: ease},
        {...frame(between, (begin + end) / 2, turning ? 0 : -20, .96), offset: .52, easing: soft},
        {...frame(destination, end, 0, .9), offset: settle},
        frame(destination, end, 0, 0)
      ]);
      run(originCue, [{opacity: .66}, {opacity: .52, offset: .18}, {opacity: 0, offset: .46}, {opacity: 0}]);
      run(destinationCue, [{opacity: 0}, {opacity: 0, offset: .52}, {opacity: .52, offset: settle}, {opacity: 0}]);
      if (movingGlyph) {
        run(movingGlyph.marker, [{opacity: 1}, {opacity: 1, offset: settle}, {opacity: 0}]);
        if (movingGlyph.mask) run(movingGlyph.mask, [{opacity: 0}, {opacity: 0, offset: settle}, {opacity: 1}]);
      }
      motion.finished.then(cleanup, () => {});
    } catch (_) { cleanup(); }
  }
  window.WorldProjection = {plainClick, save, take, animate};
})();
