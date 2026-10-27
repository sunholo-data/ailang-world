# Iteration 308 — controller verification (outside any sandbox), row 150 M7a-fix

## Pristine-base baseline (.base-world-iter308 @ 0363b3ae, AILANG_BIN pinned v0.41.0)
```
== start 2026-10-27T16:04:06Z pwd=/Users/voightkampff/dev/sunholo-data/.base-world-iter308
gofmt rc=0
vet rc=0
pycompile rc=0
pkg rc=0 ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	24.589s
race1 rc=0 
race2 rc=0 
race3 rc=1 --- FAIL: TestMarketingSupervisorFakeDaemonCleanup (12.47s)         supervisor_test.go:464: normal shell cleanup: exit status 1 {"verdict": "SupervisorError", "detail": "[Errno 1] Operation not permitted"} 
BASELINE DONE for 0363b3ae 2026-10-27T16:07:46Z
```

## Opt-in stress row, fixed tree @ 31211dbf (25 s Python wall, as executed)
```
run 1: supervisor_test.go:229: stress summary: {"trials": 159, "wrapper_calls": 477, "raw_calls": 477, "raw_eperm": 6, "raw_esrch": 471, "raw_success": 0, "wrapper_exceptions": [], "errors": ["TimeoutError('25-second stress wall deadline exceeded')", "25-second stress wall deadline exceeded"], "wrapper_exception_count": 0, "wall_seconds": 25.075696209001762}; Python exit: exit status 1
FAIL	github.com/sunholo-data/ailang-world/scripts/marketingcheck	25.341s
run 2: supervisor_test.go:229: stress summary: {"trials": 160, "wrapper_calls": 480, "raw_calls": 480, "raw_eperm": 0, "raw_esrch": 480, "raw_success": 0, "wrapper_exceptions": [], "errors": ["TimeoutError('25-second stress wall deadline exceeded')", "25-second stress wall deadline exceeded"], "wrapper_exception_count": 0, "wall_seconds": 25.144100833000266}; Python exit: exit status 1
FAIL	github.com/sunholo-data/ailang-world/scripts/marketingcheck	25.346s
run 3: supervisor_test.go:229: stress summary: {"trials": 155, "wrapper_calls": 465, "raw_calls": 465, "raw_eperm": 4, "raw_esrch": 461, "raw_success": 0, "wrapper_exceptions": [], "errors": ["TimeoutError('25-second stress wall deadline exceeded')", "25-second stress wall deadline exceeded"], "wrapper_exception_count": 0, "wall_seconds": 25.024889500004065}; Python exit: exit status 1
FAIL	github.com/sunholo-data/ailang-world/scripts/marketingcheck	25.241s
```

## Controller mutant M7a-fix-eperm (except ProcessLookupError only) — unit row and stress row @ 25 s bound
```
stress mutant run 1: supervisor_test.go:229: stress summary: {"trials": 147, "wrapper_calls": 441, "raw_calls": 441, "raw_eperm": 6, "raw_esrch": 435, "raw_success": 0, "wrapper_exceptions": [{"trial": 83, "pair": 0, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 83, "pair": 1, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 83, "pair": 2, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 88, "pair": 0, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 88, "pair": 1, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 88, "pair": 2, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}], "errors": ["TimeoutError('25-second stress wall deadline exceeded')", "25-second stress wall deadline exceeded"], "wrapper_exception_count": 6, "wall_seconds": 25.07508266599325}; Python exit: exit status 1
supervisor_test.go:252: real kill_group must raise nothing after leader reap: {"trials": 147, "wrapper_calls": 441, "raw
stress mutant run 2: supervisor_test.go:229: stress summary: {"trials": 149, "wrapper_calls": 447, "raw_calls": 447, "raw_eperm": 7, "raw_esrch": 440, "raw_success": 0, "wrapper_exceptions": [{"trial": 36, "pair": 0, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 36, "pair": 1, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 36, "pair": 2, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 49, "pair": 0, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 49, "pair": 1, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}, {"trial": 49, "pair": 2, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}], "errors": ["TimeoutError('25-second stress wall deadline exceeded')", "25-second stress wall deadline exceeded"], "wrapper_exception_count": 6, "wall_seconds": 25.127080999998725}; Python exit: exit status 1
supervisor_test.go:252: real kill_group must raise nothing after leader reap: {"trials": 149, "wrapper_calls": 447, "raw
```

