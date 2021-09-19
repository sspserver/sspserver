//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package models

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"os"
// 	"testing"

// 	"bitbucket.org/geniusrabbit/corelib/models"
// 	"geniusrabbit.dev/sspserver/internal/models/types"
// )

// func Test_Formats(t *testing.T) {
// 	var (
// 		data    []byte
// 		err     error
// 		formats []*types.Format
// 	)

// 	// Init native format
// 	{
// 		if data, err = fileData("./assets/format.native.json"); err != nil {
// 			t.Error(err)
// 			return
// 		}

// 		format := models.Format{
// 			ID:        1,
// 			Codename:  "native",
// 			Title:     "Native",
// 			Active:    1,
// 			Width:     0,
// 			Height:    0,
// 			MinWidth:  30,
// 			MinHeight: 30,
// 		}
// 		format.Config.UnmarshalJSON(data)
// 		formats = append(formats, FormatFromModel(&format))

// 		if !formats[0].Types.Has(types.FormatNativeType) {
// 			t.Error("it must be native format type", formats[0].Types.Types())
// 		}
// 	}

// 	// Init video format
// 	{
// 		if data, err = fileData("./assets/format.video.json"); err != nil {
// 			t.Error(err)
// 			return
// 		}

// 		format := models.Format{
// 			ID:        2,
// 			Codename:  "video",
// 			Title:     "Video",
// 			Active:    1,
// 			Width:     0,
// 			Height:    0,
// 			MinWidth:  30,
// 			MinHeight: 30,
// 		}
// 		format.Config.UnmarshalJSON(data)
// 		formats = append(formats, FormatFromModel(&format))

// 		if !formats[1].Types.Has(types.FormatVideoType) {
// 			t.Error("it must be video format type", formats[1].Types.Types())
// 		}
// 	}
// }

// func configByJSON(path string) (*models.FormatConfig, error) {
// 	var (
// 		conf   *models.FormatConfig
// 		f, err = os.Open(path)
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	err = json.NewDecoder(f).Decode(&conf)
// 	return conf, err
// }

// func fileData(path string) ([]byte, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	return ioutil.ReadAll(f)
// }
