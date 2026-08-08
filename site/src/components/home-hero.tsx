import { Check, Circle } from "lucide-react";
import { CopyCommand } from "@/components/copy-command";
import { homePageSeo, installCommand } from "@/lib/site";

/**
 * Port of Boxsi (home)/components/Hero.tsx — same structure, paddings, blur, pattern.
 * Quick-start block: env-sentinel.dev pattern (label → command+Copy → what this does).
 */
export function HomeHero() {
  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="relative mx-auto w-full max-w-250 border-x-2 border-zinc-100 pt-32 pb-10 md:pt-37.5 lg:pt-48 lg:pb-30">
        {/* Plain img — next/image wrapper breaks template's absolute inset-0 pattern. */}
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src="/patterns/bg-pattern.png"
          alt=""
          className="absolute inset-0 z-1 size-full"
        />

        <div className="relative space-y-7.5 text-center">
          <div className="relative z-10 mx-auto max-w-200 space-y-4 md:px-12.5">
            <div className="inline-flex items-center justify-center gap-2 rounded-full border-2 border-zinc-200 bg-white px-4 py-1.5 text-xs/none font-semibold uppercase text-zinc-900">
              <Circle
                className="size-2.5 fill-primary text-primary"
                aria-hidden="true"
              />
              <span>{homePageSeo.eyebrow}</span>
            </div>

            <h1 className="text-center text-hero leading-hero tracking-hero md:text-hero-md lg:text-7xl">
              Bring your take{" "}
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
              <span className="text-zinc-500">Ship a piece you trust.</span>
            </h1>

            <p>
              From a half-formed idea or a messy draft to a piece that{" "}
              <span className="text-zinc-950">sounds like you</span> — shaped
              with your judgment still in the chair.
            </p>

            <div className="mx-auto w-full max-w-xl space-y-4 pt-3 text-left">
              <p className="text-center text-sm text-zinc-500">
                Run this command to install Pysar:
              </p>

              <CopyCommand command={installCommand} />

              <div className="space-y-2">
                <p className="text-sm font-semibold text-zinc-800">
                  What this does:
                </p>
                <ul className="space-y-1.5 text-sm text-zinc-600">
                  <li className="flex items-start gap-2">
                    <Check
                      className="mt-0.5 size-4 shrink-0 text-primary"
                      aria-hidden="true"
                    />
                    <span>
                      Installs{" "}
                      <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[13px] text-zinc-800">
                        pysar
                      </code>{" "}
                      on your PATH
                    </span>
                  </li>
                  <li className="flex items-start gap-2">
                    <Check
                      className="mt-0.5 size-4 shrink-0 text-primary"
                      aria-hidden="true"
                    />
                    <span>macOS and Linux (amd64 / arm64)</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <Check
                      className="mt-0.5 size-4 shrink-0 text-primary"
                      aria-hidden="true"
                    />
                    <span>
                      Then run{" "}
                      <code className="rounded bg-zinc-100 px-1.5 py-0.5 font-mono text-[13px] text-zinc-800">
                        pysar init
                      </code>
                    </span>
                  </li>
                </ul>
              </div>
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
