"""Run load-bearing mutants only in the isolated candidate copy; restore byte-exactly."""
from pathlib import Path
import json, os, subprocess, hashlib
ROOT = Path('/private/tmp/ailang-world-row174-mutations-20261010')
BANK = Path(__file__).resolve().parent
ENV = dict(os.environ, AILANG_BIN='/Users/voightkampff/.pinned-ailang/ailang', WORLD_PKG_AILANG_BIN='/Users/voightkampff/.pinned-ailang/ailang', WORLD_TOOL_AILANG_BIN='/Users/voightkampff/.pinned-ailang-tools/v0.52.1/ailang')
MANIFEST = 'packages/se-tools/transitions.json'
pristine = (ROOT / MANIFEST).read_bytes()
base = json.loads((ROOT / 'host/coordinator/testdata/se-tools-descriptor-baseline.json').read_bytes())
records = []
def run(label, file, mutate, package, test, assertion):
    p = ROOT / file
    original = p.read_bytes()
    try:
        p.write_bytes(mutate(original))
        result = subprocess.run(['go', 'test', package, '-run', '^'+test+'$', '-v', '-count=1'], cwd=ROOT, env=ENV, capture_output=True, text=True)
        (BANK / ('mutation-'+label+'.log')).write_text(result.stdout+result.stderr)
        assert result.returncode != 0 and assertion in result.stdout, (label, result.returncode, result.stdout, result.stderr)
        records.append(dict(mutant=label, returncode=result.returncode, assertion=assertion))
    finally:
        p.write_bytes(original)
        assert p.read_bytes() == original

def entry_mutation(fn):
    def mutate(raw):
        entries=json.loads(raw); fn(entries)
        return (json.dumps(entries, indent=2)+'\n').encode()
    return mutate

def manifest(label, fn, assertion):
    run(label, MANIFEST, entry_mutation(fn), './host/coordinator', 'TestSeToolsDescriptorCompaction', assertion)
manifest('delete-read', lambda x:x.pop(0), 'manifest/floor census')
manifest('duplicate-read', lambda x:x.append(x[0].copy()), 'manifest/floor census')
manifest('tenth-entry', lambda x:x.append(dict(x[0],id='unexpected-tool')), 'manifest/floor census')
manifest('net-effect', lambda x:x[0]['declaredEffects'].append(dict(effect='Net',scope='worktree',cost=1)), 'complete non-description equivalence: ailang-read')
manifest('additional-properties', lambda x:x[0]['inputSchema'].update(additionalProperties=True), 'complete non-description equivalence: ailang-read')
manifest('remove-output-world', lambda x:x[0]['outputSchema']['properties'].pop('world'), 'complete non-description equivalence: ailang-read')
manifest('transition-file', lambda x:x[0].update(transitionFnFile='changed.ail'), 'complete non-description equivalence: ailang-read')
manifest('workspace-description', lambda x:x[-1].update(description='changed'), 'complete non-description equivalence: workspace-exec')
manifest('unknown-field', lambda x:x[0].update(unknown='new'), 'complete non-description equivalence: ailang-read')
manifest('old-descriptions', lambda x:[e.update(description=b['description']) for e,b in zip(x,base) if e['id']!='workspace-exec'], 'candidate byte budget: 13947/7430 exceeds 10500/3800')
run('source-byte', 'packages/se-tools/se_tools/read.ail', lambda x:x+b'\n', './host/coordinator', 'TestSeToolsDescriptorSourceEquality', 'parent transition source-byte equivalence')
run('projection-old-description', 'host/projection/mcpname.go', lambda x:x.replace(b'Description: d.Description,', ('Description: '+json.dumps(base[0]['description'])+',').encode()), './host/daemon', 'TestSeToolsDescriptorLiveExactEight', 'live descriptor equality (description/name/input/output)')
run('disable-interpreter-verification', 'host/capsule/capsule.go', lambda x:x.replace(b'if err := verifyExecutable(execPath, entry.Interpreter); err != nil {', b'if err := error(nil); err != nil {'), './host/capsule', 'TestF1PinnedInterpreterHashMismatchRefusedBeforeExec', 'want *HashMismatchError')
run('disable-tool-verification', 'host/daemon/workspace.go', lambda x:x.replace(b'if err := verifyArchivedTool(v.path, v.ref); err != nil {', b'if err := error(nil); err != nil {'), './host/daemon', 'TestWorkspaceToolBinaryIsArchivedAndVerified', 'tampered archived tool ran')
assert (ROOT / MANIFEST).read_bytes() == pristine
result = subprocess.run(['go','test','./host/coordinator','./host/daemon','./host/capsule','-run','TestSeToolsDescriptor|TestF1PinnedInterpreterHashMismatchRefusedBeforeExec|TestWorkspaceToolBinaryIsArchivedAndVerified','-v','-count=1'],cwd=ROOT,env=ENV,capture_output=True,text=True)
(BANK / 'mutation-restored-controls.log').write_text(result.stdout+result.stderr)
assert result.returncode == 0 and '--- SKIP:' not in result.stdout, result.stdout
(BANK / 'descriptor-mutations.json').write_text(json.dumps(dict(isolated_copy=str(ROOT),mutations=records,byte_exact_restore=True,restored_controls='PASS'),indent=2)+'\n')
print('PASS:',len(records),'intended mutation failures; byte-exact restores; final controls PASS')
