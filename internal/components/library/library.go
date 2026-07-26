package library

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

func Descriptor() components.Component {
	return components.Component{
		ID:          "library",
		Name:        "RGB Library",
		Description: "Publication-facing library boundary.",
	}
}
