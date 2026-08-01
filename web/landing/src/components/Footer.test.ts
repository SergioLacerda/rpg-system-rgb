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
    expect(html).toContain("GitHub · Discord · Changelog");
  });

  it("renders English footer copy", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Footer, {
      props: { lang: "en" },
    });

    expect(html).toContain("RGB System · open source");
  });
});
