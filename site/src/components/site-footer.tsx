import Link from "next/link";
import { CornerTicks } from "@/components/corner-ticks";
import { SiteLogo } from "@/components/site-logo";
import { siteConfig } from "@/lib/site";

/** Boxsi footer rhythm — dark bar, logo cell, links. */
export function SiteFooter() {
  return (
    <footer className="mt-auto size-full overflow-hidden bg-ink px-4 text-primary-foreground md:px-5 lg:px-7.5">
      <div className="site-container relative z-10 divide-y-2 divide-white/10 border-x-2 border-white/10">
        <div className="flex flex-wrap divide-y-2 divide-white/10 md:divide-y-0">
          <Link
            href="/"
            className="relative flex w-full items-center justify-center bg-ink p-6.5 md:w-auto md:border-e-2 md:border-white/10"
          >
            <SiteLogo variant="dark" />
            <CornerTicks className="text-white/35" />
          </Link>

          <nav
            aria-label="Footer"
            className="order-2 flex w-full flex-wrap items-center justify-center gap-6 px-5 py-6.5 md:order-1 md:w-auto md:grow"
          >
            <Link
              href="/docs"
              className="text-xs uppercase text-primary-foreground transition-colors duration-300 hover:text-primary"
            >
              Docs
            </Link>
            <Link
              href="/docs/install"
              className="text-xs uppercase text-primary-foreground transition-colors duration-300 hover:text-primary"
            >
              Install
            </Link>
            <a
              href={siteConfig.github}
              className="text-xs uppercase text-primary-foreground transition-colors duration-300 hover:text-primary"
              rel="noopener noreferrer"
            >
              GitHub
            </a>
          </nav>
        </div>

        <div className="flex flex-col gap-1 px-5 py-6 text-sm text-primary-foreground/60 md:flex-row md:items-center md:justify-between">
          <p>MIT licensed. Bring your take.</p>
          <p className="text-xs uppercase tracking-wide">
            © {new Date().getFullYear()} {siteConfig.name}
          </p>
        </div>
      </div>
    </footer>
  );
}
