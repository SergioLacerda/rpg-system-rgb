package core

import "fmt"

// AdvancementBudget is the fixed number of points a character spends per
// level, per docs/core/en/core/progression.md § Advancement Budget. Every
// advancement choice spends this same budget; none of them add to it.
const AdvancementBudget = 2

// SpecializationRequirementDiscount is the number of points ignored, per
// required vector, when checking an ability's Requirements against a
// character specialized in that ability's own Vector (progression.md §
// Specialization: "ignore one tier of requirements"). progression.md does
// not define a numeric formula linking Ability.Tier to Requirements point
// thresholds elsewhere in this package, so one point per vector is used
// here as the smallest deterministic unit Requirements already operates
// in.
const SpecializationRequirementDiscount = 1

// AdvancementChoice identifies which of the seven progression.md
// advancement choices a character took at a given level. Exactly one
// choice applies per level; choices do not stack within the same level.
type AdvancementChoice string

const (
	// AdvancementVectorGrowth is the default choice: +2 points distributed
	// across R, G, and B.
	AdvancementVectorGrowth AdvancementChoice = "vector_growth"
	// AdvancementNewAbility grants one ability whose Tier/Requirements the
	// character currently meets.
	AdvancementNewAbility AdvancementChoice = "new_ability"
	// AdvancementNewReaction grants one reaction-timed ability.
	AdvancementNewReaction AdvancementChoice = "new_reaction"
	// AdvancementSpecialization commits the character to one vector's
	// ability tree.
	AdvancementSpecialization AdvancementChoice = "specialization"
	// AdvancementAbilityImprovement modifies one known ability's cost,
	// duration, or limits, at the Game Master's discretion.
	AdvancementAbilityImprovement AdvancementChoice = "ability_improvement"
	// AdvancementNewResource grants a campaign-declared custom resource
	// pool.
	AdvancementNewResource AdvancementChoice = "new_resource"
	// AdvancementNewStateAccess grants access to one previously
	// unavailable tactical state or maneuver. It collapses into granting
	// an ability that produces that state (see AdvanceNewStateOrManeuverAccess)
	// — there is no separate per-character State gate in this package.
	AdvancementNewStateAccess AdvancementChoice = "new_state_or_maneuver_access"
)

// AdvancementRecord is one level's advancement choice, kept on Character
// for audit/history.
type AdvancementRecord struct {
	Level  int
	Choice AdvancementChoice
}

// recordAdvancement advances the character to the next level and appends
// an audit record for the choice taken. Every advancement method calls
// this exactly once, which is what enforces "choices do not stack within
// the same level": each call always moves to a new level.
func (character *Character) recordAdvancement(choice AdvancementChoice) {
	character.Level++
	character.AdvancementHistory = append(character.AdvancementHistory, AdvancementRecord{
		Level:  character.Level,
		Choice: choice,
	})
}

// MeetsRequirements reports whether the character's vectors satisfy the
// ability's Requirements, applying SpecializationRequirementDiscount when
// the character is specialized in the ability's own Vector.
func (character Character) MeetsRequirements(ability Ability) bool {
	specialized := character.Specializations[ability.Vector]
	for vector, required := range ability.Requirements {
		have, err := character.Vectors.Value(vector)
		if err != nil {
			return false
		}
		threshold := required
		if specialized {
			threshold -= SpecializationRequirementDiscount
			if threshold < 0 {
				threshold = 0
			}
		}
		if have < threshold {
			return false
		}
	}
	return true
}

// AdvanceVectorGrowth applies the default advancement choice: the fixed
// AdvancementBudget distributed across R, G, and B as the player chooses.
func (character *Character) AdvanceVectorGrowth(gain Vectors) error {
	if gain.R < 0 || gain.G < 0 || gain.B < 0 {
		return fmt.Errorf("vector growth must be non-negative per vector: R=%d G=%d B=%d", gain.R, gain.G, gain.B)
	}
	if total := gain.R + gain.G + gain.B; total != AdvancementBudget {
		return fmt.Errorf("vector growth must distribute exactly %d points, got %d", AdvancementBudget, total)
	}
	grown := Vectors{
		R: character.Vectors.R + gain.R,
		G: character.Vectors.G + gain.G,
		B: character.Vectors.B + gain.B,
	}
	if err := grown.Validate(); err != nil {
		return err
	}
	character.Vectors = grown
	character.recordAdvancement(AdvancementVectorGrowth)
	return nil
}

