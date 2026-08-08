import { Circle, Sparkles } from "lucide-react";
import { homePageSeo } from "@/lib/site";

/**
 * Port of Boxsi (home)/components/Hero.tsx — same structure, paddings, blur, pattern.
 * Copy is Pysar; sizes/colors follow the template (zinc + hero theme tokens).
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
              <span className="inline-flex size-8.5 items-center justify-center rounded-lg bg-primary md:size-13 lg:-translate-y-3">
                <Sparkles
                  className="size-hero-sparkle fill-white text-white"
                  aria-hidden="true"
                />
              </span>{" "}
              <span className="text-zinc-500">Ship a piece you trust.</span>
            </h1>

            <p>
              Shape an idea or rough draft into a piece you trust. CLI plus MCP
              and <span className="text-zinc-950">ps-*</span> skills for{" "}
              <span className="text-zinc-950">Claude Code</span> and{" "}
              <span className="text-zinc-950">Cursor</span> — without posting on
              your behalf.
            </p>
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
