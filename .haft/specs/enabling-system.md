<!-- DRAFT — onboarding by haft agent on 2026-07-18; operator must review and edit -->

# Enabling System Spec

## Team

Solo maintainer at this stage (repo owner). The repository is not yet a git repository, so there is no commit history to infer team size or contribution pattern from.

## Tooling snapshot (informational, not a spec-section claim)

- Build: not yet established — the repo currently contains only haft governance scaffolding (`.haft/`, `CLAUDE.md`, `.mcp.json`) and no package manifest or source tree.
- Test: none yet.
- Lint: none yet.
- Release: none yet.
- Dogfood/inspiration reference: `/Users/malyshev/Projects/Me/haft` — a standalone tool with an MCP server, AI framing, slash commands, and Claude/Cursor integration, currently on v8.1.0 (governance-substrate pivot: consumed via host-AI skills + MCP, not a standalone agent). Pysar's harness surface (skills, MCP kernel, contract docs) is expected to follow a similar shape, per the operator's stated intent.

## ES.enabling.architecture.001 Functional core / imperative shell, five one-way layers

```yaml spec-section
id: ES.enabling.architecture.001
spec: enabling-system
kind: enabling.architecture
title: Functional core / imperative shell, five one-way layers
statement_type: definition
claim_layer: description
owner: human
status: active
valid_until: 2027-01-18
depends_on: [TS.target.boundary.001]
supersedes: []
terms: [author-surface, harness-surface, editorial-engine]
target_refs: [TS.target.role.001, TS.target.boundary.001]
evidence_required: []
```

No code exists yet, so this is an intended architecture (`claim_layer: description`), not an observation of built structure. It follows `CLAUDE.md`'s existing "functional core, imperative shell" rule for this project. Five layers, one-way dependency (outer depends on inner, never reverse):

1. **Surface layer** — author-surface verbs (chat) and harness-surface (skills/CLI/MCP), presentation only. Owns: routing an author's plain-language intent to the right pass.
2. **Skills/routing layer** — maps author intent → pipeline stage (intake, staff-edit, sharpen, voice-lock, SEO-optimize) without exposing those names outward. Owns: intent-to-phase resolution.
3. **Editorial engine core (functional core)** — pure content transforms: stake/outline/angle scaffolding, edit passes, voice-lock, checklist generation. No I/O. Owns: the actual editorial logic, testable without mocking the world.
4. **Artifact/state layer** — drafts, `sources.md`, ship-ready body, checklist, cover — the project's own artifact graph (parallel to but distinct from `.haft/`'s governance graph). Owns: persistence of piece state across passes.
5. **Platform adapter layer (imperative shell)** — per-platform checklist/cover requirements, any external research/citation fetch. Owns: the only layer allowed to touch the outside world (network, files, platform APIs).

## ES.enabling.work_methods.001 Actor, trigger, and closing check per artifact kind

```yaml spec-section
id: ES.enabling.work_methods.001
spec: enabling-system
kind: enabling.work_methods
title: Actor, trigger, and closing check per artifact kind
statement_type: duty
claim_layer: description
owner: human
status: active
valid_until: 2027-01-18
depends_on: [ES.enabling.architecture.001]
supersedes: []
terms: [stake, ship-ready, voice-lock, publish-checklist]
target_refs: [TS.environment-change.001]
evidence_required: []
```

| Artifact | Actor | Trigger | Closing check |
|---|---|---|---|
| Stake/outline/angles | editorial engine core | author states an idea | author confirms the stake reads as their own words |
| Draft edit (engagement/substance) | editorial engine core | author flags a problem ("this opening is weak") or asks to "write it" | edited passage preserves author POV (no rewrite of stated thesis) |
| `sources.md` | platform adapter layer (research) | author explicitly opts in ("I need sources") | sources added; thesis/body text diff is empty outside the opt-in pass |
| Voice-lock pass | editorial engine core | author says "it sounds like AI," and only after packaging | pass runs strictly after packaging, never before |
| Ship-ready body | editorial engine core, full pass | author says "write it" | body assembled without surfacing phase names to the author |
| Publish checklist (+ cover) | platform adapter layer | author says "I'm ready to post" + names a platform | checklist scoped to the named platform; no post/publish action taken |

