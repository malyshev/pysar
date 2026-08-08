import Link from "next/link";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { rewriteDocHref } from "@/lib/docs/links";

type DocsMarkdownProps = {
  body: string;
};

export function DocsMarkdown({ body }: DocsMarkdownProps) {
  return (
    <div className="prose prose-neutral max-w-none prose-headings:font-serif prose-headings:tracking-tight prose-a:text-foreground prose-code:rounded prose-code:bg-muted prose-code:px-1 prose-code:py-0.5 prose-code:before:content-none prose-code:after:content-none">
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
