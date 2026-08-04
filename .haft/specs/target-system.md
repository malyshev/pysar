<!-- DRAFT — onboarding by haft agent on 2026-07-18; operator must review and edit -->

# Target System Spec

## TS.environment-change.001 Author leaves with a ship-ready body

```yaml spec-section
id: TS.environment-change.001
spec: target-system
kind: target.environment
title: Author leaves with a ship-ready body
statement_type: definition
claim_layer: object
owner: human
status: active
valid_until: 2027-01-18
depends_on: []
supersedes: []
terms: [author-surface, harness-surface, stake, voice-lock, ship-ready]
target_refs: []
evidence_required: []
```

Before Pysar runs, an author with an idea or a rough draft has no shaped piece and no named diagnosis of what's wrong with what they have. After Pysar runs, that same author has: (1) a stake, outline, and angles scaffolded from their own words when starting from a bare idea; (2) a draft that keeps their POV while the correct editorial pass (engagement, substance, or voice-lock) has been applied; (3) sourced material added only when they opted into research, without their thesis being overwritten; and (4) a ship-ready body they trust, finalized on their own explicit signal. The observable flip: an author who could not turn a take into a piece they trust on their own now can, without ever having to learn or use pipeline/phase vocabulary (intake, staff-edit, sharpen, SEO optimize) to get there. An opt-in SEO/discoverability checklist (citation resolution, title/subtitle tuning, tags, meta description, URL slug, scannability) is available on explicit author signal (`--seo`), packaging the body before voice-lock (dec-20260804-e3234e50, superseding dec-20260718-e9f5b5e6's earlier full removal). Cover generation and actually posting/publishing to a platform remain out of current scope -- Pysar's promise ends at a ship-ready, optionally SEO-packaged body, not a publishing pipeline.

### Use case: Shape a raw idea into a piece

**Who:** An author who has a topic or idea but hasn't shaped it into a piece yet.
**What they do:** Says something like "I have an idea about X." Pysar scaffolds a stake, outline, and angles directly from the author's own words — no intake/phase vocabulary required.
**Why it matters:** Many people who want to publish long-form pieces struggle with writing, not with lacking a CMS. The gap is going from idea to shaped piece, not tooling.
**Acceptance signal:** Given a bare idea statement, Pysar produces a stake + outline + angles that read as an extension of the author's own words, not generic AI output.

### Use case: Improve an existing draft while preserving voice

**Who:** An author who has a draft that "feels wrong" but can't name why, or knows a specific part (e.g. the opening) is weak.
**What they do:** Says "Here's my draft" (intake from file) or "This opening is weak." Pysar routes automatically to the right edit (engagement vs. substance) based on the artifact's current state, and later runs a voice-lock pass if the output "sounds like AI" — strictly after packaging, never before.
**Why it matters:** Authors think in terms of what's wrong with the piece, not in terms of which editorial stage should run. Forcing pipeline vocabulary on them is a product failure.
**Acceptance signal:** A draft intake preserves the author's POV; a vague complaint ("this opening is weak," "it sounds like AI") resolves to the correct pass without the author naming it, and voice-lock never runs before packaging.

### Use case: Add sourced research on demand, never in place of the author's thesis

**Who:** An author whose draft or idea has a gap that needs citations or competitive context.
**What they do:** Says "I need sources." Pysar runs an opt-in research pass and adds sourcing material without rewriting the thesis.
**Why it matters:** Default input is the author's idea or draft, not web research. Research fills gaps the author opts into — it never replaces or overrides their take.
**Acceptance signal:** Triggering research adds a sourcing artifact; the author's thesis/body text is unchanged by the pass unless the author separately asks for a rewrite.

### Use case: Go from idea to a finalized ship-ready body

**Who:** An author ready to take a piece all the way to something they trust and consider done.
**What they do:** Says "Write it" for a full pass (idea → draft → editorial passes → ship-ready), then "I'm ready to post" to explicitly finalize the piece as ship-ready.
**Why it matters:** This is the core promise — "Bring your take. Pysar helps you ship it." The outcome is a body the author trusts, finalized on their own explicit signal. An opt-in SEO/discoverability pass (`--seo`) may package that body for search and platform discovery before voice-lock (dec-20260804-e3234e50); cover generation remains out of current scope.
**Acceptance signal:** "Write it" ends in a ship-ready body without exposing internal phase names; "I'm ready to post" finalizes that body as ship-ready; the SEO/discoverability checklist is produced only when the author explicitly opts in via `--seo`, never auto-triggered by finalization, and cover generation is never triggered.

