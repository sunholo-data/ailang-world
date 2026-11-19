// DOM faces sit on the same projected points as the graph painted beneath them.
(() => {
  'use strict';
  if (typeof window === 'undefined') return;

  const escape = value => String(value ?? '').replace(/[&<>"']/g, ch => ({'&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;'}[ch]));
  const clamp = (value, lo, hi) => Math.max(lo, Math.min(hi, value));
  const statePath = node => node.draft ? '/surface#' + new URLSearchParams({draft:node.id}) : '/surface?' + new URLSearchParams({brane: node.lane, at: node.id}) + '#heads';
  const accent = lane => /^#[0-9a-f]{3}([0-9a-f]{3})?$/i.test(lane?.colour) ? lane.colour : '#657b91';
  const shortAge = value => String(value ?? '').replace(/\b(minutes?|mins?)\b/i, 'm')
    .replace(/\b(hours?|hrs?)\b/i, 'h').replace(/\bdays?\b/i, 'd')
    .replace(/\bweeks?\b/i, 'w').replace(/\s+(?=[mhdw]\b)|\s+ago$/gi, '');
  const exactTime = new Intl.DateTimeFormat(undefined, {year:'numeric', month:'short', day:'numeric', hour:'2-digit', minute:'2-digit', timeZone:'UTC'});

  function create(container, options) {
    const scene = document.createElement('section');
    scene.className = 'hc-scene';
    scene.setAttribute('aria-label', 'Branch heads');
    scene.hidden = true;
    scene.innerHTML = '<canvas class="hc-graph" aria-hidden="true"></canvas><div class="hc-scroll"><div class="hc-space"><div class="hc-rear-layer" hidden></div></div></div><canvas class="hc-inspection" aria-hidden="true"></canvas><div class="hc-inspect-tip" aria-hidden="true" hidden></div>';
    container.append(scene);

    const scroll = scene.querySelector('.hc-scroll'), space = scene.querySelector('.hc-space'), canvas = scene.querySelector('canvas');
    const ctx = canvas.getContext('2d');
    const inspection = scene.querySelector('.hc-inspection'), tip = scene.querySelector('.hc-inspect-tip');
    const inspectionInk = inspection.getContext('2d');
    const pickable = [];
    let inspectable = false, pointer = null;
    const inspected = () => {
      const id = hovered || focusedByUser;
      return inspectable && model?.N[id] ? id : '';
    };
    const rearLayer = space.querySelector('.hc-rear-layer');
    const cards = new Map(), rearLinks = new Map();
    let lips = new Set();
    const stamps = new Map(), paneTitles = new Map(), crossings = new Map();
    const neighborhood = new Map(), incident = new Map();
    const scratch = document.createElement('template');
    let model, layout, sidePoints, angle = 0, selected = '', hovered = '', frame = 0, silentFocus = false;
    let initialPan = true, focusedByUser = '', paintedSelection, paintWidth = 0, paintHeight = 0, pixelRatio = 0;
    let floorPhase = 0;
    const faceScales = new Map();
    let points = new Map(), headsByID = new Map(), headsByLane = new Map(), order = [], histories = new Map();

    // The top and the inward-facing wall are part of the native link's hit area.
    function facade(el, head, scale = 1) {
      const dx = (head.returnX || 0) * scale, rise = Math.max(0, -head.returnY || 0) * scale;
      el.classList.toggle('hc-left-wall', dx < 0);
      el.style.setProperty('--hc-lane', accent(model.LANES[head.lane]));
      el.style.setProperty('--hc-rise', rise + 'px');
      el.style.setProperty('--hc-return', Math.abs(dx) + 'px');
      el.style.setProperty('--hc-roof-skew', -Math.atan2(dx, rise) + 'rad');
      el.style.setProperty('--hc-wall-skew', (dx ? Math.atan(-rise / dx) : 0) + 'rad');
    }

    function setLayout(T, next) {
      if (model === T && layout === next) return;
      model = T;
      layout = next;
      headsByID = new Map(layout.heads.map(head => [head.id, head]));
      headsByLane = new Map(layout.heads.map(head => [head.lane, head]));
      order = [...layout.points.keys()].sort((a, b) => (layout.points.get(b).depth || 0) - (layout.points.get(a).depth || 0));
      histories = new Map(layout.heads.map(head => [head.lane, []]));
      for (const id of [...order].reverse()) {
        const node = model.N[id];
        if (node && headsByLane.get(node.lane)?.id !== id) histories.get(node.lane)?.push(id);
      }
      neighborhood.clear();
      incident.clear();
      crossings.clear();
      for (const edge of layout.edges) {
        if (!neighborhood.has(edge.from)) neighborhood.set(edge.from, new Set());
        if (!neighborhood.has(edge.to)) neighborhood.set(edge.to, new Set());
        neighborhood.get(edge.from).add(edge.to);
        neighborhood.get(edge.to).add(edge.from);
        for (const id of [edge.from, edge.to]) {
          if (!incident.has(id)) incident.set(id, []);
          incident.get(id).push(edge);
        }
        const from = model.N[edge.from]?.lane, to = model.N[edge.to]?.lane;
        if (from !== to) {
          const key = [from, to].sort().join('|');
          crossings.set(key, (crossings.get(key) || 0) + 1);
        }
      }
      stamps.clear();
      paneTitles.clear();
      for (const [id, entry] of cards) {
        if (headsByID.has(id)) continue;
        entry.el.remove();
        cards.delete(id);
      }
      layout.heads.forEach((head, index) => {
        const node = model.N[head.id];
        if (!node) return;
        const lane = model.LANES[head.lane];
        const label = model.laneTitle?.(head.lane) || lane?.name || head.lane;
        const title = options.title(node) || node.title || node.label || node.id;
        const stamp = Date.parse(node.at);
        const compact = head.height <= 48;
        const rel = node.rel || (Number.isFinite(stamp) ? exactTime.format(stamp) : '');
        const shownTime = compact ? shortAge(rel) : rel;
        const status = head.attention ? 'Needs you' : head.live ? 'Live' : '';
        const identity = status || node.open?.label || node.episode || (node.pr ? '#' + node.pr : '');
        const timing = Number.isFinite(stamp)
          ? `<time datetime="${escape(node.at)}" title="${escape(exactTime.format(stamp) + ' UTC')}">${escape(shownTime)}</time>`
          : `<span>${escape(shownTime)}</span>`;
        const bar = `<span class="hc-face-bar">${options.glyph(node, compact ? 10 : 12)}<span class="hc-lane-name">${escape(label)}</span><span class="hc-age">${compact && status ? escape(status) : timing}</span></span>`;
        const markup = compact ? bar : `${bar}<strong>${escape(title)}</strong>${head.height > 76 && identity ? `<span class="hc-identity">${escape(identity)}</span>` : ''}`;
        let entry = cards.get(head.id);
        if (!entry) {
          const el = document.createElement('a');
          el.className = 'hc-face';
          el.dataset.headState = head.id;
          el.dataset.headLane = head.lane;
          entry = {el};
          cards.set(head.id, entry);
        }
        const {el} = entry;
        if (entry.markup !== markup) {
          el.innerHTML = markup;
          entry.markup = markup;
        }
        el.classList.toggle('needs-human', Boolean(head.attention));
        el.classList.toggle('is-draft', Boolean(node.draft));
        el.classList.toggle('is-compact', compact);
        el.classList.toggle('is-middle', !compact && head.height <= 76);
        el.href = statePath(node);
        el.title = label + ' · ' + title;
        el.style.width = head.width + 'px';
        el.style.height = head.height + 'px';
        facade(el, head);
        el.style.zIndex = String(layout.heads.length - index);
        if (space.children[index + 1] !== el) {
          const focused = document.activeElement === el;
          space.insertBefore(el, space.children[index + 1] || null);
          if (focused && document.activeElement !== el) {
            silentFocus = true;
            el.focus({preventScroll: true});
            silentFocus = false;
          }
        }
      });
      lips = new Set();
      for (const head of layout.heads) {
        let previousY = head.y - head.height / 2;
        const limit = head.height <= 48 ? 2 : 3;
        let count = 0;
        for (const id of histories.get(head.lane)) {
          const projected = layout.points.get(id);
          if (previousY - projected.y < 20) continue;
          lips.add(id);
          previousY = projected.y;
          if (++count === limit) break;
        }
      }
      syncRearLinks();
      space.style.width = Math.ceil(layout.width) + 'px';
      space.style.height = Math.ceil(layout.height) + 'px';
      paintedSelection = null;
    }

    function glyphInfo(node) {
      const key = node.lane + '|' + node.kind;
      if (!stamps.has(key)) {
        scratch.innerHTML = options.glyph(node, 10);
        const style = scratch.content.firstElementChild?.style;
        const polygon = style?.clipPath.match(/-?\d+(?:\.\d+)?%/g)?.map(value => parseFloat(value) / 100);
        stamps.set(key, {
          colour: style?.borderTopColor || style?.backgroundColor || accent(model.LANES[node.lane]),
          hollow: Boolean(style?.borderTopWidth),
          round: style?.borderTopLeftRadius === '50%',
          polygon,
        });
      }
      return stamps.get(key);
    }

    function drawGlyph(node, point, alpha) {
      const {colour, hollow, round, polygon} = glyphInfo(node);
      const size = 2 + 7 * clamp(point.scale ?? 1, 0, 1);
      const x = point.x - size / 2, y = point.y - size / 2;
      ctx.globalAlpha = alpha;
      ctx.fillStyle = colour;
      ctx.strokeStyle = colour;
      ctx.lineWidth = 1.4;
      ctx.beginPath();
      if (polygon && polygon.length >= 6) {
        ctx.moveTo(x + size * polygon[0], y + size * polygon[1]);
        for (let i = 2; i < polygon.length; i += 2) ctx.lineTo(x + size * polygon[i], y + size * polygon[i + 1]);
        ctx.closePath();
      } else if (round) ctx.arc(point.x, point.y, size / 2, 0, Math.PI * 2);
      else ctx.rect(x, y, size, size);
      if (hollow) ctx.stroke();
      else ctx.fill();
    }

    function paneWidth(head, point) {
      return Math.max(24, head.width * clamp(point.scale ?? .45, .26, .9) * .88);
    }

    function syncRearLinks() {
      const titled = new Set();
      for (const id of order) {
        if (!lips.has(id)) continue;
        const node = model.N[id];
        const point = layout.points.get(id);
        const width = paneWidth(headsByLane.get(node.lane), point);
        if (width < 63) continue;
        const height = 17;
        const x = point.x - width / 2, y = point.y - height / 2;
        const covered = layout.heads.some(head =>
          x >= head.x - head.width / 2 && x + width <= head.x + head.width / 2 &&
          y >= head.y - head.height / 2 && y + height <= head.y + head.height / 2);
        if (covered) continue;
        titled.add(id);
        let el = rearLinks.get(id);
        if (!el) {
          el = document.createElement('a');
          el.className = 'hc-rear';
          el.append(document.createElement('span'));
          el.dataset.rearState = id;
          rearLinks.set(id, el);
        }
        const title = options.title(node) || node.title || node.label || id;
        const label = model.laneTitle?.(node.lane) || node.lane;
        el.firstElementChild.textContent = title;
        facade(el, headsByLane.get(node.lane), Math.min(.35, width / 560));
        el.title = title;
        el.setAttribute('aria-label', label + ' · ' + title);
        el.href = statePath(node);
        el.style.left = x + 'px';
        el.style.top = y + 'px';
        el.style.width = width + 'px';
        el.style.height = height + 'px';
        const index = titled.size - 1;
        if (rearLayer.children[index] !== el) rearLayer.insertBefore(el, rearLayer.children[index] || null);
      }
      for (const [id, el] of rearLinks) {
        if (titled.has(id)) continue;
        el.remove();
        rearLinks.delete(id);
      }
    }

    function drawPane(node, point, head, alpha, theme) {
      const width = paneWidth(head, point);
      const height = width >= 63 ? 17 : 9;
      const x = point.x - width / 2, y = point.y - height / 2;
      ctx.globalAlpha = alpha;
      ctx.fillStyle = theme.pane;
      ctx.strokeStyle = theme.edge;
      ctx.lineWidth = 1;
      ctx.beginPath();
      ctx.moveTo(x, y + height);
      ctx.lineTo(x + 4, y);
      ctx.lineTo(x + width - 4, y);
      ctx.lineTo(x + width, y + height);
      ctx.closePath();
      ctx.fill();
      ctx.stroke();
      if (width < 63 || angle < .86) return;
      const text = options.title(node) || node.title || node.label || '';
      ctx.font = '500 9px ' + theme.font;
      ctx.textBaseline = 'middle';
      const key = node.id + ':' + Math.floor(width / 4);
      let visible = paneTitles.get(key);
      if (visible === undefined) {
        visible = text;
        if (ctx.measureText(visible).width > width - 18) {
          let lo = 0, hi = visible.length;
          while (lo < hi) {
            const mid = Math.ceil((hi + lo) / 2);
            if (ctx.measureText(text.slice(0, mid) + '…').width > width - 18) hi = mid - 1;
            else lo = mid;
          }
          visible = text.slice(0, lo) + '…';
        }
        paneTitles.set(key, visible);
      }
      ctx.fillStyle = theme.text;
      ctx.fillText(visible, x + 8, point.y + .5);
    }

    function drawLane(head, start, trail, alpha, theme) {
      const end = points.get(trail.at(-1));
      if (!end) return;
      const reveal = clamp((angle - .55) / .45, 0, 1);
      const startWidth = head.width * (.35 + .65 * reveal) / 2;
      const endWidth = paneWidth(head, end) * .11;
      const top = start.y - head.height * .5 * reveal + 2;
      ctx.strokeStyle = theme.edge;
      ctx.globalAlpha = alpha;
      ctx.lineWidth = 1;
      ctx.beginPath();
      ctx.moveTo(start.x - startWidth, top);
      ctx.lineTo(end.x - endWidth, end.y);
      ctx.moveTo(start.x + startWidth, top);
      ctx.lineTo(end.x + endWidth, end.y);
      ctx.stroke();
    }

    // Transverse rulings share the scene's horizon. Panning keeps them
    // fixed in world space; passing one state advances them by one interval.
    function drawGround(left, right, opacity, ink) {
      const {horizon, front, pitch, street, streetY} = layout.ground;
      const center = layout.width / 2;
      ctx.strokeStyle = ink;
      ctx.lineWidth = 1;
      // These two ground lines frame the space between the banks of branes.
      // The street is world geometry, so it stays put when you pan sideways.
      if (street && streetY > horizon) {
        const bottom = Math.max(front, layout.height);
        const nearHalf = street / 2 * (bottom - horizon) / (streetY - horizon);
        const horizonFade = ctx.createLinearGradient(0, horizon, 0, streetY);
        horizonFade.addColorStop(0, 'transparent');
        horizonFade.addColorStop(.43, ink);
        horizonFade.addColorStop(1, ink);
        ctx.beginPath();
        ctx.moveTo(center, horizon);
        ctx.lineTo(center - nearHalf, bottom);
        ctx.lineTo(center + nearHalf, bottom);
        ctx.closePath();
        ctx.fillStyle = horizonFade;
        ctx.globalAlpha = opacity * .12;
        ctx.fill();
        ctx.beginPath();
        ctx.moveTo(center, horizon);
        ctx.lineTo(center - nearHalf, bottom);
        ctx.moveTo(center, horizon);
        ctx.lineTo(center + nearHalf, bottom);
        ctx.strokeStyle = horizonFade;
        ctx.globalAlpha = opacity * .62;
        ctx.lineWidth = 1.25;
        ctx.stroke();
        ctx.lineWidth = 1;
        ctx.strokeStyle = ink;
      }
      for (let row = -1; row <= 8; row++) {
        const depth = row - floorPhase;
        if (depth < -.5) continue;
        const scale = 1 / (1 + depth * .42);
        const y = horizon + (front - horizon) * scale;
        const spacing = pitch * scale;
        const fade = opacity * clamp(8 - depth, 0, 1);
        ctx.globalAlpha = fade * (.19 + .27 * scale);
        ctx.beginPath();
        ctx.moveTo(left, y);
        ctx.lineTo(right, y);
        ctx.stroke();
        ctx.globalAlpha = fade * (.28 + .30 * scale);
        ctx.beginPath();
        const first = Math.ceil((left - center) / spacing);
        const last = Math.floor((right - center) / spacing);
        for (let column = first; column <= last; column++) {
          const x = center + column * spacing;
          const advance = .022 * scale;
          ctx.moveTo(x - (x - center) * advance, y - (y - horizon) * advance);
          ctx.lineTo(x + (x - center) * advance, y + (y - horizon) * advance);
        }
        ctx.stroke();
      }
    }

    function drawFootprints(head, point, trail, activeLane, activeID, opacity, ink) {
      const active = head.lane === activeLane;
      const size = faceScales.get(head.id) || 1;
      const bottom = point.y + head.height * size / 2;
      const foot = bottom + 2 * size;
      const half = (head.width * .5 - 2) * size;
      ctx.strokeStyle = ink;
      if (active) ctx.strokeStyle = `color-mix(in srgb, ${ink} 55%, ${glyphInfo(model.N[head.id]).colour})`;
      ctx.lineWidth = active ? 1.4 : 1;
      ctx.globalAlpha = opacity * (active ? .85 : .54);
      ctx.beginPath();
      ctx.moveTo(point.x - half, foot);
      ctx.lineTo(point.x + half, foot);
      if (head.returnX) {
        const sideX = point.x + Math.sign(head.returnX) * half;
        ctx.moveTo(sideX, foot);
        ctx.lineTo(sideX + head.returnX * size, foot + head.returnY * size);
      }
      ctx.stroke();
      if (!active) return;
      let lastX = point.x, lastY = foot;
      for (const id of trail) {
        const p = points.get(id);
        const scale = clamp((p.y - layout.ground.horizon) / Math.max(1, point.y - layout.ground.horizon), 0, 1);
        const y = p.y + (foot - point.y) * scale;
        if (id !== activeID && lastY - y < 11) continue;
        ctx.globalAlpha = opacity * (.10 + .12 * scale);
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(lastX, lastY);
        ctx.lineTo(p.x, y);
        ctx.stroke();
        lastX = p.x;
        lastY = y;
        const span = 3 + 5 * scale;
        ctx.globalAlpha = opacity * (id === activeID ? .85 : .22 + .22 * scale);
        ctx.beginPath();
        ctx.moveTo(p.x - span, y);
        ctx.lineTo(p.x + span, y);
        ctx.stroke();
      }
    }

    function theme() {
      const style = getComputedStyle(scene);
      return {
        edge: style.getPropertyValue('--edge').trim() || '#8496a5',
        pane: style.getPropertyValue('--strip').trim() || '#ffffff',
        text: style.getPropertyValue('--text').trim() || '#2c3e50',
        font: style.getPropertyValue('--font').trim() || 'sans-serif',
        halo: style.getPropertyValue('--bg').trim() || '#ffffff',
      };
    }

    function shape(id, point) {
      const head = headsByID.get(id);
      if (head) return {x:point.x, y:point.y, rx:head.width / 2, ry:head.height / 2};
      if (lips.has(id)) {
        const width = paneWidth(headsByLane.get(model.N[id].lane), point);
        return {x:point.x, y:point.y, rx:width / 2, ry:width >= 63 ? 8.5 : 4.5};
      }
      const radius = 1 + 3.5 * point.scale;
      return {x:point.x, y:point.y, rx:radius, ry:radius};
    }

    function port(bounds, other) {
      const dx = other.x - bounds.x, dy = other.y - bounds.y;
      const step = Math.max(Math.abs(dx) / (bounds.rx + 2), Math.abs(dy) / (bounds.ry + 2));
      return step >= 1
        ? {x:bounds.x + dx / step, y:bounds.y + dy / step}
        : {x:bounds.x, y:bounds.y};
    }

    function drawInspection(active, colors, x, y, width, height, dpr) {
      const ink = inspectionInk;
      if (!ink) return;
      ink.setTransform(dpr, 0, 0, dpr, -x * dpr, -y * dpr);
      ink.clearRect(x, y, width, height);
      const edges = incident.get(active) || [];
      inspection.dataset.incidentCount = edges.length;
      const focused = points.get(active);
      if (!focused) return;

      const mark = (at, color, radius) => {
        ink.lineWidth = 1.6;
        ink.fillStyle = colors.halo;
        ink.strokeStyle = color;
        ink.beginPath();
        ink.arc(at.x, at.y, radius, 0, Math.PI * 2);
        ink.fill();
        ink.stroke();
      };
      const endpoints = [];
      ink.lineJoin = ink.lineCap = 'round';
      for (const {from, to} of edges) {
        const a = points.get(from), b = points.get(to);
        if (!a || !b) continue;
        const origin = port(shape(from, a), b), target = port(shape(to, b), a);
        const node = model.N[to].kind === 'shared' ? model.N[from] : model.N[to];
        const color = `color-mix(in srgb, ${colors.text} 63%, ${glyphInfo(node).colour})`;
        const dx = target.x - origin.x, dy = target.y - origin.y;
        const horizontal = Math.abs(dx) >= Math.abs(dy) * .45;
        const cx = horizontal ? dx * .33 : 0, cy = horizontal ? 0 : dy * .37;
        if (Math.hypot(dx, dy) > 5) {
          ink.beginPath();
          ink.moveTo(origin.x, origin.y);
          ink.bezierCurveTo(origin.x + cx, origin.y + cy, target.x - cx, target.y - cy, target.x, target.y);
          ink.setLineDash([]);
          ink.strokeStyle = colors.halo;
          ink.lineWidth = 5;
          ink.stroke();
          ink.setLineDash(model.N[to].draft ? [5, 4] : []);
          ink.strokeStyle = color;
          ink.lineWidth = 2;
          ink.stroke();
          ink.setLineDash([]);
          const mid = {x:(origin.x + target.x) / 2, y:(origin.y + target.y) / 2};
          const dir = Math.atan2(dy, dx);
          ink.beginPath();
          ink.moveTo(mid.x - 5 * Math.cos(dir - .55), mid.y - 5 * Math.sin(dir - .55));
          ink.lineTo(mid.x, mid.y);
          ink.lineTo(mid.x - 5 * Math.cos(dir + .55), mid.y - 5 * Math.sin(dir + .55));
          ink.stroke();
        }
        endpoints.push([origin, color, from === active ? 3.3 : 2.6], [target, color, to === active ? 3.3 : 2.6]);
      }
      for (const args of endpoints) mark(...args);
      if (!cards.has(active) && !rearLinks.has(active)) mark(focused, colors.text, lips.has(active) ? 4 : 5);
    }

    function paint() {
      frame = 0;
      if (!model || !layout || !angle || !ctx) return;
      const width = scene.clientWidth, height = scene.clientHeight;
      if (!width || !height) return;
      const dpr = Math.min(2, devicePixelRatio || 1);
      if (width !== paintWidth || height !== paintHeight || dpr !== pixelRatio) {
        canvas.width = Math.ceil(width * dpr);
        canvas.height = Math.ceil(height * dpr);
        inspection.width = canvas.width;
        inspection.height = canvas.height;
        paintWidth = width;
        paintHeight = height;
        pixelRatio = dpr;
      }
      const x = scroll.scrollLeft, y = scroll.scrollTop;
      ctx.setTransform(dpr, 0, 0, dpr, -x * dpr, -y * dpr);
      ctx.clearRect(x, y, width, height);
      const colors = theme();
      const active = inspected();
      const focus = model.N[active];
      const related = neighborhood.get(active);
      const visible = (a, margin = 0) => a.x >= x - margin && a.y >= y - margin && a.x <= x + width + margin && a.y <= y + height + margin;
      const pointOpacity = .78 + angle * .22;
      const paneOpacity = clamp((angle - .46) * 2.1, 0, 1);
      const contains = id => id === active || related?.has(id);
      ctx.lineCap = 'round';
      pickable.length = 0;
      const groundOpacity = clamp((angle - .58) / .42, 0, 1);
      if (groundOpacity && layout.ground) {
        drawGround(x, x + width, groundOpacity, colors.edge);
        const picked = focus || model.N[selected];
        for (const head of layout.heads) {
          const point = points.get(head.id);
          if (visible(point, head.width)) drawFootprints(head, point, histories.get(head.lane), picked?.lane, picked?.id, groundOpacity, colors.edge);
        }
      }
      if (paneOpacity) {
        for (const head of layout.heads) {
          const trail = histories.get(head.lane);
          const start = points.get(head.id);
          if (!trail?.length || !start) continue;
          drawLane(head, start, trail, paneOpacity * (focus?.lane === head.lane ? .52 : .25), colors);
        }
      }
      for (const {from, to} of layout.edges) {
        const a = points.get(from), b = points.get(to);
        if (!a || !b || Math.max(a.x, b.x) < x - 48 || Math.max(a.y, b.y) < y - 48 ||
            Math.min(a.x, b.x) > x + width + 48 || Math.min(a.y, b.y) > y + height + 48) continue;
        const parent = model.N[from], child = model.N[to];
        if (from === active || to === active) continue;
        const sameLane = parent.lane === child.lane;
        const dx = b.x - a.x, dy = b.y - a.y;
        const distance = Math.hypot(dx, dy);
        if (distance < 2) continue;
        const reach = clamp(Math.abs(dx) * .46, 1, 64);
        const sign = dx < 0 ? -1 : 1;
        const density = sameLane ? 1 : crossings.get([parent.lane, child.lane].sort().join('|')) || 1;
        const alpha = (focus ? (sameLane ? .07 : .12) : sameLane ? .12 : .35) / Math.sqrt(density);
        ctx.globalAlpha = pointOpacity * alpha;
        ctx.setLineDash(child.draft ? [5, 4] : []);
        ctx.strokeStyle = sameLane ? colors.edge : glyphInfo(child.kind === 'shared' ? parent : child).colour;
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(a.x, a.y);
        ctx.bezierCurveTo(a.x + reach * sign, a.y, b.x - reach * sign, b.y, b.x, b.y);
        ctx.stroke();

      }

      ctx.setLineDash([]);
      for (const id of order) {
        const point = points.get(id);
        const node = model.N[id];
        if (lips.has(id) && node && visible(point, 150) && (rearLayer.hidden || !rearLinks.has(id))) {
          drawPane(node, point, headsByLane.get(node.lane), paneOpacity * (focus && !contains(id) ? .58 : .92), colors);
          pickable.push({id, bounds:shape(id, point), lip:true});
        }
      }
      const occupied = new Set();
      for (let index = order.length - 1; index >= 0; index--) {
        const id = order[index];
        const point = points.get(id);
        const node = model.N[id];
        if (!node || !visible(point, 24) || lips.has(id) || (angle > .85 && headsByID.has(id))) continue;
        const strong = contains(id);
        const cell = Math.floor(point.x / 7) + ',' + Math.floor(point.y / 7);
        const firstInCell = !occupied.has(cell);
        if (!firstInCell && !strong) continue;
        if (firstInCell) {
          occupied.add(cell);
          pickable.push({id, bounds:shape(id, point)});
        }
        const alpha = strong ? .94 : focus ? .14 : .18 + .22 * point.scale;
        drawGlyph(node, point, pointOpacity * alpha);
      }
      ctx.globalAlpha = 1;
      drawInspection(active, colors, x, y, width, height, dpr);
    }

    function repaint() {
      if (!frame) frame = requestAnimationFrame(paint);
    }

    function updateSelection() {
      const active = inspected();
      if (paintedSelection === active + '|' + selected) return;
      paintedSelection = active + '|' + selected;
      if (active) scene.dataset.inspectState = active;
      else delete scene.dataset.inspectState;
      const related = neighborhood.get(active);
      const selectedLane = model?.N[selected]?.lane;
      const decorate = (id, el) => {
        const isSelected = id === selected, isInspected = id === active;
        const isRelated = related?.has(id) && !isInspected;
        const showsSelectedLane = Boolean(selectedLane && el.dataset.headLane === selectedLane);
        el.classList.toggle('is-selected', isSelected || showsSelectedLane);
        el.classList.toggle('is-inspected', isInspected);
        el.classList.toggle('is-related', Boolean(isRelated));
        el.classList.toggle('is-muted', Boolean(active && !isRelated && !isInspected));
        if (isSelected) el.setAttribute('aria-current', 'true');
        else el.removeAttribute('aria-current');
      };
      for (const [id, {el}] of cards) decorate(id, el);
      for (const [id, el] of rearLinks) decorate(id, el);
    }

    // Keep the departing faces in flight as their real predecessors move forward.
    function startTravel(next, direction) {
      const incoming = new Set(next.heads.map(head => head.id));
      const layer = document.createElement('div');
      layer.className = 'hc-passing';
      layer.inert = true;
      layer.setAttribute('aria-hidden', 'true');
      space.append(layer);
      for (const [id, entry] of cards) {
        if (incoming.has(id)) continue;
        const el = entry.el;
        const old = points.get(id), rear = next.points.get(id);
        if (!old) continue;
        const visible = old.x > scroll.scrollLeft - 180 && old.x < scroll.scrollLeft + scroll.clientWidth + 180;
        if (!visible) continue;
        cards.delete(id);
        el.removeAttribute('data-head-state');
        el.removeAttribute('data-head-lane');
        el.removeAttribute('href');
        el.tabIndex = -1;
        layer.append(el);
        const x = rear ? rear.x - old.x : 0;
        const y = rear ? rear.y - old.y : direction * 96;
        el.animate([{translate:'0 0', scale:1, opacity:1}, {translate:x+'px '+y+'px', scale:rear ? .65 : 1.18, opacity:0}],
          {duration:320, easing:'ease-out', fill:'forwards'}).finished.then(() => el.remove(), () => el.remove());
      }
      Promise.allSettled(layer.getAnimations({subtree:true}).map(a => a.finished)).then(() => layer.remove());
      return {points, heads:headsByID, direction};
    }

    function render(T, next, side, value, focused, travel) {
      angle = clamp(value, 0, 1);
      scene.hidden = !angle;
      if (!angle || !T || !next) {
        inspectable = false;
        hovered = '';
        focusedByUser = '';
        tip.hidden = true;
        updateSelection();
        return;
      }
      setLayout(T, next);
      selected = focused || '';
      sidePoints = side;
      const ready = angle >= .98;
      const reveal = clamp((angle - .34) / .66, 0, 1);
      scene.style.setProperty('--hc-backdrop', String(clamp((angle - .12) / .88, 0, 1)));
      scene.style.setProperty('--hc-reveal', String(reveal));
      scene.classList.toggle('is-ready', ready);
      scroll.inert = !ready;
      inspectable = angle === 1 && !travel;
      rearLayer.hidden = !inspectable;
      if (!inspectable) {
        hovered = '';
        tip.hidden = true;
      }
      if (initialPan && scroll.clientWidth) {
        initialPan = false;
        const requested = headsByID.get(selected) || headsByLane.get(model.N[selected]?.lane);
        const center = layout.width / 2;
        const centered = layout.heads.slice(0, 2).every(head => Math.abs(center - head.x) + head.width / 2 <= scroll.clientWidth / 2);
        scroll.scrollLeft = Math.max(0, (centered ? center : requested?.x ?? center) - scroll.clientWidth / 2);
      }
      points = new Map();
      faceScales.clear();
      floorPhase = travel ? travel.direction * travel.progress : 0;
      const yaw = (1 - angle) * 90;
      const scrollX = scroll.scrollLeft, scrollY = scroll.scrollTop;
      const axisX = scrollX + scroll.clientWidth / 2;
      for (const [id, endpoint] of layout.points) {
        const sidePoint = sidePoints?.get(id);
        const start = sidePoint ? {...sidePoint, x: sidePoint.x + scrollX, y: sidePoint.y + scrollY} : endpoint;
        let point = WorldHeadCamera.turn(start, endpoint, angle, axisX);
        const head = headsByID.get(id);
        let faceScale = 1;
        if (travel) {
          const past = travel.points.get(id);
          const origin = past || {x:endpoint.x, y:endpoint.y + 96, scale:endpoint.scale * 1.18};
          const t = travel.progress;
          point = {x:origin.x + (endpoint.x-origin.x)*t, y:origin.y + (endpoint.y-origin.y)*t,
            scale:origin.scale + (endpoint.scale-origin.scale)*t, depth:endpoint.depth};
          if (head) {
            const oldHead = travel.heads.get(id);
            const ratio = oldHead ? oldHead.width/head.width : past ? .58 : 1.18;
            faceScale = ratio + (1-ratio)*t;
          }
        }
        points.set(id, point);
        const entry = cards.get(id);
        if (entry) {
          faceScales.set(id, faceScale);
          entry.el.style.transform = `translate3d(${point.x}px,${point.y}px,0) translate(-50%,-50%) perspective(800px) rotateY(${yaw}deg) scale(${faceScale})`;
        }
      }
      updateSelection();
      repaint();
    }

    function face(event) {
      return event.target.closest('.hc-face, .hc-rear');
    }

    function hasVisibleFocus(el) {
      return Boolean(el && el === document.activeElement && el.getClientRects().length && !el.closest('[hidden], [inert]'));
    }

    function focusHead(id) {
      const lane = model?.N[id]?.lane;
      const el = cards.get(id)?.el || rearLinks.get(id) || cards.get(headsByLane.get(lane)?.id)?.el;
      if (!el || scene.hidden || scroll.inert) return;
      const viewport = scroll.getBoundingClientRect();
      const target = el.getBoundingClientRect();
      const gap = 9;
      const dx = target.left < viewport.left + gap ? target.left - viewport.left - gap
        : Math.max(0, target.right - viewport.right + gap);
      const dy = target.top < viewport.top + gap ? target.top - viewport.top - gap
        : Math.max(0, target.bottom - viewport.bottom + gap);
      scroll.scrollBy(dx, dy);
      silentFocus = true;
      el.focus({preventScroll:true});
      silentFocus = false;
      if (!hasVisibleFocus(el)) return;
      focusedByUser = el.dataset.headState || el.dataset.rearState;
      hovered = '';
      tip.hidden = true;
      updateSelection();
      repaint();
    }

    function updatePointer(event) {
      if (!inspectable || event.pointerType === 'touch') return;
      const {clientX, clientY} = event;
      pointer = {clientX, clientY, pointerType:event.pointerType};
      const target = document.elementFromPoint(clientX, clientY);
      const link = target?.closest('.hc-face[data-head-state], .hc-rear');
      let id = link?.dataset.headState || link?.dataset.rearState || '';
      let canvasOnly = false;
      const box = scene.getBoundingClientRect();
      if (!id && target?.closest('.hc-space') === space && !target.closest('a, button, .hc-passing')) {
        const x = clientX - box.left + scroll.scrollLeft, y = clientY - box.top + scroll.scrollTop;
        let best = Infinity;
        for (const {id: candidate, bounds, lip} of pickable) {
          const dx = Math.abs(bounds.x - x), dy = Math.abs(bounds.y - y);
          const distance = Math.hypot(dx, dy);
          if (lip ? dx <= bounds.rx && dy <= bounds.ry : distance <= bounds.rx + 3) {
            // Painted small glyphs sit above lips. Within each kind select
            // the nearest mark rather than a broad neighbor's hover radius.
            const score = distance + (lip ? 100 : 0);
            if (score >= best) continue;
            best = score;
            id = candidate;
          }
        }
        canvasOnly = Boolean(id);
      }
      if (id !== hovered) {
        hovered = id;
        updateSelection();
        repaint();
      }
      tip.hidden = !canvasOnly;
      if (canvasOnly) {
        const title = options.title(model.N[id]) || model.N[id].title || model.N[id].label || id;
        if (tip.textContent !== title) tip.textContent = title;
        tip.style.left = clamp(clientX - box.left + 11, 6, box.width - tip.offsetWidth - 6) + 'px';
        tip.style.top = clamp(clientY - box.top + 12, 6, box.height - tip.offsetHeight - 6) + 'px';
      }
      return canvasOnly ? id : null;
    }

    scene.addEventListener('focusin', event => {
      const el = face(event);
      if (!silentFocus && hasVisibleFocus(el)) {
        hovered = '';
        tip.hidden = true;
        focusedByUser = el.dataset.headState || el.dataset.rearState;
        options.select(focusedByUser);
        updateSelection();
        repaint();
      }
    });
    scene.addEventListener('click', event => {
      if (event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
      const data = face(event)?.dataset;
      const id = data?.headState || data?.rearState ||
        (event.detail && event.target.closest('.hc-space') && updatePointer(event));
      if (!id) return;
      event.preventDefault();
      event.stopPropagation();
      options.open(id);
    });
    scene.addEventListener('focusout', event => {
      if (!face(event) || event.relatedTarget?.closest?.('.hc-face, .hc-rear')) return;
      focusedByUser = '';
      updateSelection();
      repaint();
    });
    scene.addEventListener('pointermove', updatePointer);
    scene.addEventListener('pointerover', updatePointer);
    scene.addEventListener('pointerleave', () => {
      pointer = null;
      hovered = '';
      tip.hidden = true;
      updateSelection();
      repaint();
    });
    scroll.addEventListener('scroll', () => {
      if (pointer) updatePointer(pointer);
      repaint();
      options.pan?.();
    }, {passive:true});
    const resize = new ResizeObserver(repaint);
    resize.observe(scene);
    const appearance = new MutationObserver(repaint);
    appearance.observe(document.documentElement, {attributes:true, attributeFilter:['class', 'data-theme']});

    function aim(preferred) {
      const left = scroll.scrollLeft, right = left + scroll.clientWidth;
      const current = headsByLane.get(model?.N[preferred]?.lane);
      if (current && current.x >= left && current.x <= right) return preferred;
      return layout?.heads.reduce((best, h) => !best || Math.abs(h.x - (left+right)/2) < Math.abs(best.x - (left+right)/2) ? h : best, null)?.id;
    }
    function reveal(id, smooth) {
      const head = headsByLane.get(model?.N[id]?.lane);
      if (!head) return;
      const left = scroll.scrollLeft, right = left + scroll.clientWidth;
      if (head.x - head.width / 2 >= left && head.x + head.width / 2 <= right) return;
      scroll.scrollTo({left:head.x - scroll.clientWidth / 2, behavior:smooth ? 'smooth' : 'instant'});
    }
    return {render, element:scene, startTravel, aim, reveal, focus: focusHead, destroy() {
      if (frame) cancelAnimationFrame(frame);
      resize.disconnect();
      appearance.disconnect();
      scene.remove();
    }};
  }

  window.WorldHeadScene = {create};
})();
