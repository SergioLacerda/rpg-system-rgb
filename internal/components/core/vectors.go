package core

import "fmt"

// Vector identifies one of the three RGB Core axes: R, G, or B.
type Vector string

const (
	// VectorR transforms the source of pressure (teaching label: attack).
	VectorR Vector = "R"
	// VectorG changes the actor's relation to pressure (teaching label: move).
	VectorG Vector = "G"
	// VectorB preserves continuity under pressure (teaching label: defend).
	VectorB Vector = "B"
)

// Vectors holds a character's R, G, and B point values.
type Vectors struct {
	R int
	G int
	B int
}

// Validate reports whether all three vectors are non-negative.
func (vectors Vectors) Validate() error {
	if vectors.R < 0 || vectors.G < 0 || vectors.B < 0 {
		return fmt.Errorf("vectors must be non-negative: R=%d G=%d B=%d", vectors.R, vectors.G, vectors.B)
	}
	return nil
}

// Value returns the point value for the named vector.
func (vectors Vectors) Value(vector Vector) (int, error) {
	switch vector {
	case VectorR:
		return vectors.R, nil
	case VectorG:
		return vectors.G, nil
	case VectorB:
		return vectors.B, nil
	default:
		return 0, fmt.Errorf("unknown vector %q", vector)
	}
}

// Validate reports whether the vector is one of R, G, or B.
func (vector Vector) Validate() error {
	switch vector {
	case VectorR, VectorG, VectorB:
		return nil
	default:
		return fmt.Errorf("unknown vector %q", vector)
	}
}

// TeachingLabel returns the onboarding shorthand for the vector
// (attack / move / defend).
func (vector Vector) TeachingLabel() string {
	switch vector {
	case VectorR:
		return "attack"
	case VectorG:
		return "move"
	case VectorB:
		return "defend"
	default:
		return "unknown"
	}
}

// NormativeLabel returns the adjudication-grade description of the vector's
// role, used to classify ambiguous actions.
func (vector Vector) NormativeLabel() string {
	switch vector {
	case VectorR:
		return "transform pressure source"
	case VectorG:
		return "change relation to pressure"
	case VectorB:
		return "preserve continuity under pressure"
	default:
		return "unknown"
	}
}
