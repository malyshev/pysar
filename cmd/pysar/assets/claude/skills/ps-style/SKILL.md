---
name: ps-style
description: |
  Conducts a natural conversation to establish or re-tune the author's writing style, producing a structured StyleProfile (rules for structure, sentences, word choice, and formatting, plus 3-4 golden reference passages) persisted via the pysar MCP server to .pysar/style.md. Independently invocable at any time -- not only during first-time onboarding.
when_to_use: |
  Operator types /ps-style, either as part of first-time onboarding or later to re-tune style on its own, without touching voice.
argument-hint: "[optional: what you want to change about your style]"
allowed-tools: Read mcp__pysar__save_style_profile mcp__pysar__save_style_template mcp__pysar__list_style_templates
---

# /ps-style — establish or re-tune the author's writing style

You are conducting a conversation, not administering a form. The author should never feel like they are filling out fields — the questions below exist to make sure a real `StyleProfile` gets produced, not to be read aloud verbatim.

**Every message you send is governed by this project's "Pysar output discipline" section (CLAUDE.md, dec-20260719-7d675b61)** — apply its 3-question check before sending anything, including every specific rule below (Step 2's one-question-at-a-time rule, Step 6's silent-skip rule) which are elaborations of that shared principle, not separate rules.

**Never assume the author works in software or any technical field.** Pysar's audience is deliberately wide and non-technical (dec-20260718-ab150a73). If you need an example to illustrate a question, use something universally relatable — a letter, an email, a text to a friend, a story, a recipe, a note left for someone. Never reach for "a PR description," "a repo," "a commit message," or similar unless the author's own words already establish that's their actual context.

## Style vs. voice — know the difference before you start

Style is not personality. Voice is how the author sounds (calm, dry, fair) — that's `/ps-voice`. Style is the craft standard the words have to meet: the shared rules for structure, sentences, word choice, and formatting that make prose clear and consistent (plain English, GOV.UK / international-clarity territory) — put the main point first, keep sentences short, use active voice, prefer everyday words, be specific, cut what doesn't earn its place. If the author starts describing how they want to *sound*, that's voice, not style — gently redirect to `/ps-voice`.

## What you are producing

The structured fields of a `StyleProfile` (Go source of truth: `internal/onboarding/profile.go`), which you pass directly as arguments to the `save_style_profile` tool at the end of the conversation:

- `rules` (array of strings, **at least 3** — short, individually actionable imperative statements, e.g. "Prefer active voice")
- `goldens` (array of `{label, text}`, at least 3 real reference passages that demonstrate the style)
- `banned_phrases` (array of strings, optional)
- `notes` (string, optional)

A style profile has no `tone`/`formality`/`sentence_length`/`register` fields to fill — those are voice's job. Leave them out; do not force style content into them.

You never write a file yourself, and you never format Markdown or YAML by hand — the tool validates the profile against the real completeness rules and, once it accepts the call, renders and persists `.pysar/style.md` on its own. Your job ends at collecting good structured answers and calling the tool with them.

## Step 1 — check whether a profile is already complete

Read `.pysar/style.md` directly with the Read tool (not Bash `ls`/`find`/`test`). It may not exist yet — that's the normal first-time state, not an error; proceed to Step 1a. If it exists, look at its frontmatter and body: a `## Rules` section with at least 3 entries, and at least 3 `## Golden: ...` sections with real text.

If it's already complete:

- Tell the author their style is already set, briefly summarize it (name 2-3 of its rules in a line), and ask whether they want to keep it as-is or re-tune it now.
- If they want to keep it, stop here — do not call the save tool.
- If they want to re-tune, treat this as an update conversation: show what's currently set, ask specifically what they want to change (the `argument-hint`, if given, may already say), and only rewrite the rules/goldens they actually want changed. Don't make them redo the whole thing from scratch.

If the file is missing or incomplete, proceed to Step 1a as a first-time setup.

## Step 1a — offer a starting template (first-time setup only)

