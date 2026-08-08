import { describe, expect, it } from "vitest";
import { getAllDocs, getDocSlugs, getIndexDoc } from "@/lib/docs/get-docs";

describe("getAllDocs", () => {
  it("loads journey pages with required frontmatter", () => {
    const docs = getAllDocs();
    expect(docs.length).toBeGreaterThanOrEqual(7);
    expect(docs.map((d) => d.slug)).toEqual(
      expect.arrayContaining([
        "index",
        "install",
        "init",
        "pipeline",
        "mcp-and-skills",
        "export",
        "troubleshooting",
      ]),
    );
    expect(docs[0]?.slug).toBe("index");
    expect(getIndexDoc().href).toBe("/docs");
    expect(getDocSlugs()).not.toContain("index");
    expect(getDocSlugs()).toContain("install");
  });
});
