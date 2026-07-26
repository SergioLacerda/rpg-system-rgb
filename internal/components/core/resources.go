package core

// Resources tracks a character's derived durability pool: health, shield,
// and equipped armor.
type Resources struct {
	MaxHealth     int
	CurrentHealth int
	MaxShield     int
	CurrentShield int
	Armor         int
}

// DeriveResources computes starting Resources from Vectors, per the
// canonical formulas in docs/core/en/core/character_creation.md:
// Health = 4 + R + B, Shield = B * 3.
func DeriveResources(vectors Vectors) Resources {
	health := 4 + vectors.R + vectors.B
	shield := vectors.B * 3
	return Resources{
		MaxHealth:     health,
		CurrentHealth: health,
		MaxShield:     shield,
		CurrentShield: shield,
	}
}

// IsDown reports whether the character has no health remaining.
func (resources Resources) IsDown() bool {
	return resources.CurrentHealth <= 0
}
