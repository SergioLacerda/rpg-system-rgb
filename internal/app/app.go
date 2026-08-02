// Package app is the application layer: it wires component boundaries
// together and exposes use cases to the entrypoints in cmd.
package app

import (
	"fmt"
	"strings"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/bundles"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/maker"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/specialist"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/tooling"
)

// Components returns the registered component boundaries of the scaffold.
func Components() []components.Component {
	return []components.Component{
		core.Descriptor(),
		maker.Descriptor(),
		specialist.Descriptor(),
		tooling.Descriptor(),
		bundles.Descriptor(),
	}
}

// Hello reports the scaffold status and its registered components.
func Hello() string {
	var builder strings.Builder
	builder.WriteString("RGB System V2 scaffold ready")
	for _, component := range Components() {
		fmt.Fprintf(&builder, "\n- %s: %s", component.ID, component.Name)
	}
	return builder.String()
}
