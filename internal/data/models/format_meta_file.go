//
// @project GeniusRabbit AdNet 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package models

import "strings"

// Format asset defaults
const (
	FormatAssetMain = "main"
	FormatAssetIcon = "icon"
)

// FormatFileRequirement rule
type FormatFileRequirement struct {
	// ID of the assert
	ID int `json:"id,omitempty"`

	// Is required file
	Required bool `json:"required,omitempty"`

	// Name of the asset asset (optional for the first asset) by default: main
	Name string `json:"name,omitempty"`

	// AdjustSize if size is bigger
	AdjustSize bool `json:"adjust_size,omitempty"`

	// Max or exact width of image or video in pixels
	Width int `json:"width,omitempty"`

	// Max or exact height of image or video in pixels
	Height int `json:"height,omitempty"`

	// Minimal allowed width of image or video in pixels
	MinWidth int `json:"min_width,omitempty"`

	// Minimal allowed height of image or video in pixels
	MinHeight int `json:"min_height,omitempty"`

	// Is animated type of file acceptable
	Animated bool `json:"animated,omitempty"`

	// Is sound acceptable
	Sound bool `json:"sound,omitempty"`

	// Thumbs settings for autogenerator
	Thumbs []string `json:"thumbs,omitempty"`

	// AllowedTypes of files (for paticular asset)
	AllowedTypes []string `json:"allowed_types,omitempty"`
}

// MinMaxWidth sizes
func (f *FormatFileRequirement) MinMaxWidth() (mn, mx int) {
	if f.MinWidth > 0 {
		mn = f.MinWidth
		mx = f.Width
	} else {
		mn = f.Width
	}
	return
}

// MinMaxHeight sizes
func (f *FormatFileRequirement) MinMaxHeight() (mn, mx int) {
	if f.MinHeight > 0 {
		mn = f.MinHeight
		mx = f.Height
	} else {
		mn = f.Height
	}
	return
}

// IsRequired asset
func (f *FormatFileRequirement) IsRequired() bool {
	if !f.Required {
		f.Required = f.IsMain()
	}
	return f.Required
}

// IsMain asset
func (f FormatFileRequirement) IsMain() bool {
	return f.GetName() == FormatAssetMain
}

// GetName of the asset
func (f FormatFileRequirement) GetName() string {
	if f.Name == "" {
		return FormatAssetMain
	}
	return f.Name
}

// IsVideoSupport of the file
func (f FormatFileRequirement) IsVideoSupport() bool {
	return f.IsMimeTypeSupport("video")
}

// IsImageSupport of the file
func (f FormatFileRequirement) IsImageSupport() bool {
	return f.IsMimeTypeSupport("image")
}

// IsSoundSupport of the file
func (f FormatFileRequirement) IsSoundSupport() bool {
	return f.IsMimeTypeSupport("sound")
}

// IsHTML5Support of the file
func (f FormatFileRequirement) IsHTML5Support() bool {
	return f.IsMimeTypeSupport("html5")
}

// IsFixed format size
func (f FormatFileRequirement) IsFixed() bool {
	return f.Width > 0 && f.MinWidth == 0 && f.Height > 0 && f.MinHeight == 0
}

// IsMimeTypeSupport of the file
func (f FormatFileRequirement) IsMimeTypeSupport(mimePrefix string) bool {
	if len(f.AllowedTypes) == 0 {
		return true
	}
	for _, tp := range f.AllowedTypes {
		if tp == mimePrefix || strings.HasPrefix(tp, mimePrefix+"/") {
			return true
		}
	}
	return false
}

// SoftEqual compare
func (f FormatFileRequirement) SoftEqual(asset *FormatFileRequirement) (r bool) {
	if f.Name != asset.Name && (!f.IsMain() || !asset.IsMain()) {
		return
	}
	if r = !(len(f.AllowedTypes) > 0 && len(asset.AllowedTypes) > 0); !r {
		for _, mime := range asset.AllowedTypes {
			if f.IsMimeTypeSupport(mime) {
				r = true
				break
			}
		}
	}
	if r {
		switch {
		case f.MinWidth <= 0 && asset.MinWidth <= 0 && f.MinHeight <= 0 && asset.MinHeight <= 0:
			r = f.Width == asset.Width && f.Height == asset.Height
		case (f.MinWidth > 0 && asset.MinWidth <= 0) || (f.MinHeight > 0 && asset.MinHeight <= 0):
			r = false
		default:
			r = compareAspect(f.Width, f.MinWidth, asset.Width, asset.MinWidth) == 0 &&
				compareAspect(f.Height, f.MinHeight, asset.Height, asset.MinHeight) == 0
		}
	}
	return r
}
