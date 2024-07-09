package fixtures

import (
	"math/rand"
)

// randomChoice returns a random element from the given slice.
func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

// randomBool returns a random boolean value.
func randomBool() bool {
	return rand.Intn(2) == 1
}
