package stubs

import (
	"log"
	"math"
	"math/rand"
)

const (
	// TODO: some more work to ensure the full range of types is explored
	defaultMaxFloat64 = math.MaxFloat64 / 2
	defaultMinFloat64 = -math.MaxFloat64 / 2
	defaultMaxFloat32 = math.MaxFloat32 / 2
	defaultMinFloat32 = -math.MaxFloat32 / 2
	defaultMaxInt64   = math.MaxInt64>>1 - 1
	defaultMinInt64   = -math.MaxInt64 >> 1
	defaultMaxInt32   = math.MaxInt32>>1 - 1
	defaultMinInt32   = -math.MaxInt32 >> 1
	defaultMaxUint64  = math.MaxUint64>>1 - 1
	defaultMinUint64  = uint64(0)
	defaultMaxUint32  = math.MaxUint32>>1 - 1
	defaultMinUint32  = uint32(0)
)

type randomGenerator struct {
	// TODO: rand generator in struct
	defaultSeeder
	// TODO: custom boundaries
	// TODO: non uniform distributions
}

// randomO:pts describes options for the source of entropy
type randomOpts struct {
	AutoSeed bool
}

func newRandomGenerator(opts randomOpts) *randomGenerator {
	r := randomGenerator{}
	r.autoseed = opts.AutoSeed
	debugLog("autoseed: %t", r.autoseed)
	return &r
}

// IntN returns a random integer
func (r *randomGenerator) IntN(n int) int {
	if r.autoseed {
		r.AutoSeed()
	}
	return rand.Intn(n)
}

func (r *randomGenerator) Bool() bool {
	if r.autoseed {
		r.AutoSeed()
	}
	return rand.Intn(2) == 1
}

func (r *randomGenerator) Sign() int {
	if rand.Intn(2) == 1 {
		return 1
	}
	return -1
}

// TODO: random strings for Base64 and passwords

func (r *randomGenerator) Float64(min, max float64) float64 {
	if min > max {
		return 0
	}
	if min == 0 && max == 0 {
		min = defaultMinFloat64
		max = defaultMaxFloat64
	}
	if min == max {
		return min
	}
	if r.autoseed {
		r.AutoSeed()
	}
	log.Printf("DEBUG: min: %v, max: %v", min, max)
	return min + rand.Float64()*(max-min)
}

func (r *randomGenerator) Float32(min, max float32) float32 {
	if min > max {
		return 0
	}
	if min == 0 && max == 0 {
		min = defaultMinFloat32
		max = defaultMaxFloat32
	}
	if min == max {
		return min
	}
	if r.autoseed {
		r.AutoSeed()
	}
	return min + rand.Float32()*(max-min)
}

func (r *randomGenerator) Int64(min, max int64) int64 {
	if min > max {
		return 0
	}
	if min == 0 && max == 0 {
		min = defaultMinInt64
		max = defaultMaxInt64
	}
	if min == max {
		return min
	}
	if r.autoseed {
		r.AutoSeed()
	}
	return min + rand.Int63n(max-min+1)
}

func (r *randomGenerator) Int32(min, max int32) int32 {
	if min > max {
		return 0
	}
	if min == 0 && max == 0 {
		min = defaultMinInt32
		max = defaultMaxInt32
	}
	if min == max {
		return min
	}
	if r.autoseed {
		r.AutoSeed()
	}
	return min + rand.Int31n(max-min+int32(1))
}

func (r *randomGenerator) Uint64(min, max uint64) uint64 {
	if min > max {
		return 0
	}
	if min == 0 && max == 0 {
		max = defaultMaxUint64>>1 - uint64(1)
	}
	if min == max {
		return min
	}
	if r.autoseed {
		r.AutoSeed()
	}
	if max-min > defaultMaxInt64 {
		// TODO: make range > maxInt32 reachable
		debugLog("warning unsupported range")
	}
	return min + uint64(rand.Int63n(int64(max-min+uint64(1))))
}

func (r *randomGenerator) Uint32(min, max uint32) uint32 {
	if min > max {
		return 0
	}
	if min == 0 && max == 0 {
		max = defaultMaxUint32
	}
	if min == max {
		return min
	}
	if r.autoseed {
		r.AutoSeed()
	}
	if max-min > defaultMaxInt32 {
		// TODO: make range > maxInt32 reachable
		debugLog("warning unsupported range")
	}
	// TODO: make range > maxInt32 reachable
	return min + uint32(rand.Int31n(int32(max-min+uint32(1))))
}