## TS.target.role.001 Pysar plays an author-directed editorial engine, not an autonomous publisher

```yaml spec-section
id: TS.target.role.001
spec: target-system
kind: target.role
title: Pysar plays an author-directed editorial engine, not an autonomous publisher
statement_type: definition
claim_layer: object
owner: human
status: active
valid_until: 2027-01-18
depends_on: [TS.environment-change.001]
supersedes: []
terms: [author-surface, editorial-engine, stake, voice-lock]
target_refs: []
evidence_required: []
```

Role (what Pysar is assigned to be, distinct from what it can do or what it does on any given run): an **author-directed editorial engine**. Pysar is invoked by the author's own verbs ("I have an idea," "here's my draft," "write it," "I'm ready to post") and never originates a piece's thesis, never runs a pass the author didn't request or imply, and never crosses into publishing/CMS territory. The role is deliberately narrower than "autonomous ghostwriter" or "editor-in-chief": Pysar shapes and packages a take the author already owns; it does not decide what the author should think or when the piece goes out. Capability (scaffolding, editing, optional research, packaging) is available under this role, but capability existing does not mean Pysar exercises it unprompted — each pass runs only on the author's signal, per TS.environment-change.001.

## TS.target.boundary.001 In-scope: shape/edit/package a piece; out-of-scope: originate a thesis or auto-publish

```yaml spec-section
id: TS.target.boundary.001
spec: target-system
kind: target.boundary
title: In-scope shaping/editing/packaging vs. out-of-scope thesis origination and auto-publish
statement_type: definition
claim_layer: object
owner: human
status: active
valid_until: 2027-01-18
depends_on: [TS.target.role.001]
supersedes: []
terms: [author-surface, harness-surface, ship-ready, author-identity]
target_refs:
  - TS.boundary.law-definition
  - TS.boundary.admissibility-gate
  - TS.boundary.deontics-duty
  - TS.boundary.evidence-carrier
evidence_required: []
```

CHR-10 Boundary Norm Square, the four perspectives named in `target_refs`:

- **TS.boundary.law-definition** — Pysar is an editorial engine that shapes an author's own idea or draft into a ship-ready piece; it is not a CMS, not a publishing platform, and not a research/citation generator by default.
- **TS.boundary.admissibility-gate** — a Pysar session is admissible only when acting on an author-supplied idea or draft; Pysar does not originate a piece's stake on its own, and a research pass is admissible only after the author explicitly opts in.
- **TS.boundary.deontics-duty** — Pysar has a duty to preserve the author's POV through every editorial pass and to never publish or post on the author's behalf; the author has the corresponding duty to explicitly signal readiness ("I'm ready to post") before a piece is finalized as ship-ready.
- **TS.boundary.evidence-carrier** — the boundary holds when every ship-ready body traces back to author-supplied input (idea or draft), `sources.md` exists only on sessions where research was opted into, no platform post/publish action occurs -- Pysar never posts on the author's behalf -- and the SEO/discoverability checklist is generated only on the author's explicit `--seo` opt-in, never automatically (dec-20260804-e3234e50, superseding dec-20260718-e9f5b5e6's earlier full removal of this capability).

**In scope:** idea-to-stake/outline/angles scaffolding; draft intake preserving POV; routing vague complaints ("this opening is weak," "it sounds like AI") to the correct edit pass; opt-in research that adds sourcing without rewriting thesis; opt-in SEO/discoverability packaging (`--seo`) that resolves citations, tunes title/subtitle for CTR and read-completion, and fills a discoverability checklist (tags, meta description, URL slug, scannability) before voice-lock; full idea→draft→edit→ship-ready pass, finalized on explicit author signal.

**Out of scope:** originating a piece's thesis without author input; web research by default (opt-in only); actually posting/publishing to a platform (Pysar never posts on the author's behalf); cover/image generation (not part of the SEO/discoverability packaging capability); auto-triggering the SEO/discoverability checklist without the author's explicit `--seo` opt-in; forcing intake/staff-edit/sharpen/SEO-optimize vocabulary onto the author surface; a cross-project multi-author or author-registry model (dec-20260718-e84221af -- each initialized project/piece is the sole unit of authorial voice; a distinct voice means a distinct `pysar init`, not a role assignment within one project).
