package core

// MovementDistance returns the maximum distance, in meters, a character
// with the given G value can move in one turn, per
// docs/core/en/combat/movement.md: ground and aerial movement share the
// same formula, G * 2 meters, on a 1 square = 1 meter grid.
func MovementDistance(g int) int {
	return g * 2
}

// EffectiveG returns G reduced by equipped armor and shield mobility
// penalties, floored at zero, per docs/core/en/equipment/armor.md and
// docs/core/en/equipment/shields.md. Callers combine EffectiveG with
// MovementDistance when equipment applies; MovementDistance itself stays a
// pure G * 2 formula.
func EffectiveG(g int, armor ArmorCategory, shield ShieldCategory) (int, error) {
	armorPenalty, err := armor.MobilityPenalty()
	if err != nil {
		return 0, err
	}
	shieldPenalty, err := shield.MobilityPenalty()
	if err != nil {
		return 0, err
	}
	return max(g-armorPenalty-shieldPenalty, 0), nil
}
