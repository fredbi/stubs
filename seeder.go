package stubs

import (
	"math/rand"
	"time"
)

// seeder knows how to set random seeds
type seeder interface {
	SetSeed(int64)
	AutoSeed()
}

type defaultSeeder struct {
	autoseed bool
	seed     int64
}

func (d *defaultSeeder) AutoSeed() {
	rand.Seed(time.Now().UnixNano())
}

func (d *defaultSeeder) SetSeed(seed int64) {
	rand.Seed(seed)
}
