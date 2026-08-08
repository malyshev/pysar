import { DocsMarkdown } from "@/components/docs-markdown";
import { DocsNav } from "@/components/docs-nav";
import { SiteShell } from "@/components/site-shell";
import { getAllDocs, getIndexDoc } from "@/lib/docs/get-docs";
import { buildDocMetadata } from "@/lib/docs/metadata";

export function generateMetadata() {
  return buildDocMetadata(getIndexDoc());
}

export default function DocsIndexPage() {
  const docs = getAllDocs();
  const index = getIndexDoc();

  return (
    <SiteShell>
      <DocsNav docs={docs} currentSlug="index" />
      <DocsMarkdown body={index.body} />
    </SiteShell>
  );
}
