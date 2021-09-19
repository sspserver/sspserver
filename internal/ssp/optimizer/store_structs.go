package optimizer

//easyjson:json
type keyValType struct {
	CountryID  byte    `json:"cc,omitempty"`
	LanguageID uint    `json:"ln,omitempty"`
	DeviceID   uint    `json:"dv,omitempty"`
	OsID       uint    `json:"os,omitempty"`
	BrowserID  uint    `json:"br,omitempty"`
	Value      float64 `json:"vl,omitempty"`
}
