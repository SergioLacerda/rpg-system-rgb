package app

import (
	"reflect"
	"strings"
	"testing"
)

func TestComponentsExposeExpectedBoundaries(t *testing.T) {
	var ids []string
	for _, component := range Components() {
		if component.ID == "" {
			t.Fatal("component ID must be non-empty")
		}
		if component.Name == "" {
			t.Fatalf("%s component name must be non-empty", component.ID)
		}
		if component.Description == "" {
			t.Fatalf("%s component description must be non-empty", component.ID)
		}
		ids = append(ids, component.ID)
	}

	expected := []string{
		"core",
		"maker",
		"specialist",
		"tooling",
		"bundles",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected component IDs: got %v want %v", ids, expected)
	}
}

func TestHelloOutputIsDeterministic(t *testing.T) {
	output := Hello()
	if !strings.HasPrefix(output, "RGB System V2 scaffold ready") {
		t.Fatalf("unexpected hello prefix: %q", output)
	}
	for _, id := range []string{"core", "maker", "specialist", "tooling", "bundles"} {
		if !strings.Contains(output, "- "+id+":") {
			t.Fatalf("hello output missing component %s:\n%s", id, output)
		}
	}
}
