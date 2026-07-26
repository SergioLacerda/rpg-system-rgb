package corebehavior

import (
	"testing"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
)

// mirrors: tests/features/equipment/mobility_penalty.feature#Armor reduces effective G
func TestArmorReducesEffectiveGFeatureExamples(t *testing.T) {
	cases := []struct {
		armor     core.ArmorCategory
		effective int
	}{
		{armor: core.ArmorCategoryNone, effective: 6},
		{armor: core.ArmorCategoryLight, effective: 5},
		{armor: core.ArmorCategoryMedium, effective: 4},
		{armor: core.ArmorCategoryHeavy, effective: 3},
	}
	for _, tc := range cases {
		got, err := core.EffectiveG(6, tc.armor, core.ShieldCategoryNone)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.effective {
			t.Fatalf("armor=%s effective G got %d want %d", tc.armor, got, tc.effective)
		}
	}
}

// mirrors: tests/features/equipment/mobility_penalty.feature#Shield reduces effective G
func TestShieldReducesEffectiveGFeatureExamples(t *testing.T) {
	cases := []struct {
		shield    core.ShieldCategory
		effective int
	}{
		{shield: core.ShieldCategoryNone, effective: 6},
		{shield: core.ShieldCategoryLight, effective: 6},
		{shield: core.ShieldCategoryMedium, effective: 5},
		{shield: core.ShieldCategoryHeavy, effective: 4},
	}
	for _, tc := range cases {
		got, err := core.EffectiveG(6, core.ArmorCategoryNone, tc.shield)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.effective {
			t.Fatalf("shield=%s effective G got %d want %d", tc.shield, got, tc.effective)
		}
	}
}

// mirrors: tests/features/equipment/mobility_penalty.feature#Heavy armor and heavy shield together floor at zero
func TestCombinedPenaltiesFloorAtZeroFeatureExample(t *testing.T) {
	got, err := core.EffectiveG(2, core.ArmorCategoryHeavy, core.ShieldCategoryHeavy)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Fatalf("effective G got %d want 0", got)
	}
}
