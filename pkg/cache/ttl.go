package cache

import (
	"math/rand"
	"time"
)

func JitterTTL(base time.Duration) time.Duration {
	jitter := base / 10

	if jitter <= 0 {
		jitter = time.Second
	}

	return base + time.Duration(rand.Int63n(int64(jitter)))
}
