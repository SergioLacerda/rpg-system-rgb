package fixtures

import (
	"os"
	"strings"
	"testing"
)

func TestCoreV2FixtureFilesExposeReusableIdentities(t *testing.T) {
	cases := []struct {
		file     string
		required []string
	}{
		{
			file: "characters/core-v2.yaml",
			required: []string{
				"id: red-vanguard",
				"id: green-runner",
				"id: blue-warden",
				"id: rgb-balanced",
			},
		},
		{
			file: "encounters/laboratory-evacuation.yaml",
			required: []string{
				"id: laboratory-evacuation",
				"actor: blue-warden",
				"actor: green-runner",
			},
		},
		{
			file: "examples/combat-example.yaml",
			required: []string{
				"id: attack-1-a-fires-rifle-at-b",
				"id: attack-2-b-fires-smg-at-a",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			body, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			text := string(body)
			for _, required := range tc.required {
				if !strings.Contains(text, required) {
					t.Fatalf("%s missing %q", tc.file, required)
				}
			}
		})
	}
}
