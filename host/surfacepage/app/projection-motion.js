// Carry one selected brane across the two documents. URLs carry navigation;
// the short-lived geometry is only used when both pages can animate it.
(() => {
  'use strict';
  if (typeof window === 'undefined' || typeof document === 'undefined') return;
  const key = 'world.projection-motion';
  const reduced = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches ||
    document.documentElement.dataset.motion === 'reduced' ||
    document.getElementById('root')?.classList.contains('reduced');
  const plainClick = event => !event.defaultPrevented && event.button === 0 &&
    !event.altKey && !event.ctrlKey && !event.metaKey && !event.shiftKey;
  const rect = el => {
    const box = el?.getBoundingClientRect();
    if (!box?.width || !box?.height || box.right <= 0 || box.bottom <= 0 ||
        box.left >= innerWidth || box.top >= innerHeight) return null;
    return { x: box.left, y: box.top, w: box.width, h: box.height };
  };
  const here = url => {
    const address = new URL(url, location.href);
    address.searchParams.sort();
    return address.href;
  };
  function save(href, identity, glyph, face, colour, side) {
    try {
      sessionStorage.removeItem(key);
      if (reduced() || !Element.prototype.animate) return;
      const point = rect(glyph), panel = rect(face);
      if (point && panel) sessionStorage.setItem(key, JSON.stringify({
        url: here(href), identity, point, panel, side,
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
  const position = (element, box) => Object.assign(element.style, {
    left: box.x + 'px', top: box.y + 'px', width: box.w + 'px', height: box.h + 'px'
  });
  function animate(from, identity, glyph, face, side) {
    if (!from || from.identity !== identity || from.expires <= Date.now() ||
        reduced() || !Element.prototype.animate) return;
    const point = rect(glyph), panel = rect(face);
    if (!point || !panel) return;
    const carrier = document.createElement('div');
    carrier.className = 'projection-carrier';
    carrier.setAttribute('aria-hidden', 'true');
    carrier.style.setProperty('--projection-ink', from.colour);
    const edge = document.createElement('div'), marker = document.createElement('div');
    edge.className = 'projection-edge';
    marker.className = 'projection-marker';
    position(edge, panel);
    position(marker, point);
    const sameGlyph = glyph.cloneNode(true);
    sameGlyph.removeAttribute('id');
    marker.append(sameGlyph);
    carrier.append(edge, marker);
    document.body.append(carrier);
    const duration = 520, easing = 'cubic-bezier(.22,.72,.2,1)';
    const move = (a, b) => `translate(${a.x - b.x}px, ${a.y - b.y}px) scale(${a.w / b.w}, ${a.h / b.h})`;
    const frame = (box, angle, opacity) => ({
      left: box.x + 'px', top: box.y + 'px', width: box.w + 'px', height: box.h + 'px',
      transform: `perspective(1000px) rotateY(${angle}deg)`, opacity
    });
    const turn = (from.panel.x <= panel.x ? 1 : -1) * 55;
    const midpoint = {};
    for (const axis of ['x', 'y', 'w', 'h']) midpoint[axis] = (from.panel[axis] + panel[axis]) / 2;
    try {
      edge.animate([
        frame(from.panel, from.side ? turn : 0, .9),
        {'offset': .48, ...frame(midpoint, turn * (side ? -.7 : .7), .85)},
        frame(panel, side ? -turn : 0, 0)
      ], { duration, easing });
      const flight = marker.animate([
        { transform: move(from.point, point) + ' scale(1.35)', opacity: 1 },
        { opacity: 1, offset: .77 },
        { transform: 'none', opacity: 0 }
      ], { duration, easing });
      glyph.animate([
        { opacity: 0 }, { opacity: 0, offset: .77 }, { opacity: 1 }
      ], { duration, easing });
      flight.finished.then(() => carrier.remove(), () => carrier.remove());
    } catch (_) { carrier.remove(); }
  }
  window.WorldProjection = { plainClick, save, take, animate };
})();
