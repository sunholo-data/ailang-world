// End-on projection of the timeline's lanes and exact recorded state IDs.
(() => {
  'use strict';
  const DAY = 86400000;
  const clamp = (value, min, max) => Math.min(max, Math.max(min, value));
  const mix = (a, b, t) => a + (b - a) * t;
  const perspective = depth => .24 + .76 / (1 + depth * .56);

  function runningLanes(T, states) {
    const lanes = new Set(), episodes = new Map();
    const edge = T.D?.nowEdge;
    // Follow mode refreshes the Running summary even when the graph is unchanged.
    if (edge?.open) {
      for (const session of edge.open.live || []) {
        if (session.lane && T.N[session.id]?.lane === session.lane) lanes.add(session.lane);
      }
      return lanes;
    }
    for (const state of states) {
      if (state.sessionLive || state.open?.sessionLive) lanes.add(state.lane);
      if (!state.episode) continue;
      const matches = episodes.get(state.episode) || new Set();
      matches.add(state.lane);
      episodes.set(state.episode, matches);
    }
    if (!edge?.sessionsUnavailable) {
      for (const session of edge?.sessions || []) {
        const matches = episodes.get(session.episode);
        if (matches?.size === 1) lanes.add(matches.values().next().value);
      }
    }
    return lanes;
  }

  function laneStrips(T, selected, attention) {
    const records = Object.values(T.N);
    const byLane = new Map();
    for (const state of records) {
      if (!byLane.has(state.lane)) byLane.set(state.lane, []);
      byLane.get(state.lane).push({state, time: Date.parse(state.at)});
    }
    const now = Number.isFinite(T.TNOW) ? T.TNOW : Date.now();
    const ageOf = time => Number.isFinite(time) ? Math.max(0, (now - time) / DAY) : Infinity;
    const live = runningLanes(T, records);
    const strips = [...byLane].map(([lane, states]) => {
      // Stable rank ordering preserves the timeline's tip even if an
      // appended episode's timestamp predates another record in the lane.
      states.sort((a, b) => b.state.rank - a.state.rank);
      const age = states.reduce((youngest, record) => Math.min(youngest, ageOf(record.time)), Infinity);
      const priority = lane === selected || T.N[selected]?.lane === lane ? 4
        : live.has(lane) ? 3 : attention?.has(lane) ? 2 : age < 2 ? 1 : 0;
      return {lane, states, age, priority, attention:attention?.has(lane), live:live.has(lane)};
    }).sort((a, b) => b.priority - a.priority || a.age - b.age ||
      b.states[0].state.rank - a.states[0].state.rank || a.lane.localeCompare(b.lane));
    if (strips.length && !strips[0].priority) strips[0].priority = 1;
    const near = [], middle = [], far = [];
    for (const strip of strips) {
      if (strip.age < 2) near.push(strip);
      else if (strip.age < 7) middle.push(strip);
      else far.push(strip);
    }
    return {near, middle, far, ageOf};
  }

  // A level image plane. Its ground reference is the street line at Z=1.
  function project(ground, x, elevation, z) {
    return {
      x: ground.vx + x / z,
      y: ground.horizon + (ground.streetY - ground.horizon - elevation) / z
    };
  }

  function block(ground, center, width, height, length = ground?.blockDepth) {
    const {vx, horizon, streetY} = ground || {};
    if (![vx, horizon, streetY, center?.x, center?.y, width, height, length].every(Number.isFinite) ||
        !(width > 0 && height > 0 && length > 0 && streetY > horizon)) return null;
    const l = center.x - width / 2, r = center.x + width / 2;
    const t = center.y - height / 2, b = center.y + height / 2;
    const eyeHeight = streetY - horizon;
    const explicitDepth = Number.isFinite(center.depth) && center.depth > 0;
    if (!explicitDepth && b <= horizon) return null;
    const depth = explicitDepth ? center.depth : eyeHeight / (b - horizon);
    const backDepth = depth + length;
    const front = [{x:l,y:t}, {x:r,y:t}, {x:r,y:b}, {x:l,y:b}];
    const back = front.map(({x, y}) => project(ground, (x - vx) * depth, eyeHeight - (y - horizon) * depth, backDepth));
    const top = t > horizon ? [front[0], front[1], back[1], back[0]] : null;
    const side = vx < l ? [front[0], front[3], back[3], back[0]]
      : vx > r ? [front[1], back[1], back[2], front[2]] : null;
    // An elevated historical label still casts its footprint onto the
    // ground plane at the same pair of physical depths.
    const worldLeft = (l - vx) * depth, worldRight = (r - vx) * depth;
    const floor = [project(ground, worldLeft, 0, depth), project(ground, worldRight, 0, depth),
      project(ground, worldRight, 0, backDepth), project(ground, worldLeft, 0, backDepth)];
    return {front, back, top, side, floor, depth, backDepth};
  }

  function measureLane(member, ground, worldHeight) {
    const eyeHeight = ground.streetY - ground.horizon;
    const centerDrop = eyeHeight - worldHeight / 2;
    const paneHalfHeight = worldHeight * ground.paneHeightRatio / 2;
    const roof = project(ground, 0, worldHeight, member.z + ground.blockDepth).y;
    const first = Math.max(member.z + ground.blockDepth,
      (centerDrop + paneHalfHeight) / (roof - ground.horizon - 4));
    const gapRatio = 2 * paneHalfHeight / (centerDrop - paneHalfHeight);
    const historyPitch = Math.max(ground.historyPitch, first * gapRatio / (1 - 3 * gapRatio) * 1.06);
    member.historyPitch = historyPitch;
    member.firstRear = first;
    const halfWidth = member.width * member.z / 2;
    // A whole lane needs at most two bounds: its head and the run of
    // smaller cards behind it. Their projected sides converge on the axis.
    const band = (front, back, halfWidth, halfHeight) => ({
      front, back, halfWidth, topDrop: centerDrop - halfHeight, bottomDrop: centerDrop + halfHeight,
      top: (centerDrop - halfHeight) / back, bottom: (centerDrop + halfHeight) / front
    });
    member.bands = [band(member.z, member.z + ground.blockDepth, halfWidth, worldHeight / 2)];
    if (member.states.length > 1) {
      const last = first + (member.states.length - 2) * historyPitch;
      member.bands.push(band(first, last + ground.paneDepth, halfWidth * ground.paneWidthRatio, paneHalfHeight));
    }
  }

  function bankPosition(member, placed, side, ground) {
    const blocked = [];
    for (const other of placed) {
      const innerX = other.worldX * side;
      if (innerX < 0) continue;
      for (const band of member.bands) {
        for (const peer of other.bands) {
          if (band.bottom <= peer.top || peer.bottom <= band.top) continue;
          // A card fits when either its projected X bounds or its rays
          // from the vanishing point clear the neighboring band.
          const inner = innerX - peer.halfWidth, outer = innerX + peer.halfWidth;
          const left = inner * Math.max(band.topDrop / peer.bottomDrop, band.front / peer.back);
          const right = outer * Math.min(band.bottomDrop / peer.topDrop, band.back / peer.front);
          blocked.push([left - band.halfWidth - 8, right + band.halfWidth + 8]);
        }
      }
    }
    let x = ground.street / 2 + member.bands[0].halfWidth;
    // Use the first free position on this bank, including gaps between lanes.
    for (const [left, right] of blocked.sort((a, b) => a[0] - b[0])) {
      if (x <= left) break;
      x = Math.max(x, right);
    }
    return x;
  }

  function placeLanes(T, members, ground, gutter) {
    const byLane = new Map(members.map(member => [member.lane, member]));
    const shared = (T.GROUPS || []).map(group => byLane.get(group.shared)).filter(Boolean);
    // Follow's first represented repo owns the optical axis. A legacy
    // model uses its shared lane or the lane containing shared history.
    const center = shared[0] || byLane.get('shared') ||
      members.find(member => member.states.some(record => record.state.kind === 'shared'));
    ground.centerLane = center?.lane || '';
    const placed = [];
    const groupSides = new Map();
    let reach = 0;
    if (center) {
      center.worldX = 0;
      placed.push(center);
      reach = center.width / 2;
    }
    const ordered = [...shared.slice(1), ...members.filter(member => member !== center && !shared.includes(member))
      .sort((a, b) => b.priority - a.priority || a.z - b.z || a.lane.localeCompare(b.lane))];
    for (const [i, member] of ordered.entries()) {
      const repoLane = T.groupOf?.(member.lane)?.shared;
      const preference = groupSides.get(repoLane) || (i % 2 ? -1 : 1);
      const preferred = bankPosition(member, placed, preference, ground);
      const opposite = bankPosition(member, placed, -preference, ground);
      // Keep a repo's faces together where the opposite bank would only
      // save a negligible amount of screen width.
      const side = preferred <= opposite + gutter * member.z ? preference : -preference;
      member.worldX = side * (side === preference ? preferred : opposite);
      if (shared.includes(member)) groupSides.set(member.lane, side);
      placed.push(member);
      reach = Math.max(reach, Math.abs(member.worldX) / member.z + member.width / 2);
    }
    return reach;
  }

  function layout(T, options = {}) {
    const viewportWidth = Math.max(1, Math.round(options.width || 320));
    const viewportHeight = Math.max(1, Math.round(options.height || 420));
    const points = new Map(), heads = [], edges = [];
    if (!T?.N) return {points, heads, width: viewportWidth, height: viewportHeight, edges};

    const {near, middle, far, ageOf} = laneStrips(T, options.selected, options.attention);
    const narrow = viewportWidth < 560;
    const gutter = narrow ? 12 : 34;
    const gap = narrow ? 12 : 28;
    const street = clamp(viewportWidth * .18, 66, 184);
    const faceWidth = clamp((viewportWidth - street - 2 * (narrow ? 5 : gutter)) / 2, 120, 310);
    // Desktop can fit several faces along each bank without losing the
    // street. A minimum readable width keeps denser scenes pannable.
    const occupied = Math.max(1, Math.ceil(near.length / 2));
    const bankWidth = (viewportWidth - street) / 2 - gutter;
    const nearWidth = narrow ? faceWidth : clamp((bankWidth - (occupied - 1) * gap) / occupied, 172, faceWidth);
    const nearHeight = narrow ? clamp(61 + (viewportWidth - 320) * .16, 61, 82) : clamp(nearWidth * .35, 74, 96);
    const horizon = Math.max(44, viewportHeight * .19);
    const span = Math.max(nearHeight * 3.8, viewportHeight * .50);
    const tiers = [{strips: near}, {strips: middle}, {strips: far}];
    let oldestKnown = 14;
    for (const band of tiers) {
      for (const strip of band.strips) {
        for (const record of strip.states) {
          const age = ageOf(record.time);
          if (Number.isFinite(age)) oldestKnown = Math.max(age, oldestKnown);
        }
      }
    }
    const unknownDepth = Math.log1p(oldestKnown) + .3;
    const depthFor = age => Number.isFinite(age) ? Math.log1p(age) : unknownDepth;

    for (const [tier, band] of tiers.entries()) {
      const byAge = [...band.strips].sort((a, b) => a.age - b.age || a.lane.localeCompare(b.lane));
      const ageOrder = new Map(byAge.map(({lane}, i) => [lane, i]));
      band.members = band.strips.map(strip => {
        const ageDepth = depthFor(strip.age);
        const stagger = Math.min(tier === 2 ? 20 : 18, ageOrder.get(strip.lane) * (tier === 0 ? 16 : tier === 1 ? 11 : 6));
        return {...strip, ageDepth, stagger};
      });
    }

    const groundReach = ({age, stagger}) => span * mix(.98, .92, age / 2) + nearHeight / 2 - stagger;
    const eyeHeight = near.length ? Math.max(...tiers[0].members.map(groundReach)) : span * .98 + nearHeight / 2;
    const streetY = horizon + eyeHeight;
    const focal = viewportHeight / (2 * Math.tan(54 * Math.PI / 360));
    const blockDepth = .115;
    const ground = {vx: 0, horizon, streetY, street, focal, front: horizon + span * 1.35,
      pitch: faceWidth * .74, blockDepth, historyPitch: .16, paneHeightRatio: .36, paneWidthRatio: .88,
      paneDepth: blockDepth * .37};
    for (const [tier, band] of tiers.entries()) {
      for (const member of band.members) {
        const {age, ageDepth, priority, stagger} = member;
        // Place the foot in depth first. Both dimensions of the same
        // world-size facade then shrink with the camera's inverse depth.
        const reach = tier === 0 ? groundReach(member) :
          (tier === 1 ? mix(.70, .62, (age - 2) / 5) : mix(.34, .55, perspective(ageDepth))) * eyeHeight - stagger;
        const scale = reach / eyeHeight;
        const z = 1 / scale;
        const widthFactor = tier > 0 && priority ? 1.1 : 1;
        member.z = z;
        member.width = nearWidth * widthFactor / z;
        member.height = nearHeight / z;
        member.y = project(ground, 0, nearHeight / 2, z).y;
        measureLane(member, ground, nearHeight);
      }
    }

    const reach = placeLanes(T, tiers.flatMap(band => band.members), ground, gutter);
    const width = reach <= viewportWidth / 2 ? viewportWidth : Math.ceil(2 * (reach + gutter));
    ground.vx = width / 2;
    let height = viewportHeight;
    for (const band of tiers) {
      for (const member of band.members) {
        const {lane, states, priority, width: w, height: h, y, z, worldX, historyPitch} = member;
        const elevation = nearHeight / 2;
        const x = ground.vx + worldX / z;
        const scale = 1 / z, id = states[0].state.id;
        heads.push({id, lane, x, y, scale, width: w, height: h, priority, historyPitch,
          attention: member.attention, live: member.live, depth: z});
        points.set(id, {x, y, scale, depth: z});
        height = Math.max(height, Math.ceil(y + h / 2 + 18));

        for (let index = 1; index < states.length; index++) {
          const rearZ = member.firstRear + (index - 1) * historyPitch;
          points.set(states[index].state.id, {
            ...project(ground, worldX, elevation, rearZ), scale: 1 / rearZ, depth: rearZ
          });
        }
      }
    }
    const seen = new Set();
    for (const [to, parents] of Object.entries(T.parents || {})) {
      if (!points.has(to)) continue;
      for (const from of parents) {
        const edge = JSON.stringify([from, to]);
        if (from === to || !points.has(from) || seen.has(edge)) continue;
        seen.add(edge);
        edges.push({from, to});
      }
    }
    heads.sort((a, b) => b.priority - a.priority || points.get(a.id).depth - points.get(b.id).depth);
    return {points, heads, width, height, edges, ground};
  }

  function paneSize(ground, head, point) {
    const z = point.depth, scale = head.depth / z;
    const pitch = head.historyPitch;
    const drop = (head.y - ground.horizon) * scale;
    // Between this roof at Z + paneDepth and the next front at Z + pitch,
    // reserve a fifth of the projected interval for the gap and edge strokes.
    // The farther panes taper as this interval falls below their normal height.
    const clearance = drop * (pitch - ground.paneDepth) / (z + pitch);
    return {
      width: head.width * ground.paneWidthRatio * scale,
      height: Math.min(head.height * ground.paneHeightRatio * scale, clearance * .8)
    };
  }

  // The timeline supplies the longitudinal coordinate. The end-on layout
  // supplies the transverse coordinate. Yaw both about one line of sight,
  // then apply perspective at the current depth.
  function turn(a, b, progress, axisX = 0) {
    if (progress <= 0) return {x:a.x, y:a.y, scale:a.scale ?? 1, depth:a.depth ?? 0};
    if (progress >= 1) return {x:b.x, y:b.y, scale:b.scale ?? 1, depth:b.depth ?? 0};
    const yaw = progress * Math.PI / 2;
    const side = Math.cos(yaw), front = Math.sin(yaw);
    const wa = 1 / (a.scale ?? 1), wb = 1 / (b.scale ?? 1);
    const w = mix(wa, wb, front * front);
    return {
      x:axisX + (side * (a.x - axisX) * wa + front * (b.x - axisX) * wb) / w,
      y:mix(a.y * wa, b.y * wb, front * front) / w,
      scale:1 / w,
      depth:w
    };
  }

  const api = {layout, turn, project, block, paneSize};
  if (typeof window !== 'undefined') window.WorldHeadCamera = api;
  if (typeof module !== 'undefined') module.exports = api;
})();
