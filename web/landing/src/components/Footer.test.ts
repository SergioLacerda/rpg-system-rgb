import { experimental_AstroContainer as AstroContainer } from "astro/container";
import { describe, expect, it } from "vitest";
import Footer from "./Footer.astro";

describe("Footer", () => {
  it("renders Portuguese footer copy", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Footer, {
      props: { lang: "pt-br" },
    });

    expect(html).toContain("RGB System · código aberto");
    expect(html).toContain('href="https://github.com/SergioLacerda"');
    expect(html).toContain(
      'href="https://github.com/SergioLacerda/rpg-system-rgb/blob/main/CHANGELOG.md"',
    );
    expect(html).toContain(">GitHub</a>");
    expect(html).toContain(">Changelog</a>");
    expect(html).not.toContain("Discord");
  });

  it("renders English footer copy", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Footer, {
      props: { lang: "en" },
    });

    expect(html).toContain("RGB System · open source");
  });
});
