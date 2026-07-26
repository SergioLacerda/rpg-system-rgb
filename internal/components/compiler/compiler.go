package compiler

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

func Descriptor() components.Component {
	return components.Component{
		ID:          "compiler",
		Name:        "RGB Compiler",
		Description: "Future HTML and PDF generation orchestration boundary.",
	}
}
