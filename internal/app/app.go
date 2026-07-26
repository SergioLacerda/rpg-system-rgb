package app

import (
	"fmt"
	"strings"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/bundles"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/compiler"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/core"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/library"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/maker"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/specialist"
	"github.com/SergioLacerda/rpg-system-rgb/internal/components/tooling"
)

func Components() []components.Component {
	return []components.Component{
		core.Descriptor(),
		maker.Descriptor(),
		specialist.Descriptor(),
		tooling.Descriptor(),
		library.Descriptor(),
		compiler.Descriptor(),
		bundles.Descriptor(),
	}
}

func Hello() string {
	var builder strings.Builder
	builder.WriteString("RGB System V2 scaffold ready")
	for _, component := range Components() {
		builder.WriteString(fmt.Sprintf("\n- %s: %s", component.ID, component.Name))
	}
	return builder.String()
}
