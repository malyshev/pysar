import path from "node:path";

/** Repo-root `docs/` — never `site/engineering/`. */
export function getDocsDir(): string {
  return path.resolve(process.cwd(), "..", "docs");
}

export function docHref(slug: string): string {
  return slug === "index" ? "/docs" : `/docs/${slug}`;
}

/** GitHub blob URL for links that escape the docs tree (e.g. ../README.md). */
export const DOCS_GITHUB_TREE =
  "https://github.com/malyshev/pysar/blob/master";
