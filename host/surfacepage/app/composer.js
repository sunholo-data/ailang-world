// The dialog drafts a delegation against an explicitly published preset snapshot.
// Native dialog focus handling leaves the underlying state stack and pan intact.
(() => {
  if (typeof document === 'undefined') return;
  const byId = id => document.getElementById(id);
  const dialog = byId('episode-draft');
  if (!dialog?.dataset.catalog) return;
  const catalog = JSON.parse(dialog.dataset.catalog);
  if (!catalog.length) return;

  const form = byId('draft-form'), review = byId('draft-review'), title = byId('draft-title');
  const model = byId('draft-model'), ttl = byId('draft-ttl'), task = byId('draft-task');
  const hosts = {codex: 'Codex', claude: 'Claude Code'};
  let preset = catalog[0], host = preset.host;
  const text = (tag, value) => {
    const item = document.createElement(tag);
    item.textContent = value;
    return item;
  };
  const repoName = p => p.repo.split('/').filter(Boolean).pop();
  function authority(p) {
    const effects = new Set(p.grants.map(g => g.Effect));
    const actions = [];
    for (const [effect, name] of [['Workspace.Read','Read files'],['Workspace.Write','edit files'],['Workspace.Delete','delete files'],['Workspace.Exec','run profiled commands'],['Ailang.RunNet','run AILANG with network access']]) {
      if (effects.has(effect)) actions.push(name);
    }
    return actions.join(', ') || 'Use the grants listed below';
  }
  function consent(p) {
    const section = document.createElement('div');
    const effects = new Set(p.grants.map(g => g.Effect));
    const permissions = [
      ['Files', [ ['Workspace.Read','read'], ['Workspace.Write','write'], ['Workspace.Delete','delete'] ]
        .filter(([effect]) => effects.has(effect)).map(([,verb]) => verb).join(', ') || 'No file grant'],
      ['Execution', effects.has('Workspace.Exec') ? 'Commands in the '+p.profile+' profile' : 'No profiled commands'],
      ['Network', effects.has('Ailang.RunNet') ? 'Network-enabled AILANG runs' : 'No Ailang.RunNet grant']
    ];
    const table = text('dl', ''); table.className='draft-permissions';
    for (const [label,value] of permissions) table.append(text('dt',label),text('dd',value));
    section.append(table);
    const checks = p.checks || [];
    const outside = checks.filter(c => c.unsandboxed);
    if (outside.length) {
      const warning = text('p','Separate terminal approval required for unsandboxed checks: ' + outside.map(c => c.name).join(', '));
      warning.className = 'draft-caution';
      section.append(warning);
    }
    const detail = document.createElement('details');
    detail.append(text('summary',p.grants.length+' grants · '+checks.length+' finish checks'));
    const list = document.createElement('ul');
    for (const g of p.grants) list.append(text('li',g.Effect+' · '+g.Scope+' · '+g.Budget+' calls'));
    for (const c of checks) list.append(text('li',c.name+' · '+c.argv.join(' ')+' · '+c.timeout_s+' s · '+c.sandbox));
    detail.append(list);
    section.append(detail);
    return section;
  }
  function radio(parent, group, label, hint, selected, change) {
    const row = document.createElement('label');
    const input = document.createElement('input');
    input.type='radio'; input.name=group; input.checked=selected;
    const body=text('span','');
    body.append(text('strong',label),text('small',hint));
    input.addEventListener('change',change);
    row.append(input, body);
    parent.append(row);
  }
  function showHosts() {
    const area=byId('draft-hosts');
    area.replaceChildren();
    for (const [key,label] of Object.entries(hosts)) {
      radio(area,'host',label,key===preset.host?'Preset host':'Host override',key===host,()=>{
        host=key;
        model.value=host===preset.host?preset.model:'';
      });
    }
  }
  function choosePreset(p) {
    preset=p; host=p.host;
    byId('draft-parent').textContent=repoName(p)+' / '+p.base+' @ '+p.commit.slice(0,8);
    ttl.value=String(p.seconds/60); model.value=p.model;
    byId('draft-authority').replaceChildren(consent(p));
    showHosts();
  }
  function showPresets(first) {
    const area=byId('draft-presets');
    area.replaceChildren();
    for (const p of catalog.filter(p=>p.repo===first.repo && p.base===first.base)) {
      radio(area,'preset',p.name,authority(p),p===first,()=>choosePreset(p));
    }
    choosePreset(first);
  }
  const locations=new Map();
  for (const p of catalog) {
    const key=JSON.stringify([p.repo,p.base]);
    if (!locations.has(key)) locations.set(key,p);
  }
  for (const p of locations.values()) {
    radio(byId('draft-repos'),'repo',repoName(p),p.base+' @ '+p.commit.slice(0,8),p===preset,()=>showPresets(p));
  }
  showPresets(preset);
  const open=byId('compose-btn');
  open.hidden=false;
  open.onclick=()=>dialog.showModal();
  byId('draft-close').onclick=()=>dialog.close();
  function edit() {
    form.hidden=false; review.hidden=true; title.textContent='New episode'; task.focus();
  }
  form.onsubmit=event=>{
    event.preventDefault();
    if (!task.value.trim()) { task.setCustomValidity('Enter a task for this episode.'); task.reportValidity(); return; }
    const seconds=Math.round(Number(ttl.value)*60);
    if (!Number.isSafeInteger(seconds)||seconds<1) { ttl.setCustomValidity('Choose a valid session lifetime.'); ttl.reportValidity(); return; }
    form.hidden=true; review.hidden=false; title.textContent='Review episode';
    const summary=text('dl','');
    for (const [label,value] of [['Repository',preset.repo],['Base',preset.base+' @ '+preset.commit],['Preset',preset.name],['Agent',hosts[host]+' · '+(model.value.trim()||'Host default')],['Lifetime',seconds+' seconds'],['Exec profile',preset.profile||'None']]) {
      summary.append(text('dt',label),text('dd',value));
    }
    const editButton=text('button','Edit draft'); editButton.type='button'; editButton.onclick=edit;
    review.replaceChildren(text('h2',task.value.trim()),summary,consent(preset),text('p','Draft only. Starting this episode requires the attended World launcher.'),editButton);
    title.focus(); byId('draft-review').parentElement.scrollTop=0;
  };
  task.oninput=()=>task.setCustomValidity('');
  ttl.oninput=()=>ttl.setCustomValidity('');
})();
