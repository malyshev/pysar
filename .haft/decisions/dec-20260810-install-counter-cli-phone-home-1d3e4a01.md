---
id: dec-20260810-install-counter-cli-phone-home-1d3e4a01
kind: DecisionRecord
version: 5
status: deprecated
title: Phone-home from pysar binary (first-run / init)
context: site
mode: standard
valid_until: 2026-11-10T00:00:00Z
created_at: 2026-08-10T15:17:39Z
updated_at: 2026-08-10T15:23:58Z
links:
  - ref: prob-20260810-5af07d47
    type: based_on
  - ref: sol-20260810-4f52275a
    type: based_on
---

# Phone-home from pysar binary (first-run / init)

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

**Selected:** Phone-home from pysar binary (first-run / init)

**Selection policy:** Prefer the metric that measures activated product use across all install channels (curl script, go install, Windows zip), not website architecture proxies (script GETs, release download_count, Umami pageviews). Accept a Cloudflare (or equivalent) aggregate sink as Shell; keep the count trigger in the Go CLI Core. Optimize for signal quality of real installs; treat public vanity display and edge purity as secondary. Operator refinement: each ping carries pysar version and OS (and arch as needed) so the counter is segmentable, not a bare increment.

**Why selected:** V5 alone covers every install path that ends in a running pysar binary. Website-only counters (V1–V4) miss go install / Windows and conflate fetch with success. Phone-home on first-run or successful init, with version + OS in the payload, gives an honest activation counter and useful breakdown without pretending curl of install.sh equals install.


**Invariants:**
- Phone-home failure must never fail install, init, or normal CLI use
- Payload includes at least pysar version and OS (arch allowed); no requirement to send PII or machine hostname
- Default behavior must document what is sent and how to opt out
- Motherhome stays Next output:export + Pages Direct Upload; any CF write path is a narrow edge sink, not Next SSR
- Public UI must not label the metric as pure 'installs' if it is 'activations' / first-run pings — wording matches definition
- Static /install.sh remains usable without the beacon

**Pre-conditions:**
- [ ] Problem and portfolio framed (prob-20260810-5af07d47, sol-20260810-4f52275a)
- [ ] Cloudflare or other aggregate sink can be provisioned for increments

**Post-conditions:**
- [ ] CLI emits version+OS on the chosen first-run/init trigger
- [ ] Aggregate count readable by operator (and optionally site)
- [ ] Opt-out path documented

**Admissibility:**
- NOT: Silent telemetry with no docs or opt-out
- NOT: Failing user commands when the counter sink is unreachable
- NOT: Equating Umami pageviews or bare /install.sh GETs with this counter

## 3. Rationale

**Counterargument:** V1 (GitHub download_count) or V3 (install.sh beacon) ship faster with zero or minimal Go CLI change and keep trust surface on the website/script. Instrumenting pysar itself expands blast radius into the product binary and may delay the vanity counter the operator wanted on the motherhome.

**Selected variant weakest link:** User trust and opt-out: any default phone-home can be read as telemetry creep; if blocked by firewall or disabled, counts understate installs. Forgeable anonymous pings without proof-of-possession remain possible at the sink.

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Phone-home from pysar binary (first-run / init) | **Selected** | V5 alone covers every install path that ends in a running... |
| GitHub Releases download_count sum (no new infra) | Rejected | Fast stepping stone but measures asset fetches, not activated installs; omits go install; cannot attach accurate OS/version of the running binary. |
| Cloudflare log/analytics on GET /install.sh | Rejected | Counts script download only; inspect-only and bots; no version of installed binary. |
| Success beacon from install.sh to CF Worker + KV/D1 | Rejected | Closer to script success but misses non-script installs and still lives outside the binary that knows its own version. |
| Umami custom event / pixel from install.sh | Rejected | Same channel gap as V3; awkward from curl; weaker fit for version/OS from the binary. |

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| A successful first-run or init path emits one aggregate ping including pysar version and OS without failing the user command when the sink is down | Integration or manual smoke: run pysar init (or documented first-run) with sink up and with sink down | Command exit 0 in both cases; sink shows +1 with version+OS fields when up |
| Non-script install channels (at least go install or Windows zip path) can increment the same counter | Smoke without using getpysar.com/install.sh | At least one non-script path produces a counted ping with version+OS |

## 4. Consequences

**Rollback plan:**
Triggers:
- Operator or users reject default phone-home as unacceptable trust cost
- Sink or CLI ping causes install/init failures or measurable support load
- Counts remain near-zero after two weeks of known installs (broken path)
Steps:
1. Disable or make opt-in the CLI ping (build flag / env / config)
2. Remove or hide public counter UI if any
3. Tear down or idle the CF Worker/KV/D1 sink
4. Fall back to V1 download_count for ops if needed
Blast radius: Go CLI (cmd/pysar or internal telemetry), optional Cloudflare Worker/KV/D1 + secrets, site display if wired; install.sh unchanged unless we add a note

**Refresh triggers:**
- Privacy/regulatory complaint about CLI telemetry
- CF sink cost or abuse (spam pings)
- Operator wants public homepage counter wired
- Switch of first-run vs init trigger after dogfood

**Affected files:** cmd/pysar, internal/, install.sh, site/src, docs/, .github/workflows

## Deprecated (2026-08-10)

**Reason:** Operator withdrew acceptance before implementation: needs more thinking time on install-counter approach (CLI phone-home vs install.sh beacon vs other). Decision must not bind further work.

## Deprecated (2026-08-10)

**Reason:** Re-confirm withdraw after reopen: V5 is not accepted. Keep as historical record only; no implementation authority. Operator thinking; original portfolio sol-20260810-4f52275a + prob-20260810-5af07d47 remain the choice space.
