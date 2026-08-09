import Link from "next/link";
import type { ReactNode } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { CopyCommand } from "@/components/copy-command";
import { rewriteDocHref } from "@/lib/docs/links";
import { reactChildrenText } from "@/lib/react-children-text";

type DocsMarkdownProps = {
  body: string;
};

/** Prose on the About-style reading column — zinc / primary; fenced code = CopyCommand. */
export function DocsMarkdown({ body }: DocsMarkdownProps) {
  return (
    <div
      className={[
        "prose max-w-none",
        "prose-headings:font-heading prose-headings:font-bold prose-headings:tracking-tight prose-headings:text-zinc-950",
        "prose-p:text-zinc-600 prose-li:text-zinc-600 prose-strong:text-zinc-950",
        "prose-a:font-medium prose-a:text-zinc-950 prose-a:underline prose-a:decoration-primary/50 prose-a:underline-offset-4 prose-a:transition-colors hover:prose-a:text-primary hover:prose-a:decoration-primary",
        "prose-code:rounded-sm prose-code:border prose-code:border-zinc-200 prose-code:bg-zinc-100 prose-code:px-1.5 prose-code:py-0.5 prose-code:whitespace-nowrap prose-code:text-zinc-950 prose-code:before:content-none prose-code:after:content-none",
        "prose-th:px-3 prose-th:py-2 prose-th:text-left prose-td:px-3 prose-td:py-2 prose-td:align-top prose-td:text-zinc-600",
        "prose-hr:border-zinc-100",
        "prose-blockquote:border-primary prose-blockquote:text-zinc-600",
      ].join(" ")}
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          a: ({ href, children }) => {
            const nextHref = rewriteDocHref(href) ?? href ?? "#";
            const external =
              nextHref.startsWith("http://") ||
              nextHref.startsWith("https://") ||
              nextHref.startsWith("mailto:");

            if (external) {
              return (
                <a href={nextHref} rel="noopener noreferrer">
                  {children}
                </a>
              );
            }

            return <Link href={nextHref}>{children}</Link>;
          },
          // Keep flag columns from wrapping mid-token (e.g. `--` / `review`).
          table: ({ children }) => (
            <div className="my-5 overflow-x-auto">
              <table>{children}</table>
            </div>
          ),
          // Fenced blocks: soft Copy bar (same chrome as home quick start).
          pre: ({ children }) => {
            const command = reactChildrenText(children as ReactNode).replace(
              /\n$/,
              "",
            );
            return <CopyCommand command={command} className="my-5" />;
          },
        }}
      >
        {body}
      </ReactMarkdown>
    </div>
  );
}
