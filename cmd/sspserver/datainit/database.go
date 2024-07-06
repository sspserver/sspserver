//go:build dbloader
// +build dbloader

package datainit

import (
	"context"
	"net/url"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/pkg/errors"

	"github.com/geniusrabbit/adcorelib/models"
	"github.com/geniusrabbit/adcorelib/storage/loader"
	"github.com/geniusrabbit/adcorelib/storage/loader/dbloader"
	"github.com/geniusrabbit/adcorelib/storage/types/dbtypes"

	"github.com/sspserver/sspserver/internal/database"
)

func init() {
	for _, dialect := range database.ListOfDialects() {
		dataLoaderAccessor[dialect] = dbConnector
	}
}

func dbConnector(ctx context.Context, u *url.URL) DataLoaderAccessorFnk {
	db, err := database.Connect(ctx, u.String(), gocast.Bool(u.Query().Get("debug")))
	if err != nil {
		panic(err)
	}
	period, _ := time.ParseDuration(u.Query().Get("interval"))
	if period == 0 {
		period = time.Minute * 5
	}
	var formatAccessor loader.DataAccessor
	return func(ctx context.Context, dataType string) (loader.DataAccessor, error) {
		if formatAccessor == nil {
			switch dataType {
			case "format", "campaign":
				formatAccessor = loader.NewPeriodicReloader(&dbtypes.FormatList{},
					dbloader.Loader(db), period, period*10, "loader_formats")
			}
		}
		switch dataType {
		case "format":
			return formatAccessor, nil
		case "campaign":
			return loader.NewCombinedLoader(
				campaingMerge,
				loader.NewPeriodicReloader(&dbtypes.CampaignList{}, dbloader.Loader(db), period, period*10, "loader_campaigns"),
				loader.NewPeriodicReloader(&dbtypes.AdList{}, dbloader.Loader(db), period, period*10, "loader_ads"),
				loader.NewPeriodicReloader(&dbtypes.LinkList{}, dbloader.Loader(db), period, period*10, "loader_adlinks"),
				formatAccessor,
			), nil
		case "company":
			return loader.NewPeriodicReloader(&dbtypes.CompanyList{}, dbloader.Loader(db),
				period, period*10, "loader_companies"), nil
		case "zone":
			return loader.NewPeriodicReloader(&dbtypes.ZoneList{}, dbloader.Loader(db),
				period, period*10, "loader_zones"), nil
		case "rtb_source", "source":
			return loader.NewPeriodicReloader(&dbtypes.RTBSourceList{}, dbloader.Loader(db),
				period, period*10, "loader_sources"), nil
		case "rtb_access_point", "access_point":
			return loader.NewPeriodicReloader(&dbtypes.RTBAccessPointList{}, dbloader.Loader(db),
				period, period*10, "loader_access_points"), nil
		}
		return nil, errors.Wrap(ErrUnsupportedDataType, dataType)
	}
}

func campaingMerge(datas ...[]any) []any {
	camps := datas[0]
	for _, cmpBase := range camps {
		cmp := cmpBase.(*models.Campaign)
		for _, adBase := range datas[1] {
			ad := adBase.(*models.Ad)
			if ad.Status.IsApproved() && ad.Active.IsActive() && ad.CampaignID == cmp.ID {
				ad.Campaign = cmp
				for _, fBase := range datas[3] {
					format := fBase.(*models.Format)
					if ad.FormatID == format.ID {
						ad.Format = format
						break
					}
				}
				if ad.Format == nil || (ad.Format.IsNative() && len(ad.Assets) == 0) {
					continue
				}
				cmp.Ads = append(cmp.Ads, ad)
			}
		}
		for _, adBase := range datas[2] {
			link := adBase.(*models.AdLink)
			if link.Status.IsApproved() && link.Active.IsActive() && link.CampaignID == cmp.ID {
				link.Campaign = cmp
				cmp.Links = append(cmp.Links, link)
			}
		}
		if len(cmp.Ads) == 0 {
			cmp.Links = nil
		}
	}
	return camps
}
