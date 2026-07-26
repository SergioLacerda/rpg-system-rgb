package specialist

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

func Descriptor() components.Component {
	return components.Component{
		ID:          "specialist",
		Name:        "RGB Specialist",
		Description: "Rule consultation and validation boundary.",
	}
}
