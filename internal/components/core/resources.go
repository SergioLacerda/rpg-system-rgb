package core

type Resources struct {
	MaxHealth     int
	CurrentHealth int
	MaxShield     int
	CurrentShield int
	Armor         int
}

func DeriveResources(vectors Vectors) Resources {
	health := 4 + vectors.R + vectors.B
	shield := vectors.B
	return Resources{
		MaxHealth:     health,
		CurrentHealth: health,
		MaxShield:     shield,
		CurrentShield: shield,
	}
}

func (resources Resources) IsDown() bool {
	return resources.CurrentHealth <= 0
}
