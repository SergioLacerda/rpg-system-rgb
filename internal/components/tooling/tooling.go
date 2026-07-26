package tooling

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

func Descriptor() components.Component {
	return components.Component{
		ID:          "tooling",
		Name:        "RGB Tooling",
		Description: "Validators, indexers, compilers, and maintenance command boundary.",
	}
}
