# Iteration 309 — controller verification (outside any sandbox), row 150 M7c2b

Tree: `.wt-world-iter309-row150-m7c2b`, base `a96820c5`; executor commit `5e20341f`; controller test-only commit `76216725`. `AILANG_BIN` = `~/.pinned-ailang/ailang` (AILANG v0.41.0).

## Pristine-base gate baseline (a96820c5)
```
== start 2026-10-27T20:03:35Z pwd=/Users/voightkampff/dev/sunholo-data/.wt-world-iter309-row150-m7c2b head=a96820c58c678b262371bedbb8023e7b6eda6864 load=5.96 5.17 4.30
gofmt rc=0 files=[]
vet rc=0
pyparse rc=0
bash-n rc=0
version-grep shell=0 control-doc=23
compile-fence rc=0 pkgs=34
daemon-cleanup rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	6.789s
pkg rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	25.249s
race rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	84.138s 
BASELINE DONE for a96820c5 2026-10-27T20:05:57Z
```

## Executor tree (5e20341f), executor snapshot cmp-identical to the tree before commit
```
4f6e10674ff972159d82e4d82e47e6146f649272bd8550b6927225171b52daee  scripts/capture_marketing.sh
65f297fecaca61efb3846cb0f103e2a0a629e39fa34db1742fd552cf47c4bc68  scripts/marketing_capture.py
162eb12b03a561a79c1c5d1a8755f7cfc2883cf7d7e89cea037d2358c586ef56  scripts/marketingcheck/capture_test.go
== start 2026-10-27T20:18:02Z pwd=/Users/voightkampff/dev/sunholo-data/.wt-world-iter309-row150-m7c2b head=5e20341f7a0dbe3df4a6de8744bec9133c934a59 load=4.76 4.86 5.65
gofmt rc=0 files=[]
vet rc=0
pyparse rc=0
bash-n rc=0
version-grep shell=0 control-doc=23
compile-fence rc=0 pkgs=34
daemon-cleanup rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	7.004s
pkg rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	30.249s
race rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	93.172s 
BASELINE DONE for 5e20341f 2026-10-27T20:20:14Z
```

## Red-first: head test file against base production files (shell + Python from a96820c5)
```
--- FAIL: TestMarketingCaptureFakeChromeSixShots (0.46s)
    --- FAIL: TestMarketingCaptureFakeChromeSixShots/six-shots (0.01s)
    --- FAIL: TestMarketingCaptureFakeChromeSixShots/serve-argv (0.01s)
    --- FAIL: TestMarketingCaptureFakeChromeSixShots/version-changed (0.01s)
    --- FAIL: TestMarketingCaptureFakeChromeSixShots/old-chrome (0.01s)
    --- FAIL: TestMarketingCaptureFakeChromeSixShots/cleanup (0.01s)
    --- FAIL: TestMarketingCaptureFakeChromeSixShots/scheme (0.00s)
    --- FAIL: TestMarketingCaptureFakeChromeSixShots/bind-refused (0.01s)
    --- PASS: TestMarketingCaptureFakeChromeSixShots/daemon-readiness-config (0.13s)
```
`daemon-readiness-config` passes at base (base already required `--ready-file` and rejected `--ready-url` with usage rc 2), so the controller added two mutants for it below.

