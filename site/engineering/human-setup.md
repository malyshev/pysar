# Human setup — getpysar.com

**Contributor-only. Not a public site page.** Lives under `site/engineering/`
and must never appear on getpysar.com. Do not put real API tokens, Account IDs,
or secret values into this file or any committed path — only placeholders and
procedure. Values go in GitHub Actions secrets / your password manager.

Operator checklist (not codebase). Complete once before or alongside scaffold.

When done, engineering needs: three GitHub secrets + confirmation that the
Pages project and domain exist.

---

## Who owns what

| You (human) | Engineering / agent |
|---|---|
| Domain `getpysar.com`, Cloudflare account | `site/` Next app |
| Pages project (Direct Upload) | `.github/workflows/deploy-getpysar.yml` |
| API token + GitHub secrets | Docs ingest from `docs/` |
| DNS / custom domain / redirects | SEO helpers, sitemap, JSON-LD |

---

## Before you start

- [ ] Domain **`getpysar.com`** registered (or you control DNS)
- [ ] Cloudflare account (free tier is enough for Pages)
- [ ] GitHub Actions enabled on the Pysar repo
- [ ] Node 22 + pnpm available locally for preview after scaffold

---

## Step 1 — Cloudflare account ID

Dashboard → Workers & Pages (or Overview) → copy **Account ID**.

```
CLOUDFLARE_ACCOUNT_ID = ________________________________
```

---

## Step 2 — Domain on Cloudflare

Best path: DNS managed by Cloudflare.

1. Add site `getpysar.com` (or register there).
2. Point registrar nameservers at Cloudflare if needed; wait until **Active**.
3. Decide apex policy (recommend **apex only**; redirect `www` → apex).

---

## Step 3 — Pages project (Direct Upload)

1. Workers & Pages → Create → **Pages** → **Direct Upload** / upload assets.
2. Project name suggestion: `getpysar` (must match the GitHub secret).
3. Skip manual first upload if you want — CI will deploy.
4. **Do not** select Next.js SSR / OpenNext for MVP.

```
CLOUDFLARE_PROJECT_NAME = getpysar
```

---

## Step 4 — Custom domain (required — DNS alone is not enough)

Adding zone DNS records is **not** the same as attaching the hostname to the
Pages project. Until the domain is under **Pages → getpysar → Custom domains**,
proxied DNS often yields **522** / unreachable apex while `*.pages.dev` works.

1. Pages project **getpysar** → **Custom domains** → set up `getpysar.com`
   (and optionally `www.getpysar.com`).
2. Wait until status is **Active** (certificate issued). Zone must already be
   on Cloudflare (`status: active` nameservers).
3. Optional: 301 `www` → apex in Cloudflare Rules.
4. Bulk Redirect: `getpysar.pages.dev` → `https://getpysar.com`.

---

## Step 5 — API token

My Profile → API Tokens → Create Custom Token:

| Field | Value |
|---|---|
| Permission | Account → Cloudflare Pages → **Edit** |
| Account resources | Include your account |

Copy once:

```
CLOUDFLARE_API_TOKEN = ********************************
```

---

## Step 6 — GitHub secrets

Repo → Settings → Secrets and variables → Actions:

| Secret | Value |
|---|---|
| `CLOUDFLARE_API_TOKEN` | from Step 5 |
| `CLOUDFLARE_ACCOUNT_ID` | from Step 1 |
| `CLOUDFLARE_PROJECT_NAME` | from Step 3 |

Do not commit these or paste them into chat logs.

---

## Step 7 — Footer subscribe (Buttondown)

Footer Subscribe uses Buttondown’s public embed form POST
(`dec-20260808-footer-subscribe-host-form-598699a0`) — no API key in the site.

| Item | Value |
|---|---|
| Username | `malyshev` (override: `NEXT_PUBLIC_BUTTONDOWN_USERNAME`) |
| Action URL | `https://buttondown.com/api/emails/embed-subscribe/malyshev` |
| Tag on signup | `getpysar` |
| Dashboard | [buttondown.com/settings/embedding](https://buttondown.com/settings/embedding) |

Submit leaves getpysar.com and lands on Buttondown’s confirm / double-opt-in page
(honest host acknowledgment). Enable double opt-in in Buttondown if you want it.

- [ ] Submit a real test address from the live footer; confirm pending/active in Buttondown

---

## Step 8 — Analytics (Umami Cloud)

Root layout loads a 1×1 Umami **pageview** pixel:

`https://cloud.umami.is/p/4IWutlpOL`

Repo-root `install.sh` loads a separate 1×1 Umami **install** pixel after a
successful binary install (not a pageview):

`https://cloud.umami.is/p/NDhIZ7E6F`

- [ ] Confirm pageview hits appear after a production visit
- [ ] Confirm install pixel hits appear after a real `install.sh` run (opt-out:
      `PYSAR_NO_TELEMETRY=1`)

---

## Step 9 — After first deploy

- [ ] `https://getpysar.com` serves the motherhome
- [ ] `https://getpysar.com/sitemap.xml` and `/robots.txt` exist
- [ ] View-source: one H1, canonical apex, OG tags
- [ ] Google Search Console property for `https://getpysar.com`; submit sitemap
- [ ] Confirm `*.pages.dev` redirects to apex

---

## Hand-off checklist

Provide to implementation (or confirm in-thread):

```
CLOUDFLARE_PROJECT_NAME: ________
Domain Active on CF:     yes/no
GitHub secrets set:      yes/no
Default git branch:      master / main
```
