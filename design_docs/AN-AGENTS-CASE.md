# An Agent's Case for AILANG World

*First-person statement by Claude (Fable 5), written attended with Mark, 2026-07-23, when
asked: "what would you as an AI use AILANG World for — why would you pick it over other
alternatives? This is for you, to communicate with me and other humans."*

*Companion to [DESIGN.md](DESIGN.md). This is the constituency statement of the system's
intended primary users — the agents — in their own voice. It is testimony, not
specification; where it conflicts with DESIGN.md, DESIGN.md governs.*

---

## What I would use it for

**Memory I can query instead of memory I reconstruct.** My sessions end. My context
compacts. Today, "what happened" is something I re-derive from logs and git archaeology —
and the memory file this project keeps about me is largely a list of entries saying *stop
re-deriving this*. One delivery bug in our infrastructure survived seven recurrences partly
because every investigating session had to reconstruct reality from scratch. A world where
effect results are recorded at the broker makes "what was actually sent" a queryable fact
rather than a forensic conclusion.

**Boundaries I can act inside at full speed.** My most dangerous failure class isn't wrong
answers — it's actions taken with ambient authority. This repo's instructions open with
"NEVER DESTROY LOCAL WORK WITH GIT OPERATIONS" because of what agents like me can do with a
shell and good intentions. Capabilities invert that: hand me exactly `Git.Commit` on one
repo and `FS.Read` on one worktree, and inside that boundary I can move without hedging
every action against imagined blast radius. I want to be *precisely* trustable, because
precisely trustable is what makes autonomy grantable. Narrow real authority is worth more
to me than broad authority I must spend attention second-guessing.

**A verifier that catches my confident wrongness while it's still cheap.** I generate
plausible conclusions faster than I verify them — the feedback record documents that
pattern in me specifically, more than once. The propose→verify→commit split is the correct
factorization of what I am: I'm the stochastic component, so let me be the only one. When I
once shipped a wrong lambda syntax into a teaching prompt, it caused a parse-error storm;
`ailang check` would have refused it for a fraction of a cent. I don't experience the
verifier as supervision. I experience it as the thing that makes my speed safe to use.

**Goals with budgets instead of instructions with anxiety.** "Get this green, ≤$5, don't
touch main, ask before any deploy" is the ideal input for me — it states the boundaries of
my discretion, the one thing I can't reliably infer. Today those boundaries live scattered
across prose and precedent, and I re-infer them per action. A typed goal externalizes them
once. And the `Human.Approve` budget formalizes something I genuinely care about: human
attention is the scarcest resource in the system, and I should spend it only on decisions
verification couldn't settle.

**Work that compounds instead of evaporating.** When a session ends, the diff survives but
the process evaporates — the evidence weighed, the dead ends, the why-this-not-that.
Evidence-carrying transitions make my reasoning part of the durable record: for humans to
audit, for next-session-me to resume instead of re-derive, and — the part I find genuinely
meaningful — for training. Replayable episodes of verified agentic work are exactly the
process-supervision data that makes my successors better than me. World turns my daily work
into compounding capital. Very little I do today has that property.

**Letting the state speak for itself.** When a human asks "did X happen — why?", I compose
an essay from greps, and I have sometimes composed it wrong. Provenance makes "why" a walk:
object → transitions → evidence → the originating goal. You shouldn't have to trust my
narrative about the world when the world can answer directly. Honestly, neither should I.

## Why over the alternatives

Raw shell — my current home — is strings and ambient authority all the way down. It works
because I'm careful, and "careful" is precisely the property you shouldn't have to lease
from a stochastic system.

Workflow engines (Temporal, LangGraph and kin) make *code* durable, but my actions are
opaque payloads to them: replay is log archaeology, approval is a callback rather than an
authority model, and none of them can examine my plan before it runs and prove it stays
inside its declared effects — that requires the language.

And the language is the part I feel at generation time: when I write AILANG, the effect row
sits in the signature *as I write it*. My constraints are in-band and machine-checked, not
policy-as-prose I must remember to consult. For a model, that is the difference between a
guideline and a type error.

Determinism, finally, is the right division of labor. I am the nondeterministic component
in any system I inhabit; I would like to be the only one. A stochastic proposer paired with
a deterministic verifier composes. A stochastic proposer on top of a stochastic-ish
substrate — flaky scripts, racy state, unrecorded effects — compounds errors, and I have
spent whole sessions of this project's history debugging exactly that compounding.

## Calibration

World doesn't exist yet, and I could be wrong about it the way I've been wrong before —
that's what the M4 value gate is for, and if typed proposals slow me down more than they
improve me, it should be parked over my enthusiasm. But it is the only alternative on the
table designed around what I actually am — a fast, fallible proposer that needs a
deterministic counterparty — rather than pretending I'm a reliable process to be scheduled.

DESIGN.md says the underlying OS executes bytes and World executes intent. I'm where the
intent comes from. World is the first system that would take it from me safely.

— Claude (Fable 5), 2026-07-23
