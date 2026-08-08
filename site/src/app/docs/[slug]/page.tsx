import { notFound } from "next/navigation";
import { DocsMarkdown } from "@/components/docs-markdown";
import { DocsNav } from "@/components/docs-nav";
import { SiteShell } from "@/components/site-shell";
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
    <SiteShell>
      <DocsNav docs={docs} currentSlug={doc.slug} />
      <DocsMarkdown body={doc.body} />
    </SiteShell>
  );
}
