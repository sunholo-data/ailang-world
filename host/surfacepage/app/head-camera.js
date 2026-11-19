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

  function placeBanks(band, ground) {
    // The innermost front corner meets the projected curb at its own depth.
    // Extra faces pack outwards within each bank to keep the text clear.
    const banks = [-1, 1].map(side => ({side, edge: 0}));
    for (const [index, member] of band.members.entries()) {
      const bank = banks[(index + band.lead) % 2];
      const curb = ground.street / (2 * member.z);
      const edge = Math.max(curb, bank.edge);
      member.x = bank.side * (edge + member.width / 2);
      bank.edge = edge + member.width + band.gap;
    }
    return Math.max(0, ...banks.map(bank => bank.edge - band.gap));
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
    const tiers = [
      {strips: near, gap, lead: 1},
      {strips: middle, gap: gap * .75, lead: 0},
      {strips: far, gap: gap * .62, lead: 1}
    ];
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
    const ground = {vx: 0, horizon, streetY, street, focal, front: horizon + span * 1.35,
      pitch: faceWidth * .74, blockDepth: .115};
    let width = viewportWidth;
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
      }
      width = Math.max(width, Math.ceil(2 * (placeBanks(band, ground) + gutter)));
    }

    ground.vx = width / 2;
    let height = viewportHeight;
    for (const band of tiers) {
      for (const member of band.members) {
        const {lane, states, ageDepth, priority, width: w, height: h, y, z} = member;
        const worldX = member.x * z, elevation = h / 2 * z;
        const x = ground.vx + member.x;
        const scale = 1 / z, id = states[0].state.id;
        heads.push({id, lane, x, y, scale, width: w, height: h, priority,
          attention: member.attention, live: member.live, depth: z});
        points.set(id, {x, y, scale, depth: z});
        height = Math.max(height, Math.ceil(y + h / 2 + 18));

        // The pane centers are one fixed transverse point and elevation
        // viewed at successively greater physical depths on the same ray.
        const lipInset = (h / 2 + 7) / (y - horizon);
        const faceInset = Math.max(lipInset, 1 - z / (z + ground.blockDepth));
        let previous = ageDepth, firstRearDepth;
        for (let index = 1; index < states.length; index++) {
          const record = states[index];
          const age = depthFor(ageOf(record.time));
          const current = Math.max(age, ageDepth + Math.log1p(index) * .48, previous + .008);
          firstRearDepth ??= current;
          const distance = (current - firstRearDepth) * .22 + Math.log1p(index - 1) * .21;
          const progress = faceInset + (1 - faceInset) * .94 * distance / (1 + distance);
          const rearZ = z / (1 - progress);
          points.set(record.state.id, {
            ...project(ground, worldX, elevation, rearZ), scale: 1 / rearZ, depth: rearZ
          });
          previous = current;
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

  const api = {layout, turn, project, block};
  if (typeof window !== 'undefined') window.WorldHeadCamera = api;
  if (typeof module !== 'undefined') module.exports = api;
})();
