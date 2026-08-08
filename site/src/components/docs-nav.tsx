import Link from "next/link";
import type { DocPage } from "@/lib/docs/types";
import { cn } from "@/lib/utils";

type DocsNavProps = {
  docs: DocPage[];
  currentSlug: string;
};

export function DocsNav({ docs, currentSlug }: DocsNavProps) {
  return (
    <nav aria-label="Documentation" className="mb-10 border-b border-border pb-4">
      <ul className="flex flex-wrap gap-x-4 gap-y-2 text-sm">
        {docs.map((doc) => {
          const active = doc.slug === currentSlug;
          return (
            <li key={doc.slug}>
              <Link
                href={doc.href}
                className={cn(
                  "transition-colors hover:text-foreground",
                  active
                    ? "font-medium text-foreground"
                    : "text-muted-foreground",
                )}
              >
                {doc.slug === "index" ? "Guide" : doc.title}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
