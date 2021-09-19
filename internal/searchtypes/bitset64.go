//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

// Bitset64 of types
type Bitset64 uint64

// NewBitset64 from types
func NewBitset64(types ...uint) (b Bitset64) {
	b.Set(types...)
	return
}

// Set type values
func (b *Bitset64) Set(types ...uint) {
	for _, t := range types {
		*b |= 1 << Bitset64(t%64)
	}
}

// Unset type values
func (b *Bitset64) Unset(types ...uint) {
	for _, t := range types {
		*b &= ^(1 << Bitset64(t%64))
	}
}

// Has type in bitset
func (b *Bitset64) Has(t uint) bool {
	return b != nil && *b&(1<<Bitset64(t%64)) != 0
}

// Numbers from bites
func (b Bitset64) Numbers() (list []uint) {
	for i := Bitset64(0); i < 64; i++ {
		if b&(1<<i) != 0 {
			list = append(list, uint(i))
		}
	}
	return
}
