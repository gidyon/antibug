package culture

import (
	"encoding/json"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/culture"
	"github.com/jinzhu/gorm"
)

const culturesTable = "cultures"

// Culture is culture model in database
type Culture struct {
	LabTechID           string `gorm:"type:varchar(50);not null"`
	HospitalID          string `gorm:"type:varchar(50);not null"`
	CountyCode          string `gorm:"type:int(11);not null"`
	SubCountyCode       string `gorm:"type:int(11);not null"`
	PatientID           string `gorm:"type:varchar(50);not null"`
	PatientGender       string `gorm:"type:enum('male','female','all');not null;default:'all'"`
	PatientAge          int32  `gorm:"type:tinyint(4);not null"`
	CultureSource       string `gorm:"type:varchar(50);not null"`
	TestMethod          string `gorm:"type:varchar(50);not null"`
	PathogensFound      []byte `gorm:"type:json;not null"`
	PathogensIndex      string `gorm:"type:varchar(50);not null"`
	AntimicrobialsUsed  []byte `gorm:"type:json;not null"`
	AntimicrobialsIndex string `gorm:"type:varchar(50);not null"`
	Editors             []byte `gorm:"type:json;not null"`
	CultureResults      []byte `gorm:"type:json;not null"`
	ResultsTimestampSec int64  `gorm:"type:bigint(20);not null"`
	gorm.Model
}

// TableName ...
func (*Culture) TableName() string {
	return culturesTable
}

// GetCultureDB gets the database model of a culture
func GetCultureDB(culturePB *culture.Culture) (*Culture, error) {
	return getCultureDB(culturePB)
}

func getCultureDB(culturePB *culture.Culture) (*Culture, error) {
	if culturePB == nil {
		return nil, errs.NilObject("culture pb")
	}

	cultureDB := &Culture{
		LabTechID:           culturePB.LabTechId,
		HospitalID:          culturePB.HospitalId,
		CountyCode:          culturePB.CountyCode,
		SubCountyCode:       culturePB.SubCountyCode,
		PatientID:           culturePB.PatientId,
		PatientGender:       culturePB.PatientGender,
		PatientAge:          culturePB.PatientAge,
		CultureSource:       culturePB.CultureSource,
		TestMethod:          culturePB.TestMethod.String(),
		ResultsTimestampSec: culturePB.ResultsTimestampSec,
	}

	var (
		data []byte
		err  error
	)

	// Marshal pathogens found
	if len(culturePB.PathogensFound) > 0 {
		data, err = json.Marshal(culturePB.PathogensFound)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "pathogens found")
		}
		cultureDB.PathogensFound = data
	}

	// Marshal antimicrobials used
	if len(culturePB.AntimicrobialsUsed) > 0 {
		data, err = json.Marshal(culturePB.AntimicrobialsUsed)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "antimicrobials used")
		}
		cultureDB.AntimicrobialsUsed = data
	}

	// Marshal culture results
	if len(culturePB.CultureResults) > 0 {
		data, err = json.Marshal(culturePB.CultureResults)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "culture results")
		}
		cultureDB.CultureResults = data
	}

	// Marshal editors
	if len(culturePB.Editors) > 0 {
		data, err = json.Marshal(culturePB.Editors)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "editors")
		}
		cultureDB.Editors = data
	}

	return cultureDB, nil
}

// GetCulturePB gets the protobuf message equivalence of a pathogen
func GetCulturePB(cultureDB *Culture) (*culture.Culture, error) {
	return getCulturePB(cultureDB)
}

func getCulturePB(cultureDB *Culture) (*culture.Culture, error) {
	if cultureDB == nil {
		return nil, errs.NilObject("culture db")
	}

	culturePB := &culture.Culture{
		LabTechId:           cultureDB.LabTechID,
		HospitalId:          cultureDB.HospitalID,
		CountyCode:          cultureDB.CountyCode,
		SubCountyCode:       cultureDB.SubCountyCode,
		PatientId:           cultureDB.PatientID,
		PatientGender:       cultureDB.PatientGender,
		PatientAge:          cultureDB.PatientAge,
		CultureSource:       cultureDB.CultureSource,
		TestMethod:          culture.TestMethod(culture.TestMethod_value[cultureDB.TestMethod]),
		ResultsTimestampSec: cultureDB.ResultsTimestampSec,
	}

	var err error
	// Unmarshal pathogens found
	if len(cultureDB.PathogensFound) > 0 {
		err = json.Unmarshal(cultureDB.PathogensFound, &culturePB.PathogensFound)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "pathogens found")
		}
	}

	// Unmarshal antimicrobials used
	if len(cultureDB.AntimicrobialsUsed) > 0 {
		err = json.Unmarshal(cultureDB.AntimicrobialsUsed, &culturePB.AntimicrobialsUsed)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "antimicrobials used")
		}
	}

	// Unmarshal culture results
	if len(cultureDB.CultureResults) > 0 {
		err = json.Unmarshal(cultureDB.CultureResults, &culturePB.CultureResults)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "culture results")
		}
	}

	// Unmarshal editors
	if len(cultureDB.Editors) > 0 {
		err = json.Unmarshal(cultureDB.Editors, &culturePB.Editors)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "editors")
		}
	}

	return culturePB, nil
}
