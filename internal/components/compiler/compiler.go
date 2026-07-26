// Package compiler marks the HTML and PDF generation orchestration boundary.
package compiler

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "compiler",
		Name:        "RGB Compiler",
		Description: "Future HTML and PDF generation orchestration boundary.",
	}
}
