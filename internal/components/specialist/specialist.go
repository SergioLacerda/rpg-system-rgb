// Package specialist marks the rule consultation and validation component boundary.
package specialist

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "specialist",
		Name:        "RGB Specialist",
		Description: "Rule consultation and validation boundary.",
	}
}
