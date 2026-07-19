<!-- DRAFT — onboarding by haft agent on 2026-07-18; operator must review and edit -->

# Term Map

```yaml term-map
entries:
  - term: editorial-engine
    domain: target
    definition: >-
      How Pysar frames itself: not an "editorial pipeline product sold to editors,"
      but the engine underneath that helps an author who already has a take turn
      it into a shipped piece.
    not:
      - CMS
      - editorial pipeline product
    aliases: []
    owners: [human]
  - term: author-surface
    domain: target
    definition: >-
      The default, audience-facing interaction language. Authors speak in verbs
      about writing ("write it," "here's my draft," "I need sources," "this
      opening is weak") and never in pipeline/phase vocabulary.
    not:
      - harness-surface
    aliases: []
    owners: [human]
  - term: harness-surface
    domain: enabling
    definition: >-
      The maintainer- and agent-facing language: fixed pipeline order, Status
      enum, skills, rules, MCP kernel. Documented in CONTRACT.md and HOSTS.md
      (not yet present in-repo). Authors never see this surface.
    not:
      - author-surface
    aliases: []
    owners: [human]
  - term: phase-names
    domain: enabling
    definition: >-
      Harness-internal pipeline stage names (intake, staff-edit, sharpen, SEO
      optimize). Never required vocabulary for authors; exposed only through
      harness-surface docs (e.g. PIPELINE.md) and used internally to route
      author intents to the right pass.
    not:
      - author-surface vocabulary
    aliases:
      - pipeline stages
    owners: [human]
  - term: stake
    domain: target
    definition: >-
      The core claim/thesis of a piece, scaffolded from the author's own words
      before outline and angles are built out.
    not: []
    aliases: []
    owners: [human]
  - term: voice-lock
    domain: target
    definition: >-
      An editorial pass that removes the "sounds like AI" quality from a
      packaged piece. Runs after packaging, never before.
    not: []
    aliases: []
    owners: [human]
  - term: ship-ready
    domain: target
    definition: >-
      The state of a body the author trusts, finalized on the author's own
      explicit signal ("I'm ready to post"). The final state Pysar produces
      in current scope -- platform-specific publish checklists and cover
      generation are excluded (dec-20260718-e9f5b5e6).
    not: []
    aliases: []
    owners: [human]
  - term: author-identity
    domain: target
    definition: >-
      Pysar's identity/ownership model: single-voice-per-project. Each
      initialized project (a distinct `pysar init`) is the sole unit of
      authorial voice -- a distinct voice means a distinct initialized
      project, never a role or persona assigned within one project. There
      is no cross-project author/persona registry (dec-20260718-e84221af).
      Unlike Medium's POC, which assigns per-article authors from a shared
      authors/ registry serving many personas on one platform, Pysar serves
      one author's voice per initialized project -- a cross-project registry
      would be redundant scope, not a missing feature.
    not:
      - multi-author model
      - author registry
      - persona registry
    aliases: []
    owners: [human]
status: draft
```

This is a DRAFT term map. Entries carry product/architecture meaning and are a human confirmation point — review before any entry backs an active target-system claim.
