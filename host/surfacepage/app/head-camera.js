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

  function placeBanks(band) {
    // Lateral relevance fills the entrances first, with newer or quieter
    // faces still retaining their own world depth and staggering.
    const banks = [-1, 1].map(side => ({side, edge: band.street / 2}));
    for (const [index, member] of band.members.entries()) {
      const bank = banks[(index + band.lead) % 2];
      member.x = bank.side * (bank.edge + member.width / 2);
      bank.edge += member.width + band.gap;
    }
    return Math.max(...banks.map(bank => bank.edge - band.gap));
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
    const middleHeight = narrow ? clamp(50 + (viewportWidth - 320) * .12, 50, 63) : 65;
    const horizon = clamp(viewportHeight * .13, 45, 130);
    const span = Math.max(320, viewportHeight * .58);
    const tiers = [
      {strips: near, gap, street, lead: 1, width: nearWidth, height: nearHeight},
      {strips: middle, gap: gap * .75, street: street * .74, lead: 0, width: faceWidth * .76, height: middleHeight},
      {strips: far, gap: gap * .62, street: street * .56, lead: 1, width: Math.max(102, faceWidth * .57), height: 40}
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

    let width = viewportWidth;
    for (const [tier, band] of tiers.entries()) {
      const byAge = [...band.strips].sort((a, b) => a.age - b.age || a.lane.localeCompare(b.lane));
      const ageOrder = new Map(byAge.map(({lane}, i) => [lane, i]));
      band.members = band.strips.map(strip => {
        const depth = depthFor(strip.age);
        const scale = perspective(depth);
        const shrink = tier === 2 ? mix(.65, 1, scale) : 1;
        const regularWidth = band.width * shrink;
        const prominentWidth = tier === 2 ? faceWidth * .86 : tier === 1 ? faceWidth * .88 : nearWidth;
        const width = strip.priority ? Math.max(regularWidth, prominentWidth) : regularWidth;
        const stagger = Math.min(tier === 2 ? 20 : 18, ageOrder.get(strip.lane) * (tier === 0 ? 16 : tier === 1 ? 11 : 6));
        return {...strip, depth, scale, width, height: band.height, stagger};
      });
      width = Math.max(width, Math.ceil(2 * (placeBanks(band) + gutter)));
    }

    const center = width / 2;
    let height = viewportHeight;
    let streetY = near.length ? 0 : horizon + span * .98 + nearHeight / 2;
    for (const [tier, band] of tiers.entries()) {
      for (const member of band.members) {
        const {lane, states, age, depth, scale, priority, width: w, height: h} = member;
        const x = center + member.x;
        const distance = tier === 0 ? mix(.98, .92, age / 2)
          : tier === 1 ? mix(.60, .52, (age - 2) / 5) : mix(.16, .34, scale);
        const y = horizon + span * distance - member.stagger;
        const towardCenter = x < center ? 1 : -1;
        const returnDepth = tier === 0 ? clamp(w * .13, narrow ? 16 : 25, narrow ? 27 : 48) : clamp(w * .12, 7, 27);
        const returnX = towardCenter * returnDepth;
        const returnY = tier === 0 ? -clamp(w * .07, 11, 18) : -clamp(h * .15, 6, 16);
        const id = states[0].state.id;
        heads.push({id, lane, x, y, scale, width: w, height: h, priority,
          attention: member.attention, live: member.live, returnX, returnY});
        points.set(id, {x, y, scale, depth});
        height = Math.max(height, Math.ceil(y + h / 2 + 18));
        if (tier === 0) streetY = Math.max(streetY, y + h / 2);

        // The first rear state touches the face's lip. Beyond it, the lane
        // tapers toward its own point on the common horizon.
        const vanishingX = clamp(center + (x - center) * .32 - towardCenter * (tier + 1) * 9, 4, width - 4);
        const faceInset = (h / 2 + 7) / (y - horizon);
        let previous = depth, firstRearDepth;
        for (let index = 1; index < states.length; index++) {
          const record = states[index];
          const current = Math.max(depthFor(ageOf(record.time)), depth + Math.log1p(index) * .48, previous + .008);
          firstRearDepth ??= current;
          const distance = (current - firstRearDepth) * .22 + Math.log1p(index - 1) * .21;
          const progress = faceInset + (1 - faceInset) * .94 * distance / (1 + distance);
          points.set(record.state.id, {
            x: mix(x, vanishingX, progress), y: mix(y, horizon, progress),
            scale: perspective(current) * .84, depth: current
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
    return {points, heads, width, height, edges, ground:{horizon, front:horizon + span * 1.35, pitch:faceWidth * .74, street, streetY}};
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
      depth:mix(a.depth ?? 0, b.depth ?? 0, front * front)
    };
  }

  const api = {layout, turn};
  if (typeof window !== 'undefined') window.WorldHeadCamera = api;
  if (typeof module !== 'undefined') module.exports = api;
})();
