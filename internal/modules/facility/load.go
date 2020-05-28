package facility

import (
	"github.com/gidyon/antibug/pkg/api/facility"
)

func (fapi *facilityAPIServer) loadCountiesData() error {
	err := fapi.loadAndCacheCounties()
	if err != nil {
		return err
	}

	err = fapi.loadAndCacheSubCounties()
	if err != nil {
		return err
	}

	return nil
}

func (fapi *facilityAPIServer) loadAndCacheCounties() error {
	// `county` varchar(50) NOT NULL,
	// `code` int(11) NOT NULL

	countiesDB := make([]*County, 0)

	err := fapi.sqlDB.Find(&countiesDB).Error
	if err != nil {
		return err
	}

	counties := make([]*facility.County, 0)

	for _, countyDB := range countiesDB {
		counties = append(counties, &facility.County{
			County: countyDB.County,
			Code:   int32(countyDB.ID),
		})
	}

	fapi.counties = counties

	return nil
}

func (fapi *facilityAPIServer) loadAndCacheSubCounties() error {
	// `sub_county` varchar(50) NOT NULL,
	// `code` int(11) NOT NULL

	subCountiesDB := make([]*SubCounty, 0)

	err := fapi.sqlDB.Find(&subCountiesDB).Error
	if err != nil {
		return err
	}

	subCounties := make([]*facility.SubCounty, 0)

	for _, subCountyDB := range subCountiesDB {
		subCounties = append(subCounties, &facility.SubCounty{
			SubCounty: subCountyDB.SubCounty,
			Code:      int32(subCountyDB.ID),
		})
	}

	fapi.subCounties = subCounties

	return nil
}
