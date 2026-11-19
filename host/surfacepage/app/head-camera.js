// End-on projection of the timeline's lanes and exact recorded state IDs.
(() => {
  'use strict';
  const DAY = 86400000;
  const clamp = (value, min, max) => Math.min(max, Math.max(min, value));
  const mix = (a, b, t) => a + (b - a) * t;
  const unitsPerZ = 288;

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

  function timeAxis(T, options) {
    const xs = options.timeXs || T.LX[1];
    const presentIndex = T.presentIndex;
    const requested = T.timeIndex.get(options.stop);
    const stop = T.N[options.stop];
    const historical = stop && !stop.draft && Number.isInteger(requested) && requested <= presentIndex;
    const referenceIndex = historical ? requested : presentIndex;
    const referenceX = xs[referenceIndex], presentX = xs[presentIndex];
    const xAt = index => {
      const i = Math.floor(index);
      return index === i ? xs[i] : mix(xs[i], xs[i + 1], index - i);
    };
    return {
      xs, referenceIndex, presentIndex, referenceX, presentX, unitsPerZ,
      firstIndex: 0, lastIndex: referenceIndex, stop: historical ? stop.id : '', xAt,
      depthAt: index => 1 + (referenceX - xAt(index)) / unitsPerZ
    };
  }

  function laneStrips(T, time, options) {
    const byLane = new Map();
    const records = Object.values(T.N);
    const running = time.stop ? new Set() : runningLanes(T, records);
    const focus = time.stop || options.selected;
    for (const state of records) {
      const index = T.timeIndex.get(state.id);
      if (!Number.isInteger(index) || index > time.referenceIndex || (time.stop && state.draft)) continue;
      if (!byLane.has(state.lane)) byLane.set(state.lane, []);
      byLane.get(state.lane).push({state, index, z:time.depthAt(index)});
    }
    const strips = [...byLane].map(([lane, states]) => {
      states.sort((a, b) => b.index - a.index || Number(b.state.id === time.stop) - Number(a.state.id === time.stop) ||
        b.state.rank - a.state.rank);
      const idx = states[0].index;
      const age = T.undated && idx === 0 ? Infinity : Math.max(0, (T.TT[time.referenceIndex] - T.TT[idx]) / DAY);
      const live = running.has(lane), attention = !time.stop && options.attention?.has(lane);
      const priority = lane === focus || T.N[focus]?.lane === lane ? 4 : live ? 3 : attention ? 2 : age < 2 ? 1 : 0;
      return {lane, states, age, priority, attention, live, z:states[0].z};
    }).sort((a, b) => b.priority - a.priority || a.z - b.z || a.lane.localeCompare(b.lane));
    if (strips.length && !strips[0].priority) strips[0].priority = 1;
    return strips;
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

  function measureLane(member, ground, centered, base) {
    const drop = ground.streetY - ground.horizon - ground.eventElevation;
    const times = [...new Set(member.states.map(record => record.z))];
    const byTime = new Map(times.map((z, i) => [z, Math.min(z - (times[i - 1] ?? -Infinity),
      (times[i + 1] ?? Infinity) - z)]));
    const slots = [];
    let companion = 0;
    for (const [index, record] of member.states.entries()) {
      const {z} = record;
      const pitch = byTime.get(z);
      const length = index ? Math.min(ground.paneDepth, pitch * .18) : ground.blockDepth;
      const nominal = ground.worldHeight * (index ? ground.paneHeightRatio : 1);
      const height = index ? Math.min(nominal, drop * pitch / (z + pitch) * .58 || nominal) : nominal;
      const halfWidth = member.worldWidth / 2 * (index ? ground.paneWidthRatio : 1);
      const back = z + length;
      const slot = {record, z, back, halfWidth, height, length, factor:1, offset:0, direction:0,
        top:(drop - height / 2) / back, bottom:(drop + height / 2) / z};
      const peers = slots.filter(peer => verticallyOverlaps(slot, peer));
      // Keep a companion's bank when earlier slots become heads or leave the cut.
      const bankIndex = peers.length ? ++companion : index;
      slot.bankSide = base?.points.get(record.state.id)?.bankSide ?? (bankIndex % 2 ? 1 : -1);
      if (peers.length) {
        slot.direction = centered ? slot.bankSide : 1;
        for (const peer of peers) {
          if (centered && peer.direction && peer.direction !== slot.direction) continue;
          // Put the whole solid next to its overlapping neighbour. The
          // affine world-X offset clears it at any base bank position.
          slot.factor = Math.max(slot.factor, peer.factor * back / peer.z);
          slot.offset = Math.max(slot.offset, (peer.offset + peer.halfWidth + 5) * back / peer.z + halfWidth);
        }
        if (!slot.offset) slot.direction = 0;
      }
      slots.push(slot);
    }
    member.slots = slots;
  }

  const verticallyOverlaps = (a, b) => a.top < b.bottom && b.top < a.bottom;

  function bankPosition(member, placed, side, ground, anchor = 0) {
    const blocked = [];
    for (const other of placed) {
      for (const slot of member.slots) {
        for (const peer of other.slots) {
          if (!verticallyOverlaps(slot, peer)) continue;
          const x = peer.worldX * side;
          const lo = x - peer.halfWidth, hi = x + peer.halfWidth;
          if (hi <= 0) continue;
          const left = lo / (lo < 0 ? peer.z : peer.back);
          const right = hi / peer.z;
          const gap = 4 / Math.max(slot.z, peer.z);
          blocked.push([((left - gap) * slot.z - slot.halfWidth - slot.offset) / slot.factor,
            ((right + gap) * slot.back + slot.halfWidth - slot.offset) / slot.factor]);
        }
      }
    }
    let x = Math.max(anchor, ground.street / 2 + member.slots[0].halfWidth);
    for (const [left, right] of blocked.sort((a, b) => a[0] - b[0])) {
      if (x <= left) break;
      x = Math.max(x, right);
    }
    return x;
  }

  function placeLanes(T, members, ground, gutter, base, viewLane) {
    const byLane = new Map(members.map(member => [member.lane, member]));
    const shared = (T.GROUPS || []).map(group => byLane.get(group.shared)).filter(Boolean);
    const defaultCenter = shared[0] || byLane.get('shared') ||
      members.find(member => member.states.some(record => record.state.kind === 'shared')) || members[0];
    const prior = new Map(base?.heads.map(head => [head.lane, head]) || []);
    const centerLane = prior.has(viewLane) ? viewLane : base?.ground.centerLane || defaultCenter?.lane || '';
    const center = byLane.get(centerLane);
    ground.centerLane = centerLane;
    const placed = [], groupSides = new Map();
    let reach = 0;
    for (const member of members) {
      const emphasize = member.age >= 2 && member.priority ? 1.1 : 1;
      const head = prior.get(member.lane);
      member.worldWidth = head ? head.width * head.depth : ground.worldWidth * emphasize;
      measureLane(member, ground, member === center, base);
    }
    function place(member, side, anchor) {
      member.worldX = side * anchor;
      for (const slot of member.slots) {
        slot.worldX = center === member ? slot.direction * slot.offset : side * (slot.factor * anchor + slot.offset);
        reach = Math.max(reach, (Math.abs(slot.worldX) + slot.halfWidth) / slot.z);
      }
      placed.push(member);
    }
    if (center) place(center, 0, 0);
    if (base) {
      // Move along the same street through every cut. Lanes that have no records
      // here still define the order and can still be the camera's reference.
      const origin = prior.get(centerLane)?.worldX ?? 0;
      const middle = base.laneOrder.indexOf(centerLane);
      for (const [side, queue] of [[-1, base.laneOrder.slice(0, middle).reverse()], [1, base.laneOrder.slice(middle + 1)]]) {
        let last = 0;
        for (const lane of queue) {
          const member = byLane.get(lane);
          if (!member) continue;
          const anchor = Math.max(last + 1, Math.abs(prior.get(lane).worldX - origin));
          last = bankPosition(member, placed, side, ground, anchor);
          place(member, side, last);
        }
      }
      return reach;
    }
    const ordered = [...shared.slice(1), ...members.filter(member => member !== center && !shared.includes(member))];
    for (const [i, member] of ordered.entries()) {
      const repoLane = T.groupOf?.(member.lane)?.shared;
      const preference = groupSides.get(repoLane) || (i % 2 ? -1 : 1);
      const preferred = bankPosition(member, placed, preference, ground);
      const opposite = bankPosition(member, placed, -preference, ground);
      const side = preferred <= opposite + gutter * member.z ? preference : -preference;
      place(member, side, side === preference ? preferred : opposite);
      if (shared.includes(member)) groupSides.set(member.lane, side);
    }
    return reach;
  }

  function layout(T, options = {}) {
    const viewportWidth = Math.max(1, Math.round(options.width || 320));
    const viewportHeight = Math.max(1, Math.round(options.height || 420));
    const points = new Map(), heads = [], edges = [];
    if (!T?.N) return {points, heads, width: viewportWidth, height: viewportHeight, edges};

    const time = timeAxis(T, options);
    const members = laneStrips(T, time, options);
    const base = time.stop || options.viewLane
      ? options.base || layout(T, {...options, stop:undefined, viewLane:undefined}) : null;
    const narrow = viewportWidth < 560;
    const gutter = narrow ? 12 : 34;
    const gap = narrow ? 12 : 28;
    const street = clamp(viewportWidth * .18, 66, 184);
    const faceWidth = clamp((viewportWidth - street - 2 * (narrow ? 5 : gutter)) / 2, 120, 310);
    const near = members.filter(member => member.age < 2).length;
    const occupied = Math.max(1, Math.ceil(near / 2));
    const bankWidth = (viewportWidth - street) / 2 - gutter;
    const nearWidth = base?.ground.worldWidth ??
      (narrow ? faceWidth : clamp((bankWidth - (occupied - 1) * gap) / occupied, 172, faceWidth));
    const nearHeight = base?.ground.worldHeight ??
      (narrow ? clamp(61 + (viewportWidth - 320) * .16, 61, 82) : clamp(nearWidth * .35, 74, 96));
    const horizon = Math.max(44, viewportHeight * .19);
    const span = Math.max(nearHeight * 3.8, viewportHeight * .50);
    const streetY = horizon + span * .98 + nearHeight / 2;
    const focal = viewportHeight / (2 * Math.tan(54 * Math.PI / 360));
    const blockDepth = .115;
    const ground = {vx: 0, horizon, streetY, street, focal, front: horizon + span * 1.35,
      pitch: faceWidth * .74, blockDepth, paneHeightRatio: .36, paneWidthRatio: .88,
      paneDepth: blockDepth * .37, eventElevation:nearHeight / 2, worldWidth:nearWidth, worldHeight:nearHeight};
    time.yAt = index => project(ground, 0, ground.eventElevation, time.depthAt(index)).y;
    const reach = placeLanes(T, members, ground, gutter, base, options.viewLane);
    const width = reach <= viewportWidth / 2 ? viewportWidth : Math.ceil(2 * (reach + gutter));
    ground.vx = width / 2;
    let height = viewportHeight;
    for (const member of members) {
      const {lane, priority, attention, live, worldX} = member;
      for (const [index, slot] of member.slots.entries()) {
        const {state, index:timeIndex} = slot.record;
        const {z} = slot;
        const point = {...project(ground, slot.worldX, ground.eventElevation, z), worldX:slot.worldX, bankSide:slot.bankSide,
          scale:1/z, depth:z, timeIndex};
        if (!index) {
          heads.push({...point, id:state.id, lane, priority, attention, live, worldX,
            width:member.worldWidth/z, height:nearHeight/z});
          height = Math.max(height, Math.ceil(point.y + nearHeight / (2 * z) + 18));
        } else {
          point.paneHeight = slot.height / z;
          point.paneDepth = slot.length;
        }
        points.set(state.id, point);
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
    const laneOrder = base?.laneOrder || [...heads].sort((a,b) => a.worldX - b.worldX).map(head => head.lane);
    return {points, heads, width, height, edges, ground, time, laneOrder,
      travelZ:(time.presentX - time.referenceX) / unitsPerZ};
  }

  function paneSize(ground, head, point) {
    const scale = head.depth / point.depth;
    return {width:head.width * ground.paneWidthRatio * scale,
      height:point.paneHeight ?? head.height * ground.paneHeightRatio * scale};
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
