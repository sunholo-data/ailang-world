#!/usr/bin/env python3
"""Row 93 floor harness — pre-registered classification (design §4.6, V17; AC3.1).

Every task-run is exactly one of PASS, FAIL or HARNESS_FAULT:

  * PASS          — ``stdout_ok`` (whatever else happened).
  * FAIL          — an outcome. In either arm: compile/runtime/logic_error, step_exhausted,
                    non_agentic, refused, output_format, reasoning_stall. In the WORLD arm only
                    (World pays for its own faults): ``timeout``, ``world_tool_error`` (any MCP
                    or daemon error after a green pre-flight) and any UNEXPLAINED ``api_error``.
  * HARNESS_FAULT — the bar's "api_error / resource_limit / harness classes" (there is no
                    literal ``harness`` category, V17): api_error, resource_limit,
                    quota_exhausted, rate_limit, wire_drift, ``harness_setup`` (worktree,
                    grader, mock, pre-flight) and, in the SHELL arm only, ``timeout``.

An unexplained exit is ``api_error`` until a typed cause is evidenced. In the shell arm every
api_error is a harness fault (it makes the agent ineligible — never a FAIL, MUT-FAULT-AS-FAIL);
in the World arm an api_error is a harness fault only when its cause is TYPED
(``cause_typed: true``, e.g. a provider 5xx in the transcript) and evidenced, and a FAIL
otherwise. Every HARNESS_FAULT must carry a non-empty ``cause_evidence`` string.

Categories the table does not name (verify_error, constraint_violation, cost_killed,
thrash_aborted, policy_violation, ``none`` on a failed grade, …) are REFUSED: the table is
pre-registered, so an unlisted outcome stops the run rather than being silently binned.
"""
from __future__ import annotations

PASS = 'PASS'
FAIL = 'FAIL'
HARNESS_FAULT = 'HARNESS_FAULT'

FAIL_EITHER = frozenset({'compile_error', 'runtime_error', 'logic_error', 'step_exhausted',
                         'non_agentic', 'refused', 'output_format', 'reasoning_stall'})
FAULT_EITHER = frozenset({'resource_limit', 'quota_exhausted', 'rate_limit', 'wire_drift',
                          'harness_setup'})
# arm-dependent: timeout, api_error, world_tool_error
QUOTA_CATEGORIES = frozenset({'quota_exhausted', 'rate_limit'})


class ClassificationError(Exception):
    """A row the pre-registered table cannot classify; the run stops."""


def classify(row: dict) -> str:
    arm = row.get('arm')
    if arm not in ('shell', 'world'):
        raise ClassificationError(f'unknown arm {arm!r}')
    if row.get('stdout_ok') is True:
        return PASS
    if row.get('stdout_ok') is not False:
        raise ClassificationError('stdout_ok must be a bool')
    cat = row.get('error_category')
    evidence = (row.get('cause_evidence') or '').strip()
    if cat in FAIL_EITHER:
        return FAIL
    if cat == 'timeout':
        if arm == 'world':
            return FAIL
        result = HARNESS_FAULT
    elif cat == 'world_tool_error':
        if arm != 'world':
            raise ClassificationError('world_tool_error in the shell arm')
        return FAIL
    elif cat == 'api_error':
        if arm == 'world' and not (row.get('cause_typed') is True and evidence):
            return FAIL  # unexplained (no typed, evidenced cause): World pays for it
        result = HARNESS_FAULT
    elif cat in FAULT_EITHER:
        result = HARNESS_FAULT
    else:
        raise ClassificationError(f'category {cat!r} is not in the pre-registered table (§4.6)')
    if not evidence:
        raise ClassificationError(f'{arm} {cat} harness fault without cause_evidence')
    return result


def classify_all(rows):
    out = []
    for r in rows:
        r = dict(r)
        r['class'] = classify(r)
        out.append(r)
    return out
