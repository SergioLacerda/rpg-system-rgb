// Package bundles marks the machine-consumable bundle output boundary.
package bundles

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "bundles",
		Name:        "RGB Bundles",
		Description: "Future machine-consumable bundle output boundary.",
	}
}
