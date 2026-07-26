// Package tooling marks the validators, indexers, and maintenance command boundary.
package tooling

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "tooling",
		Name:        "RGB Tooling",
		Description: "Validators, indexers, compilers, and maintenance command boundary.",
	}
}
