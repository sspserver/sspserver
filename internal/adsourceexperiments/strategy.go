package adsourceexperiments

import (
	"math/rand"
	"sync/atomic"
	"time"
)

// IndexStrategy interface
type IndexStrategy interface {
	GetIndex() uint
}

// Counter strategy wrapper

type counterStrategy struct {
	counter uint32
}

// NewCounterStrategy interface
func NewCounterStrategy() IndexStrategy {
	return &counterStrategy{}
}

func (st *counterStrategy) GetIndex() uint {
	return uint(atomic.AddUint32(&st.counter, 1))
}

var _ IndexStrategy = (*counterStrategy)(nil)

// Random strategy wrapper
type randomStrategy struct {
	rnd *rand.Rand
}

// NewRandomStrategy interface
func NewRandomStrategy() IndexStrategy {
	return &randomStrategy{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (st *randomStrategy) GetIndex() uint {
	return uint(rand.Uint32())
}

var _ IndexStrategy = (*randomStrategy)(nil)
