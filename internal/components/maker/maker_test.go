package maker

import "testing"

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID != "maker" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("unexpected descriptor: %+v", descriptor)
	}
}
