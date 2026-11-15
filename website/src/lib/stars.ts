// The repository's GitHub stars, which the header shows beside its GitHub
// link, read from GitHub's API once when the site builds.
//
// GitHub answers 60 requests an hour without a token, so the build sends
// GITHUB_TOKEN when it is set. A failed request leaves the link without
// stars and never fails the build.

let count: Promise<number | undefined> | undefined;

/** repoStars is the repository's stars, or none where GitHub did not say. */
export function repoStars(): Promise<number | undefined> {
  count ??= fetchStars('nsgo-dev/nsgo');
  return count;
}

async function fetchStars(name: string): Promise<number | undefined> {
  const headers: Record<string, string> = {
    accept: 'application/vnd.github+json',
    'x-github-api-version': '2022-11-28',
  };
  const token: string | undefined = import.meta.env.GITHUB_TOKEN;
  if (token) headers.authorization = `Bearer ${token}`;
  try {
    const response = await fetch(`https://api.github.com/repos/${name}`, {
      headers,
      signal: AbortSignal.timeout(10_000),
    });
    if (!response.ok) {
      console.warn(`stars: GitHub answered ${response.status} for ${name}; showing none`);
      return undefined;
    }
    const body = (await response.json()) as { stargazers_count?: unknown };
    return typeof body.stargazers_count === 'number' ? body.stargazers_count : undefined;
  } catch (error) {
    console.warn(`stars: could not reach GitHub for ${name}; showing none: ${error}`);
    return undefined;
  }
}

/** formatStars writes a count as GitHub does, such as 6, 1.2k or 12k. */
export function formatStars(count: number): string {
  return new Intl.NumberFormat('en', {
    notation: 'compact',
    maximumFractionDigits: 1,
  })
    .format(count)
    .toLowerCase();
}
