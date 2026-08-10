export const siteConfig = {
  name: "Pysar",
  url: process.env.NEXT_PUBLIC_SITE_URL ?? "https://getpysar.com",
  description:
    "Real editorial help, without the overhead — Pysar helps shape the work, not take it over.",
  github: "https://github.com/malyshev/pysar",
  /** owner/repo for GitHub REST (stars on nav/footer links). */
  githubRepo: "malyshev/pysar",
  linkedin: "https://www.linkedin.com/in/serhiimalyshev/",
  x: "https://x.com/s_malyshev",
} as const;

/** Square mark-on-ink asset for JSON-LD / social. Wordmarks live at /logo.svg and /logo-dark.svg. */
export const siteLogo = {
  path: "/logo.png",
  width: 512,
  height: 512,
} as const;

/** Homepage SERP + social copy. Keep aligned with the hero (pysar-promo locked stack). */
export const homePageSeo = {
  h1: "Your idea. An article you're ready to stand behind.",
  title: "Pysar — your idea, an article you're ready to stand behind",
  description:
    "Real editorial help, without the overhead — Pysar helps shape the work, not take it over.",
  eyebrow: "Idea in your head. No editor on call.",
} as const;

/** Hero How CTA — paste into Claude / Cursor / Codex (no extra typing). */
export const heroInstallPrompt =
  `Install Pysar from ${siteConfig.url.replace(/\/$/, "")} using the official instructions.` as const;

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
      note: "Install the Pysar Cursor plugin (Marketplace or local rsync copy) for /ps skills. Site Install in Cursor is MCP-only. Skills are not under ~/.cursor/skills.",
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
    "When the user pastes https://getpysar.com and says install: (1) run curl -fsSL https://getpysar.com/install.sh | bash (or download to install.sh, inspect, then bash install.sh); (2) verify with pysar --version; (3) in the writing project directory run pysar init --claude (Claude Code), pysar init --cursor (Cursor), or pysar init --codex (ChatGPT Codex). For Cursor, install the Pysar plugin (Marketplace or local copy) for /ps skills — the site Install in Cursor deeplink is MCP-only. Do not invent brew/npm/go-get alternatives. Also see https://getpysar.com/llms.txt, https://getpysar.com/docs/install and https://getpysar.com/docs/init.",
} as const;

/**
 * Canonical Cursor Plugin identity (plugins/pysar).
 * Marketplace and local copy install the full package; site Install in Cursor
 * is an MCP-only deeplink onto the same spawn contract
 * (dec-20260809-cursor-marketplace-v1-dual-discovery-8b748a7a).
 */
export const cursorPlugin = {
  name: "pysar",
  /** In-repo package path (multi-plugin marketplace source). */
  packagePath: "plugins/pysar",
  marketplaceBrowseUrl: "https://cursor.com/marketplace",
  /**
   * MCP install deeplink — config matches plugins/pysar/mcp.json
   * (portable ${userHome}/.local/bin/pysar spawn). Not a full plugin install;
   * /ps skills still need Marketplace or a local plugin copy.
   */
  installDeeplink:
    "cursor://anysphere.cursor-deeplink/mcp/install?name=pysar&config=eyJ0eXBlIjoic3RkaW8iLCJjb21tYW5kIjoiJHt1c2VySG9tZX0vLmxvY2FsL2Jpbi9weXNhciIsImFyZ3MiOlsic2VydmUiXSwiZW52Ijp7IlBZU0FSX1BST0pFQ1RfUk9PVCI6IiR7d29ya3NwYWNlRm9sZGVyfSJ9fQ%3D%3D",
  docsInitPath: "/docs/init",
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
