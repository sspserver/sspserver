package info

// FieldOptionType of value
type FieldOptionType string

// InfoFieldOption types
const (
	FieldOptionString          FieldOptionType = "string"
	FieldOptionInt             FieldOptionType = "int"
	FieldOptionFloat           FieldOptionType = "float"
	FieldOptionBool            FieldOptionType = "bool"
	FieldOptionStringStringMap FieldOptionType = "map[string]string"
)

// Documentation object info
type Documentation struct {
	Title string `json:"title,omitempty"`
	Link  string `json:"link"`
}

// FieldOptionSelectItem describes variant of option
type FieldOptionSelectItem struct {
	Name        string      `json:"name,omitempty"`
	Value       interface{} `json:"value"`
	Description string      `json:"description,omitempty"`
	Link        string      `json:"link,omitempty"`
}

// FieldOption description
type FieldOption struct {
	Name        string                  `json:"name"`
	Required    bool                    `json:"required,omitempty"`
	Type        FieldOptionType         `json:"type,omitempty"`
	Description string                  `json:"description,omitempty"`
	Default     interface{}             `json:"default,omitempty"`
	Select      []FieldOptionSelectItem `json:"select,omitempty"` // Available variants
}

// Subprotocol needed for describing of subprotocols support.
// For example OpenRTB supports many subformats like VAST, OpenNative, etc.
type Subprotocol struct {
	Name        string          `json:"name,omitempty"`
	Protocol    string          `json:"protocol"`
	Versions    []string        `json:"versions,omitempty"`
	Description string          `json:"description,omitempty"`
	Docs        []Documentation `json:"docs,omitempty"`
}

// Platform info driver description
type Platform struct {
	Name         string          `json:"name,omitempty"`
	Protocol     string          `json:"protocol"`
	Versions     []string        `json:"versions,omitempty"`
	Description  string          `json:"description,omitempty"`
	Docs         []Documentation `json:"docs,omitempty"`
	Options      []*FieldOption  `json:"options,omitempty"`
	Subprotocols []Subprotocol   `json:"subprotocols,omitempty"`
}

// Short information
func (p Platform) Short() Platform {
	short := p
	short.Docs = nil
	short.Options = nil
	short.Subprotocols = nil
	return short
}
