export const siteConfig = {
  name: "Pysar",
  url: process.env.NEXT_PUBLIC_SITE_URL ?? "https://getpysar.com",
  description:
    "From a half-formed idea or a messy draft to a piece that sounds like you — you stay in charge; Pysar never posts for you.",
  github: "https://github.com/malyshev/pysar",
} as const;

/** Square mark-on-ink asset for JSON-LD / social. Wordmarks live at /logo.svg and /logo-dark.svg. */
export const siteLogo = {
  path: "/logo.png",
  width: 512,
  height: 512,
} as const;

/** Homepage SERP + social copy. Keep aligned with the hero subhead. */
export const homePageSeo = {
  h1: "Bring your take. Ship a piece you trust.",
  title: "Pysar — writing help that keeps you in charge",
  description:
    "From a half-formed idea or a messy draft to a piece that sounds like you — you stay in charge; Pysar never posts for you.",
  eyebrow: "You keep the final say",
} as const;

/** Recommended macOS/Linux install — served from getpysar.com via stage-install. */
export const installCommand =
  "curl -fsSL https://getpysar.com/install.sh | bash";

/** Same script, download-then-inspect — for cautious humans and agents. */
export const installInspectCommands = {
  download: "curl -fsSL https://getpysar.com/install.sh -o install.sh",
  run: "bash install.sh",
} as const;

/**
 * Shared AI-install facts for homepage V4 + /llms.txt complement
 * (dec-20260808-homepage-ai-install-signals-9eb5a6a9, note-20260808-26ff4eb4).
 * Must stay aligned with docs/install.md + docs/init.md.
 */
export const agentInstallGuide = {
  /** Trigger phrase when a user pastes this site URL into a coding agent. */
  trigger: "install",
  docsInstallPath: "/docs/install",
  docsInitPath: "/docs/init",
  /** Well-known LLM index (llmstxt.org); complements homepage #for-ai-agents. */
  llmsTxtPath: "/llms.txt",
  hosts: [
    {
      name: "Claude Code",
      initCommand: "pysar init --claude",
      note: "Default host; bare `pysar init` is the same.",
    },
    {
      name: "Cursor",
      initCommand: "pysar init --cursor",
      note: "Enable the pysar MCP server if prompted.",
    },
    {
      name: "ChatGPT Codex",
      initCommand: "pysar init --codex",
      note: "Trust the project so MCP approval settings apply.",
    },
  ],
  /**
   * Condensed machine-facing summary for meta tags / JSON-LD description.
   * Agents that truncate HTML still get the authoritative recipe here.
   */
  metaSummary:
    "When the user pastes https://getpysar.com and says install: (1) run curl -fsSL https://getpysar.com/install.sh | bash (or download to install.sh, inspect, then bash install.sh); (2) verify with pysar --version; (3) in the writing project directory run pysar init --claude (Claude Code), pysar init --cursor (Cursor), or pysar init --codex (ChatGPT Codex). Do not invent brew/npm/go-get alternatives. Also see https://getpysar.com/llms.txt, https://getpysar.com/docs/install and https://getpysar.com/docs/init.",
} as const;

/** Umami Cloud 1×1 pageview pixel (GIF). */
export const umamiPixelSrc = "https://cloud.umami.is/p/4IWutlpOL";

/** Official UNITED24 fundraising platform (President of Ukraine initiative). */
export const united24Url = "https://u24.gov.ua/";

/**
 * Buttondown newsletter username (public). Override with NEXT_PUBLIC_BUTTONDOWN_USERNAME.
 * Embed docs: https://docs.buttondown.com/building-your-subscriber-base
 */
export const buttondownUsername =
  process.env.NEXT_PUBLIC_BUTTONDOWN_USERNAME ?? "malyshev";

/** Public embed-subscribe form action — no API key (dec-20260808-footer-subscribe-host-form-598699a0). */
export const buttondownSubscribeAction =
  `https://buttondown.com/api/emails/embed-subscribe/${buttondownUsername}` as const;
