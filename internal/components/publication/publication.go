// Package publication owns Go-based public documentation publication.
package publication

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "publication",
		Name:        "RGB Publication",
		Description: "Go-owned HTML Library and PDF publication boundary.",
	}
}
