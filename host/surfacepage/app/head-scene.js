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
  const exactTime = new Intl.DateTimeFormat(undefined, {year:'numeric', month:'short', day:'numeric', hour:'2-digit', minute:'2-digit', second:'2-digit', timeZone:'UTC'});

  const tickDate = new Intl.DateTimeFormat('en-GB', {day:'numeric', month:'short', timeZone:'UTC'});
  const tickClock = new Intl.DateTimeFormat('en-GB', {hour:'2-digit', minute:'2-digit', timeZone:'UTC'});

  function create(container, options) {
    const scene = document.createElement('section');
    scene.className = 'hc-scene';
    scene.setAttribute('aria-label', 'Branch heads');
    scene.hidden = true;
    scene.innerHTML = '<canvas class="hc-graph" aria-hidden="true"></canvas><div class="hc-scroll"><div class="hc-space"><div class="hc-rear-layer" hidden></div></div></div><div class="hc-time-axis" role="group" aria-label="Time, UTC" hidden></div><canvas class="hc-inspection" aria-hidden="true"></canvas><div class="hc-inspect-tip" aria-hidden="true" hidden></div>';
    container.append(scene);

    const scroll = scene.querySelector('.hc-scroll'), space = scene.querySelector('.hc-space'), canvas = scene.querySelector('canvas');
    const ctx = canvas.getContext('2d');
    const timeAxis = scene.querySelector('.hc-time-axis');
    const inspection = scene.querySelector('.hc-inspection'), tip = scene.querySelector('.hc-inspect-tip');
    const preview = WorldHeadNodePreview.create(scene, options);
    const nodeLink = '.hc-face[data-head-state], .hc-rear, [data-lineage-state]';
    const stateID = el => el?.dataset.headState || el?.dataset.rearState || el?.dataset.lineageState || '';
    const inspectionInk = inspection.getContext('2d');
    let inspectable = false, pointer = null;
    const inspected = () => {
      const id = hovered || focusedByUser;
      return inspectable && model?.N[id] ? id : '';
    };
    const rearLayer = space.querySelector('.hc-rear-layer');
    const lineage = WorldHeadLineage.create(space, {title:options.title, tag:options.tag, path:statePath, accent});
    const cards = new Map(), rearLinks = new Map();
    const stamps = new Map(), paneTitles = new Map(), crossings = new Map();
    const neighborhood = new Map(), incident = new Map();
    const scratch = document.createElement('template');
    let model, layout, sidePoints, angle = 0, selected = '', hovered = '', frame = 0, silentFocus = false;
    let panAxis, panWidth = 0, panOffset = 0, focusedByUser = '', paintedSelection, paintWidth = 0, paintHeight = 0, pixelRatio = 0;
    let floorTravel = 0, cameraX = 0;
    let timeTicks = [];
    const faceSizes = new Map();
    let points = new Map(), headsByID = new Map(), headsByLane = new Map(), order = [], histories = new Map();

    let passingLayer = null;
    const svgNS = 'http://www.w3.org/2000/svg';
    const miniatures = document.createElementNS(svgNS, 'svg');
    miniatures.setAttribute('class', 'hc-miniatures');
    rearLayer.append(miniatures);

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
        ? WorldHeadCamera.block(ground, point, width, height, reveal * (point.blockDepth ?? (pane ? ground.paneDepth : ground.blockDepth)))
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
        point: {x:project(a.x,b.x), y:project(a.y,b.y), depth, scale:1/depth,
          blockDepth:a.blockDepth + t*(b.blockDepth-a.blockDepth)},
        size: {width:project(from.width,to.width), height:project(from.height,to.height)}
      };
    }

    // Continue an entering or departing front along its viewing ray. Clip
    // it short of the eye while it fades past the camera.
    function movedCamera(point, deltaZ, {vx, horizon}, deltaX) {
      const depth = Math.max(point.depth * .35, point.depth + deltaZ);
      const ratio = point.depth / depth;
      return {x:vx + (point.x-vx)*ratio + deltaX/depth, y:horizon + (point.y-horizon)*ratio,
        depth, scale:1/depth, blockDepth:point.blockDepth};
    }

    function setTimeAxis() {
      const axis = layout.time;
      const ticks = [];
      if (axis && model.dated) {
        const {firstIndex,lastIndex} = axis;
        // Keep both ends of the scale. Between them, label the nearest
        // distinct instants with enough space for a date and a clock time.
        const addTick = index => {
          if (index < firstIndex || (model.undated && index === 0)) return;
          const ms = model.TT[index], y = axis.groundYAt(index);
          if (!Number.isFinite(ms) || !Number.isFinite(y) || y < 0 || y > layout.height ||
              ticks.some(tick => Math.abs(tick.y - y) < 44)) return;
          ticks.push({index,ms,y});
        };
        addTick(lastIndex);
        addTick(firstIndex + Number(model.undated && firstIndex === 0));
        for (let i = lastIndex - 1; i > firstIndex; i--) addTick(i);
      }
      timeAxis.replaceChildren();
      timeTicks = ticks;
      if (!ticks.length) return;
      ticks.sort((a,b) => a.y - b.y);
      const top = ticks[0].y;
      timeAxis.style.top = top + 'px';
      timeAxis.style.height = (ticks[ticks.length - 1].y - top) + 'px';
      const label = document.createElement('span');
      label.className = 'hc-time-title';
      label.textContent = 'Time · UTC';
      timeAxis.append(label);
      for (const {ms,y} of ticks) {
        const tick = document.createElement('time'), date = new Date(ms);
        tick.dateTime = date.toISOString();
        tick.title = exactTime.format(date) + ' UTC';
        tick.style.top = (y - top) + 'px';
        tick.textContent = tickDate.format(date) + '\n' + tickClock.format(date);
        timeAxis.append(tick);
      }
    }

    function setLayout(T, next) {
      if (model === T && layout === next) return;
      model = T;
      layout = next;
      setTimeAxis();
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
        if (compact) el.setAttribute('aria-label', label + ' · ' + title);
        else el.removeAttribute('aria-label');
        el.style.setProperty('--hc-lane', accent(lane));
        if (space.children[index] !== el) {
          const focused = document.activeElement === el;
          space.insertBefore(el, space.children[index] || null);
          if (focused && document.activeElement !== el) {
            silentFocus = true;
            el.focus({preventScroll: true});
            silentFocus = false;
          }
        }
      });
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
      return WorldHeadCamera.paneSize(layout.ground, head, point);
    }

    function miniatureLink() {
      const el = document.createElementNS(svgNS, 'a');
      el.setAttribute('class', 'hc-rear hc-mini');
      const paths = ['body', 'front'].map(name => {
        const path = document.createElementNS(svgNS, 'path');
        path.setAttribute('class', 'hc-mini-' + name);
        el.append(path);
        return path;
      });
      return {el, paths};
    }

    function placeMiniature(entry, point, width, height) {
      const ground = layout.ground;
      const block = WorldHeadCamera.block(ground, point, width, height, point.blockDepth ?? ground.paneDepth);
      const outline = vertices => vertices ? 'M' + vertices.map(p => p.x + ',' + p.y).join('L') + 'Z' : '';
      entry.paths[0].setAttribute('d', outline(block.side) + outline(block.top));
      entry.paths[1].setAttribute('d', outline(block.front));
      // Narrow the ink as well as the faces near the vanishing point.
      entry.el.style.strokeWidth = Math.min(.85, height * .10) + 'px';
    }

    function syncRearLinks() {
      const focusedState = document.activeElement?.dataset.rearState;
      miniatures.setAttribute('width', layout.width);
      miniatures.setAttribute('height', layout.height);
      const deepestHead = layout.heads.reduce((depth, head) => Math.max(depth, head.depth), 0);
      let miniIndex = 0;
      const fullLinks = [];
      for (const id of order) {
        if (headsByID.has(id)) continue;
        const node = model.N[id], point = layout.points.get(id);
        const {width,height} = paneSize(headsByLane.get(node.lane), point);
        // Only group silhouettes behind every full-size link. This keeps
        // native hit order identical to physical occlusion on either bank.
        const mini = height < 13 && point.depth > deepestHead;
        let entry = rearLinks.get(id);
        if (!entry || mini !== Boolean(entry.paths)) {
          entry?.el.remove();
          if (mini) entry = miniatureLink();
          else {
            const el = document.createElement('a');
            el.className = 'hc-rear';
            entry = attachFront(el);
          }
          entry.el.dataset.rearState = id;
          rearLinks.set(id, entry);
        }
        const {el} = entry;
        const title = options.title(node) || node.title || node.label || id;
        const label = model.laneTitle?.(node.lane) || node.lane;
        el.style.setProperty('--hc-lane', accent(model.LANES[node.lane]));
        el.setAttribute('aria-label', label + ' · ' + title);
        el.setAttribute('href', statePath(node));
        if (mini) {
          placeMiniature(entry, point, width, height);
          if (miniatures.children[miniIndex] !== el) miniatures.insertBefore(el, miniatures.children[miniIndex] || null);
          miniIndex++;
        } else {
          entry.content.textContent = height >= 13 && width >= 50 ? title : '';
          fullLinks.push(el);
        }
      }
      // Tab reaches the current heads, then the nearby labels, before
      // entering the long run of miniatures behind them.
      fullLinks.reverse().forEach((el, index) => {
        if (rearLayer.children[index] !== el) rearLayer.insertBefore(el, rearLayer.children[index] || null);
      });
      for (const [id, {el}] of rearLinks) {
        if (layout.points.has(id) && !headsByID.has(id)) continue;
        el.remove();
        rearLinks.delete(id);
      }
      const restored = rearLinks.get(focusedState)?.el;
      if (restored && restored !== document.activeElement) {
        silentFocus = true;
        restored.focus({preventScroll:true});
        silentFocus = false;
        if (document.activeElement === restored) focusedByUser = focusedState;
      }
    }

    const colorsForPane = (base, edge, tint) => `color-mix(in srgb, ${edge} ${tint}%, ${base})`;

    function drawPane(node, point, head, alpha, theme) {
      const {width,height} = paneSize(head, point);
      const x = point.x - width / 2, y = point.y - height / 2;
      ctx.globalAlpha = alpha;
      const laneColor = accent(model.LANES[node.lane]);
      ctx.strokeStyle = colorsForPane(theme.edge, laneColor, 40);
      ctx.lineWidth = Math.min(1, height * .10);
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
      path(front, height < 13 ? colorsForPane(theme.pane, laneColor, 5) : theme.pane);
      if (width < 50 || height < 13 || angle < .86) return;
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

    function drawGround(left, right, top, bottom, opacity, ink) {
      const ground = layout.ground;
      const {vx, horizon, street, streetY, pitch} = ground;
      const eye = streetY - horizon;
      if (bottom <= horizon || eye <= 0) return;
      const project = (x, z) => WorldHeadCamera.project(ground, x, 0, z);
      const near = eye / (bottom - horizon);
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

      // Every transverse line sits on the recorded plane named by its ruler tick.
      // Continue to project those depths as the camera moves toward the cut.
      const offset = layout.travelZ - floorTravel;
      ctx.strokeStyle = ink;
      ctx.globalAlpha = opacity * .38;
      ctx.beginPath();
      for (const tick of timeTicks) {
        const y = offset ? layout.time.groundYAt(tick.index, offset) : tick.y;
        if (y < top || y > bottom) continue;
        ctx.moveTo(left, y);
        ctx.lineTo(right, y);
      }
      ctx.stroke();
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
      const {width,height} = paneSize(headsByLane.get(model.N[id].lane), point);
      return {x:point.x, y:point.y, rx:width / 2, ry:height / 2};
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
      if (!points.has(active)) return;

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
      timeAxis.style.transform = 'translateY(' + (-y) + 'px)';
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
      const groundOpacity = clamp((angle - .78) / .22, 0, 1);
      if (groundOpacity && layout.ground) {
        drawGround(x, x + width, y, y + height, groundOpacity, colors.edge);
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
      // At rest every history state is a native link. During a camera
      // move, paint the miniature solids at their interpolated depth.
      const movingMiniatures = miniatures.hasAttribute('hidden');
      const glyphOpacity = 1 - paneOpacity;
      for (const id of order) {
        const point = points.get(id), node = model.N[id];
        if (!node || !visible(point, 160)) continue;
        const entry = rearLinks.get(id);
        if (entry && paneOpacity && (rearLayer.hidden || (entry.paths && movingMiniatures))) {
          drawPane(node, point, headsByLane.get(node.lane), paneOpacity, colors);
        }
        if (glyphOpacity) drawGlyph(node, point, glyphOpacity * (contains(id) ? .94 : .40));
      }
      ctx.globalAlpha = 1;
      drawInspection(active, colors, x, y, width, height, dpr);
    }

    function repaint() {
      if (!frame) frame = requestAnimationFrame(paint);
    }

    function updateSelection() {
      const active = inspected();
      const focus = active || (points.has(selected) ? selected : layout?.time.stop);
      const origin = lineage.link(focus) ? lineage.origin : focus;
      const point = inspectable && layout?.time.stop && points.get(origin);
      lineage.update(model, layout, point ? origin : '', point ? shape(origin, point) : null);
      const key = [active, selected, inspectable].join('|');
      if (paintedSelection === key) return;
      paintedSelection = key;
      if (active) scene.dataset.inspectState = active;
      else delete scene.dataset.inspectState;
      const related = neighborhood.get(active);
      const ancestors = new Set(inspectable && layout?.time.stop ? model?.parents[focus] : []);
      const selectedLane = model?.N[selected]?.lane;
      const decorate = (id, el) => {
        const isSelected = id === selected, isInspected = id === active;
        const isRelated = related?.has(id) && !isInspected;
        const showsSelectedLane = Boolean(selectedLane && el.dataset.headLane === selectedLane);
        el.classList.toggle('is-selected', isSelected || showsSelectedLane);
        el.classList.toggle('is-inspected', isInspected);
        el.classList.toggle('is-ancestor', ancestors.has(id));
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
      // The destination can have a wider scroll surface. Carry the old
      // screen positions into its coordinates before the first animation frame.
      const dx = next.ground.vx - layout.ground.vx;
      const origins = dx ? new Map([...points].map(([id,p]) => [id, {...p, x:p.x + dx}])) : points;
      for (const id of rearLinks.keys()) {
        const pane = points.get(id), node = model.N[id];
        const head = node && headsByLane.get(node.lane);
        if (!pane || !head || sizes.has(id)) continue;
        sizes.set(id, paneSize(head, pane));
      }
      for (const [id, entry] of cards) {
        if (incoming.has(id)) continue;
        const {el} = entry;
        const old = origins.get(id), rear = next.points.get(id);
        if (!old || Math.abs(old.x - next.ground.vx) > scroll.clientWidth / 2 + 180) continue;
        const head = headsByID.get(id);
        const from = sizes.get(id) || head;
        const nextHead = nextByLane.get(head.lane);
        const end = rear || movedCamera(old, floorTravel - next.travelZ, next.ground, cameraX - next.ground.cameraX);
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
      return {points:origins, sizes, departures, floorFrom:floorTravel, cameraFrom:cameraX};
    }

    function render(T, next, side, value, focused, travel) {
      angle = clamp(value, 0, 1);
      scene.hidden = !angle;
      if (!angle || !T || !next) {
        inspectable = false;
        hovered = '';
        focusedByUser = '';
        preview.hide();
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
      miniatures.toggleAttribute('hidden', !inspectable);
      timeAxis.hidden = !inspectable || !timeAxis.childElementCount;
      if (!travel) {
        passingLayer?.remove();
        passingLayer = null;
      }
      if (!inspectable) {
        hovered = '';
        preview.hide();
      }
      const viewportWidth = scroll.clientWidth;
      if (viewportWidth && (panAxis !== layout.ground?.vx || panWidth !== viewportWidth)) {
        panWidth = viewportWidth;
        panAxis = layout.ground?.vx ?? layout.width / 2;
        scroll.scrollLeft = panAxis - viewportWidth / 2 + panOffset;
      }
      points = new Map();
      faceSizes.clear();
      const targetTravel = next.travelZ || 0;
      floorTravel = travel ? travel.floorFrom + (targetTravel - travel.floorFrom) * travel.progress : targetTravel;
      cameraX = travel ? travel.cameraFrom + (next.ground.cameraX - travel.cameraFrom)*travel.progress : next.ground.cameraX;
      const yaw = 90 * Math.max(0, 1 - angle / .94);
      const scrollX = scroll.scrollLeft, scrollY = scroll.scrollTop;
      const axisX = scrollX + scroll.clientWidth / 2;
      for (const [id, endpoint] of layout.points) {
        const sidePoint = sidePoints?.get(id);
        const start = sidePoint ? {...sidePoint, x: sidePoint.x + scrollX, y: sidePoint.y + scrollY} : endpoint;
        let point = angle === 1 ? endpoint : WorldHeadCamera.turn(start, endpoint, angle, axisX);
        const head = headsByID.get(id);
        let displayed = head && {width:head.width, height:head.height};
        if (travel) {
          const past = travel.points.get(id);
          const origin = past || movedCamera(endpoint, targetTravel - travel.floorFrom, next.ground, next.ground.cameraX - travel.cameraFrom);
          const ratio = endpoint.depth / origin.depth;
          const target = displayed || paneSize(headsByLane.get(model.N[id].lane), endpoint);
          const previous = travel.sizes.get(id) || {width:target.width * ratio, height:target.height * ratio};
          const pose = cameraTravel(origin, endpoint, previous, target, travel.progress);
          point = pose.point;
          if (head) displayed = pose.size;
        }
        if (point !== endpoint) {
          point.blockDepth ??= endpoint.blockDepth;
          if (endpoint.paneHeight !== undefined) {
            point.paneHeight = endpoint.paneHeight * endpoint.depth / point.depth;
            point.paneWidth = endpoint.paneWidth * endpoint.depth / point.depth;
          }
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
          if (entry.paths) continue;
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
      if (inspectable && !preview.element.hidden) showPreview();
      repaint();
    }

    function face(event) {
      return event.target.closest(nodeLink);
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
      focusedByUser = stateID(el);
      hovered = '';
      showPreview();
      updateSelection();
      repaint();
    }

    function updatePointer(event) {
      if (!inspectable || event.pointerType === 'touch') return;
      const {clientX, clientY} = event;
      pointer = {clientX, clientY, pointerType:event.pointerType};
      const target = document.elementFromPoint(clientX, clientY);
      const link = target?.closest(nodeLink);
      // Crossing the small gutter between a face and its forward strips
      // keeps that face's branch open until the pointer reaches the link.
      if (!link && target?.closest('.hc-forward')) return;
      const id = stateID(link);
      if (id !== hovered) {
        hovered = id;
        updateSelection();
        repaint();
      }
      showPreview();
    }

    function showPreview() {
      const id = inspected();
      const anchor = cards.get(id)?.el || rearLinks.get(id)?.el || lineage.link(id);
      if (hovered || anchor?.matches(':focus-visible')) preview.show(model, id, anchor, hovered ? pointer : null);
      else preview.hide();
    }

    scene.addEventListener('focusin', event => {
      const el = face(event);
      if (!silentFocus && hasVisibleFocus(el)) {
        hovered = '';
        focusedByUser = stateID(el);
        if (!el.dataset.lineageState) options.select(focusedByUser);
        showPreview();
        updateSelection();
        repaint();
      }
    });
    scene.addEventListener('click', event => {
      if (event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
      const id = stateID(face(event));
      if (!id) return;
      event.preventDefault();
      event.stopPropagation();
      options.open(id);
    });
    scene.addEventListener('focusout', event => {
      if (!face(event) || event.relatedTarget?.closest?.(nodeLink)) return;
      focusedByUser = '';
      showPreview();
      updateSelection();
      repaint();
    });
    scene.addEventListener('pointermove', updatePointer);
    scene.addEventListener('pointerover', updatePointer);
    scene.addEventListener('pointerleave', () => {
      pointer = null;
      hovered = '';
      showPreview();
      updateSelection();
      repaint();
    });
    scroll.addEventListener('scroll', () => {
      // A resize can clamp scrollLeft before layout runs. Remember the
      // deliberate offset only while the viewport matches the framed scene.
      if (panWidth && panWidth === scroll.clientWidth) panOffset = scroll.scrollLeft + panWidth / 2 - panAxis;
      if (pointer) updatePointer(pointer);
      else if (!tip.hidden) showPreview();
      repaint();
    }, {passive:true});
    const resize = new ResizeObserver(repaint);
    resize.observe(scene);
    const appearance = new MutationObserver(repaint);
    appearance.observe(document.documentElement, {attributes:true, attributeFilter:['class', 'data-theme']});

    function center() {
      panOffset = 0;
      scroll.scrollLeft = layout.ground.vx - scroll.clientWidth / 2;
    }
    return {render, element:scene, startTravel, center, focus: focusHead, destroy() {
      if (frame) cancelAnimationFrame(frame);
      preview.destroy();
      resize.disconnect();
      appearance.disconnect();
      scene.remove();
    }};
  }

  window.WorldHeadScene = {create};
})();
