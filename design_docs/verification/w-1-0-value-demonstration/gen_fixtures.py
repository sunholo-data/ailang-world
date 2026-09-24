#!/usr/bin/env python3
"""gen_fixtures.py — w-prove-1-0-phase-a (row 92, Phase A) commit-fixture generator.

Builds one worldd commit fixture per taken 1.0-value-demonstration question, chained:
Q1 genesis (entryIndex 0) -> Q2 (1) -> Q3 (2) -> Q4/C4 (3). Each commit carries one
incident object (semanticId world/mission/incident/<slug>, provenance w-prove-1-0-phase-a)
plus one evidence object per source (semanticId world/mission/evidence/<slug>/<n>).

The FACTS (questions, answers, evidence kind/ref/check/excerpt) are copied verbatim from
the sprint plan's harvest table (design_docs/planned/w-prove-the-1-0-bar-phase-a-sprint-plan.md
section 2) and its design doc — they are the plan's ground truths, quoted, never paraphrased.

Wire schema (host/daemon/handlers.go, decoder is DisallowUnknownFields):
  observedHead (genesis: ""), objects[]{hash,interfaceHash,semanticId,provenance,payload(base64)},
  nextWorld{ref,revision,stateRoot,logHead}, entry{header{entryIndex,semanticsEpoch,
  transitionFn,interpreter,prevEntryHash,writtenBy},entryHash,transitionRef}.

Object `hash` MUST be sha256:<hex> of the base64-decoded payload bytes (the store
content-verifies, host/store/store.go verifyObject). `header.interpreter` is computed by
reading the raw bytes of $HOME/.pinned-ailang/ailang — never hard-coded.

--check re-reads the emitted files and asserts: exact key sets (DisallowUnknownFields),
every object hash == sha256(payload), contiguous entryIndex chain, prevEntryHash chaining,
and interpreter == sha256 of the pinned binary's current bytes.
"""

import argparse
import base64
import hashlib
import json
import os
import pathlib
import sys

