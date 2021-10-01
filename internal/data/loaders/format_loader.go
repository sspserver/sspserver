package loaders

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jinzhu/gorm"

	cmodels "geniusrabbit.dev/sspserver/internal/data/models"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// FormatLoader from source like database of filesystem
func FormatLoader(source interface{}) func() ([]*types.Format, error) {
	switch src := source.(type) {
	case *gorm.DB:
		return func() ([]*types.Format, error) {
			return dbFormatLoader(src)
		}
	case string:
		return func() ([]*types.Format, error) {
			return fileFormatLoader(src)
		}
	}
	return nil
}

func dbFormatLoader(database *gorm.DB) (list []*types.Format, err error) {
	var formats []*cmodels.Format
	if err = database.Find(&formats).Error; err != nil {
		return nil, err
	}
	for _, format := range formats {
		if format == nil {
			continue
		}
		list = append(list, models.FormatFromModel(format))
	}
	return list, err
}

type formatData struct {
	Formats []*cmodels.Format `json:"formats"`
}

func fileFormatLoader(filename string) (list []*types.Format, _ error) {
	var (
		formats   formatData
		data, err = ioutil.ReadFile(filename)
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &formats); err != nil {
		return nil, err
	}
	for _, format := range formats.Formats {
		if format == nil {
			continue
		}
		list = append(list, models.FormatFromModel(format))
	}
	return list, err
}
