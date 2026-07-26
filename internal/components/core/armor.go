package core

import "fmt"

// ArmorCategory identifies a fixed physical armor tier, per
// docs/core/en/equipment/armor.md.
type ArmorCategory string

const (
	// ArmorCategoryNone means no armor is equipped.
	ArmorCategoryNone ArmorCategory = "none"
	// ArmorCategoryLight is minimal protection with a low mobility penalty.
	ArmorCategoryLight ArmorCategory = "light"
	// ArmorCategoryMedium balances protection and mobility reduction.
	ArmorCategoryMedium ArmorCategory = "medium"
	// ArmorCategoryHeavy is high protection with a significant mobility penalty.
	ArmorCategoryHeavy ArmorCategory = "heavy"
)

// Validate reports whether the category is a known armor tier.
func (category ArmorCategory) Validate() error {
	switch category {
	case ArmorCategoryNone, ArmorCategoryLight, ArmorCategoryMedium, ArmorCategoryHeavy:
		return nil
	default:
		return fmt.Errorf("unknown armor category %q", category)
	}
}

// Protection returns the category's fixed armor value, per the Armor Types
// table in docs/core/en/equipment/armor.md.
func (category ArmorCategory) Protection() (int, error) {
	switch category {
	case ArmorCategoryNone:
		return 0, nil
	case ArmorCategoryLight:
		return 2, nil
	case ArmorCategoryMedium:
		return 4, nil
	case ArmorCategoryHeavy:
		return 6, nil
	default:
		return 0, fmt.Errorf("unknown armor category %q", category)
	}
}

// MobilityPenalty returns how much the category reduces effective G, per
// the Armor Types table in docs/core/en/equipment/armor.md.
func (category ArmorCategory) MobilityPenalty() (int, error) {
	switch category {
	case ArmorCategoryNone:
		return 0, nil
	case ArmorCategoryLight:
		return 1, nil
	case ArmorCategoryMedium:
		return 2, nil
	case ArmorCategoryHeavy:
		return 3, nil
	default:
		return 0, fmt.Errorf("unknown armor category %q", category)
	}
}