# ── FACTS: the plan's harvest table, quoted verbatim (do not paraphrase) ─────────────
# Each entry: slug, question, answer, diagnosed_at, evidence[{kind,ref,check,excerpt}].
FACTS = [
    {
        "slug": "iter171-index-row",
        "question": "Why did iteration 171 have a full log entry and no index row?",
        "answer": ("The pinned v0.30.0 binary lacks the mission command group: "
                   "`ailang.v0.30.0 mission rotate-log world --keep 31` prints "
                   "`Error: unknown command 'mission'`; the index hole is visible in git "
                   "history of world-mission-index.md (the row exists today because it was "
                   "added manually after the iter-172 rule)."),
        "diagnosed_at": "2026-09-08",
        "evidence": [
            {
                "kind": "binary",
                "ref": "$HOME/.pinned-ailang/ailang.v0.30.0",
                "check": "ailang.v0.30.0 mission rotate-log world --keep 31 => Error: "
                         "unknown command 'mission'",
                "excerpt": ("The pinned v0.30.0 binary has no `ailang mission` command at "
                            "all: `~/.pinned-ailang/ailang mission rotate-log world --keep 31` "
                            "prints **`Error: unknown command 'mission'`**."),
            },
            {
                "kind": "index",
                "ref": "world-mission-index.md (git history)",
                "check": "grep -c '^| 171 |' world-mission-index.md => 1 NOW (the manual "
                         "repair, not iter-171's Gate 4)",
                "excerpt": ("iteration **171** had a full log entry and **no index row** "
                            "(`grep -c '^| 171 |'` -> **0**, control `170` -> **1**)."),
            },
            {
                "kind": "record",
                "ref": "world-mission.md ~lines 82-101; log ## 172",
                "check": "read the `mission`-command pin-version-gap bullet",
                "excerpt": ("the failure is SILENT in the direction that matters: Gate 2's "
                            "first instruction is to grep the index before picking, so a "
                            "missing row does not produce an error."),
            },
        ],
    },
    {
        "slug": "iter154-unfenced-pi",
        "question": ("Why did three designer runs execute unfenced while the rulebook said "
                     "they were sandboxed?"),
        "answer": ("The fleet's scripts/mission_pi_run.sh invokes "
                   "`pi --mode json --no-session --model \"$MODEL\" < \"$DIRECTIVE\"` with "
                   "no -e flag and no PI_FENCE_ROOT anywhere in the file, so the shared "
                   "skill's two sandbox extensions are never wired; issue "
                   "sunholo-data/ailang#1043 is OPEN."),
        "diagnosed_at": "2026-09-04",
        "evidence": [
            {
                "kind": "script",
                "ref": "scripts/mission_pi_run.sh (fleet, live via gh api)",
                "check": "grep -n 'pi --mode|PI_FENCE_ROOT| -e ' => invocation matches, ZERO "
                         "-e/PI_FENCE_ROOT hits",
                "excerpt": ("measured first-party at `scripts/mission_pi_run.sh:155`: the "
                            "invocation is `pi --mode json --no-session --model \"$MODEL\" "
                            "< \"$DIRECTIVE\"`, with **no `-e` flag at all** and no "
                            "`PI_FENCE_ROOT`."),
            },
            {
                "kind": "recipe",
                "ref": "shared-skill pi recipe (fleet-owned)",
                "check": "recipe invocation block mandates -e sandbox/index.ts and -e "
                         "worktree-fence.ts",
                "excerpt": ("The recipe's invocation block is explicit — `-e "
                            "\"$REPO/tools/pi-extensions/sandbox/index.ts\" -e "
                            "\"$REPO/tools/pi-extensions/worktree-fence.ts\"`."),
            },
            {
                "kind": "issue",
                "ref": "sunholo-data/ailang#1043",
                "check": "gh issue view 1043 --repo sunholo-data/ailang --json state --jq "
                         ".state => OPEN",
                "excerpt": ("the code half is filed upstream as "
                            "[`ailang#1043`](https://github.com/sunholo-data/ailang/issues/1043)."),
            },
        ],
    },
    {
        "slug": "iter155-faithfulness-proof",
        "question": ("Why did a green `shasum -c` faithfulness proof pass over a destroyed "
                     "commit split?"),
        "answer": ("The prescribed manifest covers the **final tree**, identical whether the "
                   "split is correct or collapsed; `git add design_docs host` staged disk "
                   "state, MS1 swallowed MS3's files, MS3 committed nothing, and `shasum -c` "
                   "returned OK on every file."),
        "diagnosed_at": "2026-09-05",
        "evidence": [
            {
                "kind": "manifest",
                "ref": "world-mission.md row 81",
                "check": "grep -n \"manifest is over the FINAL TREE\" world-mission.md",
                "excerpt": ("**That manifest is over the FINAL TREE.** Bisectability is a "
                            "property of the SPLIT, and the final tree is identical whether "
                            "the split is correct or whether commit 1 swallowed everything "
                            "and commits 2-3 are empty."),
            },
            {
                "kind": "log",
                "ref": "world-mission-log.md ## 155",
                "check": "git show --stat per commit; rebuilt row-59 split in git history",
                "excerpt": ("**My first commit reconstruction was wrong and I rebuilt it.** "
                            "`git add design_docs host` staged whatever was on disk rather "
                            "than the snapshot's named files, so MS1's commit swallowed "
                            "MS3's two doc files."),
            },
            {
                "kind": "retry",
                "ref": "world-mission-log.md ## 155",
                "check": "rebuild boundary gate; zsh word-splitting trap",
                "excerpt": ("zsh does not word-split unquoted expansions, so all three paths "
                            "arrived as ONE argument, `cp` failed, and the boundary gates "
                            "still printed `vet=0 verifygate=0`."),
            },
        ],
    },
    {
        "slug": "iter181-ci-bench-401",
        "question": ("Why did PR #141's remote CI go red on `BenchmarkRESTCommit` when every "
                     "local gate was green?"),
        "answer": ("a036062 is an ancestor of base and its diff adds "
                   "`req.Header.Set(\"Authorization\", auth)` to bench_test.go — the bench "
                   "predated the session gate and POSTed to the now-protected /v1/commit "
                   "unauthenticated, 401 (the middleware working as designed)."),
        "diagnosed_at": "2026-09-24",
        "evidence": [
            {
                "kind": "commit",
                "ref": "git a036062",
                "check": "git show a036062 -- host/daemon/bench_test.go | grep -c "
                         "'req.Header.Set(\"Authorization\", auth)' >= 1",
                "excerpt": ("`a036062` is an ancestor of base and its diff adds "
                            "`req.Header.Set(\"Authorization\", auth)` to `bench_test.go` — "
                            "the bench predated the session gate and 401'd."),
            },
            {
                "kind": "status",
                "ref": "world-mission.md STATUS 2026-09-24 (iter-181)",
                "check": "grep \"bench predates the sprint\" world-mission.md",
                "excerpt": ("remote CI on PR #141 then caught what every local gate missed — "
                            "`BenchmarkRESTCommit` POSTing to the now-protected `/v1/commit` "
                            "unauthenticated (401: the middleware working as designed; the "
                            "bench predates the sprint and was not in the plan's local gate "
                            "list)."),
            },
        ],
    },
]

