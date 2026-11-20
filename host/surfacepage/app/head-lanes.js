// Choose a Heads lane from the same repository groups as Timeline.
(() => {
  'use strict';
  if (typeof window === 'undefined') return;

  function create(container, chooseLane) {
    const strip = document.createElement('div');
    strip.className = 'hc-lanes';
    strip.hidden = true;

    const select = document.createElement('select');
    select.setAttribute('aria-label', 'Thread');
    strip.append(select);
    container.append(strip);

    select.addEventListener('change', () => chooseLane(select.value));
    // Keep type-ahead, arrows and Escape inside the native picker.
    select.addEventListener('keydown', event => event.stopPropagation());
    strip.addEventListener('wheel', event => event.stopPropagation());

    let previous = '';
    function update({active, graph, current, heads, historical, order = []}) {
      strip.hidden = !active || !heads.length;
      if (strip.hidden) return;

      const byLane = new Map(heads.map(head => [head.lane, head]));
      const sidebar = graph.ROWS.flatMap(row => graph.QUIET[row] || [row]);
      const lanes = [...new Set([...sidebar, ...order, ...byLane.keys()])].filter(lane => byLane.has(lane));
      const shared = new Set(graph.GROUPS.map(group => group.shared));
      const groups = graph.GROUPS.map(group => ({
        id: group.id,
        options: [group.shared, ...lanes.filter(lane => !shared.has(lane) && graph.groupOf(lane) === group)]
          .filter(lane => byLane.has(lane)).map(lane => {
            const {live, attention} = byLane.get(lane);
            const label = [graph.laneTitle(lane), !historical && live && 'Running', !historical && attention && 'Needs you']
              .filter(Boolean).join(', ');
            return {lane, label};
          })
      })).filter(group => group.options.length);
      const signature = JSON.stringify(groups);
      if (signature !== previous) {
        const row = document.createDocumentFragment();
        for (const {id, options} of groups) {
          const target = graph.GROUPS.length > 1 ? document.createElement('optgroup') : row;
          if (target !== row) {
            target.label = id;
            row.append(target);
          }
          for (const {lane, label} of options) target.append(new Option(label, lane));
        }
        select.replaceChildren(row);
        previous = signature;
      }
      if (select.value !== current && byLane.has(current)) select.value = current;
    }

    return {update};
  }

  window.WorldHeadLanes = {create};
})();
