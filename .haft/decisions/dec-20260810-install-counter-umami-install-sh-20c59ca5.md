---
id: dec-20260810-install-counter-umami-install-sh-20c59ca5
kind: DecisionRecord
version: 1
status: active
title: Umami custom event / pixel from install.sh
context: site
mode: standard
valid_until: 2026-11-10T00:00:00Z
created_at: 2026-08-10T15:30:59Z
updated_at: 2026-08-10T15:30:59Z
links:
  - ref: prob-20260810-5af07d47
    type: based_on
  - ref: sol-20260810-4f52275a
    type: based_on
---

# Umami custom event / pixel from install.sh

## 1. Problem Frame

**Signal:** Operator wants a count of installs visible or usable with the current motherhome: Next output:export → Cloudflare Pages Direct Upload (dec-20260808-da785647), static /install.sh that curls GitHub release assets, Umami pageview pixel already on the site. Pure static has no server-side counter; curl|bash does not execute the Umami GIF. Need a durable install metric that fits Cloudflare + this architecture without silently redefining 'install' as homepage hits.

**Constraints:**
- Keep Next output:export + GHA Pages Direct Upload as the motherhome shape unless a narrow CF edge piece is justified
- Do not treat Umami homepage pageviews as installs
- Metric definition must be honest about over/undercount (curl retries, CI mirrors, go install, Windows zip)
- No PII requirement; opt-out or privacy-light beacon preferred if phone-home
- install.sh must keep working offline-ish after download; beacon failure must not fail install

**Acceptance:** Operator can read a numeric install metric (dashboard and/or public site) whose definition is explicit (script fetches vs successful binary install vs release asset downloads); the static Pages Direct Upload path for the marketing site remains intact; a smoke can increment the metric once by a documented path.

## 2. Decision

**Selected:** Umami custom event / pixel from install.sh

**Selection policy:** Prefer reusing the already-shipped Umami analytics plane over new Cloudflare KV/D1 state or Go CLI telemetry. Optimize for ops simplicity and trust surface outside the pysar binary; accept that the metric is successful install.sh completion (with optional OS/arch from the script), not activated use from go install / Windows zip. Prior withdrawn V5 (CLI phone-home) stays deprecated — do not re-bind CLI telemetry.

**Why selected:** Umami is already on getpysar.com. A best-effort ping from install.sh after the binary lands keeps the counter out of the Go CLI, avoids a new CF write product for MVP, and still beats bare /install.sh GET counts or unlabeled GitHub download_count as a success signal for the documented curl path. Operator rejected binding V5 after reconsidering; V4 is the lightest durable path on the script channel.


**Invariants:**
- Ping failure must never fail install.sh
- No PII / hostname required; OS/arch and optional release tag are enough
- Document what is sent (privacy-light)
- Do not treat Umami homepage pageviews as this counter
- Motherhome stays static Pages export; Umami remains the analytics vendor for this metric unless superseded
- Public wording must match definition (script-success installs, not all channels)

**Pre-conditions:**
- [ ] Umami Cloud already used on site (umamiPixelSrc)
- [ ] install.sh is the documented macOS/Linux path

**Post-conditions:**
- [ ] Successful install.sh emits one Umami event (name/schema documented)
- [ ] Operator can read a count in Umami
- [ ] Opt-out or skip documented if we add an env flag

**Admissibility:**
- NOT: Failing install when Umami is down
- NOT: Embedding a long-lived private Umami API secret in the public install.sh if a public collect path exists
- NOT: Claiming the counter covers go install / Windows without those channels

## 3. Rationale

**Counterargument:** V3 (CF Worker + KV) gives a first-party sink under our control with clearer INCR semantics; V1 needs zero script change. Umami vendor lock and awkward server-side/collect-from-curl may force a second migration to V3 anyway.

**Selected variant weakest link:** curl→Umami collect is weaker than browser pixels: may need a collect URL/API shape that works from bash; if auth tokens are required in the script they become public; go install and Windows zip remain uncounted.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Umami custom event / pixel from install.sh | **Selected** | Umami is already on getpysar.com. A best-effort ping from... |
| GitHub Releases download_count sum (no new infra) | Rejected | Zero infra but measures asset fetches not script success; wrong label if shown as installs. |
| Cloudflare log/analytics on GET /install.sh | Rejected | GET of script ≠ run to completion; inspect-only noise. |
| Success beacon from install.sh to CF Worker + KV/D1 | Rejected | Stronger first-party sink but adds Functions/KV before Umami reuse is exhausted; escalate later if Umami collect fails dogfood. |
| Phone-home from pysar binary (first-run / init) | Rejected | Previously decided then withdrawn (deprecated); operator chose not to instrument the CLI for this round. |

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| After a successful install.sh run, Umami records an install (or equivalently named) event without the script exiting non-zero when Umami is unreachable | Manual smoke: run install path with network; then block Umami host and re-run | Both runs leave pysar on PATH / exit 0; Umami shows +1 event only on the reachable run |
| Event payload can include OS (and arch if trivial) from the script environment | Umami event properties on a smoke install | At least os (and arch if sent) present on the recorded event |

## 4. Consequences

**Rollback plan:**
Triggers:
- Umami collect from install.sh cannot be made to work without embedding a secret that is unacceptable
- Beacon causes install failures or support load
- Operator decides public/ops metric needs go-install coverage (reopen V5)
Steps:
1. Remove or no-op the ping in install.sh
2. Stop labeling Umami events as installs on any UI
3. Optionally move to V3 CF sink or revive V5 via new decide
Blast radius: install.sh; Umami event schema/dashboard; optional site copy; no Go CLI change required

**Refresh triggers:**
- Umami collect API change breaks bash ping
- Need go-install / Windows in the same counter
- Spam/abuse of the public collect endpoint

**Affected files:** install.sh, docs/install.md, site/src/lib/site.ts, site/engineering
