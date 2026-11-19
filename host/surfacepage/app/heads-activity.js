// The columns are UTC calendar days, ending on the snapshot's day.
(() => {
  'use strict';
  const DAY = 86400000;

  function trace(states, now) {
    const days = Array(30).fill(0);
    const firstDay = Math.floor(now / DAY) - days.length + 1;
    let last = null;
    for (const state of states) {
      const time = Date.parse(state.at);
      if (!Number.isFinite(time) || time > now) continue;
      if (last === null || time > last) last = time;
      const column = Math.floor(time / DAY) - firstDay;
      if (column >= 0) days[column]++;
    }
    return {days, last};
  }

  if (typeof window !== 'undefined') window.WorldHeadActivity = {trace};
  if (typeof module !== 'undefined') module.exports = {trace};
})();
