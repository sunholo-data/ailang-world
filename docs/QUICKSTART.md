# Quickstart — run AILANG World's daemon in 5 minutes

*Every command below was executed verbatim on 2026-07-28 against `dev` (attended demo, Mark +
coordinator). Maintained under coding-standards **S7**: if this doc drifts from the binary,
that is a defect.*

## 1. Build and start

```bash
go build -o /tmp/ailang-worldd ./cmd/ailang-worldd
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin /tmp/ailang-v0300/ailang &
```

`serve` is loopback-only (a non-loopback `--bind` is refused, no override). `--ailang-bin`
archives and pins the interpreter at startup — its content hash becomes the D1 replay pin that
every log entry carries.

Every read below is bounded: a store read that has not answered within **10 s** is abandoned and
the route answers **HTTP 503, class `Timeout`**, naming the deadline — never a hang and never a
500. A genuine internal failure answers **HTTP 500 with the fixed body `internal store failure`**;
the verbatim cause (DSN path, driver detail) goes to the daemon's **stderr**, one line per error
carrying the route. So the terminal running `serve` is where you read *why* a 500 happened — the
HTTP client is told *that* one happened, and nothing about the host.

```bash
/tmp/ailang-worldd health
```

Returns the daemon version, DB path, and the pinned `interpreter_ref` + version.

## 2. Commit the genesis transition

A commit is JSON: `observedHead` (empty for genesis) + content-addressed `objects` (payload
base64; `hash` MUST be `sha256:<hex>` of the payload bytes — the store verifies) + `nextWorld`
+ the frozen 6-field log `entry` header. Generate a valid one:

```bash
python3 - <<'EOF'
import json, hashlib, base64
def sha(b): return "sha256:" + hashlib.sha256(b).hexdigest()
payload = json.dumps({"goal": "hello, World"}).encode()
interp  = open("/tmp/ailang-v0300/ailang","rb").read()
eh = sha(b"genesis-entry-1")
c = {"observedHead": "",
     "objects": [{"hash": sha(payload), "interfaceHash": sha(b"iface-v1"),
                  "semanticId": "world/demo/genesis-goal", "provenance": "quickstart",
                  "payload": base64.b64encode(payload).decode()}],
     "nextWorld": {"ref": sha(b"world-1"), "revision": 0,
                   "stateRoot": sha(b"state-1"), "logHead": eh},
     "entry": {"header": {"entryIndex": 0, "semanticsEpoch": 1,
                          "transitionFn": sha(payload), "interpreter": sha(interp),
                          "prevEntryHash": sha(b"genesis"), "writtenBy": "quickstart"},
               "entryHash": eh, "transitionRef": sha(payload)}}
open("/tmp/genesis.json","w").write(json.dumps(c, indent=2))
EOF
/tmp/ailang-worldd commit --file /tmp/genesis.json
```

Returns `{"selectedHead": "sha256:…"}` — the world now exists.

## 3. Read everything back

```bash
/tmp/ailang-worldd head
/tmp/ailang-worldd log get 0
/tmp/ailang-worldd world get "$(/tmp/ailang-worldd head)"
```

Check `log get 0`'s `header.interpreter` against `health`'s `interpreter_ref`: **identical** —
that is the replay pin, live. Object payloads read back with
`object get <hash> --payload` (base64; the store verified content-vs-hash on write).

## 4. See the guarantees refuse things

```bash
/tmp/ailang-worldd commit --file /tmp/genesis.json
```

→ HTTP 409 `HeadConflict` with `observedHead`/`selectedHead` — the structured conflict a
caller re-plans from; stale writers get facts, not corruption.

```bash
/tmp/ailang-worldd serve --db /tmp/world-demo.db
```

→ refused at startup: *"another process already holds writer authority for this database
(single-writer is enforced, not conventional)"* — the ratified arm-A lock, fail-closed.

## 5. Stop

SIGTERM (Ctrl-C / `kill`) drains bounded and releases the writer lock.

---
**Not yet in this quickstart** (arrives with the queue): effect broker + receipts (item 4),
MCP projection — drive commits from any MCP client (item 5), the approval-inbox workbench
(item 7). Transition payloads here are opaque demo objects; typed `world/*.ail` transitions
run through the replay engine (`host/replay`) and become the commit path when the broker lands.
