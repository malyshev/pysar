import Link from "next/link";
import { CornerTicks } from "@/components/corner-ticks";
import { FooterSubscribe } from "@/components/footer-subscribe";
import { SiteLogo } from "@/components/site-logo";
import { siteConfig, united24Url } from "@/lib/site";

/**
 * Port of Boxsi Footer.tsx — logo | nav | subscribe, dark hatch, bottom bar.
 * United24 sits in the bottom frame cell (Merchanto-style support link).
 */
export function SiteFooter() {
  return (
    <footer className="mt-auto size-full overflow-hidden bg-zinc-950 px-4 md:px-5 lg:px-7.5">
      <div className="relative z-10 mx-auto w-full max-w-250 divide-y-2 divide-zinc-800 border-2 border-zinc-800 pt-16 lg:pt-30">
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
              <span className="mx-2.5 text-zinc-600" aria-hidden="true">
                |
              </span>
              <a
                href="https://qorym.com/"
                target="_blank"
                rel="noopener noreferrer"
                className="text-[#44E5FE]"
              >
                Qorym
              </a>{" "}
              Partner
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
            <a
              href={siteConfig.linkedin}
              className="p-5 text-xs text-white uppercase transition-colors duration-300 hover:bg-zinc-800 md:p-7.5"
              target="_blank"
              rel="noopener noreferrer"
            >
              LinkedIn
            </a>
            <a
              href={siteConfig.x}
              className="inline-flex items-center justify-center p-5 text-white transition-colors duration-300 hover:bg-zinc-800 md:p-7.5"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="X"
            >
              <svg
                viewBox="0 0 24 24"
                aria-hidden="true"
                className="size-3.5 fill-current"
              >
                <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
              </svg>
            </a>
            <p className="p-5 text-xs text-white uppercase md:p-7.5">
              © {new Date().getFullYear()} {siteConfig.name}
            </p>
          </div>
        </div>

        <div className="flex justify-center border-b-2 border-zinc-800 px-7.5 py-10 md:py-12">
          <a
            href={united24Url}
            target="_blank"
            rel="noopener noreferrer"
            className="group flex flex-col items-center gap-2.5 transition-opacity duration-300 hover:opacity-85"
          >
            {/* eslint-disable-next-line @next/next/no-img-element -- official U24 wordmark */}
            <img
              src="/united24.svg"
              alt="United24"
              width={422}
              height={40}
              className="h-6 w-auto md:h-7"
              loading="lazy"
            />
            <span className="text-xs text-zinc-400 uppercase transition-colors group-hover:text-white">
              Support for Ukraine
            </span>
          </a>
        </div>
      </div>
    </footer>
  );
}
