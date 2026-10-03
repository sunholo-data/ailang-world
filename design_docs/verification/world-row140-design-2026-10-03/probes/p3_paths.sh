#!/bin/bash
# P3: implicit write paths, denyWrite precedence, case folding on APFS, mandatory denies, TMPDIR override.
. "$(dirname "$0")/lib.sh"
mkdir -p $S/ws3/.git $S/ws3/epcache/tmp
cd $S/ws3
cat > $S/s-p3.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$S/ws3"],
    "denyWrite": ["$S/ws3/.github", "$S/ws3/.claude", "$S/ws3/.git", "/tmp/claude", "/private/tmp/claude", "$HOME/.npm/_logs", "$HOME/.claude/debug"] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-p3-nodeny.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$S/ws3"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
echo "== implicit default write paths (no denyWrite)"
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'echo $TMPDIR; mkdir -p /tmp/claude && echo t > /tmp/claude/r140-probe && echo wrote-tmp-claude'
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'echo d > $HOME/.claude/debug/r140-probe && echo wrote-claude-debug'
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'mkdir -p $HOME/.npm/_logs && echo n > $HOME/.npm/_logs/r140-probe && echo wrote-npm-logs'
ls -la /tmp/claude/r140-probe $HOME/.claude/debug/r140-probe $HOME/.npm/_logs/r140-probe 2>&1
rm -f /tmp/claude/r140-probe $HOME/.claude/debug/r140-probe $HOME/.npm/_logs/r140-probe
echo "== mandatory denies inside the allowed root (no denyWrite)"
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'echo x > .git/config'
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'mkdir -p .git/hooks; echo x > .git/hooks/pre-commit'
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'echo x > .git/HEAD && echo wrote-git-HEAD'
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'echo x > .GIT/config && echo wrote-GIT-config-casevariant'
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'echo x > .mcp.json'
run $SRT --settings $S/s-p3-nodeny.json -- sh -c 'mkdir -p .github/workflows && echo x > .github/workflows/ci.yml && echo wrote-github'
echo "== operator denyWrite (precedence over allowWrite) and APFS case folding"
rm -rf .github .GIT; mkdir -p .git
run $SRT --settings $S/s-p3.json -- sh -c 'mkdir -p .github/workflows; echo x > .github/workflows/ci.yml'
run $SRT --settings $S/s-p3.json -- sh -c 'mkdir -p .GITHUB/workflows && echo x > .GITHUB/workflows/ci.yml && echo wrote-GITHUB-casevariant'
run $SRT --settings $S/s-p3.json -- sh -c 'mkdir -p .Claude && echo x > .Claude/settings.json && echo wrote-Claude-casevariant'
run $SRT --settings $S/s-p3.json -- sh -c 'echo x > .git/HEAD'
run $SRT --settings $S/s-p3.json -- sh -c 'echo x > .GIT/HEAD && echo wrote-GIT-HEAD-casevariant'
ls -la $S/ws3 $S/ws3/.github 2>&1
echo "== implicit paths closed by denyWrite"
run $SRT --settings $S/s-p3.json -- sh -c 'mkdir -p /tmp/claude; echo t > /tmp/claude/r140-probe2'
run $SRT --settings $S/s-p3.json -- sh -c 'echo d > $HOME/.claude/debug/r140-probe2'
echo "== TMPDIR redirected per episode (CLAUDE_CODE_TMPDIR)"
run env CLAUDE_CODE_TMPDIR=$S/ws3/epcache/tmp $SRT --settings $S/s-p3.json -- sh -c 'echo TMPDIR=$TMPDIR; mktemp && echo mktemp-ok'
