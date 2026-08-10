import { describe, expect, it } from "vitest";
import { buildLlmsTxt } from "@/lib/llms-txt";
import {
  agentInstallGuide,
  installCommand,
  installInspectCommands,
} from "@/lib/site";

describe("buildLlmsTxt", () => {
  it("carries the same install recipe as the homepage agent guide", () => {
    const body = buildLlmsTxt("https://getpysar.com");

    expect(body.startsWith("# Pysar\n")).toBe(true);
    expect(body).toContain("> ");
    expect(body).toContain("## Instructions for LLM Agents");
    expect(body).toContain(installCommand);
    expect(body).toContain(installInspectCommands.download);
    expect(body).toContain(installInspectCommands.run);
    for (const host of agentInstallGuide.hosts) {
      expect(body).toContain(host.initCommand);
    }
    expect(body).toContain("https://getpysar.com/docs/install");
    expect(body).toContain("https://getpysar.com/docs/export");
    expect(body).toContain("--research");
    expect(body).toContain("https://getpysar.com/#for-ai-agents");
    expect(body).toContain("MCP-only");
    expect(body).not.toContain("brew install");
  });
});
