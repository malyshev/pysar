import Link from "next/link";
import type { DocPage } from "@/lib/docs/types";
import { cn } from "@/lib/utils";

type DocsNavProps = {
  docs: DocPage[];
  currentSlug: string;
};

/** Journey strip on the docs rail — Boxsi uppercase nav language. */
export function DocsNav({ docs, currentSlug }: DocsNavProps) {
  return (
    <nav
      aria-label="Documentation"
      className="border-b-2 border-zinc-100 px-5 py-5"
    >
      <ul className="flex flex-wrap items-center justify-center gap-x-1 gap-y-2">
        {docs.map((doc) => {
          const active = doc.slug === currentSlug;
          return (
            <li key={doc.slug}>
              <Link
                href={doc.href}
                className={cn(
                  "inline-flex items-center gap-2 px-3 py-2 text-xs font-medium uppercase tracking-wide transition-colors duration-300",
                  active
                    ? "text-zinc-950"
                    : "text-zinc-500 hover:text-primary",
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
