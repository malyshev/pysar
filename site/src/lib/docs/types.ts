export type DocFrontmatter = {
  title: string;
  slug: string;
  nav_order: number;
  section: string;
};

export type DocPage = DocFrontmatter & {
  /** Absolute path to the source markdown file. */
  filePath: string;
  /** Markdown body without frontmatter. */
  body: string;
  /** Site path: `/docs` for index, `/docs/<slug>` otherwise. */
  href: string;
};
