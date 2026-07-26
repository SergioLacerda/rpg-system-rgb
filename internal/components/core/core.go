// Package core marks the playable RGB model component boundary.
package core

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "core",
		Name:        "RGB Core",
		Description: "Playable RGB model boundary.",
	}
}
