package getr

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"strconv"
	"strings"

	v "gopkg.in/go-playground/validator.v9"

	"github.com/demdxx/gocast"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

var validator = v.New()

//easyjson:json
type registerRequest struct {
	TrackCode     string `json:"trackcode" form:"trackcode" query:"trackcode"`
	SubSourse     string `json:"sub_sourse" form:"sub_sourse" query:"sub_sourse"`
	Email         string `json:"email" form:"email" query:"email" validate:"required"`
	Username      string `json:"username" form:"username" query:"username"`
	Firstname     string `json:"firstname" form:"firstname" query:"firstname"`
	Lastname      string `json:"lastname" form:"lastname" query:"lastname"`
	Age           string `json:"age" form:"age" query:"age"`
	DateOfBirth   string `json:"birthday" form:"birthday" query:"birthday"`
	Gender        string `json:"gender" form:"gender" query:"gender"`
	SearchGender  string `json:"search_gender" form:"search_gender" query:"search_gender"`
	PhoneNumber   string `json:"phone" form:"phone" query:"phone"`
	MessangerType string `json:"messanger_type" form:"messanger_type" query:"messanger_type"`
	Messanger     string `json:"messanger" form:"messanger" query:"messanger"`
	IP            string `json:"ip" form:"ip" query:"ip"`
	Country       string `json:"country" form:"country" query:"country"`
	City          string `json:"city" form:"city" query:"city"`
	PostCode      string `json:"zip" form:"zip" query:"zip"`
	Password      string `json:"pw" form:"pw" query:"pw"`
	UserAgent     string `json:"ua" form:"ua" query:"ua"`
}

func (req *registerRequest) QueryMapDecode(m map[string]interface{}) (err error) {
	if err = gocast.ToStruct(req, m, "query"); err == nil {
		err = req.prepare()
	}
	return
}

func (req *registerRequest) FormMapDecode(m map[string]interface{}) (err error) {
	if err = gocast.ToStruct(req, m, "form"); err == nil {
		err = req.prepare()
	}
	return
}

func (req *registerRequest) JSONDecode(reader io.Reader) (err error) {
	if err = json.NewDecoder(reader).Decode(req); err == nil {
		err = req.prepare()
	}
	return
}

// FillBidRequest updates the original *BidRequest*
// TODO process all errors and reac on them
func (req *registerRequest) FillBidRequest(request *adsource.BidRequest) {
	var (
		user    = request.UserInfo()
		geo     = request.GeoInfo()
		browser = request.BrowserInfo()
	)

	request.ImpressionUpdate(func(imp *adsource.Impression) bool {
		if imp.FormatTypes.Has(types.FormatAutoregisterType) {
			imp.ExtID = req.TrackCode
			imp.SubID1 = req.SubSourse
			if req.SearchGender != "" {
				imp.Set("search_gender", adsource.SexByString(req.SearchGender).Code())
			}
			if req.Password != "" {
				imp.Set("new_password", req.Password)
			}
			return true
		}
		return false
	})

	var age int64
	if req.Age != "" && req.Age != "0" {
		// TODO log error
		age, _ = strconv.ParseInt(req.Age, 10, 64)
	}

	user.Email = req.Email
	user.Username = req.Username
	user.AgeStart = int(age)
	user.Birthday = req.DateOfBirth
	user.SetSexByString(req.Gender)

	if req.Firstname != "" {
		user.SetDataItem("firstname", req.Firstname)
		user.SetDataItem("lastname", req.Lastname)
	}

	if req.PhoneNumber != "" {
		user.SetDataItem("phone", req.PhoneNumber)
	}

	if req.Messanger != "" {
		user.SetDataItem("messanger_type", req.MessangerType)
		user.SetDataItem("messanger", req.Messanger)
	}

	if req.IP != "" {
		geo.IP = net.ParseIP(req.IP)
	}

	if req.Country != "" {
		geo.Country = req.Country
		geo.City = req.City
	}
	if req.PostCode != "" {
		geo.Zip = req.PostCode
	}

	if req.UserAgent != "" {
		browser.UA = req.UserAgent
	}
}

func (req *registerRequest) prepare() (err error) {
	req.Password, err = preparePassword(req.Password)
	return
}

func (req *registerRequest) Validate() error {
	if req.Country != "" {
		// ...
	}
	return validator.Struct(req)
}

func preparePassword(password string) (string, error) {
	if strings.HasPrefix(password, "base64:") {
		data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(password, "base64:"))
		if err != nil {
			return "", err
		}
		password = string(data)
	}
	return password, nil
}
