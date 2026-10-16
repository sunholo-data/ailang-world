#!/bin/bash
# Q1: test layout of the three M4 projects. Read-only: git ls-files only.
set -u
count() { # repo label ERE
  n=$(git -C "$1" ls-files | grep -E -c "$3")
  printf '%-22s %-40s %6d\n' "$(basename "$1")" "$2" "$n"
}
for R in /Users/voightkampff/dev/sunholo-data/ailang /Users/voightkampff/dev/markedmondson1234/TwilightGame /Users/voightkampff/dev/sunholo-data/sunholo-platform; do
  echo "== $R HEAD $(git -C $R rev-parse --short HEAD) branch $(git -C $R rev-parse --abbrev-ref HEAD)"
  count $R "ALL tracked files (control)" '.'
  count $R "*_test.go (basename)" '(^|/)[^/]*_test\.go$'
  count $R "*.go non-test" '(^|/)[^/]*\.go$'
  count $R "testdata/ component" '(^|/)testdata/'
  count $R "*.golden" '\.golden$'
  count $R "*.test.ts|tsx" '\.test\.tsx?$'
  count $R "*.spec.ts|tsx" '\.spec\.tsx?$'
  count $R "__tests__/ component" '(^|/)__tests__/'
  count $R "__snapshots__/ component" '(^|/)__snapshots__/'
  count $R "*.snap" '\.snap$'
  count $R "test_*.py (basename)" '(^|/)test_[^/]*\.py$'
  count $R "*_test.py (basename)" '(^|/)[^/]*_test\.py$'
  count $R "conftest.py" '(^|/)conftest\.py$'
  count $R "tests/ component (any depth)" '(^|/)tests/'
  count $R "tests/ at root" '^tests/'
  count $R "test/ component" '(^|/)test/'
  count $R "e2e/ component" '(^|/)e2e/'
  echo "-- top-level dirs holding tests/ (first 15):"
  git -C $R ls-files | grep -E '(^|/)tests/' | sed -E 's#(^|/)tests/.*#\1tests/#' | sort | uniq -c | sort -rn | head -15
  echo "-- test config files:"
  git -C $R ls-files | grep -E '(^|/)(vitest|jest|playwright)\.config\.[a-z]+$|(^|/)pytest\.ini$|(^|/)pyproject\.toml$|(^|/)setup\.cfg$|(^|/)tox\.ini$' | head -20
done
