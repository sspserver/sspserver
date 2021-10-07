package types

import (
	"sort"

	"geniusrabbit.dev/sspserver/internal/searchtypes"
)

type formatLoaderFnk func() ([]*Format, error)

// SimpleFormatAccessor implements basic inmemory storage of formats
type SimpleFormatAccessor struct {
	formatList   []*Format
	formatLoader formatLoaderFnk
	directs      *searchtypes.UIntBitset
}

// NewSimpleFormatAccessor object
func NewSimpleFormatAccessor(formats []*Format) *SimpleFormatAccessor {
	return &SimpleFormatAccessor{
		formatList: formats,
	}
}

// NewSimpleFormatAccessorWithLoader object
func NewSimpleFormatAccessorWithLoader(formatLoader formatLoaderFnk) *SimpleFormatAccessor {
	accessor := &SimpleFormatAccessor{
		formatLoader: formatLoader,
	}
	_ = accessor.Reload()
	return accessor
}

// Reload sources
func (fa *SimpleFormatAccessor) Reload() error {
	if fa.formatLoader == nil {
		return nil
	}
	formats, err := fa.formatLoader()
	if formats != nil && err == nil {
		fa.formatList = formats
	}
	return err
}

// Prepare format list
func (fa *SimpleFormatAccessor) Prepare() {
	sort.Slice(fa.formatList, func(i, j int) bool { return fa.formatList[i].ID < fa.formatList[j].ID })
}

// Formats list collection
func (fa *SimpleFormatAccessor) Formats() []*Format {
	return fa.formatList
}

// FormatsBySize returns the list of acceptable formats
func (fa *SimpleFormatAccessor) FormatsBySize(w, h, minWidth, minHeight int, fmtTypes ...FormatTypeBitset) (list []*Format) {
	for _, frmt := range fa.formatList {
		if len(fmtTypes) > 0 && !fmtTypes[0].IsEmpty() && !fmtTypes[0].IsIntersec(frmt.Types) {
			continue
		} else if frmt.IsDirect() && (len(fmtTypes) < 1 || !fmtTypes[0].Has(FormatDirectType)) {
			// * NOTE: Need to exclude *direct* type format from default search
			continue
		} else if frmt.IsAutoregister() && (len(fmtTypes) < 1 || !fmtTypes[0].Has(FormatAutoregisterType)) {
			// * NOTE: Need to exclude *direct* type format from default search
			continue
		} else if frmt.IsVideo() && (len(fmtTypes) < 1 || !fmtTypes[0].Has(FormatVideoType)) {
			// * NOTE: Need to exclude *video* type format from default search
			continue
		}

		if frmt.IsDirect() || frmt.IsAutoregister() {
			list = append(list, frmt)
			continue
		}

		// Check the format with size restrictions
		if frmt.SuitsCompare(&Format{Width: w, Height: h, MinWidth: minWidth, MinHeight: minHeight}) == 0 {
			list = append(list, frmt)
		}
	}
	return
}

// FormatByID of the model
func (fa *SimpleFormatAccessor) FormatByID(id uint64) *Format {
	i := sort.Search(len(fa.formatList), func(i int) bool { return fa.formatList[i].ID >= id })
	if i >= 0 && i < len(fa.formatList) && fa.formatList[i].ID == id {
		return fa.formatList[i]
	}
	return nil
}

// FormatByCode of the model
func (fa *SimpleFormatAccessor) FormatByCode(code string) *Format {
	for _, f := range fa.formatList {
		if f.Codename == code {
			return f
		}
	}
	return nil
}

// DirectFormatSet to search
func (fa *SimpleFormatAccessor) DirectFormatSet() *searchtypes.UIntBitset {
	if fa.directs == nil {
		fa.directs = new(searchtypes.UIntBitset)
		for _, f := range fa.formatList {
			if f.IsDirect() {
				fa.directs.Set(uint(f.ID))
			}
		}
	}
	return fa.directs
}

var _ FormatsAccessor = (*SimpleFormatAccessor)(nil)
