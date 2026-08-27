import { experimental_AstroContainer as AstroContainer } from "astro/container";
import { describe, expect, it } from "vitest";
import SkillCard from "./SkillCard.astro";

const baseProps = {
  name: "Specialist",
  role: "Consults and interprets the rules",
  badge: "v0.1 · contract defined",
  desc: "Runtime behavior is not implemented yet.",
  tags: ["runtime not implemented"],
  downloadLabel: "Download package (.zip)",
  soonLabel: "Package coming soon",
  soonNote: "follow the changelog",
};

describe("SkillCard", () => {
  it("renders a packaged skill with the download CTA", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(SkillCard, {
      props: {
        ...baseProps,
        packaged: true,
        accent: "g",
        size: "33 KB · SKILL.md + references",
        downloadUrl: "/downloads/rgb-specialist-latest.zip",
      },
    });

    expect(html).toContain("Specialist");
    expect(html).toContain("Consults and interprets the rules");
    expect(html).toContain("v0.1 · contract defined");
    expect(html).toContain('href="/downloads/rgb-specialist-latest.zip"');
    expect(html).toContain("Download package (.zip)");
    expect(html).toContain("33 KB · SKILL.md + references");
    expect(html).not.toContain("Package coming soon");
  });

  it("renders a not-yet-packaged skill with the soon state instead of a download link", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(SkillCard, {
      props: {
        ...baseProps,
        name: "Maker",
        role: "Structures raw material",
        badge: "contract defined",
        packaged: false,
        accent: "r",
      },
    });

    expect(html).toContain("Maker");
    expect(html).toContain("Package coming soon");
    expect(html).toContain("follow the changelog");
    expect(html).not.toContain('class="dl"');
    expect(html).not.toContain("href=");
  });

  it("does not claim an installed runtime through the badge or description text", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(SkillCard, {
      props: { ...baseProps, packaged: true, accent: "g" },
    });

    expect(html.toLowerCase()).not.toContain("installed");
  });
});
