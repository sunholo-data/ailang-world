"""Controller's unchanged normal gates; one attempt, finite cap, source identity on both sides."""
import datetime, hashlib, json, os, pathlib, signal, subprocess, time
BANK = pathlib.Path(__file__).resolve().parent
ROOT = BANK.parents[2]
tracked = subprocess.check_output(['git', 'ls-files'], cwd=ROOT, text=True).splitlines()
FILES = ['packages/se-tools/transitions.json', 'host/coordinator/testdata/se-tools-descriptor-baseline.json',
         'scripts/verify_ail.sh', 'scripts/verify_go.sh']
FILES += [p for p in tracked if p.endswith(('.go', '.ail'))]
assert FILES and all((ROOT / p).is_file() for p in FILES), 'Source census incomplete'
def snapshot():
    return {p: hashlib.sha256((ROOT / p).read_bytes()).hexdigest() for p in sorted(set(FILES))}
head = subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=ROOT, text=True).strip()
env = os.environ.copy()
env.update(AILANG_BIN='/Users/voightkampff/.pinned-ailang/ailang',
           WORLD_PKG_AILANG_BIN='/Users/voightkampff/.pinned-ailang/ailang',
           WORLD_TOOL_AILANG_BIN='/Users/voightkampff/.pinned-ailang-tools/v0.52.1/ailang')
report = {'product_commit': head, 'started_at_utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
          'start_sources': snapshot(), 'attempt': 1, 'gates': []}
for label, command in [('verify-ail', ['./scripts/verify_ail.sh']), ('verify-go', ['./scripts/verify_go.sh'])]:
    path = BANK / (label + '-normal.log')
    if path.exists():
        raise RuntimeError('Refusing to overwrite an earlier attempt: ' + str(path))
    print('START', label, 'head', head, flush=True)
    start = time.monotonic()
    with path.open('wb') as log:
        proc = subprocess.Popen(command, cwd=ROOT, env=env, stdout=log, stderr=subprocess.STDOUT, start_new_session=True)
        capped = False
        try:
            rc = proc.wait(timeout=1800)
        except subprocess.TimeoutExpired:
            capped = True
            os.killpg(proc.pid, signal.SIGTERM)
            try:
                rc = proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                os.killpg(proc.pid, signal.SIGKILL)
                rc = proc.wait(timeout=5)
    row = {'gate': label, 'command': command, 'returncode': rc, 'elapsed_s': round(time.monotonic()-start, 3),
           'outer_cap_s': 1800, 'capped': capped, 'log': path.name, 'log_sha256': hashlib.sha256(path.read_bytes()).hexdigest()}
    report['gates'].append(row)
    report['end_sources'] = snapshot()
    report['source_identity_unchanged'] = report['start_sources'] == report['end_sources']
    report['recorded_at_utc'] = datetime.datetime.now(datetime.timezone.utc).isoformat()
    (BANK / 'normal-gates-result.json').write_text(json.dumps(report, indent=2)+'\n')
    print(json.dumps(row), flush=True)
    if rc != 0 or capped or not report['source_identity_unchanged']:
        raise SystemExit(1)
print('BOTH NORMAL GATES PASS; source bytes unchanged; product review still required', flush=True)
