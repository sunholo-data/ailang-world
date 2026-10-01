== judge lane probe ==
cmd: claude-sub -p 'reply with exactly: ok' --model claude-sonnet-4-6
rc: rc=0
output: ok

== judge run ==
cmd: claude-sub -p "$(cat /tmp/eval218-directive.md)" --model claude-sonnet-4-6 --permission-mode bypassPermissions (cwd .eval-world-iter218 @ ba9fdc6)
rc: rc=0
wall: launched 1790898474, done 1790898590 (~116s)
final line: `RECORD-VERDICT: PASS 96/100 all claims verified first-party; doc-only diff; clause map correct; D-WORLD-49 sole open decision; 215 rotated to archive; all 16 log sections present; zero HARD FAILs; zero BLOCKING findings`
tokens: not reported by claude CLI -p (subscription lane; $0 metered)
