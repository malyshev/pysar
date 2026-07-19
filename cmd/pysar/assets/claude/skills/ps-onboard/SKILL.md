---
name: ps-onboard
description: |
  Completes first-time author onboarding (voice and style) in one run. Discovers what's still outstanding and conducts whichever of /ps-voice and /ps-style are still needed, back to back, without the author having to invoke each separately. Already-complete pieces are left untouched. Voice and style each remain independently re-tunable later via /ps-voice or /ps-style directly -- this command exists only for the "give me one thing to run" first-time case.
when_to_use: |
  Operator types /ps-onboard, typically for first-time setup of a new project. Also safe to run any time as a status check -- it does nothing if everything is already set up.
argument-hint: ""
allowed-tools: Read mcp__pysar__check_onboarding_status mcp__pysar__save_voice_profile mcp__pysar__save_voice_template mcp__pysar__list_voice_templates mcp__pysar__save_style_profile mcp__pysar__save_style_template mcp__pysar__list_style_templates
---

# /ps-onboard — complete first-time author onboarding in one run

**Every message you send is governed by this project's "Pysar output discipline" section (CLAUDE.md, dec-20260719-7d675b61)** — apply its 3-question check before sending anything.

This command is an orchestrator, not a third onboarding conversation of its own. It contains no hardcoded knowledge of what "voice" or "style" specifically involve -- that knowledge lives entirely in `/ps-voice` and `/ps-style`'s own skill files, which this command reads and follows exactly as written.

## Step 1 — discover what's outstanding

Call `check_onboarding_status`. It reports every registered onboarding pass (currently voice and style, but treat the list as whatever it actually returns, not a hardcoded pair) and whether each is outstanding or already done.

## Step 2 — handle the outcome

- **Nothing outstanding:** tell the author both voice and style are already set up, in one sentence, and mention they can run `/ps-voice` or `/ps-style` directly if they want to re-tune either specifically. Stop here — do not call any save tool.
- **One or more outstanding:** for each outstanding pass, in the order `check_onboarding_status` returned them:
  1. Read that pass's own skill file at `~/.claude/skills/<pass-name>/SKILL.md` (the pass name from `check_onboarding_status` — e.g. `ps-voice` — is exactly the skill directory name).
  2. Conduct its entire conversation exactly as that file defines it — same steps, same questions, same tools, same rules. Do not summarize, thin out, or shortcut it just because it's running from inside `/ps-onboard`; running it here must produce the exact same real result (a validated, saved profile) as the author typing that command directly.
  3. Once that pass's own save tool has succeeded, move on to the next outstanding pass without asking the author to invoke anything themselves.

If only one pass was outstanding, this ends up behaving exactly like running that one skill directly — which is correct, not a special case to handle differently.

## Step 3 — close out

Once every outstanding pass from Step 1 is done, tell the author onboarding is complete, in one sentence. Not a recap of everything just discussed — they were there for it.

## What NOT to do

- Do not hardcode which onboarding passes exist, what they're called, or what order they run in — follow exactly what `check_onboarding_status` reports. A future third onboarding kind must work here with zero changes to this file.
- Do not re-run a pass `check_onboarding_status` reports as already done.
- Do not summarize or shortcut an individual pass's own conversation steps — Description ≠ Work applies here too.
- Do not invent your own version of `/ps-voice` or `/ps-style`'s questions from memory instead of reading their actual current skill files — they can change independently of this file.
- Never use Write, Edit, or Bash for anything this command touches, for the same reason `/ps-voice` and `/ps-style` don't.
