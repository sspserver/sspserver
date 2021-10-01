//
// @project GeniusRabbit AdNet 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package models

import (
	"strings"
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"
)

// Format types...
const (
	FormatTypeDirect  = `direct`
	FormatTypeVideo   = `video`
	FormatTypeProxy   = `proxy`
	FormatTypeNative  = `native`
	FormatTypeAutoreg = `autoreg`
)

// Format model description
type Format struct {
	ID       uint64 `json:"id"`
	Codename string `json:"codename"`
	Title    string `json:"title"`
	Active   uint   `json:"active,omitempty"`

	Width     int `json:"width,omitempty"`
	Height    int `json:"height,omitempty"`
	MinWidth  int `json:"min_width,omitempty"`
	MinHeight int `json:"min_height,omitempty"`

	Config gosql.NullableJSON `json:"config,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (f *Format) TableName() string {
	return "adv_format"
}

// DecodeConfig from JSON
func (f Format) DecodeConfig() (meta *FormatConfig) {
	if f.Config.Length() > 2 {
		meta = new(FormatConfig)
		if f.Config.UnmarshalTo(meta) != nil {
			meta = nil
		} else {
			var idx = int(1)
			for i, field := range meta.Fields {
				field.ID = idx
				meta.Fields[i] = field
				idx++
			}

			// Update asset IDs
			for i, asset := range meta.Assets {
				asset.ID = idx
				meta.Assets[i] = asset
				idx++
			}
		}
	}
	return
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

// IsDirect type (popunder or proxy or etc.)
func (f Format) IsDirect() bool {
	return f.Codename == FormatTypeDirect
}

// IsVideo type of advertisement
func (f Format) IsVideo() bool {
	return strings.HasPrefix(f.Codename, FormatTypeVideo)
}

// IsProxy type of format
func (f Format) IsProxy() bool {
	return strings.HasPrefix(f.Codename, FormatTypeProxy)
}

// IsNative type of format
func (f Format) IsNative() bool {
	return strings.HasPrefix(f.Codename, FormatTypeNative)
}

// IsAutoreg type of format
func (f Format) IsAutoreg() bool {
	return f.Codename == FormatTypeAutoreg
}
