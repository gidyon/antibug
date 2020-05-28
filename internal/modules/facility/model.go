package facility

import (
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/facility"
	"github.com/jinzhu/gorm"
)

const facilitiesTable = "facilities"

// Facility struct contains basic information about facility
type Facility struct {
	FacilityName  string `gorm:"type:varchar(100);unique_index;not null"`
	County        string `gorm:"type:varchar(100);not null"`
	CountyCode    int32  `gorm:"type:int(10);not null"`
	SubCounty     string `gorm:"type:varchar(100);not null"`
	SubCountyCode int32  `gorm:"type:int(10);not null"`
	gorm.Model
}

// TableName ...
func (*Facility) TableName() string {
	return facilitiesTable
}

func getFacilityDB(
	facilityPB *facility.Facility,
) (*Facility, error) {
	// handle nil
	if facilityPB == nil {
		return nil, errs.NilObject("FacilityPB")
	}

	facilityDB := &Facility{
		FacilityName:  facilityPB.FacilityName,
		County:        facilityPB.County,
		CountyCode:    facilityPB.CountyCode,
		SubCounty:     facilityPB.SubCounty,
		SubCountyCode: facilityPB.SubCountyCode,
	}

	return facilityDB, nil

}

func getFacilityPB(
	facilityDB *Facility,
) (*facility.Facility, error) {
	// handle nil
	if facilityDB == nil {
		return nil, errs.NilObject("FacilityDB")
	}

	facilityPB := &facility.Facility{
		FacilityId:    int64(facilityDB.ID),
		FacilityName:  facilityDB.FacilityName,
		County:        facilityDB.County,
		CountyCode:    facilityDB.CountyCode,
		SubCounty:     facilityDB.SubCounty,
		SubCountyCode: facilityDB.SubCountyCode,
	}

	return facilityPB, nil
}

// County contains name and code of county or state
type County struct {
	County string `gorm:"type:varchar(100);not null"`
	Code   int32  `gorm:"type:int(10);not null"`
	gorm.Model
}

// TableName ...
func (*County) TableName() string {
	return "counties"
}

// SubCounty contain name and code of subcounty region
type SubCounty struct {
	SubCounty string `gorm:"type:varchar(100);not null"`
	Code      int32  `gorm:"type:int(10);not null"`
	gorm.Model
}

// TableName ...
func (*SubCounty) TableName() string {
	return "sub-counties"
}
