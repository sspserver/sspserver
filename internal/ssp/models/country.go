package models

// CountryCode type
type CountryCode uint16

const (
	CountryCodeUndefined CountryCode = 0
)

// CountryCodeFromString convert ISO-CC into the number
func CountryCodeFromString(cc string) CountryCode {
	if cc == "**" || cc == "*" {
		return CountryCodeUndefined
	}
	switch len(cc) {
	case 2:
		return (CountryCode)(cc[0]) | (CountryCode)(cc[1])<<8
	}
	return CountryCodeUndefined
}

func (cc CountryCode) String() string {
	if cc == CountryCodeUndefined {
		return "**"
	}
	return string([]byte{byte(cc), byte(cc >> 8)})
}