## Controller mutation drill on 5e20341f (one mutant at a time, named subtest only, restored and byte-compared)
Script: `/tmp/world-iter309/mutants.py` (anchors asserted unique). The X-* rows are diff-anchored controller mutants.
```
M7c2b-version-changed: subtest=version-changed rc=1 KILLED restored=True
     capture_test.go:349: a Chrome version change during capture must refuse: rc=0 stderr="progress: 1/6 workbench-light.png ok\nprogress: 2/6 workbench-dark.png ok\nprogress: 3/6 packet.png ok\nprogress: 4/6 graph.png ok\nprogress: 5/6 provenance.png ok\nprogress: 6/6 phone.png ok\nprogress: 1/6 workben
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (1.69s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/version-changed (1.39s)
M7c-scheme: subtest=scheme rc=1 KILLED restored=True
     capture_test.go:385: every shot explicitly pins correct scheme: pass-a/workbench-dark.png argv=["--headless=new" "--no-first-run" "--no-default-browser-check" "--disable-background-networking" "--user-data-dir=/private/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/marketing-chrome-x68xnzm4" "--scr
     capture_test.go:385: every shot explicitly pins correct scheme: pass-b/workbench-dark.png argv=["--headless=new" "--no-first-run" "--no-default-browser-check" "--disable-background-networking" "--user-data-dir=/private/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/marketing-chrome-p_cvasx1" "--scr
     capture_test.go:397: capture must succeed: rc=1 stderr="progress: 1/6 workbench-light.png ok\nprogress: 2/6 workbench-dark.png ok\nprogress: 3/6 packet.png ok\nprogress: 4/6 graph.png ok\nprogress: 5/6 provenance.png ok\nprogress: 6/6 phone.png ok\nprogress: 1/6 workbench-light.png ok\nprogress: 2/6
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (1.97s)
M7c2b-identities: subtest=serve-argv rc=1 KILLED restored=True
     capture_test.go:365: both seed identities must be configured: ["serve" "--db" "/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestMarketingCaptureFakeChromeSixShotsserve-argv313131645/001/db" "--bind" "127.0.0.1:57275" "--workbench-still" "--ailang-bin" "/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt0000
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (1.98s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/serve-argv (1.59s)
A2-drop-loopback: subtest=bind-refused rc=1 KILLED restored=True
     capture_test.go:335: non-loopback bind must refuse before daemon starts: rc=0 stderr="progress: 1/6 workbench-light.png ok\nprogress: 2/6 workbench-dark.png ok\nprogress: 3/6 packet.png ok\nprogress: 4/6 graph.png ok\nprogress: 5/6 provenance.png ok\nprogress: 6/6 phone.png ok\nprogress: 1/6 workben
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (1.97s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/bind-refused (1.57s)
X-drop-diagnostics: subtest=six-shots rc=1 KILLED restored=True
     capture_test.go:417: open /var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestMarketingCaptureFakeChromeSixShotssix-shots3604223237/001/out/diagnostics.txt: no such file or directory
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (2.05s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/six-shots (1.65s)
X-drop-check: subtest=six-shots rc=1 KILLED restored=True
     capture_test.go:415: open /var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestMarketingCaptureFakeChromeSixShotssix-shots1628925710/001/out/capture-manifest.json: no such file or directory
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (2.00s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/six-shots (1.58s)
X-single-pass-continues: subtest=ALL rc=0 SURVIVED restored=True
X-lifetime-default: subtest=ALL rc=0 SURVIVED restored=True
X-ready-ignores-exit: subtest=ALL rc=0 SURVIVED restored=True
X-started-despite-stop: subtest=ALL rc=0 SURVIVED restored=True
X-old-chrome-continues: subtest=old-chrome rc=1 KILLED restored=True
     capture_test.go:343: old Chrome must refuse before any shot: rc=1 stderr=""
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (0.67s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/old-chrome (0.29s)
X-cleanup-no-term: subtest=cleanup rc=1 KILLED restored=True
     capture_test.go:394: fake worldd after capture must be dead after bounded cleanup (pid 69841)
     capture_test.go:397: capture must succeed: rc=-1 stderr="progress: 1/6 workbench-light.png ok\nprogress: 2/6 workbench-dark.png ok\nprogress: 3/6 packet.png ok\nprogress: 4/6 graph.png ok\nprogress: 5/6 provenance.png ok\nprogress: 6/6 phone.png ok\nprogress: 1/6 workbench-light.png ok\nprogress: 2/
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (125.43s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/cleanup (125.02s)
X-readiness-group-not-required: subtest=daemon-readiness-config rc=1 KILLED restored=True
     capture_test.go:430: exactly one readiness source required: rc=1 stderr=""
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (0.41s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/daemon-readiness-config (0.12s)
X-readiness-group-not-exclusive: subtest=daemon-readiness-config rc=1 KILLED restored=True
     capture_test.go:430: exactly one readiness source required: rc=1 stderr=""
     capture_test.go:430: exactly one readiness source required: rc=1 stderr=""
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (0.44s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/daemon-readiness-config (0.14s)
X-ready-url-ignored: subtest=six-shots rc=1 KILLED restored=True
     capture_test.go:397: capture must succeed: rc=1 stderr="Traceback (most recent call last):\n  File \"/Users/voightkampff/dev/sunholo-data/.wt-world-iter309-row150-m7c2b/scripts/marketing_capture.py\", line 367, in <module>\n    raise SystemExit(main())\n                     ~~~~^^\n  File \"/Users/v
     --- FAIL: TestMarketingCaptureFakeChromeSixShots (0.72s)
     --- FAIL: TestMarketingCaptureFakeChromeSixShots/six-shots (0.37s)
final restore True
```

## Survivor closure (test-only, 76216725): single-pass and daemon-exits
```
X-single-pass-continues: subtest=single-pass rc=1 KILLED — capture_test.go:359: without --verify-twice the run stops after pass A with rc 0 and no check: rc=0 pass-a=6 unexpected=[pass-b capture-manifest.json diagnostics.txt]
X-ready-ignores-exit: subtest=daemon-exits rc=1 KILLED — capture_test.go:345: a daemon that exits before HTTP readiness must refuse as NonzeroExit before any pass: verdict=Interrupted (29.5 s)
final restore True; shasum -c ctl-final.sha256: 3/3 OK
```
Still surviving (carried, non-blocking): X-lifetime-default (A1 untested), X-started-despite-stop (signal before HTTP readiness).

## Final tree (76216725)
```
== start 2026-10-27T20:21:43Z pwd=/Users/voightkampff/dev/sunholo-data/.wt-world-iter309-row150-m7c2b head=76216725417303b304228d8d43c77c37ca11fdb0 load=3.04 4.02 5.12
gofmt rc=0 files=[]
vet rc=0
pyparse rc=0
bash-n rc=0
version-grep shell=0 control-doc=23
compile-fence rc=0 pkgs=34
daemon-cleanup rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	6.118s
pkg rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	31.522s
race rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	102.017s 
BASELINE DONE for 76216725 2026-10-27T20:24:05Z
TestSubprocessSinkGate (host/verifygate) ok; TestProductionGoSurface + TestProductionContextRoots (host/store) ok
```

## Outside-worktree writes
Read the diff for `>`/`>>`/`tee`/`rm`/`mv`: every write in the new shell flow is under `--out` (`daemon-verdict.json`, `daemon-started.json`, `routes.json`, `pass-a`, `pass-b`, `diagnostics.txt`, `capture-manifest.json`); the test's `--out` is a `t.TempDir()`. Nothing under `$HOME` or `~/.ailang/state`.

## Judge findings reproduced by the controller
```
capture: missing --manifest
missing-flag rc=1
capture: unknown flag: --bogus
unknown-flag rc=1
  --bind 127.0.0.1:PORT      Loopback daemon address
```