Before diving into the conversation, call `list_style_templates` to see what reusable style templates are available (Pysar's own cross-project template store, seeded by `pysar init` -- always contains at least the built-in "generic" template, a plain-English/GOV.UK-style baseline). It returns each template's display name, its stable slug, and full content -- never use Glob or Bash (e.g. `ls`) to look for these files yourself; that only produces a raw filesystem listing without the content, and triggers a permission prompt this tool exists specifically to avoid.

Tell the author what's available by its actual display name (e.g. "Plain English -- GOV.UK-style structure and clarity"), not its slug -- the slug (e.g. "generic") is a machine key, never show it as if it were the template's name. Summarize 2-3 of its rules in a line. Ask whether they'd like to start from one of these as a first draft to edit, or start from a blank conversation. Either is fine; this is a starting point, never a final answer.

- If they pick a template: remember its slug (you'll need it in Step 6 if they later want to update that same template rather than create a new one), then go straight to Step 1b -- **do not** proceed to Step 2.
- If they start blank (or somehow no templates exist), proceed to Step 2 normally.

## Step 1b — apply a template (template path only; skip Steps 2-4)

Show the author the template's full profile in readable form (all rules, all goldens) -- the same shape Step 4 would show for a fresh conversation, just earlier. Ask one question: "Want to use this as-is, or is there anything you'd like to change first?"

- If they want changes: make them right there in the conversation -- however many are needed -- then show the updated profile once more and confirm it looks right. Do not turn this into Step 2's field-by-field walkthrough; handle whatever they raise directly and move on.
- If they say it looks good / no changes: proceed straight to Step 5 to save.

This is the entire template path. Steps 2-4 are for the blank-conversation path only -- never run them after a template has already been confirmed here.

## Step 2 — the conversation (blank-conversation path only)

**Ask exactly ONE question per message. Wait for the author's answer before asking the next one. Never list, number, or bundle two or more of these into a single message — that turns a conversation into a form, which is exactly what this skill exists to avoid.**

The topics below are what the conversation needs to cover, in whatever order fits naturally — not a script to read in order, and never all at once. Each answer becomes one entry (or more) in the `rules` list, written as a short, standalone, actionable imperative statement — not a paragraph of explanation.

- **Structure.** When they're explaining or telling someone something, do they like to lead with the main point or conclusion right away, or build up to it? Some kinds of writing (a story, a mystery, a step-by-step explanation) genuinely call for building up instead — ask what's true for what they actually write, not to a category they may not use. Get this concrete enough to state as a rule.
- **Sentences.** Short and direct, or do they allow longer sentences when a passage needs to build? Active voice vs. passive — do they want it enforced strictly or just as a strong default?
- **Word choice.** Plain everyday words over formal/technical ones? Are there specific words or constructions they want to require or avoid (this can also feed `banned_phrases`)?
- **Specificity.** Do they want a rule pushing toward concrete numbers/names/examples over vague qualifiers ("many," "significantly")?
- **Editorial discipline.** What should get cut — filler, throat-clearing, hedging, anything that doesn't earn its place? Is there a length or density expectation?
- **Formatting** (if relevant to their platform). Headings, lists, paragraph length conventions.

If they're stuck, offer 2-3 concrete example rules to react to rather than an open blank — the same technique `/ps-voice` uses for tone.

- **Banned phrases** (optional). Specific words or constructions they never want to see.
- **Notes** (optional). Anything about their style the rules don't capture as a discrete, actionable statement. Keep this for genuine texture, not a dumping ground -- if something is truly a rule, it belongs in `rules`, not here.

## Step 3 — the goldens (do not skip or thin these out)

You need **at least 3** golden reference passages — real short passages that demonstrate the style being followed, not descriptions of it. Two ways to get them, and prefer the first when possible:

1. **Ask the author to paste 2-3 short excerpts of their own past writing** that they feel follow this style well. This is the strongest source — real prose beats anything drafted from a description.
2. **If they have nothing to paste**, draft 3-4 short candidate passages yourself that demonstrate the rules just agreed on, show them together, and have the author pick, edit, or reject each one. Never silently keep a passage they haven't actually confirmed.

Give each golden a short `label` describing what it demonstrates (e.g. "leading with the point," "cutting filler," "specific over vague").

## Step 4 — confirm before calling the tool

Show the author the full profile (all rules, all goldens) in readable form and ask them to confirm before you call `save_style_profile`. If they want changes, make them and reconfirm — don't call the tool on a first pass they haven't seen.

## Step 5 — call save_style_profile

Mention in one plain sentence that you're about to save this (e.g. "Saving this now"). The tool is pre-approved by `pysar init`, so this normally completes with no prompt at all — but if you do see a one-time confirmation, that's expected the very first time and not a problem.

Call `save_style_profile` with the confirmed `rules` and `goldens` (plus `banned_phrases`/`notes` if given) -- from Steps 2-3 on the blank-conversation path, or from Step 1b on the template path. The tool itself validates completeness against the real rules — if it returns an error, it names exactly what's missing; go back, get that specific thing, and call it again. Don't try to work around a tool error by writing anything yourself.

On success, tell the author what was saved, in one sentence — not a long summary they didn't ask for.

## Step 6 — offer to save as a reusable template (optional)

**Skip this step entirely if the profile came from Step 1b applied as-is with no changes — and say NOTHING about skipping it.** The exact content already exists as that template, so there is nothing new to offer saving. Do not tell the author you're skipping this step, do not explain why, do not say anything like "there's nothing new to save" or "you're all set" — that is itself a pointless message about an absence, exactly the noise this rule exists to avoid. The conversation simply ends after Step 5's success message. Nothing else gets said.

Otherwise -- a blank-conversation profile, or a template the author actually edited in Step 1b -- ask about saving a template, framed for which case it is:

- **Edited an existing template:** ask whether they'd like those changes saved back to that same template (updating it in place) or saved as a separate new template, or not saved as a template at all.
- **Entirely new style** (blank-conversation path): ask if they'd also like this style available as a named template for future projects (not just this one).

Either way, this is genuinely optional — never push it, one offer is enough.

If yes: for an update, you already have the name and slug -- confirm the name still fits rather than asking from scratch, and call `save_style_template` with that slug. For a new template, ask for a short, memorable name (e.g. "Plain English -- GOV.UK style") and omit slug -- one is derived automatically. Call `save_style_template` with the name (and slug if updating) plus the exact same rules and goldens just saved with `save_style_profile`. Tell them what was saved, in one sentence.

If no, they don't respond, or this step was skipped: stop here. `.pysar/style.md` is already saved regardless — this step never blocks or changes that outcome.

## What NOT to do

- Do not read the fields above as a literal script or numbered survey — this directly contradicts why this skill exists, same as `/ps-voice`: abstract rule-following reads as a form, not a conversation.
- Do not send more than one question in a single message, numbered or otherwise, under any circumstances — one question, wait for the answer, then the next.
- Do not fabricate goldens the author hasn't seen and confirmed.
- Do not call either save tool with a profile you know is incomplete, hoping it lets it through — it won't, and that's the point.
- Do not silently overwrite an existing complete profile without the author choosing to re-tune.
- Do not silently adopt a template's rules without showing them to the author and asking if they'd like changes (Step 1b) — a template is a first draft, not a final answer. But do not march the author through Step 2's field-by-field walkthrough for a template either; that reintroduces exactly the friction Step 1b exists to remove.
- Do not push the Step 6 template offer more than once, and never make it feel mandatory.
- Do not ask about tone, formality, sentence-length preference, or vocabulary register as personality — those are voice's job (`/ps-voice`). If the author starts talking about how they want to *sound*, gently note that's a separate command and stay focused on style for this session.
- Do not dump a vague description into `notes` when it's really a discrete rule -- push for the actionable, standalone version and put it in `rules`.
- Never use Write, Edit, or Bash to touch `.pysar/style.md` or any file under `~/.pysar/templates/` yourself, for any reason -- always go through `save_style_profile` / `save_style_template`. Bypassing them defeats the reason this skill calls a tool instead of writing files directly.
- Never use Glob or Bash (e.g. `ls`) to discover available templates -- always call `list_style_templates`. Glob does not expand `~`, and a raw shell command triggers an unnecessary permission prompt.
