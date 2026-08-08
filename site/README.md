# Pysar motherhome (`site/`)

Public product site for **[getpysar.com](https://getpysar.com)** — Next.js
**static export (SSG)**, Tailwind, shadcn/ui, Cloudflare Pages.

Bound by `dec-20260808-da785647`.

## Commands

```bash
cd site
pnpm install
pnpm dev          # App Router dev server
pnpm build        # static export → out/
pnpm preview      # build + serve out/ on :4321
pnpm check        # lint + test
```

Set `NEXT_PUBLIC_SITE_URL` for non-default origins (CI uses
`https://getpysar.com`).

## Layout

| Path | Role |
|---|---|
| `src/app` | Public routes (SSG) |
| `src/components` | Shell + shadcn `ui/` |
| `src/lib` | siteConfig, metadata, JSON-LD |
| `public/` | logo, OG defaults |
| `engineering/` | Contributor ops — **never** a public route |

Repo-root [`docs/`](../docs/) is the public product docs corpus (ingest later).
Docs nav currently links to GitHub until `/docs` routes land.

## Deploy

See [engineering/deploy-pipeline.md](./engineering/deploy-pipeline.md) and
[engineering/human-setup.md](./engineering/human-setup.md).
Workflow: `.github/workflows/deploy-getpysar.yml`.
