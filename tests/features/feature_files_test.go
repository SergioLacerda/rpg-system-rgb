package features

import (
	"os"
	"strings"
	"testing"
)

func TestCoreV2FeatureFilesExist(t *testing.T) {
	files := []string{
		"combat/attack_evasion.feature",
		"damage/armor_shield_absorption.feature",
		"core/rgb_obstacle_approaches.feature",
		"encounters/laboratory_evacuation.feature",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			body, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			text := string(body)
			for _, required := range []string{"Feature:", "Scenario"} {
				if !strings.Contains(text, required) {
					t.Fatalf("%s missing %q", file, required)
				}
			}
		})
	}
}
