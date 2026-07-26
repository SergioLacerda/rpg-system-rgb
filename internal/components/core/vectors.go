package core

import "fmt"

type Vector string

const (
	VectorR Vector = "R"
	VectorG Vector = "G"
	VectorB Vector = "B"
)

type Vectors struct {
	R int
	G int
	B int
}

func (vectors Vectors) Validate() error {
	if vectors.R < 0 || vectors.G < 0 || vectors.B < 0 {
		return fmt.Errorf("vectors must be non-negative: R=%d G=%d B=%d", vectors.R, vectors.G, vectors.B)
	}
	return nil
}

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

func (vector Vector) Validate() error {
	switch vector {
	case VectorR, VectorG, VectorB:
		return nil
	default:
		return fmt.Errorf("unknown vector %q", vector)
	}
}

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
