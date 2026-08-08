import { JsonLd } from "@/components/json-ld";
import { SiteShell } from "@/components/site-shell";
import { buttonVariants } from "@/components/ui/button";
import { buildHomeJsonLd } from "@/lib/json-ld";
import { buildHomeMetadata } from "@/lib/metadata";
import { docsRepoUrl, homePageSeo, installDocsUrl } from "@/lib/site";
import { cn } from "@/lib/utils";

export const metadata = buildHomeMetadata();

export default function HomePage() {
  return (
    <SiteShell>
      <JsonLd data={buildHomeJsonLd()} />
      <section className="flex max-w-xl flex-col gap-6">
        <p className="font-serif text-4xl font-semibold tracking-tight text-foreground md:text-5xl">
          Pysar
        </p>
        <h1 className="text-2xl font-medium tracking-tight text-foreground md:text-3xl">
          {homePageSeo.h1}
        </h1>
        <p className="text-lg leading-relaxed text-muted-foreground">
          An author-directed editorial engine for your writing projects. Shape
          an idea or rough draft into a piece you trust — CLI, MCP, and{" "}
          <code className="rounded bg-muted px-1.5 py-0.5 text-[0.9em]">
            ps-*
          </code>{" "}
          skills for Claude Code and Cursor. It does not post on your behalf.
        </p>
        <div className="flex flex-wrap items-center gap-3 pt-2">
          <a
            href={installDocsUrl}
            className={cn(buttonVariants({ size: "lg" }))}
          >
            Install
          </a>
          <a
            href={docsRepoUrl}
            className={cn(buttonVariants({ variant: "outline", size: "lg" }))}
          >
            Read the docs
          </a>
        </div>
      </section>
    </SiteShell>
  );
}
