package loaders

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	// "github.com/hashicorp/hcl"
	// "gopkg.in/yaml.v2"

	imodels "bitbucket.org/geniusrabbit/corelib/models"
	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/eventstream"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/ssp/platform"
)

var errInvalidFileFormat = errors.New(`[SourceReloader] invalid file format`)

type companyGetter func(id uint64) *models.Company

// SourceReloader accessor
func SourceReloader(logger *zap.Logger, data interface{}, companyGetter companyGetter, eventStream eventstream.Stream, metrics, winNotify nc.Publisher) func() ([]adsource.Source, error) {
	logger = logger.With(zap.String("module", "SourceReloader"))
	switch d := data.(type) {
	case *gorm.DB:
		return DBSourceReloader(logger, d, companyGetter, eventStream, metrics, winNotify)
	case string:
		return FSSourceReloader(logger, d, companyGetter, eventStream, metrics, winNotify)
	}
	return nil
}

// DBSourceReloader accessor
func DBSourceReloader(logger *zap.Logger, database *gorm.DB, companyGetter companyGetter, eventStream eventstream.Stream, metrics, winNotify nc.Publisher) func() ([]adsource.Source, error) {
	return func() ([]adsource.Source, error) {
		var sourceList []*imodels.RTBSource
		if err := database.Find(&sourceList).Error; err != nil {
			return nil, err
		}
		return reload(sourceList, logger, companyGetter, eventStream, metrics, winNotify)
	}
}

// FSSourceReloader accessor
func FSSourceReloader(logger *zap.Logger, filename string, companyGetter companyGetter, eventStream eventstream.Stream, metrics, winNotify nc.Publisher) func() ([]adsource.Source, error) {
	return func() (sources []adsource.Source, err error) {
		sourceList, err := readSources(filename)
		if err != nil {
			return
		}
		return reload(sourceList, logger, companyGetter, eventStream, metrics, winNotify)
	}
}

func reload(sourceList []*imodels.RTBSource, logger *zap.Logger, companyGetter companyGetter, eventStream eventstream.Stream, metrics, winNotify nc.Publisher) (sources []adsource.Source, err error) {
	for _, baseSource := range sourceList {
		var (
			err    error
			src    = models.RTBSourceFromModel(baseSource, companyGetter(baseSource.CompanyID))
			source adsource.Source
		)

		if src == nil {
			logger.Error("invalid RTB source object",
				zap.Uint64("source_id", baseSource.ID),
				zap.String("protocol", baseSource.Protocol))
			continue
		}

		if fact := platform.ByProtocol(src.Protocol); fact != nil {
			source, err = fact.New(src,
				metrics,
				platform.WinNotifications(winNotify),
				eventStream,
				logger.With(
					zap.Uint64("platform_id", src.ID),
					zap.String("protocol", src.Protocol),
				))
		} else {
			logger.Error("invalid RTB client",
				zap.Uint64("source_id", src.ID),
				zap.String("protocol", src.Protocol))
			continue
		}

		if err != nil {
			logger.Error("invalid RTB source",
				zap.Uint64("source_id", src.ID),
				zap.String("protocol", src.Protocol),
				zap.Error(err))
		} else {
			sources = append(sources, source)
		}
	}
	return sources, err
}

type sourcesData struct {
	Sources []*imodels.RTBSource `json:"sources" yaml:"sources"`
}

func readSources(filename string) ([]*imodels.RTBSource, error) {
	var (
		sourcesData sourcesData
		data, err   = ioutil.ReadFile(filename)
	)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".json":
		err = json.Unmarshal(data, &sourcesData)
	// case ".yml", ".yaml":
	// 	err = yaml.Unmarshal(data, &sourcesData)
	// case ".xml":
	// 	err = xml.Unmarshal(data, &sourcesData)
	// case ".hcl":
	// 	err = hcl.Unmarshal(data, &sourcesData)
	default:
		err = errInvalidFileFormat
	}
	return sourcesData.Sources, err
}
