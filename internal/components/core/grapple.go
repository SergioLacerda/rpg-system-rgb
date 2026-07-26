package core

// ResolveGrapple reports whether a grapple attempt succeeds, per
// docs/core/en/combat/attack_and_defense.md: the attacker must match or
// exceed the target in both G (mobility) and R (pressure).
func ResolveGrapple(attackerG, attackerR, targetG, targetR int) bool {
	return attackerG >= targetG && attackerR >= targetR
}
