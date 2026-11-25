(function () {
  const presetsEl = document.getElementById('presets');
  const profilesEl = document.getElementById('profiles');
  const statusEl = document.getElementById('status');
  const refreshBtn = document.getElementById('refresh');

  function el(tag, cls, text) {
    const n = document.createElement(tag);
    if (cls) n.className = cls;
    if (text !== undefined) n.textContent = text;
    return n;
  }

  function badgeClass(status) {
    if (status === 'loaded' || status === 'active') return 'pill ok';
    if (status === 'excluded') return 'pill err';
    return 'pill warn';
  }

  function renderPresets(presets, profileByName) {
    presetsEl.textContent = '';
    if (!presets.length) {
      presetsEl.appendChild(el('div', 'status', 'No presets configured. Run ailang-worldd preset init ...'));
      return;
    }
    for (const p of presets) {
      const card = el('div', 'card');
      const top = el('div', 'row');
      top.appendChild(el('span', 'title', p.name));
      const pills = el('div', 'pills');
      if (p.agent) pills.appendChild(el('span', 'pill', p.agent));
      if (p.mode) pills.appendChild(el('span', 'pill', p.mode));
      const prof = p.profile ? profileByName.get(p.profile) : null;
      if (p.profile) {
        const st = prof ? prof.probe_status : 'missing';
        pills.appendChild(el('span', badgeClass(st), p.profile + ': ' + st));
      }
      top.appendChild(pills);
      card.appendChild(top);

      if (p.error) {
        card.appendChild(el('div', 'pill err', p.error));
      } else {
        card.appendChild(el('div', 'meta', (p.repo || '') + (p.base ? ' · base ' + p.base : '') + ' · ' + p.grants_count + ' grants'));
        const cmdText = 'ailang-worldd episode start <ep> --preset ' + p.name + ' --task <file> --launch';
        const cmdRow = el('div', 'row');
        const cmdBox = el('div', 'cmd', cmdText);
        cmdBox.style.flex = '1';
        const copyBtn = el('button', '', 'Copy');
        copyBtn.type = 'button';
        copyBtn.addEventListener('click', () => {
          if (navigator.clipboard) navigator.clipboard.writeText(cmdText);
          copyBtn.textContent = 'Copied';
          setTimeout(() => { copyBtn.textContent = 'Copy'; }, 1200);
        });
        cmdRow.appendChild(cmdBox);
        cmdRow.appendChild(copyBtn);
        card.appendChild(cmdRow);
      }
      presetsEl.appendChild(card);
    }
  }

  function renderProfiles(profiles) {
    profilesEl.textContent = '';
    if (!profiles.length) {
      profilesEl.appendChild(el('div', 'status', 'No exec profiles configured. Run ailang-worldd profile init ...'));
      return;
    }
    for (const pr of profiles) {
      const card = el('div', 'card');
      const top = el('div', 'row');
      top.appendChild(el('span', 'title', pr.project));
      const pills = el('div', 'pills');
      pills.appendChild(el('span', badgeClass(pr.state), pr.state));
      pills.appendChild(el('span', badgeClass(pr.probe_status), pr.probe_status));
      top.appendChild(pills);
      card.appendChild(top);

      if (pr.error) {
        card.appendChild(el('div', 'pill err', pr.error));
      } else {
        card.appendChild(el('div', 'meta', 'root: ' + (pr.root || '.') + (pr.digest ? ' · ' + pr.digest.slice(0, 19) : '')));
        if (pr.commands && pr.commands.length) {
          const cmds = el('div', 'pills');
          for (const c of pr.commands) cmds.appendChild(el('span', 'pill', c));
          card.appendChild(cmds);
        }
      }
      if (pr.exclude_reason) {
        card.appendChild(el('div', 'meta', 'Excluded: ' + pr.exclude_reason));
      }
      profilesEl.appendChild(card);
    }
  }

  async function load() {
    statusEl.textContent = 'Loading…';
    try {
      const [presRes, profRes] = await Promise.all([
        fetch('/surface/presets.json', { cache: 'no-store' }),
        fetch('/surface/profiles.json', { cache: 'no-store' }),
      ]);
      if (!presRes.ok || !profRes.ok) {
        throw new Error('HTTP ' + presRes.status + '/' + profRes.status);
      }
      const [presData, profData] = await Promise.all([presRes.json(), profRes.json()]);
      const profiles = profData.profiles || [];
      const presets = presData.presets || [];
      const byName = new Map();
      for (const pr of profiles) {
        if (pr.state === 'active') byName.set(pr.project, pr);
      }
      renderPresets(presets, byName);
      renderProfiles(profiles);
      statusEl.textContent = presets.length + ' presets · ' + profiles.length + ' profiles';
    } catch (err) {
      statusEl.textContent = 'Error: ' + err.message;
    }
  }

  refreshBtn.addEventListener('click', load);
  load();
})();