## Opt-in stress row after controller bound fix 7475ec3c (60 s Python wall, 70 s Go context)
```
fixed run 1: supervisor_test.go:229: stress summary: {"trials": 200, "wrapper_calls": 600, "raw_calls": 600, "raw_eperm": 7, "raw_esrch": 593, "raw_success": 0, "wrapper_exceptions": [], "errors": [], "wrapper_exception_count": 0, "wall_seconds": 33.28274495799997}; Python exit: <nil>
ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	33.783s
fixed run 2: supervisor_test.go:229: stress summary: {"trials": 200, "wrapper_calls": 600, "raw_calls": 600, "raw_eperm": 3, "raw_esrch": 597, "raw_success": 0, "wrapper_exceptions": [], "errors": [], "wrapper_exception_count": 0, "wall_seconds": 32.64467537499877}; Python exit: <nil>
ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	33.028s
fixed run 3: supervisor_test.go:229: stress summary: {"trials": 200, "wrapper_calls": 600, "raw_calls": 600, "raw_eperm": 7, "raw_esrch": 593, "raw_success": 0, "wrapper_exceptions": [], "errors": [], "wrapper_exception_count": 0, "wall_seconds": 33.33403408400045}; Python exit: <nil>
ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	33.538s
mutant run 1: supervisor_test.go:229: stress summary: {"trials": 200, "wrapper_calls": 600, "raw_calls": 600, "raw_eperm": 0, "raw_esrch": 600, "raw_success": 0, "wrapper_exceptions": [], "errors": [], "wrapper_exception_count": 0, "wall_seconds": 31.47881333300029}; Python exit: <nil>
supervisor_test.go:231: UNINFORMATIVE: raw EPERM count is 0; EPERM not provoked on this host; this run is not evidence the fix works
mutant run 2: supervisor_test.go:229: stress summary: {"trials": 200, "wrapper_calls": 600, "raw_calls": 600, "raw_eperm": 1, "raw_esrch": 599, "raw_success": 0, "wrapper_exceptions": [{"trial": 128, "pair": 0, "type": "PermissionError", "errno": 1, "detail": "[Errno 1] Operation not permitted"}], "errors": [], "wrapper_exception_count": 1, "wall_seconds": 32.54994224999973}; Python exit: exit status 1
FAIL	github.com/sunholo-data/ailang-world/scripts/marketingcheck	32.840s
mutant run 3: supervisor_test.go:229: stress summary: {"trials": 200, "wrapper_calls": 600, "raw_calls": 600, "raw_eperm": 1, "raw_esrch": 599, "raw_success": 0, "wrapper_exceptions": [], "errors": [], "wrapper_exception_count": 0, "wall_seconds": 31.30023433399765}; Python exit: <nil>
ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	31.540s
```

## C11 protected-file observation
before (12 records) vs after: marketingcheck compare rc 0; files byte-identical.

## -race campaign 1 (stop at first red; concurrent with the judge)
```
campaign start 2026-10-27T16:32:26Z head=7475ec3c 17:32  up 1 day, 27 mins, 4 users, load averages: 7.45 6.80 7.90
M7a-fix macOS race trial 1/25 rc=0 117s 
M7a-fix macOS race trial 2/25 rc=0 112s 
M7a-fix macOS race trial 3/25 rc=0 126s 
M7a-fix macOS race trial 4/25 rc=1 186s --- FAIL: TestMarketingSupervisorFakeChromeTimeoutKillsGroup (8.39s) 
CAMPAIGN DONE for 7475ec3c: 3/25 green 2026-10-27T16:41:27Z
trial 4 failure:
--- FAIL: TestMarketingSupervisorFakeChromeTimeoutKillsGroup (8.39s)
    supervisor_test.go:495: fake must record child and grandchild: []
FAIL
```

## -race campaign 2 (all 25, reds classified; after the judge finished)
```
campaign2 start 2026-10-27T16:47:12Z head=7475ec3c 17:47  up 1 day, 41 mins, 4 users, load averages: 6.21 30.93 35.98
trial 1/25 rc=0 62s class=green load=6.84 26.62 34.02 
trial 2/25 rc=0 63s class=green load=6.61 22.70 32.00 
trial 3/25 rc=0 63s class=green load=6.92 19.91 30.33 
trial 4/25 rc=0 62s class=green load=4.77 16.82 28.40 
trial 5/25 rc=0 61s class=green load=5.18 14.76 26.84 
trial 6/25 rc=0 62s class=green load=4.25 12.83 25.30 
trial 7/25 rc=0 59s class=green load=3.09 10.92 23.73 
trial 8/25 rc=0 58s class=green load=2.57 9.30 22.25 
trial 9/25 rc=0 59s class=green load=2.46 8.14 21.00 
trial 10/25 rc=0 59s class=green load=2.66 7.18 19.77 
trial 11/25 rc=0 59s class=green load=2.82 6.41 18.62 
trial 12/25 rc=0 60s class=green load=2.57 5.71 17.53 
trial 13/25 rc=0 59s class=green load=3.49 5.37 16.59 
trial 14/25 rc=0 62s class=green load=9.44 6.50 16.23 
trial 15/25 rc=0 60s class=green load=6.80 6.40 15.54 
trial 16/25 rc=0 60s class=green load=4.15 5.74 14.67 
trial 17/25 rc=0 62s class=green load=3.56 5.22 13.82 
trial 18/25 rc=0 59s class=green load=3.73 5.00 13.19 
trial 19/25 rc=0 60s class=green load=3.45 4.74 12.54 
trial 20/25 rc=0 59s class=green load=3.07 4.40 11.88 
trial 21/25 rc=0 57s class=green load=2.79 4.10 11.26 
trial 22/25 rc=0 59s class=green load=3.19 3.98 10.76 
trial 23/25 rc=0 60s class=green load=3.69 3.94 10.28 
trial 24/25 rc=0 57s class=green load=2.58 3.59 9.72 
trial 25/25 rc=0 60s class=green load=2.46 3.40 9.23 
CAMPAIGN2 DONE for 7475ec3c: green=25 eperm=0 other=0 of 25 2026-10-27T17:12:14Z
```
