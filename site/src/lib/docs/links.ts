import { DOCS_GITHUB_TREE, docHref } from "@/lib/docs/paths";

/**
 * Rewrite relative markdown/repo links to site routes or GitHub.
 * Absolute http(s) and mailto links pass through unchanged.
 */
export function rewriteDocHref(href: string | undefined): string | undefined {
  if (!href) {
    return href;
  }

  if (
    href.startsWith("http://") ||
    href.startsWith("https://") ||
    href.startsWith("mailto:") ||
    href.startsWith("#")
  ) {
    return href;
  }

  const [pathPart, hash = ""] = href.split("#");
  const suffix = hash ? `#${hash}` : "";
  const normalized = pathPart.replace(/^\.\//, "");

  if (normalized.startsWith("../")) {
    const repoPath = normalized.replace(/^\.\.\//, "");
    return `${DOCS_GITHUB_TREE}/${repoPath}${suffix}`;
  }

  const mdMatch = normalized.match(/^([a-z0-9-]+)\.md$/i);
  if (mdMatch) {
    return `${docHref(mdMatch[1])}${suffix}`;
  }

  if (normalized.startsWith("/")) {
    return `${normalized}${suffix}`;
  }

  return href;
}
