import type { Metadata } from "next";
import type { DocPage } from "@/lib/docs/types";
import { getDefaultOgImageUrl, getSiteUrl } from "@/lib/metadata";
import { siteConfig } from "@/lib/site";

function firstParagraph(body: string): string {
  const block = body
    .split(/\n\s*\n/)
    .map((p) => p.replace(/^#+\s+.*$/m, "").trim())
    .find((p) => p.length > 0 && !p.startsWith("```"));
  if (!block) {
    return siteConfig.description;
  }
  const plain = block
    .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1")
    .replace(/[`*_]/g, "")
    .replace(/\s+/g, " ")
    .trim();
  return plain.length > 160 ? `${plain.slice(0, 157)}…` : plain;
}

export function buildDocMetadata(doc: DocPage): Metadata {
  const url = `${getSiteUrl()}${doc.href}`;
  const description = firstParagraph(doc.body);
  const image = getDefaultOgImageUrl();
  const title = doc.slug === "index" ? "Docs" : doc.title;

  return {
    title,
    description,
    alternates: {
      canonical: url,
    },
    openGraph: {
      title: doc.slug === "index" ? `${title} · ${siteConfig.name}` : doc.title,
      description,
      url,
      siteName: siteConfig.name,
      locale: "en_US",
      type: "article",
      images: [{ url: image, alt: `${doc.title} — ${siteConfig.name}` }],
    },
    twitter: {
      card: "summary_large_image",
      title: doc.title,
      description,
      images: [image],
    },
  };
}
