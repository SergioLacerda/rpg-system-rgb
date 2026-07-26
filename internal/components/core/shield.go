package core

import "fmt"

// ShieldCategory identifies a fixed physical shield tier, per
// docs/core/en/equipment/shields.md. This is distinct from the energy/
// magical shield reserve (Resources.CurrentShield, derived from B — see
// resources.go); the physical shields table is the only one of the two
// with a mobility penalty.
type ShieldCategory string

const (
	// ShieldCategoryNone means no physical shield is equipped.
	ShieldCategoryNone ShieldCategory = "none"
	// ShieldCategoryLight adds minimal protection with no mobility penalty.
	ShieldCategoryLight ShieldCategory = "light"
	// ShieldCategoryMedium adds moderate protection with a small mobility penalty.
	ShieldCategoryMedium ShieldCategory = "medium"
	// ShieldCategoryHeavy adds high protection with a significant mobility penalty.
	ShieldCategoryHeavy ShieldCategory = "heavy"
)

// Validate reports whether the category is a known physical shield tier.
func (category ShieldCategory) Validate() error {
	switch category {
	case ShieldCategoryNone, ShieldCategoryLight, ShieldCategoryMedium, ShieldCategoryHeavy:
		return nil
	default:
		return fmt.Errorf("unknown shield category %q", category)
	}
}

// Protection returns the category's fixed shield value, per the Physical
// Shields table in docs/core/en/equipment/shields.md.
func (category ShieldCategory) Protection() (int, error) {
	switch category {
	case ShieldCategoryNone:
		return 0, nil
	case ShieldCategoryLight:
		return 1, nil
	case ShieldCategoryMedium:
		return 2, nil
	case ShieldCategoryHeavy:
		return 3, nil
	default:
		return 0, fmt.Errorf("unknown shield category %q", category)
	}
}

// MobilityPenalty returns how much the category reduces effective G, per
// the Physical Shields table in docs/core/en/equipment/shields.md.
func (category ShieldCategory) MobilityPenalty() (int, error) {
	switch category {
	case ShieldCategoryNone, ShieldCategoryLight:
		return 0, nil
	case ShieldCategoryMedium:
		return 1, nil
	case ShieldCategoryHeavy:
		return 2, nil
	default:
		return 0, fmt.Errorf("unknown shield category %q", category)
	}
}
