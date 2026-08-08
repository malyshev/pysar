import { describe, expect, it } from "vitest";
import { rewriteDocHref } from "@/lib/docs/links";

describe("rewriteDocHref", () => {
  it("maps peer .md links to /docs routes", () => {
    expect(rewriteDocHref("./install.md")).toBe("/docs/install");
    expect(rewriteDocHref("init.md")).toBe("/docs/init");
    expect(rewriteDocHref("./index.md")).toBe("/docs");
  });

  it("preserves hash fragments", () => {
    expect(rewriteDocHref("./install.md#script")).toBe("/docs/install#script");
  });

  it("maps ../ escapes to GitHub blob URLs", () => {
    expect(rewriteDocHref("../README.md")).toBe(
      "https://github.com/malyshev/pysar/blob/master/README.md",
    );
  });

  it("leaves absolute and fragment links alone", () => {
    expect(rewriteDocHref("https://example.com/x")).toBe("https://example.com/x");
    expect(rewriteDocHref("#section")).toBe("#section");
    expect(rewriteDocHref("mailto:a@b.c")).toBe("mailto:a@b.c");
  });
});
