---
id: dec-20260808-footer-subscribe-host-form-598699a0
kind: DecisionRecord
version: 6
status: active
title: Newsletter host form POST (Buttondown / ConvertKit)
context: site
mode: standard
valid_until: 2026-11-08
created_at: 2026-08-08T15:36:27Z
updated_at: 2026-08-08T16:26:29Z
links:
  - ref: prob-20260808-c0712044
    type: based_on
  - ref: sol-20260808-c5d18d99
    type: based_on
---

# Newsletter host form POST (Buttondown / ConvertKit)

## 1. Problem Frame

**Signal:** Footer has a Boxsi-style Subscribe UI but Next output:export + Cloudflare Pages Direct Upload has no server runtime in the Next app. Submits currently clear the field and store nothing. Need a way to collect emails that fits the locked static motherhome shape (dec-20260808-da785647).

**Constraints:**
- Keep Next output:export (no Next API routes / SSR)
- Do not invent fake success when nothing was stored
- No secrets in the static client bundle beyond what the chosen provider documents as public (form keys / list tokens)
- Prefer reversible ops — can change provider without rewriting the whole site

**Acceptance:** A visitor can submit an email from the footer; the address is durably captured in a chosen list/store; double-opt-in or clear consent path exists if the provider requires it; failure is visible (not silent success); static export + Pages deploy path remains intact.

## 2. Decision

**Selected:** Newsletter host form POST (Buttondown / ConvertKit)

**Selection policy:** Prefer the smallest change that keeps Next output:export + Pages Direct Upload intact (dec-20260808-da785647), captures email durably with built-in list/consent tooling, and defers Cloudflare Functions until a concrete need (secrets, same-origin, Turnstile) appears. Optimize for time-to-working capture and ops simplicity; treat bot/spam risk as observation until measured.

**Why selected:** V1 wires the existing Boxsi footer form to a newsletter host’s public form/API without adding Pages Functions or leaving static export. List storage, double-opt-in, and unsubscribe live in the host. Matches motherhome constraints and is reversible by swapping the endpoint or removing the POST.


**Invariants:**
- Next remains output:export — no Next server/API routes for subscribe
- Never show success unless the host acknowledged the signup (or pending confirm)
- No private API tokens in the static client bundle — only provider-documented public form endpoints/keys
- Motherhome deploy stays Cloudflare Pages Direct Upload of site/out

**Pre-conditions:**
- [ ] Operator has (or creates) a newsletter host account with a public subscribe form/API
- [ ] Form endpoint URL (and any public form key) available for NEXT_PUBLIC_ or build-time config

**Post-conditions:**
- [ ] FooterSubscribe POSTs to the host and surfaces failure visibly
- [ ] Docs/engineering note names the provider and where the list lives

**Admissibility:**
- NOT: Silent clear-on-submit with no storage
- NOT: Embedding third-party scripts that break Boxsi footer layout without operator opt-in
- NOT: Putting private Resend/API secrets in client JS

## 3. Rationale

**Counterargument:** A Cloudflare Pages Function (V2) keeps secrets off the client and enables Turnstile/same-origin control from day one; choosing V1 may force a second migration if spam or key exposure becomes real before the list proves useful.

**Selected variant weakest link:** Public form endpoint is exposed to bots/spam and any required form key is client-visible; quality depends on the host’s abuse controls and on correctly surfacing HTTP success/failure (no silent fake success).

**Rejected alternatives:**
| Variant | Verdict | Reason |
|---------|---------|--------|
| Newsletter host form POST (Buttondown / ConvertKit) | **Selected** | V1 wires the existing Boxsi footer form to a newsletter h... |
| Cloudflare Pages Function + Resend/Audience or D1 | Rejected | Correct when secrets or bot protection are load-bearing; today that adds functions/, secrets, and deploy complexity before we know list demand. V1 is the stepping stone; escalate to V2 if spam or secret needs appear. |
| Embed provider widget / hosted form page | Rejected | Breaks the Boxsi footer visual lock or hops off-site for a primary CTA; worse fit than a quiet form POST that keeps the template cell. |
| No email list — announce via GitHub Releases / RSS only | Rejected | Operator kept the Subscribe affordance and asked how to make subscription work; refusing capture leaves the UI dishonest. |

**Predictions:**
| Claim | Observable | Threshold |
|-------|------------|-----------|
| Footer submit with a valid email returns a non-error response from the chosen host and the address appears in that host’s list (or pending double-opt-in) within one minute | Manual submit on getpysar.com + host dashboard subscriber/pending count | 1 successful end-to-end capture on production |
| Static export and Pages Direct Upload path remain unchanged (no Next API routes required) | site/next.config.ts still output:export; deploy-getpysar.yml still wrangler pages deploy out | Both true after subscribe wiring lands |

## 4. Consequences

**Rollback plan:**
Triggers:
- Provider form rejects browser POSTs (CORS/auth) with no workable public form path
- Spam/bot signups dominate legitimate ones after basic host controls
- Operator decides email list is not wanted
Steps:
1. Point FooterSubscribe back to no-op or remove Subscribe cell
2. Cancel/delete the host list or rotate the public form key
3. If already migrated to V2, disable the Function route and revert footer fetch URL
Blast radius: site/src/components/footer-subscribe.tsx, public env for form endpoint/key, newsletter host account only — no Go CLI/MCP impact

**Refresh triggers:**
- Spam volume makes public form unusable
- Provider changes or retires public form API
- Need for Turnstile/secret-backed subscribe (reopen V2)

**Affected files:** site/src/components/footer-subscribe.tsx, site/src/lib/site.ts, site/.env.example, site/engineering/human-setup.md

## Impact Measurement (2026-08-08)

**Verdict:** partial

**Findings:**
/h-verify: Prediction 2 holds — static export + Pages Direct Upload unchanged; native Buttondown form POST on production. Prediction 1 incomplete — live form correctly targets Buttondown embed-subscribe/malyshev, but successful list/pending capture in Buttondown dashboard was not confirmed this session.

**Criteria met:**
- [x] Static export and Pages Direct Upload path remain unchanged
- [x] Footer POSTs to public Buttondown embed endpoint (no silent fake success)
- [x] Engineering note names Buttondown

**Criteria NOT met:**
- [ ] 1 successful end-to-end capture on production confirmed in host list/pending

**Measurements:**
- next.config output:export
- wrangler pages deploy out
- live form action buttondown embed-subscribe/malyshev method=post
- E2E Buttondown dashboard: not observed

## Impact Measurement (2026-08-08)

**Verdict:** accepted

**Findings:**
Both predictions hold. Static export + Pages path unchanged. Operator confirmed live footer → Buttondown success acknowledgment (subscribed to Serhii Malyshev newsletter).

**Criteria met:**
- [x] 1 successful end-to-end capture on production
- [x] Static export and Pages Direct Upload path remain unchanged
- [x] Footer POSTs to Buttondown; failure/success is host-visible

**Measurements:**
- live form action embed-subscribe/malyshev
- operator Buttondown success modal confirmed
- output:export + wrangler pages deploy out
