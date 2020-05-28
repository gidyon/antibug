package pathogen

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
)

var _ = Describe("Creating Pathogen #create", func() {
	var (
		createReq *pathogen.CreatePathogenRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &pathogen.CreatePathogenRequest{
			Pathogen: newPathogen(),
		}
		ctx = context.Background()
	})

	Describe("Creating pathogen with nil request", func() {
		It("should fail when request is nil", func() {
			createReq = nil
			createRes, err := PathogenAPI.CreatePathogen(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
	})

	When("Creating pathogen with some missing pathogen fields", func() {
		It("should fail when pathogen name is missing", func() {
			createReq.Pathogen.PathogenName = ""
			createRes, err := PathogenAPI.CreatePathogen(context.Background(), createReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(createRes).Should(BeNil())
		})
		It("should fail when general information is missing", func() {
			createReq.Pathogen.GeneralInformation = ""
			createRes, err := PathogenAPI.CreatePathogen(context.Background(), createReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(createRes).Should(BeNil())
		})
		It("should fail when category is missing", func() {
			createReq.Pathogen.Category = ""
			createRes, err := PathogenAPI.CreatePathogen(context.Background(), createReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(createRes).Should(BeNil())
		})
		It("should fail when epidemology is nil or with len of values zero", func() {
			createReq.Pathogen.Epidemology = nil
			createRes, err := PathogenAPI.CreatePathogen(context.Background(), createReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(createRes).Should(BeNil())
		})
		It("should fail when symptoms associated is nil or with len of values zero", func() {
			createReq.Pathogen.Symptoms = nil
			createRes, err := PathogenAPI.CreatePathogen(context.Background(), createReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(createRes).Should(BeNil())
		})
		It("should fail when additional information is nil or with len of values zero", func() {
			createReq.Pathogen.AdditionalInfo = nil
			createRes, err := PathogenAPI.CreatePathogen(context.Background(), createReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(createRes).Should(BeNil())
		})
		It("should fail when general susceptibilities is nil or with len of values zero", func() {
			createReq.Pathogen.GeneralSusceptibilities = nil
			createRes, err := PathogenAPI.CreatePathogen(context.Background(), createReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(createRes).Should(BeNil())
		})
	})

	Context("Creating pathogen with a name that already exists", func() {
		pathogenName := randomdata.SillyName()
		It("should succeed if pathogen name does not exist", func() {
			createReq.Pathogen.PathogenName = pathogenName
			createRes, err := PathogenAPI.CreatePathogen(ctx, createReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(createRes).ToNot(BeNil())
		})
		It("should fail with already exists error if name exists", func() {
			createReq.Pathogen.PathogenName = pathogenName
			createRes, err := PathogenAPI.CreatePathogen(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.AlreadyExists))
			Expect(createRes).To(BeNil())
		})
	})

	When("Creating pathogen with valid request", func() {
		It("should create pathogen in database without error", func() {
			createRes, err := PathogenAPI.CreatePathogen(ctx, createReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(createRes).ToNot(BeNil())
		})
	})
})

var pathogenCategories = []string{"virus", "fungi", "protozoa", "rickettsia", "helminth"}

// generates new pathogen agent
func newPathogen() *pathogen.Pathogen {
	return &pathogen.Pathogen{
		PathogenName:       randomdata.SillyName() + " " + randomdata.SillyName(),
		GeneralInformation: randomdata.Paragraph(),
		Category:           pathogenCategories[rand.Intn(len(pathogenCategories))],
		Epidemology: &pathogen.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		Symptoms: &pathogen.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		AdditionalInfo: &pathogen.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		Editors: &pathogen.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		GeneralSusceptibilities: &pathogen.Susceptibilities{
			Susceptibilities: []*pathogen.Susceptibility{
				&pathogen.Susceptibility{
					Title: "Susceptible",
					Antibiotics: []string{
						randomdata.SillyName(), randomdata.SillyName(), randomdata.SillyName(),
					},
				},
				&pathogen.Susceptibility{
					Title: "Susceptible",
					Antibiotics: []string{
						randomdata.SillyName(), randomdata.SillyName(), randomdata.SillyName(),
					},
				},
				&pathogen.Susceptibility{
					Title: "Intermediate",
					Antibiotics: []string{
						randomdata.SillyName(), randomdata.SillyName(), randomdata.SillyName(),
					},
				},
			},
		},
		UpdateTimeSec: time.Now().Unix(),
	}
}
