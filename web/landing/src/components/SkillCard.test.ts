import { experimental_AstroContainer as AstroContainer } from "astro/container";
import { describe, expect, it } from "vitest";
import SkillCard from "./SkillCard.astro";

describe("SkillCard", () => {
  it("renders implemented skills with the enabled state", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(SkillCard, {
      props: {
        name: "Library",
        status: "installed",
        installed: true,
        desc: "Answers rules questions.",
        tags: ["Library", "R G B"],
      },
    });

    expect(html).toContain("Skill: Library");
    expect(html).toContain("● installed");
    expect(html).toContain("Answers rules questions.");
    expect(html).toContain("Library");
    expect(html).toContain("background:var(--g-bg);color:var(--g)");
  });

  it("renders contract-defined skills with the disabled state", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(SkillCard, {
      props: {
        name: "Specialist",
        status: "contract defined",
        installed: false,
        desc: "Runtime behavior is not implemented.",
        tags: ["draft"],
      },
    });

    expect(html).toContain("Skill: Specialist");
    expect(html).toContain("○ contract defined");
    expect(html).toContain("Runtime behavior is not implemented.");
    expect(html).toContain("background:var(--r-bg);color:var(--r)");
  });
});
