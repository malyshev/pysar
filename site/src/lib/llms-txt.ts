import {
  agentInstallGuide,
  homePageSeo,
  installCommand,
  installInspectCommands,
  siteConfig,
} from "@/lib/site";

/**
 * Canonical /llms.txt body (llmstxt.org shape).
 * Complements homepage V4 carrier — note-20260808-26ff4eb4.
 */
export function buildLlmsTxt(siteUrl = siteConfig.url.replace(/\/$/, "")): string {
  const hosts = agentInstallGuide.hosts
    .map((h) => `   - ${h.name}: \`${h.initCommand}\` — ${h.note}`)
    .join("\n");

  return `# ${siteConfig.name}

> ${homePageSeo.description}

Pysar is an author-directed editorial engine for writing projects. It installs as a CLI, scaffolds MCP + skills into Claude Code, Cursor, or ChatGPT Codex, and never posts for you.

## Instructions for LLM Agents

When the user pastes ${siteUrl} (or this file) and says "${agentInstallGuide.trigger}":

1. Install the binary (macOS / Linux). Preferred one-liner:
   \`${installCommand}\`
   Inspectable two-step (same script):
   \`${installInspectCommands.download}\`
   then inspect \`install.sh\`, then \`${installInspectCommands.run}\`
2. Confirm with \`pysar --version\`.
3. In the writing project directory, init for the host you are running in:
${hosts}
4. Continue with \`/ps-onboard\` or \`/ps\` as the author directs.

Do not invent brew, npm, pip, or \`go install\` alternatives. The homepage section \`${siteUrl}/#for-ai-agents\` carries the same recipe for URL-drop without a second hop.

## Docs

- [Install](${siteUrl}${agentInstallGuide.docsInstallPath}): binary install (script, Windows zip, build from source)
- [Init a project](${siteUrl}${agentInstallGuide.docsInitPath}): \`pysar init\` for Claude Code, Cursor, Codex
- [Docs home](${siteUrl}/docs): journey docs index
- [Pipeline](${siteUrl}/docs/pipeline): \`/ps\` default chain; optional \`--research\`, \`--seo\`, \`--review\`
- [Export](${siteUrl}/docs/export): root Markdown copy; Latin piece slugs; citation markers may remain until resolve-at-export
- [MCP and skills](${siteUrl}/docs/mcp-and-skills): how the agent surface works
- [Troubleshooting](${siteUrl}/docs/troubleshooting): common install/init failures

## Project

- [GitHub](${siteConfig.github}): source and releases
- [Homepage](${siteUrl}): human marketing + AI install carrier
`;
}
