package types

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"geniusrabbit.dev/sspserver/internal/data/models"
)

var mockFormats = []*Format{
	{
		ID:       1,
		Codename: "direct",
		Types:    *NewFormatTypeBitset(FormatDirectType),
	},
	{
		ID:        2,
		Codename:  "proxy",
		Types:     *NewFormatTypeBitset(FormatProxyType),
		Width:     0,
		Height:    0,
		MinWidth:  10,
		MinHeight: 10,
	},
	{
		ID:       3,
		Codename: "proxy_200x200",
		Types:    *NewFormatTypeBitset(FormatProxyType),
		Width:    200,
		Height:   200,
	},
	{
		ID:       4,
		Codename: "banner_200x200",
		Types:    *NewFormatTypeBitset(FormatBannerType),
		Width:    200,
		Height:   200,
	},
	{
		ID:       5,
		Codename: "banner_300x300",
		Types:    *NewFormatTypeBitset(FormatBannerType),
		Width:    300,
		Height:   300,
	},
}

// MockFormats for testing
func MockFormats() []*Format {
	if len(mockFormats) != 7 {
		// Video format
		videoConfig, err := configByJSON("../assets/format.video.json")
		if err != nil {
			panic(err)
		}
		videoFormat := Format{
			ID:       6,
			Codename: "video",
			Types:    *NewFormatTypeBitset(FormatVideoType),
			Config:   videoConfig,
		}
		mockFormats = append(mockFormats, &videoFormat)

		// Native format
		nativeConfig, err := configByJSON("../assets/format.native.json")
		if err != nil {
			panic(err)
		}
		nativeFormat := Format{
			ID:       7,
			Codename: "native",
			Types:    *NewFormatTypeBitset(FormatNativeType),
			Config:   nativeConfig,
		}
		mockFormats = append(mockFormats, &nativeFormat)
	}
	return mockFormats
}

func configByJSON(path string) (*models.FormatConfig, error) {
	var (
		_, fileName, _, _ = runtime.Caller(1)
		conf              *models.FormatConfig
		f, err            = os.Open(filepath.Dir(fileName) + "/" + path)
	)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&conf)
	return conf, err
}

func fileData(path string) ([]byte, error) {
	_, fileName, _, _ := runtime.Caller(1)
	f, err := os.Open(filepath.Dir(fileName) + "/" + path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
