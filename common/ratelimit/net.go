package ratelimit

import (
	"net"
)

// NewPipes adding or sending `rate` B every second, holding max capacity B
func NewPipes(rate float64, capacity int64) (net.Conn, net.Conn) {
	return NewPipesWithClock(rate, capacity, nil)
}

func NewPipesWithClock(rate float64, capacity int64, clock Clock) (net.Conn, net.Conn) {
	bucket1 := NewBucketWithRateAndClock(rate, capacity, clock)
	bucket2 := NewBucketWithRateAndClock(rate, capacity, clock)

	p1, p2 := net.Pipe()

	ratedPipe1 := Conn(p1, bucket1)
	ratedPipe2 := Conn(p2, bucket2)

	return ratedPipe1, ratedPipe2
}

func NewPipesWithRates(rate1, rate2 float64, capacity1, capacity2 int64, clock1, clock2 Clock) (net.Conn, net.Conn) {
	bucket1 := NewBucketWithRateAndClock(rate1, capacity1, clock1)
	bucket2 := NewBucketWithRateAndClock(rate2, capacity2, clock2)

	p1, p2 := net.Pipe()

	ratedPipe1 := Conn(p1, bucket1)
	ratedPipe2 := Conn(p2, bucket2)

	return ratedPipe1, ratedPipe2
}
