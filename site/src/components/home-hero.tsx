import Link from "next/link";
import { CopyCommand } from "@/components/copy-command";
import { heroInstallPrompt, homePageSeo } from "@/lib/site";

/**
 * Port of Boxsi (home)/components/Hero.tsx — What? Why? How?
 * How = one paste prompt for AI editors (pysar-promo locked stack).
 */
export function HomeHero() {
  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="relative mx-auto w-full max-w-250 border-x-2 border-zinc-100 pt-28 pb-8 md:pt-32 lg:pt-40 lg:pb-16">
        {/* Plain img — next/image wrapper breaks template's absolute inset-0 pattern. */}
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src="/patterns/bg-pattern.png"
          alt=""
          className="absolute inset-0 z-1 size-full"
        />

        <div className="relative space-y-6 text-center">
          <div className="relative z-10 mx-auto max-w-200 space-y-4 md:px-12.5">
            <div className="inline-flex items-center justify-center gap-2 rounded-full border-2 border-zinc-200 bg-white px-5 py-2 text-xs/none font-semibold uppercase text-zinc-900">
              <span
                className="size-2.5 shrink-0 bg-primary"
                aria-hidden="true"
              />
              <span>{homePageSeo.eyebrow}</span>
            </div>

            <h1 className="text-center text-hero leading-hero tracking-hero md:text-hero-md lg:text-7xl">
              Your idea.{" "}
              <span
                className="inline-flex size-8.5 items-center justify-center md:size-13 lg:-translate-y-3"
                aria-hidden="true"
              >
                {/* eslint-disable-next-line @next/next/no-img-element -- pixel mark must stay crisp */}
                <img
                  src="/logo-mark.svg"
                  alt=""
                  width={304}
                  height={304}
                  className="pixelated size-full"
                />
              </span>{" "}
              <span className="text-zinc-500">
                An article you&apos;re ready to stand behind.
              </span>
            </h1>

            <p>{homePageSeo.description}</p>

            <div className="mx-auto w-full max-w-xl space-y-3 pt-2 text-left">
              <p className="text-center text-sm text-zinc-500">
                Paste into Claude, Cursor, or Codex:
              </p>

              <CopyCommand command={heroInstallPrompt} />

              <p className="text-center text-sm text-zinc-500">
                Or{" "}
                <Link
                  href="/docs"
                  className="font-medium text-zinc-800 underline decoration-zinc-300 underline-offset-4 transition-colors hover:text-primary hover:decoration-primary"
                >
                  read the full documentation
                </Link>
                ?
              </p>
            </div>
          </div>

          <div
            className="absolute inset-0 z-1 h-hero-blur bg-white opacity-50 filter backdrop-blur-2xl"
            aria-hidden="true"
          />
        </div>
      </div>
    </section>
  );
}
