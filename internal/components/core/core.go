package core

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

func Descriptor() components.Component {
	return components.Component{
		ID:          "core",
		Name:        "RGB Core",
		Description: "Playable RGB model boundary.",
	}
}
