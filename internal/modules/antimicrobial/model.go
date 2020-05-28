package antimicrobial

import (
	"encoding/json"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"github.com/jinzhu/gorm"
)

const antimicrobialsTable = "antimicrobials"

// Antimicrobial is a model for antimicrobial resource
type Antimicrobial struct {
	AntimicrobialName     string `gorm:"varchar(100);unique_index;not null"`
	CDiff                 string `gorm:"varchar(100);not null;default:'NA'"`
	OralBioavailability   string `gorm:"varchar(30);not null;default:'NA'"`
	ApproximateCost       string `gorm:"varchar(14);not null;default:'NA'"`
	GeneralUsage          []byte `gorm:"type:json;not null"`
	DrugMonitoring        []byte `gorm:"type:json"`
	AdverseEffects        []byte `gorm:"type:json;not null"`
	MajorInteractions     []byte `gorm:"type:json"`
	Pharmacology          []byte `gorm:"type:json;not null"`
	AdditionalInformation []byte `gorm:"type:json"`
	ActivitySpectrum      []byte `gorm:"type:json;not null"`
	Editors               []byte `gorm:"type:json;not null"`
	gorm.Model
}

// TableName ...
func (*Antimicrobial) TableName() string {
	return antimicrobialsTable
}

func getAntimicrobialDB(
	antimicrobialPB *antimicrobial.Antimicrobial,
) (*Antimicrobial, error) {
	// handle nil
	if antimicrobialPB == nil {
		return nil, errs.NilObject("AntimicrobialPB")
	}

	antimicrobialDB := &Antimicrobial{
		AntimicrobialName:   antimicrobialPB.AntimicrobialName,
		CDiff:               antimicrobialPB.CDiff,
		OralBioavailability: antimicrobialPB.OralBioavailability,
		ApproximateCost:     antimicrobialPB.ApproximateCost,
	}

	var (
		err  error
		data []byte
	)

	if len(antimicrobialPB.GetGeneralUsage().GetValues()) > 0 {
		data, err = json.Marshal(antimicrobialPB.GeneralUsage)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.GeneralUsage")
		}
		antimicrobialDB.GeneralUsage = data
	}

	if len(antimicrobialPB.GetDrugMonitoring().GetValues()) > 0 {
		data, err = json.Marshal(antimicrobialPB.DrugMonitoring)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.DrugMonitoring")
		}
		antimicrobialDB.DrugMonitoring = data
	}

	if len(antimicrobialPB.GetAdverseEffects().GetValues()) > 0 {
		data, err = json.Marshal(antimicrobialPB.AdverseEffects)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.AdverseEffects")
		}
		antimicrobialDB.AdverseEffects = data
	}

	if len(antimicrobialPB.GetMajorInteractions().GetValues()) > 0 {
		data, err = json.Marshal(antimicrobialPB.MajorInteractions)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.MajorIntercations")
		}
		antimicrobialDB.MajorInteractions = data
	}

	if len(antimicrobialPB.GetPharmacology().GetPharmacologyInfos()) > 0 {
		data, err = json.Marshal(antimicrobialPB.Pharmacology)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.Pharmacology")
		}
		antimicrobialDB.Pharmacology = data
	}

	if len(antimicrobialPB.GetAdditionalInformation().GetValues()) > 0 {
		data, err = json.Marshal(antimicrobialPB.AdditionalInformation)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.AdditionalInformation")
		}
		antimicrobialDB.AdditionalInformation = data
	}

	if len(antimicrobialPB.GetActivitySpectrum().GetSpectrum()) > 0 {
		data, err = json.Marshal(antimicrobialPB.ActivitySpectrum)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.ActivitySpectrum")
		}
		antimicrobialDB.ActivitySpectrum = data
	}

	if len(antimicrobialPB.GetEditors().GetValues()) > 0 {
		data, err = json.Marshal(antimicrobialPB.Editors)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Antimicrobial.Editors")
		}
		antimicrobialDB.Editors = data
	}

	return antimicrobialDB, nil
}

func getAntimicrobialPB(
	antimicrobialDB *Antimicrobial,
) (*antimicrobial.Antimicrobial, error) {
	// handle nil
	if antimicrobialDB == nil {
		return nil, errs.NilObject("AntimicrobialDB")
	}

	antimicrobialPB := &antimicrobial.Antimicrobial{
		AntimicrobialId:       int64(antimicrobialDB.ID),
		AntimicrobialName:     antimicrobialDB.AntimicrobialName,
		CDiff:                 antimicrobialDB.CDiff,
		OralBioavailability:   antimicrobialDB.OralBioavailability,
		ApproximateCost:       antimicrobialDB.ApproximateCost,
		GeneralUsage:          createRepeatedString(),
		DrugMonitoring:        createRepeatedString(),
		AdverseEffects:        createRepeatedString(),
		MajorInteractions:     createRepeatedString(),
		AdditionalInformation: createRepeatedString(),
		Editors:               createRepeatedString(),
		Pharmacology: &antimicrobial.Pharmacology{
			PharmacologyInfos: make([]*antimicrobial.PharmacologyInfo, 0),
		},
		ActivitySpectrum: &antimicrobial.SpectrumOfActivity{
			Spectrum: make([]*antimicrobial.Spectrum, 0),
		},
	}

	var (
		err error
	)

	// Because json.Unmarshal is not nil safe, we check that the data to unmarshal is safe
	if len(antimicrobialDB.GeneralUsage) > 0 {
		err = json.Unmarshal(antimicrobialDB.GeneralUsage, &antimicrobialPB.GeneralUsage)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.GeneralUsage")
		}
	}

	if len(antimicrobialDB.DrugMonitoring) > 0 {
		err = json.Unmarshal(antimicrobialDB.DrugMonitoring, &antimicrobialPB.DrugMonitoring)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.DrugMonitoring")
		}
	}

	if len(antimicrobialDB.AdverseEffects) > 0 {
		err = json.Unmarshal(antimicrobialDB.AdverseEffects, &antimicrobialPB.AdverseEffects)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.AdverseEffects")
		}
	}

	if len(antimicrobialDB.MajorInteractions) > 0 {
		err = json.Unmarshal(antimicrobialDB.MajorInteractions, &antimicrobialPB.MajorInteractions)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.MajorIntercations")
		}
	}

	if len(antimicrobialDB.Pharmacology) > 0 {
		err = json.Unmarshal(antimicrobialDB.Pharmacology, &antimicrobialPB.Pharmacology)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.Pharmacology")
		}
	}

	if len(antimicrobialDB.AdditionalInformation) > 0 {
		err = json.Unmarshal(antimicrobialDB.AdditionalInformation, &antimicrobialPB.AdditionalInformation)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.AdditionalInformation")
		}
	}

	if len(antimicrobialDB.ActivitySpectrum) > 0 {
		err = json.Unmarshal(antimicrobialDB.ActivitySpectrum, &antimicrobialPB.ActivitySpectrum)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.ActivitySpectrum")
		}
	}

	if len(antimicrobialDB.Editors) > 0 {
		err = json.Unmarshal(antimicrobialDB.Editors, &antimicrobialPB.Editors)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Antimicrobial.Editors")
		}
	}

	return antimicrobialPB, nil
}

func createRepeatedString() *antimicrobial.RepeatedString {
	return &antimicrobial.RepeatedString{Values: make([]string, 0)}
}
