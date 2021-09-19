//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"sort"
)

// UIntBitset any numbers
type UIntBitset struct {
	values []uint
	mask   uint64
}

// NewUIntBitset from numbers
func NewUIntBitset(vals ...uint) (b *UIntBitset) {
	return (&UIntBitset{}).Set(vals...)
}

// Len of the elements
func (b *UIntBitset) Len() int {
	if b == nil {
		return 0
	}
	return len(b.values)
}

// Mask of the set
func (b *UIntBitset) Mask() uint64 {
	return b.mask
}

// Values list
func (b *UIntBitset) Values() []uint {
	if b == nil {
		return nil
	}
	return b.values
}

// Set type values
func (b *UIntBitset) Set(vals ...uint) *UIntBitset {
	var updated = false
	for _, v := range vals {
		if !b.Has(v) {
			b.mask |= 1 << uint64(v%64)
			b.values = append(b.values, v)
			updated = true
		}
	}
	if updated {
		sort.Slice(b.values, func(i, j int) bool { return b.values[i] < b.values[j] })
	}
	return b
}

// Unset type values
func (b *UIntBitset) Unset(vals ...uint) *UIntBitset {
	newVals := b.values
	for _, v := range vals {
		idx := sort.Search(len(newVals), func(i int) bool {
			return newVals[i] >= v
		})

		if idx >= 0 && idx < len(newVals) && newVals[idx] == v {
			i := idx + 1
			for ; i < len(newVals); i++ {
				if newVals[i] != v {
					break
				}
			}
			if idx > 0 {
				if i < len(newVals) {
					newVals = append(newVals[:idx], newVals[i:]...)
				} else if idx < len(newVals)-1 {
					newVals = newVals[:idx]
				}
			} else if i < len(newVals)-1 {
				newVals = newVals[i:]
			}
		}
	}

	if len(newVals) == len(b.values) {
		return b
	}
	return NewUIntBitset(newVals...)
}

// Has type in bitset
func (b *UIntBitset) Has(v uint) bool {
	if b != nil && b.mask&(1<<uint64(v%64)) != 0 {
		idx := sort.Search(b.Len(), func(i int) bool {
			return b.values[i] >= v
		})
		return idx >= 0 && idx < b.Len() && b.values[idx] == v
	}
	return false
}

// Reset bitset value
func (b *UIntBitset) Reset() *UIntBitset {
	b.mask = 0
	if b.values != nil {
		b.values = b.values[:0]
	}
	return b
}

// ContainsOnly items from the set
func (b *UIntBitset) ContainsOnly(set *UIntBitset) (res bool) {
	if set != nil && b.mask&set.mask == b.mask {
		res = true
		for _, v := range b.values {
			if res = set.Has(v); !res {
				break
			}
		}
	}
	return
}
