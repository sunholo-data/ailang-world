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
    let floorTravel = 0;
    const faceSizes = new Map();
    let points = new Map(), headsByID = new Map(), headsByLane = new Map(), order = [], histories = new Map();

    let passingLayer = null;
    const svgNS = 'http://www.w3.org/2000/svg';

    function attachFront(el) {
      const content = document.createElement('span');
      content.className = 'hc-front-content';
      const sides = ['roof', 'cheek'].map(side => {
        const plane = document.createElementNS(svgNS, 'svg');
        plane.setAttribute('class', 'hc-plane hc-' + side);
        plane.setAttribute('aria-hidden', 'true');
        const polygon = document.createElementNS(svgNS, 'polygon');
        plane.append(polygon);
        el.append(plane);
        return {plane, polygon};
      });
      el.append(content);
      return {el, content, sides};
    }

    function blockAt(point, width, height, pane) {
      const ground = layout.ground;
      const reveal = clamp((angle - .94) / .06, 0, 1);
      return reveal && ground && WorldHeadCamera.block
        ? WorldHeadCamera.block(ground, point, width, height, ground.blockDepth * reveal * (pane ? .37 : 1))
        : null;
    }

    function placePlane({plane, polygon}, vertices, left, top, scaleX, scaleY) {
      if (!vertices) {
        plane.style.display = 'none';
        return;
      }
      const xs = vertices.map(point => point.x), ys = vertices.map(point => point.y);
      const minX = Math.min(...xs), minY = Math.min(...ys);
      const width = (Math.max(...xs) - minX) / scaleX, height = (Math.max(...ys) - minY) / scaleY;
      if (width <= .08 || height <= .08) {
        plane.style.display = 'none';
        return;
      }
      const points = vertices.map(point => [(point.x - minX) / scaleX, (point.y - minY) / scaleY]);
      plane.style.display = '';
      plane.style.left = (minX - left) / scaleX + 'px';
      plane.style.top = (minY - top) / scaleY + 'px';
      plane.style.width = width + 'px';
      plane.style.height = height + 'px';
      polygon.setAttribute('points', points.map(([x,y]) => x + ',' + y).join(' '));
    }

    function projectFace(entry, point, width, height, {displayed = {width,height}, yaw = 0, pane = false} = {}) {
      const {el, sides} = entry;
      const scaleX = displayed.width / width, scaleY = displayed.height / height;
      const left = point.x - displayed.width / 2, top = point.y - displayed.height / 2;
      const block = blockAt(point, displayed.width, displayed.height, pane);
      entry.block = block;
      el.style.width = width + 'px';
      el.style.height = height + 'px';
      el.style.transform = `translate3d(${point.x}px,${point.y}px,0) translate(-50%,-50%) perspective(800px) rotateY(${yaw}deg) scale(${scaleX}, ${scaleY})`;
      // A pane can be elevated above the street. Its true camera depth,
      // independent of the projected y or focus, determines what it occludes.
      const depth = point.depth > 0 ? point.depth : block?.depth;
      el.style.zIndex = String(Math.round(depth ? 1e6 / depth : (point.y + displayed.height / 2) * 100));
      el.classList.toggle('hc-lit-cheek', Boolean(block?.side && point.x > layout.ground.vx));
      placePlane(sides[0], block?.top, left, top, scaleX, scaleY);
      placePlane(sides[1], block?.side, left, top, scaleX, scaleY);
    }

    function cameraTravel(a, b, from, to, t) {
      const oldDepth = a.depth > 0 ? a.depth : 1 / (a.scale || 1);
      const newDepth = b.depth > 0 ? b.depth : 1 / (b.scale || 1);
      const aWeight = (1 - t) * oldDepth, bWeight = t * newDepth;
      const depth = aWeight + bWeight;
      const project = (x, y) => (x * aWeight + y * bWeight) / depth;
      return {
        point: {x:project(a.x,b.x), y:project(a.y,b.y), depth, scale:1/depth},
        size: {width:project(from.width,to.width), height:project(from.height,to.height)}
      };
    }

    // Continue an entering or departing front along its viewing ray. Clip
    // it short of the eye while it fades past the camera.
    function movedCamera(point, deltaZ) {
      const {vx, horizon} = layout.ground;
      const depth = Math.max(point.depth * .35, point.depth + deltaZ);
      const ratio = point.depth / depth;
      return {x:vx + (point.x-vx)*ratio, y:horizon + (point.y-horizon)*ratio, depth, scale:1/depth};
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
          entry = attachFront(el);
          cards.set(head.id, entry);
        }
        const {el} = entry;
        if (entry.markup !== markup) {
          entry.content.innerHTML = markup;
          entry.markup = markup;
        }
        el.classList.toggle('needs-human', Boolean(head.attention));
        el.classList.toggle('is-draft', Boolean(node.draft));
        el.classList.toggle('is-compact', compact);
        el.classList.toggle('is-distant', compact && head.height < 35);
        el.classList.toggle('is-middle', !compact && head.height <= 76);
        el.classList.toggle('is-short', !compact && head.height < 58);
        el.href = statePath(node);
        el.title = label + ' · ' + title;
        if (compact) el.setAttribute('aria-label', el.title);
        else el.removeAttribute('aria-label');
        el.style.setProperty('--hc-lane', accent(lane));
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

    function paneSize(head, point) {
      const scale = head.depth / Math.max(.01, point.depth);
      return {width:head.width * .88 * scale, height:head.height * .36 * scale};
    }

    function syncRearLinks() {
      const titled = new Set();
      for (const id of order) {
        if (!lips.has(id)) continue;
        const node = model.N[id];
        const point = layout.points.get(id);
        const {width,height} = paneSize(headsByLane.get(node.lane), point);
        if (width < 63 || height < 13) continue;
        const x = point.x - width / 2, y = point.y - height / 2;
        const covered = layout.heads.some(head =>
          layout.points.get(head.id).depth < point.depth && x >= head.x - head.width / 2 && x + width <= head.x + head.width / 2 &&
          y >= head.y - head.height / 2 && y + height <= head.y + head.height / 2);
        if (covered) continue;
        titled.add(id);
        let entry = rearLinks.get(id);
        if (!entry) {
          const el = document.createElement('a');
          el.className = 'hc-rear';
          el.dataset.rearState = id;
          entry = attachFront(el);
          rearLinks.set(id, entry);
        }
        const {el, content} = entry;
        const title = options.title(node) || node.title || node.label || id;
        const label = model.laneTitle?.(node.lane) || node.lane;
        content.textContent = title;
        el.style.setProperty('--hc-lane', accent(model.LANES[node.lane]));
        el.title = title;
        el.setAttribute('aria-label', label + ' · ' + title);
        el.href = statePath(node);
        const index = titled.size - 1;
        if (rearLayer.children[index] !== el) rearLayer.insertBefore(el, rearLayer.children[index] || null);
      }
      for (const [id, {el}] of rearLinks) {
        if (titled.has(id)) continue;
        el.remove();
        rearLinks.delete(id);
      }
    }

    const colorsForPane = (base, edge, tint) => `color-mix(in srgb, ${edge} ${tint}%, ${base})`;

    function drawPane(node, point, head, alpha, theme) {
      const {width,height} = paneSize(head, point);
      const x = point.x - width / 2, y = point.y - height / 2;
      ctx.globalAlpha = alpha;
      ctx.fillStyle = theme.pane;
      ctx.strokeStyle = theme.edge;
      ctx.lineWidth = 1;
      const block = blockAt(point, width, height, true);
      const path = (vertices, fill) => {
        ctx.beginPath();
        vertices.forEach(({x, y}, index) => index ? ctx.lineTo(x,y) : ctx.moveTo(x,y));
        ctx.closePath();
        ctx.fillStyle = fill;
        ctx.fill();
        ctx.stroke();
      };
      const front = block?.front || [{x,y}, {x:x+width,y}, {x:x+width,y:y+height}, {x,y:y+height}];
      if (block?.side) path(block.side, colorsForPane(theme.pane, theme.edge, 36));
      if (block?.top) path(block.top, colorsForPane(theme.pane, theme.edge, 16));
      path(front, theme.pane);
      if (width < 63 || height < 13 || angle < .86) return block;
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
      return block;
    }

    function drawGround(left, right, bottom, opacity, ink) {
      const ground = layout.ground;
      const {vx, horizon, street, streetY, pitch, focal} = ground;
      const eye = streetY - horizon;
      if (bottom <= horizon || eye <= 0) return;
      const project = (x, z) => WorldHeadCamera.project(ground, x, 0, z);
      const near = eye / (bottom - horizon);
      // A paving slab is square in the world. X is measured in pixels at
      // Z=1, so the same length along the optical axis is pitch / focal.
      const slab = pitch / focal;
      const far = Math.max(near + slab, Math.sqrt(eye * slab / 3));
      const leftCurb = project(-street / 2, near), rightCurb = project(street / 2, near);
      const fade = ctx.createLinearGradient(0, horizon, 0, streetY);
      fade.addColorStop(0, 'transparent');
      fade.addColorStop(.48, ink);
      fade.addColorStop(1, ink);
      ctx.fillStyle = ctx.strokeStyle = fade;
      ctx.beginPath();
      ctx.moveTo(vx, horizon);
      ctx.lineTo(leftCurb.x, leftCurb.y);
      ctx.lineTo(rightCurb.x, rightCurb.y);
      ctx.closePath();
      ctx.globalAlpha = opacity * .11;
      ctx.fill();

      ctx.beginPath();
      ctx.moveTo(vx, horizon);
      ctx.lineTo(leftCurb.x, leftCurb.y);
      ctx.moveTo(vx, horizon);
      ctx.lineTo(rightCurb.x, rightCurb.y);
      ctx.globalAlpha = opacity * .63;
      ctx.lineWidth = 1.25;
      ctx.stroke();

      ctx.lineWidth = 1;
      ctx.globalAlpha = opacity * .28;
      ctx.beginPath();
      for (const side of [-1, 1]) {
        for (let i = 1; i <= 3; i++) {
          const x = side * (street / 2 + i * pitch);
          const a = project(x, near);
          ctx.moveTo(a.x, a.y);
          ctx.lineTo(vx, horizon);
        }
      }
      ctx.stroke();

      const first = Math.ceil((near + floorTravel) / slab);
      const last = Math.floor((far + floorTravel) / slab);
      ctx.strokeStyle = ink;
      for (let i = first; i <= last; i++) {
        const z = i * slab - floorTravel;
        const y = project(0, z).y;
        const spacing = y - project(0, z + slab).y;
        ctx.globalAlpha = opacity * .38 * clamp((spacing - 2) / 13, 0, 1);
        ctx.beginPath();
        ctx.moveTo(left, y);
        ctx.lineTo(right, y);
        ctx.stroke();
      }
    }

    function drawFootprint(head, point, trail, picked, opacity, colors) {
      const ground = layout.ground;
      const {width, height} = faceSizes.get(head.id) || head;
      const solid = cards.get(head.id)?.block;
      if (!solid) return;
      const {depth:z, backDepth:back, floor} = solid;
      const project = (x, depth) => WorldHeadCamera.project(ground, x, 0, depth);
      const worldX = (point.x - ground.vx) * z;
      const left = worldX - width * z / 2, right = worldX + width * z / 2;
      // Parallel light above the left foreground casts every shadow to the
      // right and back on the ground, including the narrow end of each block.
      const roofElevation = ground.streetY - ground.horizon - (point.y - height/2 - ground.horizon) * z;
      const dx = roofElevation * .35, dz = roofElevation * .5 / ground.focal;
      const outline = [floor[0], floor[1], project(right + dx, z + dz),
        project(right + dx, back + dz), project(left + dx, back + dz), floor[3]];
      ctx.fillStyle = colors.shadow;
      ctx.globalAlpha = .12 * opacity;
      ctx.beginPath();
      ctx.moveTo(outline[0].x, outline[0].y);
      for (let i = 1; i < outline.length; i++) ctx.lineTo(outline[i].x, outline[i].y);
      ctx.closePath();
      ctx.fill();

      if (picked?.lane !== head.lane) return;
      let preceding = project(worldX, z);
      ctx.strokeStyle = colors.edge;
      ctx.lineWidth = 1;
      for (const id of trail) {
        const p = points.get(id);
        if (!(p.depth > 0)) continue;
        // A mark sits on the ground immediately below the actual state.
        const foot = project((p.x - ground.vx) * p.depth, p.depth);
        if (id !== picked.id && preceding.y - foot.y < 12) continue;
        ctx.globalAlpha = opacity * .28;
        ctx.beginPath();
        ctx.moveTo(preceding.x, preceding.y);
        ctx.lineTo(foot.x, foot.y);
        ctx.stroke();
        preceding = foot;
        const half = Math.min(9, 7 / p.depth);
        ctx.globalAlpha = opacity * (id === picked.id ? .9 : .57);
        ctx.beginPath();
        ctx.moveTo(foot.x - half, foot.y);
        ctx.lineTo(foot.x + half, foot.y);
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
        shadow: style.getPropertyValue('--ink').trim() || '#243441',
      };
    }

    function shape(id, point) {
      const head = headsByID.get(id);
      if (head) return {x:point.x, y:point.y, rx:head.width / 2, ry:head.height / 2};
      if (lips.has(id)) {
        const {width,height} = paneSize(headsByLane.get(model.N[id].lane), point);
        return {x:point.x, y:point.y, rx:width / 2, ry:height / 2};
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
      const groundOpacity = clamp((angle - .78) / .22, 0, 1);
      if (groundOpacity && layout.ground) {
        drawGround(x, x + width, y + height, groundOpacity, colors.edge);
        const picked = focus || model.N[selected];
        for (const head of layout.heads) {
          const point = points.get(head.id);
          if (visible(point, head.width)) drawFootprint(head, point, histories.get(head.lane), picked, groundOpacity, colors);
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
          const block = drawPane(node, point, headsByLane.get(node.lane), paneOpacity * (focus && !contains(id) ? .58 : .92), colors);
          pickable.push({id, bounds:shape(id, point), lip:true, block});
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
      for (const [id, {el}] of rearLinks) decorate(id, el);
    }

    // Departures use the same projection as arrivals while their anchor
    // moves away and their face shrinks into the event history.
    function startTravel(next) {
      const incoming = new Set(next.heads.map(head => head.id));
      const nextByLane = new Map(next.heads.map(head => [head.lane, head]));
      passingLayer?.remove();
      const layer = document.createElement('div');
      layer.className = 'hc-passing';
      layer.inert = true;
      layer.setAttribute('aria-hidden', 'true');
      space.append(layer);
      passingLayer = layer;
      const departures = [], sizes = new Map(faceSizes);
      for (const id of lips) {
        const pane = points.get(id), node = model.N[id];
        const head = node && headsByLane.get(node.lane);
        if (!pane || !head || sizes.has(id)) continue;
        sizes.set(id, paneSize(head, pane));
      }
      for (const [id, entry] of cards) {
        if (incoming.has(id)) continue;
        const {el} = entry;
        const old = points.get(id), rear = next.points.get(id);
        if (!old || old.x < scroll.scrollLeft - 180 || old.x > scroll.scrollLeft + scroll.clientWidth + 180) continue;
        const head = headsByID.get(id);
        const from = sizes.get(id) || head;
        const nextHead = nextByLane.get(head.lane);
        const end = rear || movedCamera(old, floorTravel - (next.travelZ || 0));
        const ratio = old.depth / end.depth;
        const to = rear && nextHead ? paneSize(nextHead, rear)
          : {width:from.width * ratio, height:from.height * ratio};
        cards.delete(id);
        el.removeAttribute('data-head-state');
        el.removeAttribute('data-head-lane');
        el.removeAttribute('href');
        el.tabIndex = -1;
        layer.append(el);
        departures.push({entry, head, old, end, from, to});
      }
      return {points, sizes, departures, floorFrom:floorTravel};
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
      rearLayer.hidden = angle < .94;
      rearLayer.inert = !inspectable;
      if (!travel) {
        passingLayer?.remove();
        passingLayer = null;
      }
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
      faceSizes.clear();
      const targetTravel = next.travelZ || 0;
      floorTravel = travel ? travel.floorFrom + (targetTravel - travel.floorFrom) * travel.progress : targetTravel;
      const yaw = 90 * Math.max(0, 1 - angle / .94);
      const scrollX = scroll.scrollLeft, scrollY = scroll.scrollTop;
      const axisX = scrollX + scroll.clientWidth / 2;
      for (const [id, endpoint] of layout.points) {
        const sidePoint = sidePoints?.get(id);
        const start = sidePoint ? {...sidePoint, x: sidePoint.x + scrollX, y: sidePoint.y + scrollY} : endpoint;
        let point = WorldHeadCamera.turn(start, endpoint, angle, axisX);
        const head = headsByID.get(id);
        let displayed = head && {width:head.width, height:head.height};
        if (travel) {
          const past = travel.points.get(id);
          const origin = past || movedCamera(endpoint, targetTravel - travel.floorFrom);
          const ratio = endpoint.depth / origin.depth;
          const glyphWidth = 2 + 7 * clamp(origin.scale, 0, 1);
          const previous = travel.sizes.get(id) || (head
            ? (past ? {width:glyphWidth, height:glyphWidth} : {width:head.width * ratio, height:head.height * ratio})
            : {width:0, height:0});
          const pose = cameraTravel(origin, endpoint, previous, displayed || previous, travel.progress);
          point = pose.point;
          if (head) displayed = pose.size;
        }
        points.set(id, point);
        const entry = cards.get(id);
        if (entry) {
          faceSizes.set(id, displayed);
          projectFace(entry, point, head.width, head.height, {displayed, yaw});
        }
      }
      if (!rearLayer.hidden) {
        for (const [id, entry] of rearLinks) {
          const point = points.get(id), head = headsByLane.get(model.N[id].lane);
          const target = paneSize(head, layout.points.get(id));
          const prior = travel?.points.get(id);
          const displayed = prior
            ? cameraTravel(prior, layout.points.get(id), travel.sizes.get(id) || target, target, travel.progress).size
            : target;
          projectFace(entry, point, target.width, target.height, {displayed, pane:true, yaw});
        }
      }
      if (travel) {
        const t = travel.progress;
        for (const {entry, head, old, end, from, to} of travel.departures) {
          const pose = cameraTravel(old, end, from, to, t);
          entry.el.style.opacity = String(1 - t);
          projectFace(entry, pose.point, head.width, head.height, {displayed:pose.size, yaw});
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
      const el = cards.get(id)?.el || rearLinks.get(id)?.el || cards.get(headsByLane.get(lane)?.id)?.el;
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

    function containsPoint(vertices, x, y) {
      if (!vertices) return false;
      let inside = false;
      for (let i = 0, j = vertices.length - 1; i < vertices.length; j = i++) {
        const a = vertices[i], b = vertices[j];
        if ((a.y > y) !== (b.y > y) && x < (b.x - a.x) * (y - a.y) / (b.y - a.y) + a.x) inside = !inside;
      }
      return inside;
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
        for (const {id: candidate, bounds, lip, block} of pickable) {
          const dx = Math.abs(bounds.x - x), dy = Math.abs(bounds.y - y);
          const distance = Math.hypot(dx, dy);
          const insidePane = (dx <= bounds.rx && dy <= bounds.ry) ||
            containsPoint(block?.top, x, y) || containsPoint(block?.side, x, y);
          if (lip ? insidePane : distance <= bounds.rx + 3) {
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
