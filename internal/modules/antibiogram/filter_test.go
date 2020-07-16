package antibiogram

import (
	"github.com/gidyon/antibug/internal/modules/culture"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	"math/rand"
)

func pastDuration() antibiogram.Duration {
	return antibiogram.Duration(
		antibiogram.Duration_value[antibiogram.Duration_name[int32(rand.Intn(len(antibiogram.Duration_value)))]],
	)
}

func regionScope() antibiogram.RegionScope {
	return antibiogram.RegionScope(
		antibiogram.RegionScope_value[antibiogram.RegionScope_name[int32(rand.Intn(len(antibiogram.RegionScope_value)))]],
	)
}

func gender() antibiogram.Gender {
	return antibiogram.Gender(
		antibiogram.Gender_value[antibiogram.Gender_name[int32(rand.Intn(len(antibiogram.Gender_value)))]],
	)
}

const (
	subjectPathogen      = "pathogens"
	subjectAntimicrobial = "antimicrobial"
)

func fakeFilter(subjectType string) *antibiogram.Filter {
	regionScope := regionScope()
	scopeValues := make([]string, 0, 3)
	switch regionScope {
	case antibiogram.RegionScope_COUNTRY:
	case antibiogram.RegionScope_COUNTY:
		scopeValues = []string{culture.CountyCode(), culture.CountyCode(), culture.CountyCode()}
	case antibiogram.RegionScope_SUB_COUNTY:
		scopeValues = []string{culture.SubCountyCode(), culture.SubCountyCode(), culture.SubCountyCode()}
	case antibiogram.RegionScope_FACILITY:
		scopeValues = []string{culture.HospitalID(), culture.HospitalID(), culture.HospitalID()}
	}

	var inputValues []*antibiogram.Value

	switch subjectType {
	case subjectPathogen:
		inputValues = []*antibiogram.Value{
			{Id: culture.Pathogen(), Name: culture.Pathogen()},
			{Id: culture.Pathogen(), Name: culture.Pathogen()},
			{Id: culture.Pathogen(), Name: culture.Pathogen()},
		}
	case subjectAntimicrobial:
		inputValues = []*antibiogram.Value{
			{Id: culture.Antimicrobial(), Name: culture.Antimicrobial()},
			{Id: culture.Antimicrobial(), Name: culture.Antimicrobial()},
			{Id: culture.Antimicrobial(), Name: culture.Antimicrobial()},
		}
	default:
		inputValues = []*antibiogram.Value{
			{Id: culture.Pathogen(), Name: culture.Pathogen()},
			{Id: culture.Pathogen(), Name: culture.Pathogen()},
			{Id: culture.Pathogen(), Name: culture.Pathogen()},
		}
	}

	return &antibiogram.Filter{
		PastDuration: pastDuration(),
		RegionScope:  regionScope,
		ScopeValues:  scopeValues,
		InputValues:  inputValues,
		Advanced:     true,
		Advance: &antibiogram.AdvancedFilter{
			Gender:     gender(),
			AgeMinDays: 365 * 1,
			AgeMaxDays: 365 * 100,
		},
	}
}