This is a project-local work-methods spec, distinct from `.haft/workflow.md` (which governs how *haft governance itself* is worked, not how Pysar's editorial passes are worked).

## ES.enabling.effect_boundaries.001 Only the owning layer may mutate its resource class; no auto-publish network call exists

```yaml spec-section
id: ES.enabling.effect_boundaries.001
spec: enabling-system
kind: enabling.effect_boundaries
title: Only the owning layer may mutate its resource class; no auto-publish network call exists
statement_type: duty
claim_layer: description
owner: human
status: active
valid_until: 2027-01-18
depends_on: [ES.enabling.architecture.001]
supersedes: []
terms: [editorial-engine, harness-surface, publish-checklist]
target_refs:
  - ES.effect.law-definition
  - ES.effect.admissibility-gate
  - ES.effect.deontics-duty
  - ES.effect.evidence-carrier
evidence_required: []
```

- **ES.effect.law-definition** — each resource class has exactly one owning layer: session/piece state → editorial engine core; on-disk piece artifacts (draft, `sources.md`, ship-ready body, checklist, cover) → artifact/state layer; network + platform requirement lookups → platform adapter layer. `.haft/` is haft's own governance graph, owned by the haft CLI/agent, never by Pysar's runtime.
- **ES.effect.admissibility-gate** — a mutation is admissible only when issued by its owning layer: the editorial engine core (functional core) never performs filesystem or network I/O directly; the platform adapter layer never authors or rewrites thesis text; nothing in Pysar's runtime writes to `.haft/`.
- **ES.effect.deontics-duty** — the platform adapter layer has a duty to never call an actual publish/post action; it produces a checklist, not a submission. Actually wiring a post/publish action is a later, separate one-way-door decision requiring explicit human authorization, not something this spec pre-approves.
- **ES.effect.evidence-carrier** — the boundary holds when: filesystem/network calls in the codebase exist only inside artifact/state and platform-adapter modules respectively; no publish/post API call exists anywhere in the codebase; and `.haft/` has zero writes originating from Pysar's own runtime code (only from the haft CLI/agent/MCP server).

## ES.enabling.agent_policy.001 Claude Code is the supported host; haft's existing auto/manual skill split governs delegation

```yaml spec-section
id: ES.enabling.agent_policy.001
spec: enabling-system
kind: enabling.agent_policy
title: Claude Code is the supported host; haft's existing auto/manual skill split governs delegation
statement_type: duty
claim_layer: description
owner: human
status: active
valid_until: 2027-01-18
depends_on: [ES.enabling.architecture.001]
supersedes: []
terms: [harness-surface]
target_refs: []
evidence_required: []
```

Supported host agent: **Claude Code** (`.claude/` present, `.mcp.json` wires the `haft` MCP server). `.cursor/mcp.json` is also present but Cursor is experimental/deferred per haft's own `PROJECT_ONBOARDING_CONTRACT.md` host table — not a v1 acceptance target for Pysar either.

This project does not need a separate agent-autonomy policy invented from scratch: `CLAUDE.md` already declares it, and this section binds to that existing declaration rather than restating it:

- Agents may invoke `haft_*` tools freely for framing, exploring, comparing, diagnosing, verifying, and status — these auto-fire on matching signal (CLAUDE.md skill catalog).
- `/h-decide` (binding a DecisionRecord) and `/h-commission` (creating a WorkCommission) are **manual-only** — an agent must never fire these autonomously, per the Transformer Mandate already stated in `CLAUDE.md` critical reminder #3.
- No commits without explicit operator permission (`CLAUDE.md` critical reminder #2) — this is the corresponding duty for any coding agent working Pysar's own source, once it exists.
- This section governs delegation into *Pysar's own* editorial-engine work, once code exists; it does not re-govern haft's own MCP tool surface, which `CLAUDE.md`'s installed haft section and the MCP server instructions already own.

## ES.enabling.commission_policy.001 Human-only commission creation, narrow default scope, no commission without a fresh governing decision

```yaml spec-section
id: ES.enabling.commission_policy.001
spec: enabling-system
kind: enabling.commission_policy
title: Human-only commission creation, narrow default scope, no commission without a fresh governing decision
statement_type: duty
claim_layer: description
owner: human
status: active
valid_until: 2027-01-18
depends_on: [ES.enabling.effect_boundaries.001, ES.enabling.agent_policy.001]
supersedes: []
terms: []
target_refs: []
evidence_required: []
```

- **Who may create:** only the human principal, via manual `/h-commission` (per `ES.enabling.agent_policy.001` — an agent never self-authorizes a WorkCommission).
- **Default scope:** `allowed_paths` narrows to whatever source tree exists under Pysar's five architecture layers (`ES.enabling.architecture.001`) once code lands; `forbidden_paths` always includes `.haft/**` — Pysar runtime work never has commission authority to edit haft's own governance graph.
- **Freshness gates before a commission is runnable:** every spec section its scope touches must be `active` with a current (non-drifted) baseline; a drifted or unbaselined governing section blocks the commission the same way `haft spec check` blocks readiness here.
- **Retirement:** a commission closes on evidence review (`ES.enabling.evidence_policy` — next phase), not on the agent's self-report that work is done.

This inherits `.haft/workflow.md`'s existing project defaults (`require_decision: true`, `require_verify: true`, `allow_autonomy: false`) rather than overriding them.

## ES.enabling.runtime_policy.001 CLI owns runtime lifecycle; the MCP/Claude Code surface only creates, observes, and reviews

```yaml spec-section
id: ES.enabling.runtime_policy.001
spec: enabling-system
kind: enabling.runtime_policy
title: CLI owns runtime lifecycle; the MCP/Claude Code surface only creates, observes, and reviews
statement_type: duty
claim_layer: description
owner: human
status: active
valid_until: 2027-01-18
depends_on: [ES.enabling.commission_policy.001]
supersedes: []
terms: []
target_refs: []
evidence_required: []
```

No harness runtime exists for Pysar's own editorial-engine code yet (no code exists). This section fixes the intended lifecycle split for when it does:

- **Lifecycle ownership:** the `haft` CLI (`haft harness run` / `prepare` / `status` / `result` / `apply` / `cancel`) starts and stops RuntimeRuns for Pysar WorkCommissions. The MCP plugin (this Claude Code session) may create WorkCommissions and observe/review RuntimeRuns and Evidence, but does not itself own long-running execution — consistent with haft's own `PROJECT_ONBOARDING_CONTRACT.md` actor table ("MCP host agent... may create drafts/notes/commissions only inside explicit tool contracts").
- **Isolation:** one RuntimeRun per WorkCommission scope; overlapping `allowed_paths` across concurrent runs is not permitted.
- **Observability:** each RuntimeRun writes RuntimeRun/PhaseOutcome/Evidence back into the haft artifact graph (per `SpecCoverage` edges in haft's `SPECIFICATION_ONTOLOGY.md`), not into a separate untracked log.

## ES.enabling.evidence_policy.001 Mechanical claims need same-repo tests; POV/voice claims need human review — never averaged together

```yaml spec-section
id: ES.enabling.evidence_policy.001
spec: enabling-system
kind: enabling.evidence_policy
title: Mechanical claims need same-repo tests; POV/voice claims need human review — never averaged together
statement_type: duty
claim_layer: description
owner: human
status: active
valid_until: 2027-01-18
depends_on: [ES.enabling.runtime_policy.001]
supersedes: []
terms: [voice-lock, ship-ready]
target_refs: []
evidence_required:
  - kind: manual
    description: Human confirms this evidence policy still matches how Pysar actually verifies claims once code and tests exist.
```

Two distinct claim classes, deliberately not pooled into one score (per CLAUDE.md's own R_eff = min, never average, weakest-link principle):

- **Mechanical/structural claims** (layer isolation per `ES.enabling.architecture.001`, effect boundaries per `ES.enabling.effect_boundaries.001`, voice-lock ordering "never before packaging," "no publish/post API call exists") — admissible evidence is an automated test in this repo (unit or contract test through the public interface, per `CLAUDE.md` reminder #5 "test contracts, not implementation"). Congruence level: CL3 (same-context, this repo's own tests) required; no external/borrowed test suite substitutes.
- **Qualitative claims** (POV preserved, "reads as the author's own words," voice-lock actually removed the "sounds like AI" quality) — no automated test can fully verify these. Admissible evidence is a human/manual review by the operator against a real sample piece. These claims are never scored as "passing" by test coverage alone.
- **R_eff composition:** a feature that satisfies its mechanical tests but fails manual POV review is **not** verified — the weaker of the two evidence classes governs, never an average across classes.
- **Refresh triggers:** any change to the five architecture layers or to platform-adapter checklist rules invalidates prior evidence for the sections it touches; default evidence `valid_until` is 90 days unless a section states otherwise, per `CLAUDE.md`'s evidence-decay glossary entry.
