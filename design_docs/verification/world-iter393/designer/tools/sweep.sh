#!/usr/bin/env bash
# usage: sweep.sh <repo> <ailang-bin> ; prints one line per .ail: file passed errcodes@file:line verified/total statuses
repo="$1"; bin="$2"
cd "$repo" || exit 2
find . -iname '*.ail' -not -path './.git/*' | sed 's|^\./||' | sort | while read -r f; do
  modname="$(grep -m1 -E '^module ' "$f" | awk '{print $2}')"
  base="${f%"$modname.ail"}"; base="${base%/}"; [ -z "$base" ] && base="."
  rel="${f#"$base"/}"; [ "$base" = "." ] && rel="$f"
  out="$( cd "$base" && "$bin" ai-check -timeout 5s "$rel" 2>/dev/null )"
  printf '%s\n' "$out" | python3 -c '
import json,sys
f=sys.argv[1]
try: d=json.loads(sys.stdin.read())
except Exception as e: print("%-55s PARSE_FAIL %s"%(f,e)); sys.exit()
c=d.get("check",{}); v=d.get("verify",{}) or {}
errs=["%s@%s:%s"%(e.get("code"),e.get("file"),e.get("line")) for e in c.get("errors",[])]
res=v.get("results",[]) or []
st={}
for r in res: st[r.get("status")]=st.get(r.get("status"),0)+1
print("%-55s passed=%-5s errs=%s verify=%s"%(f,c.get("passed"),errs or "-",st or "-"))
' "$f"
done
