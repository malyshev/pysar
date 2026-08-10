/** Format stargazer count for compact nav/footer (e.g. 1.2k). */
export function formatStarCount(n: number): string {
  if (n < 1000) return String(n);
  const k = n / 1000;
  const rounded = k >= 10 ? Math.round(k) : Math.round(k * 10) / 10;
  return `${rounded}k`;
}

type StarsCache = {
  value: number | null | undefined;
  inflight: Promise<number | null> | null;
};

const cache: StarsCache = { value: undefined, inflight: null };

/** One shared fetch per page load; null = unavailable (keep UI as GitHub only). */
export function fetchGithubStars(repo: string): Promise<number | null> {
  if (cache.value !== undefined) return Promise.resolve(cache.value);
  if (!cache.inflight) {
    cache.inflight = fetch(`https://api.github.com/repos/${repo}`, {
      headers: { Accept: "application/vnd.github+json" },
    })
      .then(async (res) => {
        if (!res.ok) return null;
        const data = (await res.json()) as { stargazers_count?: unknown };
        return typeof data.stargazers_count === "number"
          ? data.stargazers_count
          : null;
      })
      .catch(() => null)
      .then((n) => {
        cache.value = n;
        return n;
      });
  }
  return cache.inflight;
}
