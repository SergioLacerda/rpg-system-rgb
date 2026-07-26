package core

import "sort"

// InitiativeEntry pairs an actor ID with its G value and whether it acts
// with surprise, for use with InitiativeOrder.
type InitiativeEntry struct {
	ActorID  string
	G        int
	Surprise bool
}

// InitiativeOrder returns actor IDs in acting order, per
// docs/core/en/combat/attack_and_defense.md: a surprise attack ignores
// initiative and acts first; otherwise the higher G value acts first.
func InitiativeOrder(entries []InitiativeEntry) []string {
	sorted := make([]InitiativeEntry, len(entries))
	copy(sorted, entries)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Surprise != sorted[j].Surprise {
			return sorted[i].Surprise
		}
		return sorted[i].G > sorted[j].G
	})

	order := make([]string, len(sorted))
	for i, entry := range sorted {
		order[i] = entry.ActorID
	}
	return order
}
