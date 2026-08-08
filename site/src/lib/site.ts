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

/** Recommended macOS/Linux install — served from getpysar.com via stage-install. */
export const installCommand =
  "curl -fsSL https://getpysar.com/install.sh | bash";

/**
 * Buttondown newsletter username (public). Override with NEXT_PUBLIC_BUTTONDOWN_USERNAME.
 * Embed docs: https://docs.buttondown.com/building-your-subscriber-base
 */
export const buttondownUsername =
  process.env.NEXT_PUBLIC_BUTTONDOWN_USERNAME ?? "malyshev";

/** Public embed-subscribe form action — no API key (dec-20260808-footer-subscribe-host-form-598699a0). */
export const buttondownSubscribeAction =
  `https://buttondown.com/api/emails/embed-subscribe/${buttondownUsername}` as const;
