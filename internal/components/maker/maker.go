package maker

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

func Descriptor() components.Component {
	return components.Component{
		ID:          "maker",
		Name:        "RGB Maker",
		Description: "Campaign and content structuring boundary.",
	}
}
