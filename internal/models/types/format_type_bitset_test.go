//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package types

import (
	"reflect"
	"testing"
)

func TestFormatBitset(t *testing.T) {
	bits := NewFormatTypeBitset(FormatProxyType, FormatDirectType)

	bits.Set(FormatCustomType)
	if !bits.Has(FormatCustomType) {
		t.Error("bits Set/Has not working")
	}

	bits.Unset(FormatProxyType)
	if bits.Has(FormatProxyType) {
		t.Error("bits Unet/Has not working")
	}

	if !reflect.DeepEqual([]FormatType{FormatDirectType, FormatCustomType}, bits.Types()) {
		t.Errorf("bits Types problems: %v", bits.Types())
	}
}
