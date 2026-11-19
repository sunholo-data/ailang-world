// World branch heads use the same lane and state addresses as the timeline.
(function () {
  'use strict';

  function runningStatus(edge, repos, follow) {
    const messages = [];
    if (follow) {
      const open = edge.open;
      if (open && open.note) messages.push(open.note);
      else if (!open || open.status === 'unavailable') messages.push('Running status unavailable');
      else if (open.status !== 'complete') messages.push('Running status incomplete');
      for (const repo of repos) {
        if (repo.state !== 'fresh' && repo.line) messages.push(repo.line);
      }
    } else if (edge.sessionsUnavailable) {
      messages.push(edge.sessionsUnavailable);
    }
    return messages.join(' · ');
  }

  function project(model) {
    const follow = model.schema === 'world/surface-model/v2';
    const repos = model.repos || [];
    const sourceByID = new Map(repos.map(repo => [repo.id, repo]));
    const edge = model.nowEdge || {};
    const heads = [], unlinked = [];
    const byLane = new Map(), byEpisode = new Map();

    for (const lane of model.lanes || []) {
      if (!lane.states.length) continue;
      // Follow appends open states in overlay order, after all landed ranks.
      // The timeline uses that rank for its tip. History uses recorded times.
      const latest = lane.states.reduce((tip, state) => state.rank > tip.rank ? state : tip);
      const states = lane.states.slice().sort((a, b) => Date.parse(b.at) - Date.parse(a.at) || b.rank - a.rank);
      const head = {
        key: lane.key, name: lane.name, repo: lane.repo, colour: lane.colour, glyph: lane.glyph,
        states, latest, live: [], source: sourceByID.get(lane.repo),
      };
      heads.push(head);
      byLane.set(lane.key, {head, ids: new Set(states.map(state => state.id))});
      if (!follow) {
        for (const state of states) {
          if (!state.episode) continue;
          const matches = byEpisode.get(state.episode) || [];
          matches.push({head, state});
          byEpisode.set(state.episode, matches);
        }
      }
    }

    if (follow) {
      for (const session of (edge.open && edge.open.live) || []) {
        const lane = byLane.get(session.lane);
        if (lane && lane.ids.has(session.id)) lane.head.live.push(session);
        else unlinked.push({episode: session.episode, reason: session.reason || 'No indexed state for this session'});
      }
    } else if (!edge.sessionsUnavailable) {
      for (const session of edge.sessions || []) {
        const matches = byEpisode.get(session.episode) || [];
        if (matches.length === 1) {
          const {head, state} = matches[0];
          head.live.push({...session, id: state.id, lane: head.key, label: 'landed; session live'});
        } else {
          unlinked.push({
            episode: session.episode,
            reason: matches.length ? 'Multiple indexed states match this episode' : 'No indexed state for this session',
          });
        }
      }
    }

    return {heads, unlinked, status: runningStatus(edge, repos, follow)};
  }

  if (typeof window !== 'undefined') window.WorldHeadsData = {project};
  if (typeof module !== 'undefined') module.exports = {project};
}());
