import Link from "next/link";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { rewriteDocHref } from "@/lib/docs/links";

type DocsMarkdownProps = {
  body: string;
};

export function DocsMarkdown({ body }: DocsMarkdownProps) {
  return (
    <div
      className={[
        "prose max-w-none",
        "prose-neutral",
        "prose-headings:font-heading prose-headings:font-bold prose-headings:tracking-[-0.02em] prose-headings:text-ink",
        "prose-p:text-muted-foreground prose-li:text-muted-foreground prose-strong:text-ink",
        "prose-a:font-medium prose-a:text-ink prose-a:underline prose-a:decoration-primary/50 prose-a:underline-offset-4 prose-a:transition-colors hover:prose-a:text-primary hover:prose-a:decoration-primary",
        "prose-code:rounded-sm prose-code:border prose-code:border-frame prose-code:bg-muted prose-code:px-1.5 prose-code:py-0.5 prose-code:text-ink prose-code:before:content-none prose-code:after:content-none",
        "prose-pre:border-2 prose-pre:border-frame prose-pre:bg-ink prose-pre:text-primary-foreground",
        "prose-pre:code:border-0 prose-pre:code:bg-transparent prose-pre:code:text-inherit",
        "prose-hr:border-frame",
        "prose-blockquote:border-primary prose-blockquote:text-muted-foreground",
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
        }}
      >
        {body}
      </ReactMarkdown>
    </div>
  );
}
