package loaders

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jinzhu/gorm"

	cmodels "geniusrabbit.dev/sspserver/internal/data/models"
	"geniusrabbit.dev/sspserver/internal/models"
)

// TargetReloader function
func TargetReloader(source interface{}) func() ([]models.Target, error) {
	switch src := source.(type) {
	case *gorm.DB:
		return func() ([]models.Target, error) {
			return dbTargetReloader(src)
		}
	case string:
		return func() ([]models.Target, error) {
			return fileTargetReloader(src)
		}
	}
	return nil
}

func dbTargetReloader(database *gorm.DB) (list []models.Target, err error) {
	var zones []*cmodels.Zone
	if err = database.Find(&zones).Error; err != nil {
		return nil, err
	}
	for _, zone := range zones {
		if zone == nil {
			continue
		}
		list = append(list, models.TargetFromModel(*zone))
	}
	return list, err
}

type targetData struct {
	Zones []*cmodels.Zone `json:"zones"`
}

func fileTargetReloader(filename string) (list []models.Target, _ error) {
	var (
		targets   targetData
		data, err = ioutil.ReadFile(filename)
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &targets); err != nil {
		return nil, err
	}
	for _, zone := range targets.Zones {
		if zone == nil {
			continue
		}
		list = append(list, models.TargetFromModel(*zone))
	}
	return list, err
}
