"use client";

import { useEffect, useState } from "react";
import { fetchGithubStars, formatStarCount } from "@/lib/github-stars";
import { siteConfig } from "@/lib/site";
import { cn } from "@/lib/utils";

type GitHubLinkProps = {
  className?: string;
  onClick?: () => void;
};

/** GitHub repo link; appends live stargazer count when the API responds. */
export function GitHubLink({ className, onClick }: GitHubLinkProps) {
  const [stars, setStars] = useState<number | null>(null);

  useEffect(() => {
    let cancelled = false;
    void fetchGithubStars(siteConfig.githubRepo).then((n) => {
      if (!cancelled) setStars(n);
    });
    return () => {
      cancelled = true;
    };
  }, []);

  const label =
    stars === null ? "GitHub" : `GitHub, ${stars} stars`;

  return (
    <a
      href={siteConfig.github}
      className={cn(className)}
      rel="noopener noreferrer"
      aria-label={label}
      onClick={onClick}
    >
      GitHub
      {stars !== null ? (
        <span className="ms-1.5 font-normal normal-case tracking-normal text-inherit opacity-70">
          ★ {formatStarCount(stars)}
        </span>
      ) : null}
    </a>
  );
}
