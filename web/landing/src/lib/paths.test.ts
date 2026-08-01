import { describe, expect, it } from "vitest";
import {
  landingPath,
  localizedLandingPath,
  localizedRestPath,
  normalizeLandingBase,
  pathWithoutLandingBase,
} from "./paths";

describe("landing paths", () => {
  it("normalizes configured bases before composing URLs", () => {
    expect(normalizeLandingBase("/")).toBe("");
    expect(normalizeLandingBase("/rpg-system-rgb/")).toBe("/rpg-system-rgb");
    expect(normalizeLandingBase("/rpg-system-rgb")).toBe("/rpg-system-rgb");
  });

  it("prefixes paths with the configured landing base", () => {
    expect(landingPath("/")).toBe("/rpg-system-rgb/");
    expect(landingPath("/library/")).toBe("/rpg-system-rgb/library/");
    expect(landingPath("downloads/file.pdf")).toBe(
      "/rpg-system-rgb/downloads/file.pdf",
    );
  });

  it("removes the configured landing base from current URLs", () => {
    expect(pathWithoutLandingBase("/rpg-system-rgb")).toBe("/");
    expect(pathWithoutLandingBase("/rpg-system-rgb/en/")).toBe("/en/");
    expect(pathWithoutLandingBase("/external/en/")).toBe("/external/en/");
  });

  it("builds locale-aware home links with Portuguese as the unprefixed default", () => {
    expect(localizedLandingPath("pt-br")).toBe("/rpg-system-rgb/");
    expect(localizedLandingPath("en")).toBe("/rpg-system-rgb/en/");
  });

  it("preserves nested page paths when switching locales", () => {
    expect(localizedRestPath("/rpg-system-rgb/pt-br/skills/pathfinder")).toBe(
      "skills/pathfinder",
    );
    expect(localizedRestPath("/rpg-system-rgb/en/skills/pathfinder")).toBe(
      "skills/pathfinder",
    );
    expect(localizedLandingPath("pt-br", "skills/pathfinder")).toBe(
      "/rpg-system-rgb/pt-br/skills/pathfinder",
    );
    expect(localizedLandingPath("en", "skills/pathfinder")).toBe(
      "/rpg-system-rgb/en/skills/pathfinder",
    );
  });
});
