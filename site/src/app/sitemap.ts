import type { MetadataRoute } from "next";
import { getAllDocs } from "@/lib/docs/get-docs";
import { getSiteUrl } from "@/lib/metadata";

export const dynamic = "force-static";

export default function sitemap(): MetadataRoute.Sitemap {
  const siteUrl = getSiteUrl();
  const docs = getAllDocs();

  return [
    {
      url: siteUrl,
      changeFrequency: "weekly",
      priority: 1,
    },
    ...docs.map((doc) => ({
      url: `${siteUrl}${doc.href}`,
      changeFrequency: "monthly" as const,
      priority: doc.slug === "index" ? 0.9 : 0.8,
    })),
  ];
}
