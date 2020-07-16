package culture

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/culture"
	"math/rand"
	"time"
)

var (
	// Pathogens ...
	Pathogens []string
	// Antimicrobials ...
	Antimicrobials []string
	// Samples ...
	Samples = []string{"blood", "urine", "saliva", "blood", "faeces", "fluid"}
	letters = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "y", "z"}
)

// InitCulture initializes some global variables
func InitCulture() {
	// seed random generator
	rand.Seed(time.Now().UnixNano())

	// init data
	Pathogens = make([]string, 0)
	Antimicrobials = make([]string, 0)

	// Initialize Pathogens and Antimicrobials
	for _, letter := range letters {
		Pathogens = append(Pathogens, fmt.Sprintf("pathogen-%s", letter))
		Antimicrobials = append(Antimicrobials, fmt.Sprintf("antimicrobial-%s", letter))
	}
}

// CountyCode returns a random county code in range of 1 - 5
func CountyCode() string {
	// number between 1 - 47
	return fmt.Sprintf("%d", rand.Intn(5)+1)
}

// SubCountyCode returns a random sub-county code in range of 1 - 10
func SubCountyCode() string {
	// number between 1 - 200
	return fmt.Sprintf("%d", rand.Intn(10)+1)
}

// HospitalID returns a random hospital id in range of 1 - 10
func HospitalID() string {
	// number between 10-20
	return fmt.Sprintf("hospital-%d", rand.Intn(10)+1)
}

// PatientID returns a random patient id in range of 1 - 15
func PatientID() string {
	// number between 30-100
	return fmt.Sprintf("patient-%d", rand.Intn(15)+1)
}

// LabTechID returns a random lab tech id in range of 1- 5
func LabTechID() string {
	// number between 30-100
	return fmt.Sprintf("labtech-%d", rand.Intn(5)+1)
}

// PatientGender returns a random gender either male or female
func PatientGender() string {
	if rand.Intn(50)%2 == 0 {
		return "male"
	}
	return "female"
}

// Pathogen returns a random pathogen from list of Pathogens
func Pathogen() string {
	return Pathogens[rand.Intn(len(Pathogens))]
}

// Antimicrobial returns a random antimicrobial from list of Antimicrobials
func Antimicrobial() string {
	return Antimicrobials[rand.Intn(len(Antimicrobials))]
}

// Sample returns a random culture source
func Sample() string {
	return Samples[rand.Intn(len(Samples))]
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

// FakeCulture creates a fake culture
func FakeCulture() *culture.Culture {
	labTechID := LabTechID()
	return updateLabTestResults(&culture.Culture{
		CultureId:     randomdata.RandStringRunes(20),
		LabTechId:     labTechID,
		HospitalId:    HospitalID(),
		CountyCode:    CountyCode(),
		SubCountyCode: SubCountyCode(),
		PatientId:     PatientID(),
		PatientGender: PatientGender(),
		PatientAge:    int32(randomdata.Number(1, 100)),
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
