package core

// MovementDistance returns the maximum distance, in meters, a character
// with the given G value can move in one turn, per
// docs/core/en/combat/movement.md: ground and aerial movement share the
// same formula, G * 2 meters, on a 1 square = 1 meter grid.
func MovementDistance(g int) int {
	return g * 2
}
