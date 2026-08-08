export const siteConfig = {
  name: "Pysar",
  url: process.env.NEXT_PUBLIC_SITE_URL ?? "https://getpysar.com",
  description:
    "An author-directed editorial engine for your writing projects.",
  github: "https://github.com/malyshev/pysar",
} as const;

/** Square mark-on-ink asset for JSON-LD / social. Wordmarks live at /logo.svg and /logo-dark.svg. */
export const siteLogo = {
  path: "/logo.png",
  width: 512,
  height: 512,
} as const;

/** Homepage SERP + social copy. Does not replace siteConfig.description elsewhere. */
export const homePageSeo = {
  h1: "Bring your take. Ship a piece you trust.",
  title: "Pysar — author-directed editorial engine for writers",
  description:
    "Shape an idea or rough draft into a piece you trust. CLI plus MCP and ps-* skills for Claude Code and Cursor — without posting on your behalf.",
  eyebrow: "Author-directed editorial engine",
} as const;