// grantAbility validates and appends an ability granted by an advancement
// choice, gated by MeetsRequirements. AdvanceNewAbility, AdvanceNewReaction,
// and AdvanceNewStateOrManeuverAccess all share this path — they differ
// only in the AdvancementChoice recorded and, for reactions, the required
// Timing.
func (character *Character) grantAbility(ability Ability, choice AdvancementChoice) error {
	if err := ability.Validate(); err != nil {
		return fmt.Errorf("advancement ability invalid: %w", err)
	}
	if !character.MeetsRequirements(ability) {
		return fmt.Errorf("character %s does not meet requirements for ability %q", character.ID, ability.ID)
	}
	character.Abilities = append(character.Abilities, ability)
	character.recordAdvancement(choice)
	return nil
}

// AdvanceNewAbility grants one ability whose Tier/Requirements the
// character currently meets.
func (character *Character) AdvanceNewAbility(ability Ability) error {
	return character.grantAbility(ability, AdvancementNewAbility)
}

// AdvanceNewReaction grants one reaction-timed ability whose
// Tier/Requirements the character currently meets.
func (character *Character) AdvanceNewReaction(ability Ability) error {
	if ability.Timing != TimingReaction {
		return fmt.Errorf("new-reaction advancement requires an ability with Timing %q, got %q", TimingReaction, ability.Timing)
	}
	return character.grantAbility(ability, AdvancementNewReaction)
}

// AdvanceNewStateOrManeuverAccess grants access to one previously
// unavailable tactical state or maneuver. Per this mission's U-03
// resolution, it collapses into granting an ability that produces the
// target state via its documented procedure (see state_lifecycle.go) —
// there is no separate per-character State gate; every State already has
// an unconditional producing procedure available to any character that
// holds the right ability.
func (character *Character) AdvanceNewStateOrManeuverAccess(ability Ability) error {
	return character.grantAbility(ability, AdvancementNewStateAccess)
}

// AdvanceSpecialization commits the character to one vector's ability
// tree. Future AdvanceNewAbility/AdvanceNewReaction/
// AdvanceNewStateOrManeuverAccess calls for abilities of the same Vector
// apply SpecializationRequirementDiscount via MeetsRequirements.
// Specialization does not itself grant an ability.
func (character *Character) AdvanceSpecialization(vector Vector) error {
	if err := vector.Validate(); err != nil {
		return err
	}
	if character.Specializations[vector] {
		return fmt.Errorf("character %s is already specialized in %s", character.ID, vector)
	}
	if character.Specializations == nil {
		character.Specializations = map[Vector]bool{}
	}
	character.Specializations[vector] = true
	character.recordAdvancement(AdvancementSpecialization)
	return nil
}

// ReplaceAbility swaps out an existing ability (matched by ID) for a
// modified version, per the "Ability improvement" advancement choice. Per
// this mission's U-01 resolution, this is intentionally free-form: the
// engine does not validate how much the modified ability differs from the
// original — reducing cost, extending duration, or raising a limits
// ceiling is "at the Game Master's discretion" (progression.md) and is
// trusted entirely to the caller. The replacement must still pass
// Ability.Validate() as a normal, well-formed ability.
func (character *Character) ReplaceAbility(modified Ability) error {
	if err := modified.Validate(); err != nil {
		return fmt.Errorf("improved ability invalid: %w", err)
	}
	for i, existing := range character.Abilities {
		if existing.ID == modified.ID {
			character.Abilities[i] = modified
			character.recordAdvancement(AdvancementAbilityImprovement)
			return nil
		}
	}
	return fmt.Errorf("character %s does not have ability %q to improve", character.ID, modified.ID)
}

// ResourceDefinition is a campaign-declared custom resource pool, per the
// "New resource" advancement choice (progression.md: a resource pool
// "sized by the Game Master"). Per this mission's U-02 resolution, a
// campaign must declare a ResourceDefinition before it can be granted to
// a character — this mirrors how Ability itself is a declared contract,
// rather than an untyped escape hatch.
type ResourceDefinition struct {
	Name string
	Max  int
}

// Validate reports whether the resource definition is well-formed.
func (definition ResourceDefinition) Validate() error {
	if definition.Name == "" {
		return fmt.Errorf("resource definition name must be non-empty")
	}
	if definition.Max <= 0 {
		return fmt.Errorf("resource definition %q max must be positive", definition.Name)
	}
	return nil
}

// GrantResource grants the character a custom resource pool at full
// value, per a campaign-declared ResourceDefinition.
func (character *Character) GrantResource(definition ResourceDefinition) error {
	if err := definition.Validate(); err != nil {
		return fmt.Errorf("cannot grant undeclared resource: %w", err)
	}
	if character.CustomResources == nil {
		character.CustomResources = map[string]int{}
	}
	character.CustomResources[definition.Name] = definition.Max
	character.recordAdvancement(AdvancementNewResource)
	return nil
}
