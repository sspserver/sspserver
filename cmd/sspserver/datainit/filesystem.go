//go:build fsloader
// +build fsloader

package datainit

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"geniusrabbit.dev/adcorelib/storage/loader"
	"geniusrabbit.dev/adcorelib/storage/loader/fsloader"
	"geniusrabbit.dev/adcorelib/storage/types/fstypes"
)

func init() {
	dataLoaderAccessor["fs"] = func(ctx context.Context, u *url.URL) DataLoaderAccessorFnk {
		rootDir := strings.Split(u.String()[5:], "?")[0]
		period, _ := time.ParseDuration(u.Query().Get("interval"))
		if period == 0 {
			period = time.Minute * 5
		}
		getLoader := func(pattern string) loader.LoaderFnk {
			return fsloader.PatternLoader(rootDir, pattern)
		}
		return func(ctx context.Context, dataType string) (loader.DataAccessor, error) {
			switch dataType {
			case "format":
				return loader.NewPeriodicReloader(&fstypes.FormatData{}, getLoader("format*"),
					period, period*10, "loader_formats"), nil
			case "campaign":
				return loader.NewPeriodicReloader(&fstypes.CampaignData{}, getLoader("campaign*"),
					period, period*10, "loader_campaigns"), nil
			case "company":
				return loader.NewPeriodicReloader(&fstypes.CompanyData{}, getLoader("compan*"),
					period, period*10, "loader_companies"), nil
			case "zone":
				return loader.NewPeriodicReloader(&fstypes.ZoneData{}, getLoader("zone*"),
					period, period*10, "loader_zones"), nil
			case "rtb_source", "source":
				return loader.NewPeriodicReloader(&fstypes.RTBSourceData{}, getLoader("rtb_source*"),
					period*10, period, "loader_sources"), nil
			case "rtb_access_point", "access_point":
				return loader.NewPeriodicReloader(&fstypes.RTBAccessPointData{}, getLoader("rtb_access_point*"),
					period*10, period, "loader_access_points"), nil
			}
			return nil, errors.Wrap(ErrUnsupportedDataType, dataType)
		}
	}
}
