// Shared by the timeline, the selected node and the end-on view.
const WorldBrane = (() => {
  // A lane's label: its name and, for a row whose name is not already "row N", the row.
  const rowLabel = (name, row) => row && name !== 'row ' + row ? `${name} · row ${row}` : String(name || '');
  // ---------- Glyphs (one shape per brane; colour from the data) ----------
  const CLIP = { square: 'none', triangle: 'polygon(50% 4%, 100% 96%, 0 96%)', diamond: 'polygon(50% 0, 100% 50%, 50% 100%, 0 50%)',
    hexagon: 'polygon(25% 4%, 75% 4%, 100% 50%, 75% 96%, 25% 96%, 0 50%)', pentagon: 'polygon(50% 0, 100% 38%, 81% 100%, 19% 100%, 0 38%)',
    star: 'polygon(50% 0, 62% 35%, 99% 36%, 69% 58%, 80% 95%, 50% 73%, 20% 95%, 31% 58%, 1% 36%, 38% 35%)' };
  const SHARED_INK = '#526677';
  function gl(shape, colour, size, extra = '') {
    if (!shape) return `<span class="g" style="width:${size - 2}px;height:${size - 2}px;border-radius:50%;background:#8496a5;margin:1px;${extra}"></span>`;
    if (shape === 'circle') return `<span class="g" style="width:${size}px;height:${size}px;border-radius:50%;border:2px solid ${colour || SHARED_INK};${extra}"></span>`;
    return `<span class="g" style="width:${size}px;height:${size}px;background:${colour};clip-path:${CLIP[shape] || 'none'};border-radius:${shape === 'square' ? '2px' : '0'};${extra}"></span>`;
  }

  const timelineHref = (brane, at) => '/surface#' + new URLSearchParams(at ? {brane, at} : {brane});
  function graphEdges(model) {
    const lanes = new Map((model.lanes || []).map(lane => [lane.key, lane]));
    const ids = new Set([...lanes.values()].flatMap(lane => (lane.states || []).map(state => state.id)));
    const edges = (model.links || []).filter(({from, to}) => ids.has(from) && ids.has(to));
    const repos = Array.isArray(model.repos) ? model.repos.filter(repo => repo.shared && lanes.has(repo.shared)) : [];
    const shared = repos.length ? repos.map(repo => repo.shared) : ['shared'];
    for (const key of shared) {
      const line = (lanes.get(key)?.states || []).filter(state => state.kind === 'shared').sort((a, b) => a.rank - b.rank);
      for (let i = 1; i < line.length; i++) edges.push({from: line[i - 1].id, to: line[i].id});
    }
    return edges;
  }

  return {glyph: gl, label: rowLabel, sharedInk: SHARED_INK, timelineHref, graphEdges};
})();

if (typeof module !== 'undefined') module.exports = WorldBrane;
