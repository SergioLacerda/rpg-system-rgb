// Package library marks the publication-facing library component boundary.
package library

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "library",
		Name:        "RGB Library",
		Description: "Publication-facing library boundary.",
	}
}
