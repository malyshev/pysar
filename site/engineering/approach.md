# Motherhome approach — getpysar.com

**Contributor-only.** This file is under `site/engineering/` and must **never**
be exposed as a public URL on getpysar.com (not a Next route, not copied to
`public/`, not in sitemap).

**Bound:** `dec-20260808-da785647` — V1 isolated `site/` + GHA Direct Upload.  
Problem `prob-20260808-72fab46c` · portfolio `sol-20260808-8aa0e850`.

---

## Rendering mode (SSG) — locked

| Requirement | Decision |
|---|---|
| Next config | `output: 'export'` → write static files to `site/out/` |
| Images | `images.unoptimized: true` (required with static export) |
| Production runtime | Static HTML/CSS/JS on Cloudflare Pages only |
| Forbidden for MVP | Node server, Server Actions in prod, Route Handlers that need a server, Cloudflare Next SSR / OpenNext adapter |

This **is** SSG (static site generation / static export). It is not SSR.

---

## Goals

1. Ship a public motherhome at `https://getpysar.com`.
2. Keep marketing/SSG code out of the Go tool tree.
3. Use latest Next.js App Router, Tailwind, and shadcn/ui with a simple UI.
4. Static hosting only (no Node server in production) — SSG as above.
5. Reuse the verified static-export → Cloudflare Pages deploy pattern.
6. Honor the docs corpus decision: one repo-root `docs/` tree; site may ingest **that** tree only.

Out of scope for MVP: blog/MDX publishing pipeline, auth, server actions,
Cloudflare Next SSR adapter, design-system overhaul.

---

## What must not leak onto getpysar.com

| Content | Public site? |
|---|---|
| Motherhome marketing (`src/app` pages you author) | Yes |
| Repo-root `docs/*.md` (product user guide), if ingest is enabled | Yes — product docs only |
| `site/engineering/**` (this folder: deploy, Cloudflare, secrets setup) | **No — never** |
| Secret values, Account IDs, API tokens | **No** — GitHub Actions secrets / password manager only |
| Go internal packages, release matrix, CI token wiring | **No** |

Ingest path is **only** repo-root `../docs` (from `site/`). Never read
`site/engineering/` into page content.

---

## Repo boundary

```
pysar/                          # Go module root (unchanged product)
├── cmd/  internal/  go.mod     # tool
├── docs/                       # single user-docs corpus (journey + frontmatter)
├── .github/workflows/
│   ├── release.yml             # existing GoReleaser (untouched)
│   └── deploy-getpysar.yml     # NEW — site only
└── site/                       # NEW — Node package root
    ├── package.json
    ├── pnpm-lock.yaml
    ├── next.config.ts          # output: 'export'  ← SSG
    ├── components.json         # shadcn
    ├── public/                 # logo, og-default, favicons
    ├── out/                    # build artifact (gitignored)
    ├── src/
    │   ├── app/                # public routes only
    │   ├── components/         # shell + shadcn ui/
    │   └── lib/                # siteConfig, metadata, docs ingest, json-ld
    └── engineering/            # contributor ops — NEVER a public route source
```

**Rules:**

| Rule | Why |
|---|---|
| No root `package.json` for MVP | Keeps Go vs site ownership crisp (V1 vs V2) |
| CI `working-directory: site` (or `cd site`) | Isolated install/build |
| Path-filter deploy on `site/**` + `docs/**` + workflow | Product docs changes rebuild doc routes |
| Do not copy repo-root `docs/` into `site/` permanently | Single corpus invariant |
| Never route or bundle `site/engineering/**` | Keeps deploy/secrets procedures off the CDN |

---

## Stack (pins)

| Layer | Choice |
|---|---|
| Runtime target | Static HTML/CSS/JS on Cloudflare Pages |
| Framework | Next.js App Router (current latest major — verified pattern uses 16.x) |
| UI | React 19, Tailwind CSS 4, shadcn/ui (RSC + lucide) |
| Package manager | pnpm |
| Node in CI | 22 |
| Deploy CLI | `wrangler` as `site` devDependency; `pnpm exec wrangler pages deploy` |
| Config | `output: 'export'`; `images.unoptimized: true` (required for static export) |

Scaffold notes:

- `pnpm create next-app` (or equivalent) **inside `site/`** with App Router, TS, Tailwind, `src/`.
- Init shadcn into `site/` (`components.json` aliases under `@/`).
- Prefer a small set of primitives first: `Button`, `Card`, typography via Tailwind — no large component dump.

