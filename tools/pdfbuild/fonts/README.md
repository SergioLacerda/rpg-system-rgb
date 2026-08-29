# PDF Authoring Fonts

Local font files for the manual PDF authoring step. Referenced from
`docs/styles/rgb-pdf.css` via relative, non-CDN `@font-face` URLs (per
`pdf_premium.txt` P-03). No CDN reference is used anywhere in the build.

## Vendored fonts

| Family | Directory | Version | License | Source |
|---|---|---|---|---|
| Source Serif 4 | `source-serif-4/` | 4.005R | SIL OFL 1.1 (`LICENSE.md`) | https://github.com/adobe-fonts/source-serif/releases/tag/4.005R |
| Space Grotesk | `space-grotesk/` | 2.0.0 | SIL OFL 1.1 (`OFL.txt`) | https://github.com/floriankarsten/space-grotesk/releases/tag/2.0.0 |
| JetBrains Mono | `jetbrains-mono/` | v2.304 | SIL OFL 1.1 (`OFL.txt`) | https://github.com/JetBrains/JetBrainsMono/releases/tag/v2.304 |

Each directory contains the static Regular/Bold/Italic/Bold-Italic TTFs
actually used by `docs/styles/rgb-pdf.css`'s `@font-face` rules (Space
Grotesk ships no italic upstream — headings never render italic, so only
Regular/Bold are vendored for it).

Adding a new font: update `docs/styles/rgb-pdf.css` with relative
`@font-face` URLs and record the font's license in a subdirectory here,
following the same pattern.
