import { Children, isValidElement, type ReactNode } from "react";

/** Flatten React children to plain text (fenced-code extraction for docs). */
export function reactChildrenText(node: ReactNode): string {
  if (node == null || typeof node === "boolean") return "";
  if (typeof node === "string" || typeof node === "number") return String(node);
  if (Array.isArray(node)) {
    return node.map(reactChildrenText).join("");
  }
  if (isValidElement<{ children?: ReactNode }>(node)) {
    return reactChildrenText(node.props.children);
  }
  return Children.toArray(node).map(reactChildrenText).join("");
}