---

## Information architecture (MVP)

| URL | Role |
|---|---|
| `/` | Motherhome — brand, one H1, short promise, primary CTA (Install / Get started) |
| `/docs` | Docs index — generated from repo `docs/index.md` |
| `/docs/<slug>` | Journey pages from `docs/*.md` |
| `/sitemap.xml` | All public routes |
| `/robots.txt` | Allow `/`, point at sitemap |

Optional soon after: dedicated `/install` marketing page that deep-links into docs
(only if home CTA is not enough). Reserved path list lives in code next to
docs-slug validation.

---

## Docs ingest (product corpus only)

- Read Markdown + YAML frontmatter from repo-root `../docs` at **build time**
  (Node `fs` from `site` lib — path must resolve to **repo-root `docs/`**, not
  `site/engineering/`).
- Use existing keys: `title`, `slug`, `nav_order`, `section` (already on the corpus).
- Render to static routes under `/docs/...`.
- Link rewrite: relative links between docs pages must resolve to site routes.
- **Never** maintain a second prose copy under `site/content/docs`.
- **Never** ingest or expose `site/engineering/**`.

---

## SEO spine (baseline)

Port these practices into `site/src/lib/` (names can match intent, not an external repo):

1. **`siteConfig`** — `name`, `url` from `NEXT_PUBLIC_SITE_URL` (default `https://getpysar.com`), short description.
2. **Root `metadata`** — `metadataBase`, title template (`%s · Pysar`), default description, default OG/Twitter card + `og-default` image + alt.
3. **Per-page metadata helpers** — absolute canonical, OG/Twitter mirrors; home uses an **absolute** title (avoid template doubling).
4. **`robots.ts` / `sitemap.ts`** — `export const dynamic = "force-static"`; sitemap lists home + every docs slug (+ any marketing routes).
5. **JSON-LD** — home `@graph` with `WebSite` + `Organization` (logo in `public/`); docs pages can add `TechArticle` / `WebPage` later.
6. **Canonical host** — CI sets `NEXT_PUBLIC_SITE_URL=https://getpysar.com` (never `*.pages.dev`, never localhost in prod builds).
7. **Cloudflare** — Bulk Redirect `*.pages.dev` → apex; pick apex vs `www` and 301 the other.
8. **Home H1** — exactly one visible H1; SERP title ~50–60 chars; meta description ~140–160 chars; distinct from bare brand.

CI verify after build:

```bash
test -d out
test -f out/index.html
test -f out/sitemap.xml
test -f out/robots.txt
```

---

## UI / design (start simple)

- One shell: header (wordmark + nav) / main / footer.
- Max-width content column; Tailwind spacing; shadcn `Button` for CTAs.
- No dashboard chrome, no card grids in the hero, no purple-glow defaults.
- Brand first on `/`: product name as the dominant signal, one headline, one short
  supporting sentence, one CTA group.

Copy must describe **shipped** behavior only (same rule as `docs/`).

---

## Relationship to existing decisions

| Artifact | Implication |
|---|---|
| `dec-20260808-journey-docs-corpus-e8b5483b` | `docs/` stays the only user-docs corpus; site ingests; SSG choice was deferred — this plan is that choice |
| GoReleaser release workflow | Unrelated path filters; do not couple site deploy to `v*` tags |
| Portfolio `sol-20260808-8aa0e850` | V1 recommended; V2–V5 documented as rivals |

---

## Implementation phases (after `/h-decide`)

| Phase | Outcome |
|---|---|
| 0 | Human setup: domain DNS, Pages project, GitHub secrets ([human-setup.md](./human-setup.md)) — private ops, not a site page |
| 1 | Scaffold `site/` Next + Tailwind + shadcn; `output: 'export'` (SSG); shell + home |
| 2 | SEO helpers + robots/sitemap + JSON-LD + default OG/logo assets |
| 3 | Product docs ingest from repo-root `docs/` → `/docs` (govern detail later; backbone first) |
| 4 | Workflow deploy ([deploy-pipeline.md](./deploy-pipeline.md)); first production push |
| 5 | Search Console property for `https://getpysar.com`; submit sitemap |

Decision already bound: `dec-20260808-da785647`.  
Operator note (`note-20260808-b80b85fa`): repo-root `docs/` is public tool docs for the site; ingest IA governed later — **current focus is site backbone**.
