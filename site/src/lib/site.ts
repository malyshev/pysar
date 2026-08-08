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
    "From a half-formed idea or a messy draft to a piece that sounds like you — shaped with your judgment still in the chair.",
  eyebrow: "Author-directed editorial engine",
} as const;
