//
// @project GeniusRabbit rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package types

import (
	"fmt"
	"math"

	"geniusrabbit.dev/sspserver/internal/data/models"
)

// Format struct
type Format struct {
	ID        uint64               `json:"id"`
	Codename  string               `json:"codename"`
	Types     FormatTypeBitset     `json:"type"`
	Width     int                  `json:"w"`
	Height    int                  `json:"h"`
	MinWidth  int                  `json:"mw,omitempty"`
	MinHeight int                  `json:"mh,omitempty"`
	Config    *models.FormatConfig `json:"config,omitempty"`
	cloned    bool
}

// Code hash
func (f Format) Code() uint64 {
	return uint64(f.Width<<32) | uint64(uint32(f.Height))
}

// CloneWithSize of current format object with the new params
func (f *Format) CloneWithSize(width, height, minWidth, minHeight int) *Format {
	if f == nil {
		return nil
	}
	newFormat := *f
	if width > 0 {
		newFormat.Width = width
	}
	if height > 0 {
		newFormat.Height = height
	}
	if minWidth > 0 {
		newFormat.MinWidth = minWidth
	}
	if minHeight > 0 {
		newFormat.MinHeight = minHeight
	}
	newFormat.cloned = true
	return &newFormat
}

// String from format
func (f Format) String() string {
	switch {
	case f.Width > 0 && f.Height > 0:
		return fmt.Sprintf("%dx%d", f.Width, f.Height)
	case f.Width > 0:
		return fmt.Sprintf("%dx0", f.Width)
	default:
		if one := f.Types.HasOneType(); !one.IsInvalid() {
			return one.Name()
		}
	}
	return "stretch"
}

// GetConfig of the format object
func (f *Format) GetConfig() *models.FormatConfig {
	if f == nil {
		return nil
	}
	return f.Config
}

// Compare formats
func (f Format) Compare(format *Format) (r int) {
	if f.ID > 0 && f.ID == format.ID {
		return 0
	}
	if r = compareWH(f.Width, f.Height, format.Width, format.Height); r == 0 {
		r = compareWH(f.MinWidth, f.MinHeight, format.MinWidth, format.MinHeight)
	}
	if r == 0 {
		if f.ID > format.ID {
			r = 1
		} else if f.ID < format.ID {
			r = -1
		}
	}
	return r
}

// SuitsCompare formats for premetive test
func (f Format) SuitsCompare(format *Format) int {
	return f.SuitsCompareSize(format.Width, format.Height, format.MinWidth, format.MinHeight)
}

// SuitsCompareSize for premetive test
func (f Format) SuitsCompareSize(w, h, mw, mh int) int {
	if true &&
		compareAspect(f.Width, f.MinWidth, w, mw) == 0 &&
		compareAspect(f.Height, f.MinHeight, h, mh) == 0 {
		return 0
	}

	if f.Width < w || (f.Width == w && f.Height < h) {
		return -1
	}
	return 1
}

// IsZero format
func (f *Format) IsZero() bool {
	return f != nil &&
		f.Width == 0 && f.MinWidth == 0 &&
		f.Height == 0 && f.MinHeight == 0
}

// IsStretch format
func (f *Format) IsStretch() bool {
	return f != nil && !f.IsFixed()
}

// IsFixed format
func (f *Format) IsFixed() bool {
	return f != nil &&
		f.Width > 0 && (f.MinWidth <= 0 || f.Width == f.MinWidth) &&
		f.Height > 0 && (f.MinHeight <= 0 || f.Height == f.MinHeight)
}

// Diagonal of square
func (f Format) Diagonal() float64 {
	return math.Sqrt(pow2(prepSize(f.Width)) + pow2(prepSize(f.Height)))
}

// IsBanner format type
func (f Format) IsBanner() bool {
	return f.Types.Is(FormatBannerType)
}

// IsDirect format type
func (f Format) IsDirect() bool {
	return f.Types.Is(FormatDirectType)
}

// IsAutoregister format type
func (f Format) IsAutoregister() bool {
	return f.Types.Is(FormatAutoregisterType)
}

// IsProxy format type
func (f Format) IsProxy() bool {
	return f.Types.Has(FormatProxyType)
}

// IsVideo format type
func (f Format) IsVideo() bool {
	return f.Types.Has(FormatVideoType)
}

// IsNative format type
func (f Format) IsNative() bool {
	return f.Types.Has(FormatNativeType)
}

// IsCloned format object
func (f *Format) IsCloned() bool {
	return f != nil && f.cloned
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func pow2(v float64) float64 {
	return v * v
}

func prepSize(v int) float64 {
	if v < 1 {
		return float64(math.MaxInt32)
	}
	return float64(v)
}

func compareAspect(v, mv, targetV, targetMV int) int {
	if v == 0 && mv <= 0 {
		return 0
	}

	if mv == -1 {
		mv = v / 2
	}

	if targetMV == -1 {
		targetMV = targetV / 2
	}

	if mv > targetV || (mv == 0 && v > targetV) {
		return 1
	}

	if v != 0 && ((targetMV > 0 && v < targetMV) || (targetMV == 0 && v < targetV)) {
		return -1
	}

	return 0
}

func compareWH(w1, h1, w2, h2 int) int {
	if w1 == w2 && h1 == h2 {
		return 0
	}
	if (w1 < w2) || (w1 == w2 && h1 < h2) {
		return -1
	}
	return 1
}
