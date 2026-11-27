// End-on projection of the timeline's lanes and exact recorded state IDs.
(() => {
  'use strict';
  const DAY = 86400000;
  const clamp = (value, min, max) => Math.min(max, Math.max(min, value));
  const mix = (a, b, t) => a + (b - a) * t;
  const unitsPerZ = 216;

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
    const timeline = options.timeXs || T.LX[1];
    // At distant Timeline zooms multiple instants fit inside one pixel. Give
    // distinct rows enough depth for a readable foreground face in Heads.
    const xs = options.base?.time.xs || timeline.reduce((axis, x, i) => {
      axis.push(i ? axis[i-1] + Math.max(42, x-timeline[i-1]) : x);
      return axis;
    }, []);
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
      states.sort((a, b) => b.index - a.index ||
        Number(b.state.id === time.stop) - Number(a.state.id === time.stop) ||
        Number(b.state.kind === 'shared') - Number(a.state.kind === 'shared') ||
        b.state.rank - a.state.rank || a.state.id.localeCompare(b.state.id));
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

  function block(ground, center, width, height, length = center?.blockDepth ?? ground?.blockDepth) {
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

  // Assign a physical track once, from the whole graph. A merge runs along
  // its lane axis; simultaneous source states occupy adjacent tracks.
  function streetPlan(T, members, ground, gap, viewportWidth) {
    const tracks = new Map(), lanes = new Map();
    const shared = (T.GROUPS || []).map(group => members.find(m => m.lane === group.shared)).filter(Boolean);
    const center = shared[0] || members.find(m => m.states.some(r => r.state.kind === 'shared')) || members[0];
    let left = 0, right = 0;
    const groupSide = new Map();
    for (const member of new Set([center, ...shared, ...members].filter(Boolean))) {
      let instant, pos = 0, lo = 0, hi = 0;
      for (const {state,index} of member.states) {
        if (instant !== index) { instant = index; pos = 0; }
        const offset = pos % 2 ? (pos + 1) / 2 : -pos / 2;
        tracks.set(state.id, offset);
        lo = Math.min(lo, offset);
        hi = Math.max(hi, offset);
        pos++;
      }
      const repo = T.groupOf?.(member.lane)?.shared;
      const preferred = T.preferredSides?.get(member.lane);
      const side = (preferred === 1 || preferred === -1) ? preferred : (groupSide.get(repo) || (right <= -left ? 1 : -1));
      const origin = member === center ? 0 : side > 0 ? right + 1 - lo : left - 1 - hi;
      left = Math.min(left, origin + lo);
      right = Math.max(right, origin + hi);
      if (member !== center && shared.includes(member)) groupSide.set(member.lane, side);
      const width = ground.worldWidth * (member.age >= 2 && member.priority ? 1.1 : 1);
      // Leave room beside a centered merge for its first source on a phone.
      const peerWidth = clamp(viewportWidth/2 - width/2 - gap, 44, width);
      lanes.set(member.lane, {origin, width, peerWidth});
    }
    const pitch = members.reduce((w,m) => Math.max(w, lanes.get(m.lane).width), ground.worldWidth) + gap;
    const order = [...lanes.keys()].sort((a,b) => lanes.get(a).origin - lanes.get(b).origin);
    return {tracks, lanes, pitch, order, center:center?.lane || ''};
  }

  function measureCards(members, plan, ground, gap) {
    const cards = [], drop = ground.streetY - ground.horizon - ground.eventElevation;
    for (const member of members) {
      const lane = plan.lanes.get(member.lane);
      member.slots = member.states.map((record, i) => {
        const z = record.z, front = i === 0;
        const track = plan.tracks.get(record.state.id);
        const offset = track*(lane.peerWidth + gap) + Math.sign(track)*(lane.width-lane.peerWidth)/2;
        const x = lane.origin * plan.pitch - ground.cameraX + offset;
        const width = track ? lane.peerWidth : lane.width;
        const halfW = width / (2*z) * (front ? 1 : ground.paneWidthRatio);
        const halfH = ground.worldHeight / (2*z) * (front ? 1 : ground.paneHeightRatio);
        const y = drop/z;
        // The roof of a distant lateral track stays within its clearance
        // from the next track at the same instant.
        const length = Math.min(front ? ground.blockDepth : ground.paneDepth,
          gap / (Math.abs(x)/z + halfW));
        const ratio = z/(z+length), l = x/z-halfW, r = x/z+halfW;
        const card = {record, front, x, y, halfW, halfH, length, up:Infinity, down:Infinity,
          left:Math.min(l,l*ratio), right:Math.max(r,r*ratio), top:(y-halfH)*ratio, bottom:y+halfH};
        cards.push(card);
        return card;
      });
    }
    // If silhouettes would meet along the same sightline, fit their height
    // and extrusion between the two timestamp rows. No record moves sideways.
    cards.sort((a,b) => a.top - b.top);
    for (let i = 0; i < cards.length; i++) {
      const a = cards[i];
      for (let j = i+1; j < cards.length && cards[j].top < a.bottom; j++) {
        const b = cards[j];
        if (a.left >= b.right || b.left >= a.right || a.y === b.y) continue;
        const near = a.y > b.y ? a : b, far = near === a ? b : a;
        const space = near.y - far.y;
        const needed = near.halfH + far.halfH;
        const share = near.front === far.front ? near.halfH/needed : near.front ? .86 : .14;
        const above = space >= needed ? space - far.halfH : clamp(space*share, space-far.halfH, near.halfH);
        near.up = Math.min(near.up, above);
        far.down = Math.min(far.down, space-above);
      }
    }
    let panRoom = 0;
    for (const card of cards) {
      const z = card.record.z;
      card.halfH = Math.min(card.halfH, card.up*.99, card.down*.99);
      if (card.up < card.y) {
        card.length = Math.min(card.length, z*(card.up-card.halfH)/(card.y-card.up));
      }
      // A grab can bring any recorded face to the viewport centre, including
      // a distant lane that was already visible at the edge.
      panRoom = Math.max(panRoom, Math.abs(card.x)/z);
    }
    return panRoom;
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
    // Leave space for the forward stack when the viewport can fit it without
    // compressing the recorded faces into unreadable strips.
    const span = Math.min(Math.max(nearHeight * 3.8, viewportHeight * .50),
      Math.max(nearHeight * 2.55, (viewportHeight - horizon - nearHeight / 2 - 84) / .98));
    const streetY = horizon + span * .98 + nearHeight / 2;
    const focal = viewportHeight / (2 * Math.tan(54 * Math.PI / 360));
    const blockDepth = .115;
    const ground = {vx: 0, horizon, streetY, street, focal, front: horizon + span * 1.35,
      pitch: faceWidth * .74, blockDepth, paneHeightRatio: .36, paneWidthRatio: .88,
      paneDepth: blockDepth * .37, eventElevation:nearHeight / 2, worldWidth:nearWidth, worldHeight:nearHeight};
    time.yAt = index => project(ground, 0, ground.eventElevation, time.depthAt(index)).y;
    time.groundYAt = (index, offset = 0) => {
      const depth = time.depthAt(index) + offset;
      return depth > 0 ? project(ground, 0, 0, depth).y : Infinity;
    };
    const plan = base?.plan || streetPlan(T, members, ground, gap, viewportWidth);
    ground.centerLane = plan.lanes.has(options.viewLane) ? options.viewLane : plan.center;
    ground.cameraX = (plan.lanes.get(ground.centerLane)?.origin || 0) * plan.pitch;
    const margin = Math.ceil(measureCards(members, plan, ground, gap));
    const width = viewportWidth + 2*margin;
    ground.vx = width / 2;
    let height = viewportHeight;
    for (const member of members) {
      const {lane, priority, attention, live} = member;
      for (const [index, slot] of member.slots.entries()) {
        const {state, index:timeIndex} = slot.record;
        const z = slot.record.z;
        const point = {...project(ground, slot.x, ground.eventElevation, z), worldX:slot.x, blockDepth:slot.length,
          scale:1/z, depth:z, timeIndex};
        if (!index) {
          heads.push({...point, id:state.id, lane, priority, attention, live,
            width:slot.halfW*2, height:slot.halfH*2});
          height = Math.max(height, Math.ceil(point.y + slot.halfH + 18));
        } else {
          point.paneWidth = slot.halfW*2;
          point.paneHeight = slot.halfH*2;
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
    return {points, heads, width, height, edges, ground, time, plan, laneOrder:plan.order,
      travelZ:(time.presentX - time.referenceX) / unitsPerZ};
  }

  function paneSize(ground, head, point) {
    const scale = head.depth / point.depth;
    return {width:point.paneWidth ?? head.width * ground.paneWidthRatio * scale,
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
