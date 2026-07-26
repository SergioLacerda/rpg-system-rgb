// Package components holds the shared contract types that every RGB
// component boundary exposes to the application layer.
package components

// Component describes one bounded component of the RGB System V2 scaffold.
type Component struct {
	ID          string
	Name        string
	Description string
}