TOPLVL_KEYS = {"observedHead", "objects", "nextWorld", "entry"}
OBJECT_KEYS = {"hash", "interfaceHash", "semanticId", "provenance", "payload"}
WORLD_KEYS = {"ref", "revision", "stateRoot", "logHead"}
ENTRY_KEYS = {"header", "entryHash", "transitionRef"}
HEADER_KEYS = {"entryIndex", "semanticsEpoch", "transitionFn", "interpreter",
               "prevEntryHash", "writtenBy"}


def pinned_binary() -> pathlib.Path:
    path = pathlib.Path(os.environ.get("HOME", "")) / ".pinned-ailang" / "ailang"
    if not path.is_file():
        sys.exit("fatal: pinned binary not found at %s" % path)
    return path


def sha256_bytes(data: bytes) -> str:
    return "sha256:" + hashlib.sha256(data).hexdigest()


def sha256_text(text: str) -> str:
    return sha256_bytes(text.encode("utf-8"))


def interpreter_ref(binary: pathlib.Path) -> str:
    # Replay pin: serve --ailang-bin archives the interpreter by hashing the RAW file
    # bytes. self check: never hard-code; read the file.
    return sha256_bytes(binary.read_bytes())


def incident_payload(fact) -> str:
    sources = [e["sha"] for e in fact["evidence"]]
    transport = {
        "question": fact["question"],
        "answer": fact["answer"],
        "diagnosed_at": fact["diagnosed_at"],
        "sources": sources,
    }
    return json.dumps(transport, separators=(",", ":"), ensure_ascii=False)


def evidence_payload(evidence) -> str:
    transport = {k: evidence[k] for k in ("kind", "ref", "check", "excerpt")}
    return json.dumps(transport, separators=(",", ":"), ensure_ascii=False)


def build_commits(interpreter: str):
    commits = []
    prev_entry_hash = None     # prior commit's entryHash
    prev_world_ref = None      # prior commit's nextWorld.ref
    for idx, fact in enumerate(FACTS):
        slug = fact["slug"]

        # Evidence objects first: their payload EMBEDS nothing chained, but the incident
        # payload's sources[] references their hashes, so hashes must be stable here.
        objects = []
        for n, evidence in enumerate(fact["evidence"]):
            payload_bytes = evidence_payload(evidence).encode("utf-8")
            eh = sha256_bytes(payload_bytes)
            evidence["sha"] = eh  # stable; consumed by incident sources[]
            objects.append({
                "hash": eh,
                "interfaceHash": sha256_text("interface:world/mission/evidence/%s/%d" % (slug, n)),
                "semanticId": "world/mission/evidence/%s/%d" % (slug, n),
                "provenance": "w-prove-1-0-phase-a",
                "payload": base64.b64encode(payload_bytes).decode("ascii"),
            })

        # Incident object: payload contains sources[] (= evidence hashes).
        incident_bytes = incident_payload(fact).encode("utf-8")
        incident_hash = sha256_bytes(incident_bytes)
        incident_obj = {
            "hash": incident_hash,
            "interfaceHash": sha256_text("interface:world/mission/incident/%s" % slug),
            "semanticId": "world/mission/incident/%s" % slug,
            "provenance": "w-prove-1-0-phase-a",
            "payload": base64.b64encode(incident_bytes).decode("ascii"),
        }
        objects.insert(0, incident_obj)

        transition_ref = incident_hash
        prev_entry_hash_this = sha256_text("w-1-0-value-demo/genesis-prev") if idx == 0 \
            else prev_entry_hash
        entry_hash = sha256_text("entry:%d:%s" % (idx, slug))
        world_ref = sha256_text("world:%d:%s" % (idx, slug))

        observed_head = prev_world_ref if prev_world_ref is not None else ""

        commit = {
            "observedHead": observed_head,
            "objects": objects,
            "nextWorld": {
                "ref": world_ref,
                "revision": idx,
                "stateRoot": sha256_text("state:%d:%s" % (idx, slug)),
                "logHead": entry_hash,
            },
            "entry": {
                "header": {
                    "entryIndex": idx,
                    "semanticsEpoch": 1,
                    "transitionFn": sha256_text("transition:provenance-walk"),
                    "interpreter": interpreter,
                    "prevEntryHash": prev_entry_hash_this,
                    "writtenBy": "w-prove-1-0-phase-a",
                },
                "entryHash": entry_hash,
                "transitionRef": transition_ref,
            },
        }
        commits.append(commit)
        prev_entry_hash = entry_hash
        prev_world_ref = world_ref
    return commits


