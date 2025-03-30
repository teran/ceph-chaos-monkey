package monkey

import (
	"math/rand"
	"time"
)

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func getRandomChoice(max int) int {
	return rng.Intn(max)
}
