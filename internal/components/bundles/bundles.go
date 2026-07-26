package bundles

import "github.com/SergioLacerda/rpg-system-rgb/internal/components"

func Descriptor() components.Component {
	return components.Component{
		ID:          "bundles",
		Name:        "RGB Bundles",
		Description: "Future machine-consumable bundle output boundary.",
	}
}
