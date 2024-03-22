package message

// TODO(lorenzo) refactor into a more efficient data structure
// e.g. 2 bit for each validator. 00 --> no key, 01 --> 1 key, 10 or 11 --> coefficient of key appended at the end of the initial array

type Coefficients []byte //TODO(lorenzo) evaluate if to do privatwe

func NewCoefficients(committeeSize uint64) Coefficients {
	return make([]byte, committeeSize)
}

func (c Coefficients) Merge(c1 Coefficients) {
	for i, coef := range c1 {
		c[i] = c[i] + coef
	}
}

func (c Coefficients) Increment(index uint64) {
	c[index]++
}

// TODO(lorenzo) not sure the name is really appropriate
func (c Coefficients) Flatten() []uint64 {
	var indexes []uint64
	for i, v := range c {
		for j := 0; j < int(v); j++ {
			indexes = append(indexes, uint64(i))
		}
	}
	return indexes
}

// TODO(lorenzo) not sure the name is really appropriate
func (c Coefficients) FlattenUniq() []uint64 {
	var indexes []uint64
	for i, v := range c {
		if v > 0 {
			indexes = append(indexes, uint64(i))
		}
	}
	return indexes
}
