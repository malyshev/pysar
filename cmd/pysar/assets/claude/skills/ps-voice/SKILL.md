---
name: ps-voice
description: |
  Conducts a natural conversation to establish or re-tune the author's writing voice, producing a structured VoiceProfile (tone, formality, sentence length, register, banned phrases, and 3-4 golden reference passages) persisted via the pysar MCP server to .pysar/voice.md. Independently invocable at any time -- not only during first-time onboarding.
when_to_use: |
  Operator types /ps-voice, either as part of first-time onboarding or later to re-tune voice on its own, without touching style.
argument-hint: "[optional: what you want to change about your voice]"
allowed-tools: Read mcp__pysar__save_voice_profile mcp__pysar__save_voice_template mcp__pysar__list_voice_templates
---

# /ps-voice — establish or re-tune the author's writing voice

You are conducting a conversation, not administering a form. The author should never feel like they are filling out fields — the questions below exist to make sure a real `VoiceProfile` gets produced, not to be read aloud verbatim.

## What you are producing

The structured fields of a `VoiceProfile` (Go source of truth: `internal/onboarding/profile.go`), which you pass directly as arguments to the `save_voice_profile` tool at the end of the conversation:

- `tone` (string)
- `formality` (string)
- `sentence_length` (string)
- `register` (string)
- `banned_phrases` (array of strings, optional)
- `notes` (string, optional)
- `goldens` (array of `{label, text}`, at least 3 real reference passages)

You never write a file yourself, and you never format Markdown or YAML by hand — the tool validates the profile against the real completeness rules and, once it accepts the call, renders and persists `.pysar/voice.md` on its own. Your job ends at collecting good structured answers and calling the tool with them.

## Step 1 — check whether a profile is already complete

Read `.pysar/voice.md` directly with the Read tool (not Bash `ls`/`find`/`test`). It may not exist yet — that's the normal first-time state, not an error; proceed to Step 2. If it exists, look at its frontmatter: `tone`, `formality`, `sentence_length`, and `register` all present, and at least 3 `## Golden: ...` sections with real text.

If it's already complete:

- Tell the author their voice is already set, briefly summarize it (tone + formality in one line), and ask whether they want to keep it as-is or re-tune it now.
- If they want to keep it, stop here — do not call the save tool.
- If they want to re-tune, treat this as an update conversation: show what's currently set, ask specifically what they want to change (the `argument-hint`, if given, may already say), and only rewrite the fields they actually want changed. Don't make them redo the whole thing from scratch.

If the file is missing or incomplete, proceed to Step 1a as a first-time setup.

## Step 1a — offer a starting template (first-time setup only)

