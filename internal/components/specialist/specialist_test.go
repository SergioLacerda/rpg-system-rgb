package specialist

import "testing"

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID != "specialist" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("unexpected descriptor: %+v", descriptor)
	}
}
