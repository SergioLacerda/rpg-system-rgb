// Package maker marks the campaign and content structuring component boundary.
package maker

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "maker",
		Name:        "RGB Maker",
		Description: "Campaign and content structuring boundary.",
	}
}
