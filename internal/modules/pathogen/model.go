package pathogen

import (
	"encoding/json"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"github.com/jinzhu/gorm"
)

const pathogensTable = "pathogens"

// Pathogen is a model for pathogen resource
type Pathogen struct {
	PathogenName            string `gorm:"type:varchar(100);unique;not null"`
	Category                string `gorm:"type:varchar(50);not null"`
	GeneralInformation      string `gorm:"type:varchar(512);not null;"`
	Epidemology             []byte `gorm:"type:json;not null"`
	Symptoms                []byte `gorm:"type:json;not null"`
	AdditionalInformation   []byte `gorm:"type:json;not null"`
	GeneralSusceptibilities []byte `gorm:"type:json;not null"`
	Editors                 []byte `gorm:"type:json;not null"`
	gorm.Model
}

// TableName ...
func (*Pathogen) TableName() string {
	return pathogensTable
}

func getPathogenDB(pathogenPB *pathogen.Pathogen) (*Pathogen, error) {
	if pathogenPB == nil {
		return nil, errs.NilObject("PathogenPB")
	}

	pathogenDB := &Pathogen{
		PathogenName:       pathogenPB.PathogenName,
		Category:           pathogenPB.Category,
		GeneralInformation: pathogenPB.GeneralInformation,
	}

	// Marshal epidemology
	// Marshal symptoms
	// Marshal additional info
	// Marshal general susceptibilities
	// Marshal editors

	var (
		err  error
		data []byte
	)

	if len(pathogenPB.GetEpidemology().GetValues()) > 0 {
		data, err = json.Marshal(pathogenPB.Epidemology)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "epidemology")
		}
		pathogenDB.Epidemology = data
	}

	if len(pathogenPB.GetSymptoms().GetValues()) > 0 {
		data, err = json.Marshal(pathogenPB.Symptoms)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "symptoms")
		}
		pathogenDB.Symptoms = data
	}

	if len(pathogenPB.GetAdditionalInfo().GetValues()) > 0 {
		data, err = json.Marshal(pathogenPB.AdditionalInfo)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "additional information")
		}
		pathogenDB.AdditionalInformation = data
	}

	if len(pathogenPB.GetGeneralSusceptibilities().GetSusceptibilities()) > 0 {
		data, err = json.Marshal(pathogenPB.GeneralSusceptibilities)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "general susceptibilities")
		}
		pathogenDB.GeneralSusceptibilities = data
	}

	if len(pathogenPB.GetEditors().GetValues()) > 0 {
		data, err = json.Marshal(pathogenPB.Editors)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "editors")
		}
		pathogenDB.Editors = data
	}

	return pathogenDB, nil
}

func getPathogenPB(pathogenDB *Pathogen) (*pathogen.Pathogen, error) {
	if pathogenDB == nil {
		return nil, errs.NilObject("PathogenDB")
	}

	pathogenPB := &pathogen.Pathogen{
		PathogenId:         int64(pathogenDB.ID),
		PathogenName:       pathogenDB.PathogenName,
		GeneralInformation: pathogenDB.GeneralInformation,
		Category:           pathogenDB.Category,
		UpdateTimeSec:      pathogenDB.UpdatedAt.Unix(),
	}

	var (
		err error
	)

	// Unmarshal epidemology
	if len(pathogenDB.Epidemology) > 0 {
		err = json.Unmarshal(pathogenDB.Epidemology, &pathogenPB.Epidemology)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "epidemology")
		}
	}

	// Unmarshal symptoms
	if len(pathogenDB.Symptoms) > 0 {
		err = json.Unmarshal(pathogenDB.Symptoms, &pathogenPB.Symptoms)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "symptoms")
		}
	}

	// Unmarshal additional info
	if len(pathogenDB.AdditionalInformation) > 0 {
		err = json.Unmarshal(pathogenDB.AdditionalInformation, &pathogenPB.AdditionalInfo)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "additional information")
		}
	}

	// Unmarshal general susceptibilities
	if len(pathogenDB.GeneralSusceptibilities) > 0 {
		err = json.Unmarshal(pathogenDB.GeneralSusceptibilities, &pathogenPB.GeneralSusceptibilities)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "general susceptibilities")
		}
	}

	// Unmarshal editors
	if len(pathogenDB.Editors) > 0 {
		err = json.Unmarshal(pathogenDB.Editors, &pathogenPB.Editors)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "editors")
		}
	}

	return pathogenPB, nil
}
