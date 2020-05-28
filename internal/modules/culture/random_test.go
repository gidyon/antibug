package culture

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/culture"
	"math/rand"
	"time"
)

var (
	pathogens      []string
	antimicrobials []string
	samples        = []string{"blood", "urine", "saliva", "blood", "faeces", "fluid"}
	letters        = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "y", "z"}
)

func initCulture() {
	// seed random generator
	rand.Seed(time.Now().UnixNano())

	// init data
	pathogens = make([]string, 0)
	antimicrobials = make([]string, 0)

	// Initialize pathogens and antimicrobials
	for _, letter := range letters {
		pathogens = append(pathogens, fmt.Sprintf("pathogen-%s", letter))
		antimicrobials = append(antimicrobials, fmt.Sprintf("antimicrobial-%s", letter))
	}
}

// CountyCode returns a random county code in range of 1 - 47
func CountyCode() string {
	// number between 1 - 47
	return fmt.Sprintf("%d", rand.Intn(47))
}

// SubCountyCode returns a random sub-county code in range of 0 - 200
func SubCountyCode() string {
	// number between 1 - 200
	return fmt.Sprintf("%d", rand.Intn(200))
}

// HospitalID returns a random hospital id in range of 50 -100
func HospitalID() string {
	// number between 10-20
	return fmt.Sprintf("hospital-%d", rand.Intn(50)+50)
}

// PatientID returns a random patient id in range of 30 - 100
func PatientID() string {
	// number between 30-100
	return fmt.Sprintf("patient-%d", rand.Intn(30)+70)
}

// PatientGender returns a random gender either male or female
func PatientGender() string {
	if rand.Intn(50)%2 == 0 {
		return "male"
	}
	return "female"
}

// Pathogen returns a random pathogen from list of pathogens
func Pathogen() string {
	return pathogens[rand.Intn(len(pathogens))]
}

// Antimicrobial returns a random antimicrobial from list of antimicrobials
func Antimicrobial() string {
	return antimicrobials[rand.Intn(len(antimicrobials))]
}

// Sample returns a random culture source
func Sample() string {
	return samples[rand.Intn(len(samples))]
}

// Label returns a random label for a culture result
func Label() culture.Label {
	return culture.Label(
		culture.Label_value[culture.Label_name[rand.Int31n(int32(len(culture.Label_name)))]],
	)
}

// DiskDiameter returns a random disk diameter for a culture result between 0 - 30
func DiskDiameter() string {
	return fmt.Sprintf("%d mm", rand.Intn(30))
}

func updateLabTestResults(culturePB *culture.Culture) *culture.Culture {
	if culturePB.CultureResults == nil {
		culturePB.CultureResults = make([]*culture.LabTestResult, 0, 3)
	}

	for i := 0; i < len(culturePB.AntimicrobialsUsed); i++ {
		culturePB.CultureResults = append(culturePB.CultureResults, &culture.LabTestResult{
			PathogenName:      culturePB.PathogensFound[i],
			PathogenId:        culturePB.PathogensFound[i],
			AntimicrobialId:   culturePB.AntimicrobialsUsed[i],
			AntimicrobialName: culturePB.AntimicrobialsUsed[i],
			DiskDiameter:      DiskDiameter(),
			Label:             Label(),
			ResultComment:     randomdata.Paragraph(),
		})
	}

	return culturePB
}

func fakeCulture() *culture.Culture {
	labTechID := randomdata.RandStringRunes(20)
	return updateLabTestResults(&culture.Culture{
		CultureId:     randomdata.RandStringRunes(20),
		LabTechId:     randomdata.RandStringRunes(20),
		HospitalId:    HospitalID(),
		CountyCode:    CountyCode(),
		SubCountyCode: SubCountyCode(),
		PatientId:     PatientID(),
		PatientGender: PatientGender(),
		PatientAge:    fmt.Sprint(randomdata.Number(10, 40)),
		Editors:       []string{labTechID},
		TestMethod:    culture.TestMethod_DISK_DIFFUSION,
		CultureSource: Sample(),
		PathogensFound: []string{
			Pathogen(), Pathogen(), Pathogen(),
		},
		AntimicrobialsUsed: []string{
			Antimicrobial(), Antimicrobial(), Antimicrobial(),
		},
		CultureResults:      []*culture.LabTestResult{},
		ResultsTimestampSec: time.Now().Unix(),
	})
}
