// A local draft starts at an exact recorded state and lives in this tab.
(() => {
  'use strict';
  const storageKey = 'world.head-drafts/v1';
  const draftLabel = 'Draft · this tab';
  const titleLimit = 160;
  const draftID = /^draft:[a-z0-9][a-z0-9-]{5,79}$/i;
  const taskTitle = value => typeof value === 'string' ? value.trim().replace(/\s+/g, ' ') : '';

  function restore(raw) {
    const drafts = new Map();
    if (!Array.isArray(raw)) return drafts;
    for (const item of raw) {
      const title = taskTitle(item?.title);
      if (typeof item?.id !== 'string' || !draftID.test(item.id) || drafts.has(item.id) ||
          typeof item.parent !== 'string' || !item.parent.trim() || draftID.test(item.parent) ||
          !title || title.length > titleLimit || typeof item.at !== 'string' || !Number.isFinite(Date.parse(item.at))) continue;
      drafts.set(item.id, {id: item.id, parent: item.parent, title, at: new Date(item.at).toISOString()});
    }
    return drafts;
  }

  function project(raw, drafts) {
    const parents = new Map();
    if (!Array.isArray(raw?.lanes)) return {model: raw, parents};
    const lanes = [], prior = new Set();
    let lastRank = -1;
    for (const lane of raw.lanes) {
      if (lane.draft === true && lane.key.startsWith('local:draft:')) {
        for (const state of lane.states || []) prior.add(state.id);
        continue;
      }
      lanes.push(lane);
      for (const state of lane.states || []) {
        if (state.draft !== true && state.kind !== 'draft') parents.set(state.id, {lane, state});
        if (Number.isFinite(state.rank)) lastRank = Math.max(lastRank, state.rank);
      }
    }
    const overlay = [], links = [];
    for (const draft of drafts.values()) {
      const anchor = parents.get(draft.parent);
      if (!anchor || parents.has(draft.id)) continue;
      const {lane} = anchor;
      const state = {...draft, kind: 'draft', draft: true, label: draftLabel, rank: ++lastRank, badge: 0};
      overlay.push({
        key: 'local:' + draft.id, name: draftLabel, title: draft.title,
        draft: true, repo: lane.repo, row: null, colour: lane.colour, glyph: lane.glyph,
        stepsLanded: 0, states: [state]
      });
      links.push({from: draft.parent, to: draft.id, colour: lane.colour});
    }
    if (!overlay.length && !prior.size) return {model: raw, parents};
    const recorded = (raw.links || []).filter(link => !prior.has(link.from) && !prior.has(link.to));
    return {model: {...raw, lanes: lanes.concat(overlay), links: recorded.concat(links)}, parents};
  }

  function create(onChange) {
    let drafts;
    try { drafts = restore(JSON.parse(sessionStorage.getItem(storageKey))); }
    catch (_) { drafts = new Map(); }
    let parents = new Map(), editor, sequence = 0;

    function nextID() {
      let id;
      do {
        const unique = globalThis.crypto?.randomUUID?.() ??
          [Date.now().toString(36), (++sequence).toString(36), Math.random().toString(36).slice(2)].join('-');
        id = 'draft:' + unique;
      } while (drafts.has(id) || parents.has(id));
      return id;
    }

    function changed(draft, discarded) {
      try { sessionStorage.setItem(storageKey, JSON.stringify([...drafts.values()])); }
      catch (_) { /* The draft remains available in this view. */ }
      onChange?.({id: draft.id, parent: draft.parent, discarded});
    }

    function createEditor() {
      const dialog = document.createElement('dialog');
      dialog.className = 'head-draft-dialog';
      dialog.setAttribute('aria-label', draftLabel);
      dialog.innerHTML = `
        <form>
          <h2>Draft · this tab</h2>
          <p>Branch from <strong></strong></p>
          <label>Task title<input name="task" type="text" required maxlength="${titleLimit}" autofocus></label>
          <div class="head-draft-actions">
            <button type="button" data-discard hidden>Discard draft</button>
            <div class="head-draft-controls">
              <button type="button" data-cancel>Cancel</button>
              <button type="submit">Create draft</button>
            </div>
          </div>
        </form>`;
      const form = dialog.querySelector('form');
      const title = form.elements.task;
      const origin = form.querySelector('strong');
      const discard = form.querySelector('[data-discard]');
      const submit = form.querySelector('[type="submit"]');
      let context;

      function close() {
        context = null;
        dialog.close();
      }
      form.querySelector('[data-cancel]').addEventListener('click', close);
      dialog.addEventListener('cancel', () => { context = null; });
      title.addEventListener('input', () => title.setCustomValidity(''));
      discard.addEventListener('click', () => {
        const draft = context?.draft;
        if (!draft) return;
        drafts.delete(draft.id);
        close();
        changed(draft, true);
      });
      form.addEventListener('submit', event => {
        event.preventDefault();
        const value = taskTitle(title.value);
        if (!value || value.length > titleLimit) {
          title.setCustomValidity(value ? 'Shorten the task title.' : 'Enter a task title.');
          title.reportValidity();
          return;
        }
        if (!context) return;
        const {parent, draft: existing} = context;
        const draft = existing ? {...existing, title: value}
          : {id: nextID(), parent: parent.id, title: value, at: new Date().toISOString()};
        drafts.set(draft.id, draft);
        close();
        changed(draft, false);
      });
      document.body.append(dialog);
      return (parent, draft) => {
        context = {parent, draft};
        origin.textContent = parent.title || parent.summary || parent.label || parent.id;
        title.setCustomValidity('');
        title.value = draft?.title || '';
        discard.hidden = !draft;
        submit.textContent = draft ? 'Save draft' : 'Create draft';
        if (!dialog.open) dialog.showModal();
        title.focus();
        if (draft) title.select();
      };
    }

    return {
      project(raw) {
        const current = project(raw, drafts);
        parents = current.parents;
        return current.model;
      },
      edit(parentNode, existingID) {
        const draft = existingID ? drafts.get(existingID) : null;
        if (existingID && !draft) return false;
        const parent = parents.get(draft?.parent ?? parentNode?.id)?.state;
        if (!parent) return false;
        editor ??= createEditor();
        editor(parent, draft);
        return true;
      },
      createDraft(parentID, rawTitle) {
        const value = taskTitle(rawTitle).slice(0, titleLimit);
        const parent = parents.get(parentID)?.state;
        if (!parent || !value) return null;
        const draft = {id: nextID(), parent: parent.id, title: value, at: new Date().toISOString()};
        drafts.set(draft.id, draft);
        changed(draft, false);
        return draft.id;
      },
      has(id) { return drafts.has(id); }
    };
  }

  if (typeof window !== 'undefined') window.WorldHeadDrafts = {create};
  if (typeof module !== 'undefined') module.exports = {project, restore};
})();
