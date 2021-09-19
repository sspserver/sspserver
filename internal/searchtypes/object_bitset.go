//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"sort"
)

// ObjectBitset any numbers
type ObjectBitset struct {
	less   func(o1, o2 interface{}) bool
	value  func(o1 interface{}) uint
	values []interface{}
	mask   uint64
}

// NewObjectBitset from numbers
func NewObjectBitset(less func(o1, o2 interface{}) bool, value func(o1 interface{}) uint, vals ...interface{}) (b *ObjectBitset) {
	return (&ObjectBitset{less: less, value: value}).Set(vals...)
}

// Len of the elements
func (b *ObjectBitset) Len() int {
	if b == nil {
		return 0
	}
	return len(b.values)
}

// Values list
func (b *ObjectBitset) Values() []interface{} {
	if b == nil {
		return nil
	}
	return b.values
}

// Set type values
func (b *ObjectBitset) Set(vals ...interface{}) *ObjectBitset {
	var updated = false
	for _, v := range vals {
		if !b.Has(v) {
			b.mask |= 1 << uint64(b.value(v)%64)
			b.values = append(b.values, v)
			updated = true
		}
	}
	if updated {
		sort.Slice(b.values, func(i, j int) bool {
			return b.less(b.values[i], b.values[j])
		})
	}
	return b
}

// Unset type values
func (b *ObjectBitset) Unset(vals ...interface{}) *ObjectBitset {
	newVals := b.values
	for _, v := range vals {
		idx := sort.Search(len(newVals), func(i int) bool {
			return b.less(v, newVals[i])
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
	return NewObjectBitset(b.less, b.value, newVals...)
}

// Has type in bitset
func (b *ObjectBitset) Has(v interface{}) bool {
	if b != nil && b.mask&(1<<uint64(b.value(v)%64)) != 0 {
		idx := sort.Search(b.Len(), func(i int) bool {
			return b.less(v, b.values[i])
		})
		return idx >= 0 && idx < b.Len() && b.values[idx] == v
	}
	return false
}
