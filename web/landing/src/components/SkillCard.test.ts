import { experimental_AstroContainer as AstroContainer } from "astro/container";
import { describe, expect, it } from "vitest";
import SkillCard from "./SkillCard.astro";

describe("SkillCard", () => {
  it("renders installed skills with the enabled state", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(SkillCard, {
      props: {
        name: "Pathfinder",
        status: "installed",
        installed: true,
        desc: "Answers rules questions.",
        tags: ["Library", "R G B"],
      },
    });

    expect(html).toContain("Skill: Pathfinder");
    expect(html).toContain("● installed");
    expect(html).toContain("Answers rules questions.");
    expect(html).toContain("Library");
    expect(html).toContain("background:var(--g-bg);color:var(--g)");
  });

  it("renders unavailable skills with the disabled state", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(SkillCard, {
      props: {
        name: "Maker",
        status: "not installed",
        installed: false,
        desc: "Drafts structured documents.",
        tags: ["draft"],
      },
    });

    expect(html).toContain("Skill: Maker");
    expect(html).toContain("○ not installed");
    expect(html).toContain("Drafts structured documents.");
    expect(html).toContain("background:var(--r-bg);color:var(--r)");
  });
});
