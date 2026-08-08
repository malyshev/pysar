import Link from "next/link";
import { docsRepoUrl } from "@/lib/site";

export function SiteHeader() {
  return (
    <header className="border-b border-border/80">
      <div className="mx-auto flex h-14 w-full max-w-3xl items-center justify-between px-4 md:px-6">
        <Link
          href="/"
          className="text-lg font-semibold tracking-tight text-foreground"
        >
          Pysar
        </Link>
        <nav className="flex items-center gap-6 text-sm text-muted-foreground">
          <a
            href={docsRepoUrl}
            className="transition-colors hover:text-foreground"
          >
            Docs
          </a>
          <a
            href="https://github.com/malyshev/pysar"
            className="transition-colors hover:text-foreground"
          >
            GitHub
          </a>
        </nav>
      </div>
    </header>
  );
}
