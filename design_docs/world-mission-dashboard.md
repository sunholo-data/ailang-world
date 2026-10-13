# World mission dashboard — iteration 256 (2026-10-13)

- **Latest landing:** row 172 `w-decision-packet-production-projection` → `d36b148` (PR #249). The coordinator mints a canonical `world/decision-packet/v1` on the effectful commit that requests `Human.Approve`; the object page validates it (bytes, record, plan, request, anchor, lifecycle, pair-scoped uniqueness) before projecting it; REST refuses the packet label.
- **Held:** the daemon's `Human.Approve` binding → new row 176, PARKED on **D-WORLD-73** (default: hold). Until then no production path stores a packet.
- **Judge:** sonnet PASS 86 → 92/100 (fix round), zero blocking; CI PR 2/2; `verify_go.sh` + `verify_ail.sh` green on pinned v0.41.0.
- **Clauses:** 1/2/3/6/7 MET; 4 UNMET (attended 93 FINAL); 5 UNMET (114 → 139).
- **Next loop pick:** D-WORLD-69 path 26 ✔ → 105 ✔ → 172 ✔ → 150, which still waits on fleet row 173 (PNG gate). If 173 has not moved, the loop has no routable row.
- **Parked on Mark:** **D-WORLD-73** (bind `Human.Approve` in the daemon? A = bind, B = hold; default B). Ledger 60 rows, ONE OPEN.
- **Routing this fire:** controller opus; Agent-tool roles: opus planner, haiku executor (2 runs + fix round), sonnet judge (r1 + r2); design from orphan 255 (opus).
- **Quota:** codex over ration all fire.
- **Metered spend:** $0 this fire (orphan 255 quorum $0.28).
- **Orphans:** 255 (STALL-killed at gate-3, foreground planner > 40 min), credited in 256.
- **Evidence:** `design_docs/verification/world-iter256/` (judge r1/r2 reports, quorum r1/r2 JSON).
