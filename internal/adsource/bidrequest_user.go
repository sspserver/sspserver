//
// @project geniusrabbit::rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package adsource

import (
	"strings"
	"time"

	"github.com/bsm/openrtb"
	"github.com/sspserver/udetect"
	uopenrtb "github.com/sspserver/udetect/openrtb"
)

// TypeSex type
type TypeSex uint8

// User sex enum
const (
	UserSexUndefined TypeSex = 0
	UserSexMale      TypeSex = 1
	UserSexFemale    TypeSex = 2
)

func (s TypeSex) String() string {
	switch s {
	case UserSexMale:
		return "Male"
	case UserSexFemale:
		return "Female"
	}
	return "Undefined"
}

// Code from the sex
func (s TypeSex) Code() string {
	switch s {
	case UserSexMale:
		return "M"
	case UserSexFemale:
		return "F"
	}
	return ""
}

// SexByString returns the type Sex
func SexByString(sex string) (t TypeSex) {
	switch strings.ToLower(sex) {
	case "male", "m":
		t = UserSexMale
	case "female", "f":
		t = UserSexFemale
	default:
		t = UserSexUndefined
	}
	return
}

// Segment item value
type Segment struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Data item segment
type Data struct {
	Name    string    `json:"name,omitempty"`
	Segment []Segment `json:"segment,omitempty"`
}

// User information
type User struct {
	ID            string       `json:"id,omitempty"`        // Unique User ID
	Email         string       `json:"email,omitempty"`     // In some cases it's able the use email, and we are gonna use it
	Username      string       `json:"username,omitempty"`  // User profile name from the external service or potentional username
	SessionID     string       `json:"sessid,omitempty"`    // Unique session ID
	FingerPrintID string       `json:"fpid,omitempty"`      //
	ETag          string       `json:"etag,omitempty"`      //
	Birthday      string       `json:"birthday,omitempty"`  // * Prefer do not use such personal information in alghoritm
	AgeStart      int          `json:"age_start,omitempty"` // Year of birth from
	AgeEnd        int          `json:"age_end,omitempty"`   // Year of birth from
	Gender        string       `json:"gender,omitempty"`    // Gender ("M": male, "F" female, "O" Other)
	Keywords      string       `json:"keywords,omitempty"`  // Comma separated list of keywords, interests, or intent
	Geo           *udetect.Geo `json:"geo,omitempty"`
	Data          []Data       `json:"data,omitempty"`
	birthday      time.Time
}

// AvgAge number
func (u User) AvgAge() byte {
	return byte(u.AgeEnd - u.AgeStart)
}

// SetSexByString text value
func (u *User) SetSexByString(sex string) {
	u.Gender = SexByString(sex).Code()
}

// Sex constant
func (u User) Sex() TypeSex {
	switch u.Gender {
	case "M":
		return UserSexMale
	case "F":
		return UserSexFemale
	}
	return UserSexUndefined
}

// RTBObject of User
func (u User) RTBObject() *openrtb.User {
	var data []openrtb.Data
	for _, it := range u.Data {
		dataItem := openrtb.Data{Name: it.Name}
		for i := 0; i < len(it.Segment); i++ {
			dataItem.Segment = append(dataItem.Segment, openrtb.Segment{
				Name:  it.Segment[i].Name,
				Value: it.Segment[i].Value,
			})
		}
		data = append(data, dataItem)
	}

	return &openrtb.User{
		ID:         u.ID,       // Unique consumer ID of this user on the exchange
		BuyerID:    "",         // Buyer-specific ID for the user as mapped by the exchange for the buyer. At least one of buyeruid/buyerid or id is recommended. Valid for OpenRTB 2.3.
		BuyerUID:   "",         // Buyer-specific ID for the user as mapped by the exchange for the buyer. Same as BuyerID but valid for OpenRTB 2.2.
		YOB:        0,          // Year of birth as a 4-digit integer.
		Gender:     u.Gender,   // Gender ("M": male, "F" female, "O" Other)
		Keywords:   u.Keywords, // Comma separated list of keywords, interests, or intent
		CustomData: "",         // Optional feature to pass bidder data that was set in the exchange's cookie. The string must be in base85 cookie safe characters and be in any format. Proper JSON encoding must be used to include "escaped" quotation marks.
		Geo:        uopenrtb.GeoFrom(u.Geo),
		Data:       data,
		Ext:        nil,
	}
}

// SetDataItem with simple *key*, *value*
func (u *User) SetDataItem(name, value string) {
	for i, data := range u.Data {
		if data.Name == name {
			if data.Segment != nil {
				data.Segment = data.Segment[:0]
			}
			data.Segment = append(data.Segment, Segment{Value: value})
			u.Data[i] = data
			return
		}
	}

	u.Data = append(u.Data, Data{Name: name, Segment: []Segment{{Value: value}}})
}

// GetDataItem simple value by key
func (u *User) GetDataItem(name string) (v string, ok bool) {
	for _, it := range u.Data {
		if it.Name == name {
			if len(it.Segment) < 1 {
				return "", true
			}
			return it.Segment[0].Value, true
		}
	}
	return "", false
}

// GetDataItemOrDefault item
func (u *User) GetDataItemOrDefault(name, def string) string {
	if v, ok := u.GetDataItem(name); ok {
		return v
	}
	return def
}

// BirthdayTime parsed by Birthday string
func (u *User) BirthdayTime() time.Time {
	if u.birthday.IsZero() && u.Birthday != "" {
		u.birthday = parseBirthday(u.Birthday)
	}
	return u.birthday
}

func parseBirthday(date string) (t time.Time) {
	var err error
	switch {
	case strings.ContainsRune(date, '/'):
		if t, err = time.Parse("2006/01/02", date); err != nil {
			t, _ = time.Parse("02/01/2006", date)
		}
	case strings.ContainsRune(date, '-'):
		if t, err = time.Parse("2006-01-02", date); err != nil {
			t, _ = time.Parse("02-01-2006", date)
		}
	case strings.ContainsRune(date, '.'):
		if t, err = time.Parse("2006.01.02", date); err != nil {
			t, _ = time.Parse("02.01.2006", date)
		}
	}
	return
}
