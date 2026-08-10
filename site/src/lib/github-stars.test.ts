import { describe, expect, it } from "vitest";
import { formatStarCount } from "@/lib/github-stars";

describe("formatStarCount", () => {
  it("keeps small counts plain", () => {
    expect(formatStarCount(0)).toBe("0");
    expect(formatStarCount(42)).toBe("42");
    expect(formatStarCount(999)).toBe("999");
  });

  it("compacts thousands", () => {
    expect(formatStarCount(1000)).toBe("1k");
    expect(formatStarCount(1200)).toBe("1.2k");
    expect(formatStarCount(10500)).toBe("11k");
  });
});
