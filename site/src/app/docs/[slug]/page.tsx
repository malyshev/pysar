import { notFound } from "next/navigation";
import { DocsLayout } from "@/components/docs-layout";
import { DocsMarkdown } from "@/components/docs-markdown";
import { DocsNav } from "@/components/docs-nav";
import { getAllDocs, getDocBySlug, getDocSlugs } from "@/lib/docs/get-docs";
import { buildDocMetadata } from "@/lib/docs/metadata";

type PageProps = {
  params: Promise<{ slug: string }>;
};

export function generateStaticParams() {
  return getDocSlugs().map((slug) => ({ slug }));
}

export async function generateMetadata({ params }: PageProps) {
  const { slug } = await params;
  const doc = getDocBySlug(slug);
  if (!doc) {
    return {};
  }
  return buildDocMetadata(doc);
}

export default async function DocSlugPage({ params }: PageProps) {
  const { slug } = await params;
  if (slug === "index") {
    notFound();
  }

  const doc = getDocBySlug(slug);
  if (!doc) {
    notFound();
  }

  const docs = getAllDocs();

  return (
    <DocsLayout nav={<DocsNav docs={docs} currentSlug={doc.slug} />}>
      <DocsMarkdown body={doc.body} />
    </DocsLayout>
  );
}