def check_fixtures(outdir: pathlib.Path, interpreter: str) -> int:
    errors = 0
    files = sorted(outdir.glob("q*.json"))
    if len(files) != len(FACTS):
        print("check: expected %d fixtures, found %d" % (len(FACTS), len(files)))
        errors += 1
    for f in files:
        try:
            doc = json.loads(f.read_text())
        except Exception as exc:  # noqa: BLE001
            print("check: %s unreadable: %s" % (f.name, exc))
            errors += 1
            continue
        if set(doc.keys()) != TOPLVL_KEYS:
            print("check: %s top-level keys %r != %r" % (f.name, sorted(doc), sorted(TOPLVL_KEYS)))
            errors += 1
        for i, obj in enumerate(doc["objects"]):
            if set(obj.keys()) != OBJECT_KEYS:
                print("check: %s objects[%d] keys %r != %r" % (f.name, i, sorted(obj), sorted(OBJECT_KEYS)))
                errors += 1
            try:
                payload = base64.b64decode(obj["payload"])
            except Exception as exc:  # noqa: BLE001
                print("check: %s objects[%d] bad base64: %s" % (f.name, i, exc))
                errors += 1
                continue
            expect = sha256_bytes(payload)
            if obj["hash"] != expect:
                print("check: %s objects[%d].hash %s != sha256(payload) %s"
                      % (f.name, i, obj["hash"], expect))
                errors += 1
        if set(doc["nextWorld"].keys()) != WORLD_KEYS:
            print("check: %s nextWorld keys mismatch" % f.name)
            errors += 1
        if set(doc["entry"].keys()) != ENTRY_KEYS:
            print("check: %s entry keys mismatch" % f.name)
            errors += 1
        if set(doc["entry"]["header"].keys()) != HEADER_KEYS:
            print("check: %s header keys mismatch" % f.name)
            errors += 1
        if doc["entry"]["header"]["interpreter"] != interpreter:
            print("check: %s interpreter != pinned binary current sha256" % f.name)
            errors += 1
    # Chain assertions: contiguous entryIndex and prevEntryHash chaining.
    for f in files:
        doc = json.loads(f.read_text())
        idx = doc["entry"]["header"]["entryIndex"]
        if idx != files.index(f):
            print("check: %s entryIndex %d not contiguous (position %d)"
                  % (f.name, idx, files.index(f)))
            errors += 1
    for prev, cur in zip(files, files[1:]):
        pdoc = json.loads(prev.read_text())
        cdoc = json.loads(cur.read_text())
        if cdoc["entry"]["header"]["prevEntryHash"] != pdoc["entry"]["entryHash"]:
            print("check: %s prevEntryHash != %s entryHash" % (cur.name, prev.name))
            errors += 1
        if cdoc["observedHead"] != pdoc["nextWorld"]["ref"]:
            print("check: %s observedHead != %s nextWorld.ref" % (cur.name, prev.name))
            errors += 1
    return errors


def main(argv):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--out",
                        default=str(pathlib.Path(__file__).resolve().parent / "commits"),
                        help="output directory (default: the commits/ directory next to this script)")
    parser.add_argument("--check", action="store_true",
                        help="re-read emitted fixtures and verify schema/hashes/chain")
    args = parser.parse_args(argv)

    binary = pinned_binary()
    interpreter = interpreter_ref(binary)

    if args.check:
        out_dir = pathlib.Path(args.out)
        errors = check_fixtures(out_dir, interpreter)
        if errors:
            print("check: %d error(s)" % errors)
            return 1
        print("check: OK — %d fixture(s), schema + content-hashes + chain + interpreter "
              "asserted against %s" % (len(FACTS), binary))
        return 0

    out_dir = pathlib.Path(args.out)
    out_dir.mkdir(parents=True, exist_ok=True)
    for idx, commit in enumerate(build_commits(interpreter)):
        slug = FACTS[idx]["slug"]
        path = out_dir / ("q%d-%s.json" % (idx + 1, slug))
        path.write_text(json.dumps(commit, indent=2, ensure_ascii=False) + "\n")
        print("wrote %s" % path)
    # Sanity self-check the just-written tree before returning (matches --check logic).
    errors = check_fixtures(out_dir, interpreter)
    print("self-check: %d error(s); interpreter pin %s" % (errors, interpreter[:20] + "..."))
    return 0 if errors == 0 else 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))