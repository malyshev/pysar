import { DocsLayout } from "@/components/docs-layout";
import { DocsMarkdown } from "@/components/docs-markdown";
import { DocsNav } from "@/components/docs-nav";
import { getAllDocs, getIndexDoc } from "@/lib/docs/get-docs";
import { buildDocMetadata } from "@/lib/docs/metadata";

export function generateMetadata() {
  return buildDocMetadata(getIndexDoc());
}

export default function DocsIndexPage() {
  const docs = getAllDocs();
  const index = getIndexDoc();

  return (
    <DocsLayout nav={<DocsNav docs={docs} currentSlug="index" />}>
      <DocsMarkdown body={index.body} />
    </DocsLayout>
  );
}
