import { Sparkles } from "lucide-react";
import Image from "next/image";
import { homePageSeo } from "@/lib/site";

/** Boxsi Hero.tsx composition — same structure, Pysar copy. */
export function HomeHero() {
  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="site-container relative pt-32 pb-10 md:pt-37.5 lg:pt-48 lg:pb-30">
        <Image
          src="/patterns/bg-pattern.png"
          alt=""
          width={2000}
          height={1434}
          priority
          className="absolute inset-0 z-1 h-full w-full object-cover"
        />

        <div className="relative space-y-7.5 text-center">
          <div className="relative z-10 mx-auto max-w-200 space-y-4 md:px-12.5">
            <div className="inline-flex items-center justify-center gap-2 rounded-full border-2 border-frame bg-background px-4 py-1.5 text-xs/none font-semibold uppercase text-ink">
              <span
                className="size-2.5 shrink-0 rounded-full bg-primary"
                aria-hidden="true"
              />
              <span>{homePageSeo.eyebrow}</span>
            </div>

            <h1 className="text-center text-[36px] leading-[1.1em] tracking-[-0.03em] md:text-[58px] lg:text-7xl">
              Bring your take{" "}
              <span className="inline-flex size-8.5 items-center justify-center rounded-lg bg-primary align-middle md:size-13 lg:-translate-y-3">
                <Sparkles
                  className="size-[60%] fill-primary-foreground text-primary-foreground"
                  aria-hidden="true"
                />
              </span>{" "}
              <span className="text-muted-foreground">
                Ship a piece you trust.
              </span>
            </h1>

            <p>
              Shape an idea or rough draft into a piece you trust. CLI plus MCP
              and <span className="text-ink">ps-*</span> skills for{" "}
              <span className="text-ink">Claude Code</span> and{" "}
              <span className="text-ink">Cursor</span> — without posting on your
              behalf.
            </p>
          </div>

          <div
            className="absolute inset-0 z-1 h-[330px] bg-background opacity-50 backdrop-blur-2xl"
            aria-hidden="true"
          />
        </div>
      </div>
    </section>
  );
}
