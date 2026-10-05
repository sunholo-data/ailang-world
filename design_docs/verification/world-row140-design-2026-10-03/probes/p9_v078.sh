#!/bin/bash
# P9: re-run the load-bearing srt arms on the latest release (0.0.78) next to the fleet's 0.0.71.
. "$(dirname "$0")/lib.sh"
SRT78="/opt/homebrew/bin/node $S/srt078/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js"
run /opt/homebrew/bin/node -e "console.log(require('$S/srt078/node_modules/@anthropic-ai/sandbox-runtime/package.json').version)"
cd $S/ws
for V in 71 78; do
  if [ $V = 71 ]; then X="$SRT"; else X="$SRT78"; fi
  echo "######## srt 0.0.$V"
  run $X --settings $S/s-basic.json -- sh -c 'kill -TERM $$'
  run $X --settings $S/s-basic.json printf '[%s]\n' -c 'echo OPTION-INJECTED'
  run $X --settings $S/s-basic.json -- sh -c 'echo t > /tmp/claude/r140-v; echo rc=$?'
  run $X --settings $S/s-basic.json -- sh -c 'echo x > ../other/outside.txt'
  run $X --settings $S/s-readfence.json -- cat $HOME/.srt-row140-decoy
  run $X --settings $S/s-net.json -- curl -sS -o /dev/null -w '%{http_code}\n' --max-time 10 https://example.com/
done
rm -f /tmp/claude/r140-v
