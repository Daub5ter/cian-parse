package random

import (
	"math/rand"
	"time"
)

func RandomSec(min int, max int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(max-min+1) + min)
}
