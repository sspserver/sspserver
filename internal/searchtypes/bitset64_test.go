//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"reflect"
	"testing"
)

func TestBitset64(t *testing.T) {
	bits := NewBitset64(1, 3)

	bits.Set(2)
	if !bits.Has(2) {
		t.Error("bits Set/Has not working")
	}

	bits.Unset(3)
	if bits.Has(3) {
		t.Error("bits Unet/Has not working")
	}

	if !reflect.DeepEqual([]uint{1, 2}, bits.Numbers()) {
		t.Errorf("bits Numbers problems: %v", bits.Numbers())
	}
}
