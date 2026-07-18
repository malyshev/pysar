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
      The state of a body the author trusts — the point at which the publish
      checklist can run.
    not: []
    aliases: []
    owners: [human]
  - term: publish-checklist
    domain: target
    definition: >-
      A platform-specific list of pre-post requirements, plus a cover image if
      the target platform requires one. Produced when the author signals
      they're ready to post.
    not: []
    aliases: []
    owners: [human]
status: draft
```

This is a DRAFT term map. Entries carry product/architecture meaning and are a human confirmation point — review before any entry backs an active target-system claim.
