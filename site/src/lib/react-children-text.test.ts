import { createElement } from "react";
import { describe, expect, it } from "vitest";
import { reactChildrenText } from "./react-children-text";

describe("reactChildrenText", () => {
  it("joins nested code children for fenced blocks", () => {
    const tree = createElement(
      "pre",
      null,
      createElement("code", { className: "language-bash" }, "pysar init\n"),
    );
    expect(reactChildrenText(tree)).toBe("pysar init\n");
  });
});
