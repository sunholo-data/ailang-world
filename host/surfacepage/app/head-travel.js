// Move the historical cut through the retained end-on scene.
(() => {
  'use strict';

  function laneOrder(T, lane) {
    return new Map((T.LANES?.[lane]?.states || []).map(({id}, index) => [id, index]));
  }

  function project(T, base, lane, id) {
    const stop = T?.N?.[id];
    if (!base || stop?.lane !== lane || !Number.isFinite(stop.rank) ||
        stop.draft || stop.kind === 'draft' || !base.points.has(id)) return base;

    const chosenOrder = laneOrder(T, lane);
    const cutoff = chosenOrder.get(id) ?? Infinity;
    const strips = new Map(base.heads.map(head => [head.lane, {head, slots: [], visible: []}]));
    for (const [key, point] of base.points) {
      const record = T.N[key];
      const strip = strips.get(record?.lane);
      if (!strip) continue;
      strip.slots.push({key, point});
      if (!Number.isFinite(record.rank) || record.rank > stop.rank || record.draft || record.kind === 'draft') continue;
      if (record.lane === lane && record.rank === stop.rank && (chosenOrder.get(key) ?? Infinity) > cutoff) continue;
      strip.visible.push(record);
    }
    if (!strips.get(lane)?.visible.some(record => record.id === id)) return base;

    const heads = [], points = new Map();
    for (const [key, strip] of strips) {
      if (!strip.visible.length) continue;
      const order = key === lane ? chosenOrder : laneOrder(T, key);
      strip.visible.sort((a, b) => b.rank - a.rank || Number(b.id === id) - Number(a.id === id) ||
        (order.get(b.id) ?? 0) - (order.get(a.id) ?? 0));
      strip.slots.sort((a, b) => Number(b.key === strip.head.id) - Number(a.key === strip.head.id) ||
        a.point.depth - b.point.depth);
      strip.face = {...strip.head, id: key === lane ? id : strip.visible[0].id, live: false, attention: false};
      if (key === lane) heads.unshift(strip.face);
      else heads.push(strip.face);
    }

    strips.get(lane).face.priority = 4;
    for (const head of heads) {
      const strip = strips.get(head.lane);
      for (const [index, record] of strip.visible.entries()) {
        // Retain the lane and move each eligible record forward into the
        // corresponding camera slot, including the previous-state lip.
        points.set(record.id, strip.slots[index].point);
      }
    }

    const edges = base.edges.filter(({from, to}) => points.has(from) && points.has(to));
    // The recorded stop's original camera depth drives the ground travel.
    const travelZ = base.points.get(id).depth - points.get(id).depth;
    return {...base, points, heads, edges, travelZ};
  }

  const api = {project};
  if (typeof window !== 'undefined') window.WorldHeadTravel = api;
  if (typeof module !== 'undefined') module.exports = api;
})();
