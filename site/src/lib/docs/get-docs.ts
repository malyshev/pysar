import fs from "node:fs";
import path from "node:path";
import matter from "gray-matter";
import { docHref, getDocsDir } from "@/lib/docs/paths";
import type { DocFrontmatter, DocPage } from "@/lib/docs/types";

function parseFrontmatter(data: Record<string, unknown>): DocFrontmatter {
  const title = typeof data.title === "string" ? data.title.trim() : "";
  const slug = typeof data.slug === "string" ? data.slug.trim() : "";
  const section = typeof data.section === "string" ? data.section.trim() : "";
  const nav_order =
    typeof data.nav_order === "number"
      ? data.nav_order
      : Number.parseInt(String(data.nav_order ?? ""), 10);

  if (!title || !slug || !section || Number.isNaN(nav_order)) {
    throw new Error(
      `Invalid docs frontmatter (need title, slug, nav_order, section): ${JSON.stringify(data)}`,
    );
  }

  return { title, slug, nav_order, section };
}

function loadDocFile(filePath: string): DocPage {
  const raw = fs.readFileSync(filePath, "utf8");
  const { data, content } = matter(raw);
  const fm = parseFrontmatter(data as Record<string, unknown>);

  return {
    ...fm,
    filePath,
    body: content.trimStart(),
    href: docHref(fm.slug),
  };
}

/** All journey docs from repo-root docs/, sorted by nav_order. */
export function getAllDocs(): DocPage[] {
  const docsDir = getDocsDir();
  if (!fs.existsSync(docsDir)) {
    throw new Error(`Docs corpus missing at ${docsDir}`);
  }

  const files = fs
    .readdirSync(docsDir)
    .filter((name) => name.endsWith(".md"))
    .map((name) => path.join(docsDir, name));

  const docs = files.map(loadDocFile);
  const slugs = new Set(docs.map((d) => d.slug));
  if (slugs.size !== docs.length) {
    throw new Error("Duplicate docs slug in corpus");
  }
  if (!slugs.has("index")) {
    throw new Error("docs/index.md (slug: index) is required");
  }

  return docs.sort((a, b) => a.nav_order - b.nav_order || a.slug.localeCompare(b.slug));
}

export function getDocBySlug(slug: string): DocPage | undefined {
  return getAllDocs().find((doc) => doc.slug === slug);
}

export function getIndexDoc(): DocPage {
  const index = getDocBySlug("index");
  if (!index) {
    throw new Error("docs index missing");
  }
  return index;
}

/** Non-index pages for `/docs/[slug]`. */
export function getDocSlugs(): string[] {
  return getAllDocs()
    .filter((doc) => doc.slug !== "index")
    .map((doc) => doc.slug);
}
