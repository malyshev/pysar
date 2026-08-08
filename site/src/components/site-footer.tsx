import Link from "next/link";
import { CornerTicks } from "@/components/corner-ticks";
import { FooterSubscribe } from "@/components/footer-subscribe";
import { SiteLogo } from "@/components/site-logo";
import { siteConfig } from "@/lib/site";

/**
 * Port of Boxsi Footer.tsx — logo | nav | subscribe, dark hatch, bottom bar.
 * Copy/links are Pysar; no invented social networks.
 */
export function SiteFooter() {
  return (
    <footer className="mt-auto size-full overflow-hidden bg-zinc-950 px-4 md:px-5 lg:px-7.5">
      <div className="relative z-10 mx-auto w-full max-w-250 divide-y-2 divide-zinc-800 border-2 border-zinc-800 py-16 lg:py-30">
        <div className="flex flex-wrap divide-y-2 divide-zinc-800 border-t-2 border-b-0 border-zinc-800 md:border-b-2">
          <Link
            href="/"
            className="relative flex w-full items-center justify-center bg-zinc-900 p-6.5 md:w-auto"
          >
            <SiteLogo variant="dark" />
            <CornerTicks className="text-zinc-500" />
          </Link>

          <nav
            aria-label="Footer"
            className="order-2 flex h-full w-full grow flex-wrap items-center justify-center gap-6 self-center px-5 py-6.5 md:order-1 md:w-auto lg:border-0"
          >
            <Link
              href="/"
              className="text-xs text-white uppercase transition-colors duration-300 hover:text-primary"
            >
              Home
            </Link>
            <Link
              href="/docs"
              className="text-xs text-white uppercase transition-colors duration-300 hover:text-primary"
            >
              Docs
            </Link>
            <Link
              href="/docs/install"
              className="text-xs text-white uppercase transition-colors duration-300 hover:text-primary"
            >
              Install
            </Link>
            <a
              href={siteConfig.github}
              className="text-xs text-white uppercase transition-colors duration-300 hover:text-primary"
              rel="noopener noreferrer"
            >
              GitHub
            </a>
          </nav>

          <FooterSubscribe />
        </div>

        <div className="min-h-10 min-w-full bg-stripe-dark" aria-hidden="true" />

        <div className="flex flex-wrap divide-y-2 divide-zinc-800 border-b-2 border-zinc-800">
          <div className="grow p-7.5 md:border-b-0">
            <p className="self-center text-center text-xs text-white uppercase md:text-start">
              Made with{" "}
              <a
                href="https://haft.tools/"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary"
              >
                Haft
              </a>
            </p>
          </div>

          <div className="mx-auto flex h-full flex-wrap items-center justify-center self-center divide-x-2 divide-zinc-800 border-e-2 border-s-2 border-zinc-800 md:border-e-0">
            <a
              href={siteConfig.github}
              className="p-5 text-xs text-white uppercase transition-colors duration-300 hover:bg-zinc-800 md:p-7.5"
              rel="noopener noreferrer"
            >
              GitHub
            </a>
            <Link
              href="/docs"
              className="p-5 text-xs text-white uppercase transition-colors duration-300 hover:bg-zinc-800 md:p-7.5"
            >
              Docs
            </Link>
            <p className="p-5 text-xs text-white uppercase md:p-7.5">
              © {new Date().getFullYear()} {siteConfig.name}
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
}