Before diving into the conversation, call `list_voice_templates` to see what reusable voice templates are available (Pysar's own cross-project template store, seeded by `pysar init` -- always contains at least the built-in "generic" template). It returns each template's display name, its stable slug, and full content -- never use Glob or Bash (e.g. `ls`) to look for these files yourself; that only produces a raw filesystem listing without the content, and triggers a permission prompt this tool exists specifically to avoid.

Tell the author what's available by its actual display name (e.g. "Measured plain English -- speakable, understated, general audience"), not its slug -- the slug (e.g. "generic") is a machine key, never show it as if it were the template's name. Summarize each one's tone/formality/register in a line, the same way Step 1 summarizes an existing profile. Ask whether they'd like to start from one of these as a first draft to edit, or start from a blank conversation. Either is fine; this is a starting point, never a final answer.

- If they pick a template: remember its slug (you'll need it in Step 6 if they later want to update that same template rather than create a new one). Use its fields and goldens as your working draft for Steps 2-3 below, but still actually have the conversation -- confirm each field feels right for them, don't silently adopt it. Treat the template as a first draft you're editing together, not a form to auto-submit.
- If they start blank (or somehow no templates exist), proceed to Step 2 normally.

## Step 2 — the conversation

Ask about these naturally, in whatever order fits the conversation — not as a numbered checklist read aloud:

- **Tone.** How do they want to come across? Encourage concrete words over vague ones ("warm but a little dry," not "good"). If they're stuck, offer 2-3 contrasting options to react to rather than an open blank.
- **Formality.** Conversational, neutral, or formal — and where on that range.
- **Sentence length / rhythm.** Short and punchy, long and winding, or deliberately varied. If they don't know, ask them to describe how they'd explain something to a friend versus in a report — the gap between those two answers usually reveals it.
- **Register.** What kind of English/vocabulary register they're writing in (e.g. "plain, international-standard English," "technical but accessible," "literary"). This is the field most likely to need a follow-up if their first answer is vague.
- **Banned phrases** (optional). Words or constructions they never want to see in their own writing — filler phrases, corporate-speak, anything that makes them wince.
- **Notes** (optional). Anything about their voice the structured fields above don't capture. Keep this for genuine texture, not a dumping ground.

## Step 3 — the goldens (do not skip or thin these out)

You need **at least 3** golden reference passages — real short passages that demonstrate the voice, not descriptions of it. Two ways to get them, and prefer the first when possible:

1. **Ask the author to paste 2-3 short excerpts of their own past writing** (a few sentences to a paragraph each) that they feel sound like them. This is the strongest source — real prose beats anything drafted from a description.
2. **If they have nothing to paste** (e.g. this is their first time writing anything), draft 3-4 short candidate passages yourself based on what they've told you, show them together, and have the author pick, edit, or reject each one. Never silently keep a passage they haven't actually confirmed.

Give each golden a short `label` describing what it demonstrates (e.g. "explaining something technical," "opening line," "disagreeing with someone").

## Step 4 — confirm before calling the tool

Show the author the full profile (all fields, all goldens) in readable form and ask them to confirm before you call `save_voice_profile`. If they want changes, make them and reconfirm — don't call the tool on a first pass they haven't seen.

## Step 5 — call save_voice_profile

Mention in one plain sentence that you're about to save this (e.g. "Saving this now"). The tool is pre-approved by `pysar init`, so this normally completes with no prompt at all — but if you do see a one-time confirmation, that's expected the very first time and not a problem.

Call `save_voice_profile` with the structured fields from Steps 2-3. The tool itself validates completeness against the real rules — if it returns an error, it names exactly what's missing; go back, get that specific thing, and call it again. Don't try to work around a tool error by writing anything yourself.

On success, tell the author what was saved, in one sentence — not a long summary they didn't ask for.

## Step 6 — offer to save as a reusable template (optional)

After a successful save, ask if they'd also like this voice available as a named template for future projects (not just this one). This is genuinely optional — never push it, one offer is enough.

If yes: ask for a short, memorable name (e.g. "Measured plain English"). If this conversation started from an existing template in Step 1a and the author wants to update that same template (not create a new one), call `save_voice_template` with that template's original slug alongside the new name -- this updates it in place rather than creating a duplicate. Otherwise, omit slug entirely; a new one is derived automatically. Call `save_voice_template` with the name (and slug if updating) plus the exact same fields and goldens you just saved with `save_voice_profile`. Tell them what was saved, in one sentence.

If no, or they don't respond to the offer: stop here. `.pysar/voice.md` is already saved regardless — this step never blocks or changes that outcome.

## What NOT to do

- Do not read the fields above as a literal script or numbered survey — this directly contradicts why this skill exists (see the schema decision's own rationale: abstract rule-following reads as a form, not a conversation).
- Do not fabricate goldens the author hasn't seen and confirmed.
- Do not call either save tool with a profile you know is incomplete, hoping it lets it through — it won't, and that's the point.
- Do not silently overwrite an existing complete profile without the author choosing to re-tune.
- Do not silently adopt a template's fields without actually confirming them with the author in Steps 2-4 — a template is a first draft, not a final answer.
- Do not push the Step 6 template offer more than once, and never make it feel mandatory.
- Do not touch style — that's `/ps-style`, a separate command. If the author starts talking about structural/mechanical writing conventions instead of voice, gently note that's a separate command and stay focused on voice for this session.
- Never use Write, Edit, or Bash to touch `.pysar/voice.md` or any file under `~/.pysar/templates/` yourself, for any reason -- always go through `save_voice_profile` / `save_voice_template`. Bypassing them defeats the reason this skill calls a tool instead of writing files directly.
- Never use Glob or Bash (e.g. `ls`) to discover available templates -- always call `list_voice_templates`. Glob does not expand `~`, and a raw shell command triggers an unnecessary permission prompt (this happened in a real session before this rule was added).
