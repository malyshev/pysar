import Link from "next/link";
import type { DocPage } from "@/lib/docs/types";
import { cn } from "@/lib/utils";

type DocsNavProps = {
  docs: DocPage[];
  currentSlug: string;
};

export function DocsNav({ docs, currentSlug }: DocsNavProps) {
  return (
    <nav
      aria-label="Documentation"
      className="mb-10 border-b-2 border-frame pb-5"
    >
      <p className="mb-3 text-[0.65rem] font-medium uppercase tracking-[0.14em] text-muted-foreground">
        Documentation
      </p>
      <ul className="flex flex-wrap gap-x-1 gap-y-2">
        {docs.map((doc) => {
          const active = doc.slug === currentSlug;
          return (
            <li key={doc.slug}>
              <Link
                href={doc.href}
                className={cn(
                  "inline-flex items-center gap-2 px-3 py-2 text-xs font-medium uppercase tracking-wide transition-colors",
                  active
                    ? "bg-muted text-ink"
                    : "text-muted-foreground hover:text-primary",
                )}
                aria-current={active ? "page" : undefined}
              >
                {active ? (
                  <span
                    className="size-1.5 shrink-0 rounded-full bg-primary"
                    aria-hidden="true"
                  />
                ) : null}
                {doc.slug === "index" ? "Guide" : doc.title}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
