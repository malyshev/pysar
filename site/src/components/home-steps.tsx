import Image from "next/image";
import { CornerTicks } from "@/components/corner-ticks";

const steps = [
  {
    n: "01",
    title: "Start with the idea",
    body: (
      <>
        Turn a half-formed thought into a clear brief — or drop in a{" "}
        <span className="text-ink">rough draft</span> and keep going.
      </>
    ),
  },
  {
    n: "02",
    title: "Get a full draft",
    body: (
      <>
        A complete piece from what you already decided — not a blank chat
        guessing your point.
      </>
    ),
    muted: true,
  },
  {
    n: "03",
    title: "Tighten what matters",
    body: "Structure, stakes, and wording get a hard look. Earlier versions stay, so nothing useful is lost.",
    primary: true,
  },
  {
    n: "04",
    title: "Make it sound human",
    body: "Strip the AI tells, then export a draft you're ready to share yourself.",
    dark: true,
  },
] as const;

/** Boxsi StepByStep.tsx grid — same cells, Pysar pipeline. */
export function HomeSteps() {
  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="site-container relative z-10 border-frame py-16 lg:py-30">
        <div className="mb-12.5 px-5 md:px-12.5">
          <div className="grid gap-5 md:grid-cols-2">
            <div className="space-y-2.5 text-center md:text-start">
              <div className="inline-flex items-center justify-center rounded-full border-2 border-frame px-4 py-1.5 text-xs font-semibold uppercase text-primary">
                One step at a time
              </div>
              <h2 className="text-[36px] md:text-[42px] lg:text-[52px]">
                From idea to a piece{" "}
                <span className="text-muted-foreground">you trust</span>
              </h2>
            </div>
            <div className="mx-auto max-w-62 text-center md:me-0 md:place-self-end md:text-end">
              Each step leaves a version you can open and change —{" "}
              <span className="text-ink">nothing ships until you say so</span>.
            </div>
          </div>
        </div>

        <div className="grid grid-cols-2 divide-x-2 divide-y-2 divide-frame border-y-2 border-frame md:grid-cols-3">
          {steps.slice(0, 2).map((step) => (
            <div
              key={step.n}
              className={`relative p-5 md:p-7.5 lg:p-12.5 ${
                "muted" in step && step.muted ? "bg-muted" : ""
              }`}
            >
              <div className="relative mb-22.5 inline-flex items-center justify-center p-1.5">
                <span className="inline-flex size-10 items-center justify-center rounded bg-ink text-xl text-primary-foreground">
                  {step.n}
                </span>
                <CornerTicks />
              </div>
              <h3 className="mb-2.5 text-xl">{step.title}</h3>
              <p>{step.body}</p>
              <CornerTicks />
            </div>
          ))}

          <div className="relative hidden md:flex">
            <Image
              src="/patterns/bg-pattern-small.png"
              alt=""
              width={400}
              height={400}
              className="size-full object-cover"
            />
            <CornerTicks />
          </div>

          <div className="relative hidden border-b-0 md:flex">
            <div className="min-h-full min-w-full border-x border-frame bg-background bg-[repeating-linear-gradient(-45deg,var(--muted)_0,var(--muted)_1px,transparent_0,transparent_50%)] bg-size-[10px_10px]" />
            <CornerTicks />
          </div>

          {steps.slice(2).map((step) => (
            <div
              key={step.n}
              className={`relative border-b-0 p-5 md:p-7.5 lg:p-12.5 ${
                "primary" in step && step.primary
                  ? "border-e-transparent bg-primary"
                  : "bg-ink"
              }`}
            >
              <div className="relative mb-22.5 inline-flex items-center justify-center p-1.5">
                <span className="inline-flex size-10 items-center justify-center rounded bg-background text-xl text-ink">
                  {step.n}
                </span>
                <CornerTicks className="text-primary-foreground" />
              </div>
              <h3
                className={`mb-2.5 text-xl ${
                  "primary" in step && step.primary
                    ? "text-primary-foreground"
                    : "text-primary-foreground"
                }`}
              >
                {step.title}
              </h3>
              <p
                className={
                  "primary" in step && step.primary
                    ? "text-primary-foreground"
                    : "text-white/50"
                }
              >
                {step.body}
              </p>
              <CornerTicks
                className={
                  "primary" in step && step.primary
                    ? "text-primary-foreground/80"
                    : "text-white/35"
                }
              />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
