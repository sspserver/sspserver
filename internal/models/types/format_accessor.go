//
// @project GeniusRabbit rotator 2016 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2018
//

package types

import "geniusrabbit.dev/sspserver/internal/searchtypes"

// FormatsAccessor object interface
type FormatsAccessor interface {
	// Format list collection
	Formats() []*Format

	// FormatsBySize returns the list of acceptable formats
	FormatsBySize(w, h, minWidth, minHeight int, ftypes ...FormatTypeBitset) []*Format

	// FormatByID of the model
	FormatByID(id uint64) *Format

	// FormatByCode of the model
	FormatByCode(code string) *Format

	// DirectFormatSet to search
	DirectFormatSet() *searchtypes.UIntBitset
}
