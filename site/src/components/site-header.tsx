"use client";

import { Menu, X } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { CornerTicks } from "@/components/corner-ticks";
import { GitHubLink } from "@/components/github-link";
import { SiteLogo } from "@/components/site-logo";

/**
 * Port of Boxsi Navbar.tsx — same paddings, borders, stripe fill, CTA cell.
 * Copy/links are Pysar; no Preline overlay.
 */
export function SiteHeader() {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <header className="fixed inset-x-0 top-0 z-50 bg-white px-4 md:px-5 lg:px-7.5">
      <div className="mx-auto flex w-full max-w-250 items-center border-x-2 border-y-2 border-zinc-100">
        <Link
          href="/"
          className="relative flex items-center border-zinc-100 p-5 md:border-e-2 md:p-6.5"
          onClick={() => setMenuOpen(false)}
        >
          <SiteLogo priority className="h-5.5" />
          <CornerTicks className="text-zinc-500" />
        </Link>

        <nav
          aria-label="Primary"
          className="hidden h-full items-center justify-center self-center p-5 md:flex"
        >
          <Link
            href="/docs"
            className="px-4 py-2 uppercase text-zinc-950 transition-colors duration-300 hover:text-primary"
          >
            Docs
          </Link>
          <GitHubLink className="px-4 py-2 uppercase text-zinc-950 transition-colors duration-300 hover:text-primary" />
        </nav>

        <div className="grow self-center border-x-2 border-zinc-100">
          <div className="flex h-full min-h-16 justify-end gap-4 bg-white md:bg-stripe">
            <div className="relative hidden bg-zinc-100 p-2.5 sm:flex md:p-4">
              <Link
                href="/docs/install"
                className="inline-flex items-center justify-center gap-2.5 bg-zinc-900 px-5 py-3.5 font-medium uppercase text-white"
              >
                <span
                  className="size-3 shrink-0 rounded-full bg-primary"
                  aria-hidden="true"
                />
                Install
              </Link>
              <CornerTicks className="text-zinc-500" />
            </div>

            <div className="flex items-center p-1 pe-4 md:hidden">
              <button
                type="button"
                className="inline-flex size-10 items-center justify-center bg-zinc-900 text-white"
                aria-expanded={menuOpen}
                aria-controls="navbar-menu"
                onClick={() => setMenuOpen((open) => !open)}
              >
                {menuOpen ? (
                  <X className="size-6" aria-hidden="true" />
                ) : (
                  <Menu className="size-6" aria-hidden="true" />
                )}
                <span className="sr-only">
                  {menuOpen ? "Close menu" : "Open menu"}
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>

      {menuOpen ? (
        <div
          id="navbar-menu"
          className="w-full rounded-b-lg bg-white px-6 py-2 md:hidden"
          role="dialog"
          aria-label="Mobile"
        >
          <div className="flex flex-col divide-y divide-dashed divide-zinc-200">
            <Link
              href="/docs"
              className="flex items-center py-4 text-2xl/none font-medium text-zinc-900 transition-all hover:text-primary"
              onClick={() => setMenuOpen(false)}
            >
              Docs
            </Link>
            <GitHubLink
              className="flex items-center py-4 text-2xl/none font-medium text-zinc-900 transition-all hover:text-primary"
              onClick={() => setMenuOpen(false)}
            />
            <Link
              href="/docs/install"
              className="flex items-center py-4 text-2xl/none font-medium text-zinc-900 transition-all hover:text-primary"
              onClick={() => setMenuOpen(false)}
            >
              Install
            </Link>
            <div className="ms-auto py-4 text-sm text-zinc-500">
              © {new Date().getFullYear()} Pysar
            </div>
          </div>
        </div>
      ) : null}
    </header>
  );
}
