//
// @project GeniusRabbit AdNet 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package models

// FormatConfig description
type FormatConfig struct {
	// Assets list of the files and configs
	Assets []FormatFileRequirement `json:"assets,omitempty"`

	// By default empty, so in requirements only the link from Ad object
	Fields []FormatField `json:"fields,omitempty"`
}

// Intersec with other format config
func (c *FormatConfig) Intersec(conf *FormatConfig) bool {
	if c.IsEmpty() && conf.IsEmpty() {
		return true
	}

	// Check assets which must be required
	if len(c.Assets) > 0 && len(conf.Assets) > 0 {
		for _, asset := range conf.Assets {
			if asset.IsRequired() && !c.ContainsAsset(&asset) {
				return false
			}
		}

		// Check if required assets does not demands
		for _, asset := range c.Assets {
			if asset.IsRequired() && !conf.ContainsAsset(&asset, true) {
				return false
			}
		}
	} else if len(c.Assets) > 0 {
		for _, asset := range c.Assets {
			if asset.IsRequired() {
				return false
			}
		}
	}

	// Check fields which must be required
	if len(c.Fields) > 0 && len(conf.Fields) > 0 {
		for _, field := range conf.Fields {
			if field.Required && c.SimilarField(&field) == nil {
				return false
			}
		}

		// Check if required field does not demands
		for _, field := range c.Fields {
			if field.Required && conf.SimilarField(&field, true) == nil {
				return false
			}
		}
	} else if len(c.Fields) > 0 {
		for _, field := range c.Fields {
			if field.Required {
				return false
			}
		}
	}
	return true
}

// AssetByName from config
func (c *FormatConfig) AssetByName(name string) *FormatFileRequirement {
	for _, asset := range c.Assets {
		if asset.Name == name || (name == "" && asset.IsMain()) {
			return &asset
		}
	}
	return nil
}

// ContainsAsset in the list
func (c *FormatConfig) ContainsAsset(asset *FormatFileRequirement, revers ...bool) bool {
	if c == nil {
		return false
	}

	revCheck := len(revers) > 0 && revers[0]
	for _, a := range c.Assets {
		if revCheck {
			if asset.SoftEqual(&a) {
				return true
			}
		} else if a.SoftEqual(asset) {
			return true
		}
	}
	return false
}

// SimilarField in the list
func (c *FormatConfig) SimilarField(field *FormatField, revers ...bool) *FormatField {
	if c == nil {
		return nil
	}

	revCheck := len(revers) > 0 && revers[0]
	for _, f := range c.Fields {
		if revCheck {
			if field.SoftEqual(&f) {
				return &f
			}
		} else if f.SoftEqual(field) {
			return &f
		}
	}
	return nil
}

// SimpleAsset returns the main asset in case of one required asset
func (c *FormatConfig) SimpleAsset() (as *FormatFileRequirement) {
	if c == nil {
		return nil
	}

	for i, asset := range c.Assets {
		if asset.IsMain() {
			as = &c.Assets[i]
		} else if asset.Required {
			as = nil
			break
		}
	}
	return
}

// MainAsset returns the main asset if exists
func (c *FormatConfig) MainAsset() (as *FormatFileRequirement) {
	if c == nil {
		return nil
	}

	for i, asset := range c.Assets {
		if asset.IsMain() {
			return &c.Assets[i]
		}
	}
	return
}

// GetField by name
func (c *FormatConfig) GetField(name string) *FormatField {
	if c == nil {
		return nil
	}

	for i, fl := range c.Fields {
		if fl.Name == name {
			return &c.Fields[i]
		}
	}
	return nil
}

// RequiredField have in the config
func (c *FormatConfig) RequiredField(fields ...string) *FormatField {
	if c == nil {
		return nil
	}

	if len(fields) > 0 {
		var haveRequired = false
		for _, fl := range fields {
			for i, field := range c.Fields {
				if !haveRequired {
					haveRequired = field.Required
				}
				if field.Required && fl == field.Name {
					return &c.Fields[i]
				}
			}
			if !haveRequired {
				return nil
			}
		}
	} else {
		for i, field := range c.Fields {
			if field.Required {
				return &c.Fields[i]
			}
		}
	}
	return nil
}

// RequiredFieldExcept have any required field
// Fileds in param must be optional
func (c *FormatConfig) RequiredFieldExcept(fields ...string) *FormatField {
	if c == nil {
		return nil
	}

	found := false
	for i, field := range c.Fields {
		if field.Required {
			found = false
			for _, fl := range fields {
				if fl == field.Name {
					found = true
					break
				}
			}
			if !found {
				return &c.Fields[i]
			}
		}
	}
	return nil
}

// IsEmpty config
func (c *FormatConfig) IsEmpty() bool {
	return c == nil || (len(c.Assets) == 0 && len(c.Fields) == 0)
}

func compareAspect(v, mv, targetV, targetMV int) int {
	if v == 0 {
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

	if (targetMV > 0 && v < targetMV) || (targetMV == 0 && v < targetV) {
		return -1
	}

	return 0
}
