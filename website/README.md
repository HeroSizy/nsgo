# NSGo website

The site at <https://nsgo.dev/>: static Astro pages in English
and Chinese, with no UI framework, and the docs built from the repository's
`docs/public/`, whose marked code blocks run in the page. Its Live demo
links open the gallery's web build, compiled by `nsgo`, full page in a new
tab.

```sh
bun install
bun run dev          # build the gallery and the docs' apps, then serve the site at /
bun run build        # build them, typecheck, then build dist/ and its search index
bun run preview      # serve dist/
bun run gallery      # rebuild only the gallery
bun run live         # rebuild only the docs' apps
bun run screenshots  # retake the home page's screens of the gallery
bun run share-cards  # retake the cards a shared link shows
```

`astro.config.mjs` fixes the base at `/` for nsgo.dev on GitHub Pages, so the
dev server also runs at <http://localhost:4321/>. A build for another path
needs that value changed.

## Pages

| Page | Path | Published |
| --- | --- | --- |
| Home | `/` | Yes |
| Docs | `/docs/` | Yes |
| Benchmarks | `/benchmarks/` | Yes |

`src/site.ts` lists the pages; one marked `dev` is built, linked from the
header, the footer and the home page's next steps, and marked with a dot
only by the dev server, until it is filled and published.

Each page but the docs is a file in `src/pages/[...lang]/` whose
`getStaticPaths` is `langPaths(id)`: English at the base, Chinese under
`zh/`. A first visit from a Chinese browser to an English page goes to its
Chinese twin, and the language menu remembers the choice in `nsgo-lang`.

Every page's words are in `src/copy.ts`; add English and Chinese together.
Chinese headings, titles and kickers omit the trailing `。`; body sentences
keep it. `src/layouts/Layout.astro` gives every page the header, the call to
action and the footer.
The header's GitHub link shows the repository's stars, which
`src/lib/stars.ts` reads from GitHub's API as the site builds, with
`GITHUB_TOKEN` when it is set; when GitHub does not answer, the link shows
none. The website workflow rebuilds daily to keep them current.

`scripts/web-build.mjs` builds an app with `go tool nsgo build -target web`
from the repository root into a folder of `public/`, keeping its wasm only
as a `.gzip` copy. `scripts/build-gallery.mjs` writes `public/gallery/`,
served at `/gallery/`, behind `scripts/gallery-loader.html`, which reports
its download, decompression, compilation and start progress and follows the
site's appearance; `loader.html`, its first address, redirects there. All of `public/`'s
builds are generated and ignored.

## Benchmarks

The Benchmarks page shows `benchmark/results/`, which `benchmark/runner`
writes on each desktop ([benchmark](../benchmark/README.md)): bar charts a
case, with a switch beside each case's title that shows every case on one
desktop, and a table of where each framework runs, whose rows are in
`src/lib/benchmarks.ts`. A desktop without results says so.

## Docs

The docs are the Markdown in the repository's `docs/public/`, whose README
says how to write them. `src/content.config.ts` loads `en/` and `zh/`;
`src/pages/docs/` and `src/pages/zh/docs/` render every English page in each
language, the Chinese one where it exists and the English one with a notice
where it does not (`src/lib/docs.ts`). `src/lib/doc-links.mjs`, a Sätteri
plugin, turns the relative `.md` links between pages into the site's URLs.
`src/lib/live.mjs` puts a box before each code block marked
`preview=<demo>` or `program` (`docs/public/README.md`), and
`scripts/build-live.mjs` builds what the boxes run as one app in
`public/live/`, so a reader downloads NSGo once for the whole site:
`examples/previews`, the component examples, and each program block's
code, built once when a translation shares it, each a package whose
`Root` the app runs when its launch parameters name it. Programs thus
share one storage. `src/scripts/live.ts` mounts a box through the tool's
loader once it nears the window, compiling the wasm once a page, in the
site's appearance, which it follows, and the page's language. A Go program cannot be stopped, so leaving a page whose
apps started loads the next page whole.
A page has the sidebar of its language's groups, under the pages without one, its contents, the previous
and next pages, an edit link and copy buttons on its code; the styles are in
`src/styles/docs.css`.

`src/lib/llms.ts` publishes the English pages for language models: `/llms.txt` lists them by group, `/llms-full.txt` holds them all and `src/pages/docs/[...slug].md.ts` serves each at `/docs/<page>.md`, its relative links made absolute and its code blocks' marks dropped. A docs page's head names its Markdown, and every page's `llms.txt`.

Search (`src/components/DocSearch.astro`, opened from the sidebar, ⌘K or
`/`) runs on [Pagefind](https://pagefind.app/): `bun run build` indexes each
page's article after Astro into `dist/pagefind/`, one index per `<html lang>`,
so each language searches its own pages, and Chinese is split into words.
The dev server has no index, so search works only in `bun run preview`.

## Look

White by default; the header's appearance menu chooses light, dark or the
system's, stored in `nsgo-theme`, which the gallery's loader also reads. The mark gives the
site its look: its blue gradient (`#4da3ff` to `#007aff`) for the primary
buttons, the lit platform, the engine and the call to action, and its skewed speed lines
(`SpeedLines.astro`) before section names and behind the hero. Cards are
soft and rounded, code takes GitHub's light or dark colors with the
appearance, and behind the hero a WebGL shader (`HeroField.astro`) drifts
faint blue through its dot grid, faint under the copy's words, around a glow
on the platforms' board, from whose lit platform a ring of dots spreads; it
is still when less motion is asked for. Real
screens of the gallery show what NSGo makes. The tokens and layout are in
`src/styles/global.css`.

The home page's screens are `src/assets/screens/`, the gallery's desktop
and phone layouts in each of `src/site.ts`'s languages and in light and
dark, named `<layout>-<lang>-<theme>.jpg`; the page shows the pair in the
site's appearance. `bun run screenshots` retakes them all from
`public/gallery/` (run `bun run gallery` first) with Google Chrome and
macOS's `sips`; retake them when the gallery's look changes or a language
is added, giving it a locale in `scripts/screenshots.mjs`.

A shared link shows a large card, `src/assets/share/<lang>.jpg`: the home
page's hero in the page's language at 1200 by 630 in light, with the mark,
the title, the lead and the platforms' board. `bun run share-cards` retakes
them from its own `astro dev` with Google Chrome and `sips`; retake them
when the hero changes or a language is added. A build fails without its
language's card.

## Docs: goals

- **Versions from 1.0.** Until 1.0 the docs show `main`. From 1.0 each
  minor release keeps its docs, built from its release tag so
  `docs/public/` stays the only copy: `/docs/` is the latest release,
  `/docs/1.2/` an older one and `/docs/next/` `main`, with a version switch
  that keeps the page when it exists.
- **Translations kept current.** The build lists Chinese pages that are
  missing or older than the English.
- **Search by version**, once versions exist.
- **Checks.** The build fails on a broken link or anchor, and code blocks
  marked as programs compile through `go tool nsgo check`.

## Deployment

The [Website workflow](../.github/workflows/website.yml) builds `website/dist`
and publishes it to GitHub Pages.
