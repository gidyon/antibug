package antimicrobial

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var _ = Describe("Creating Antimicrobial #create", func() {
	var (
		createReq *antimicrobial.CreateAntimicrobialRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &antimicrobial.CreateAntimicrobialRequest{
			Antimicrobial: newAntimicrobial(),
		}
		ctx = context.Background()
	})

	Describe("Creating antimicrobial with nil request", func() {
		It("should fail when request is nil", func() {
			createReq = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
	})

	When("Creating antimicrobial with some missing antimicrobial fields", func() {
		It("should fail when antimicrobial name is missing", func() {
			createReq.Antimicrobial.AntimicrobialName = ""
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when antimicrobial cdiff is missing", func() {
			createReq.Antimicrobial.CDiff = ""
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when oral bioavailability is missing", func() {
			createReq.Antimicrobial.OralBioavailability = ""
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when approximate cost is missing", func() {
			createReq.Antimicrobial.ApproximateCost = ""
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when general usage has zero len or is nil", func() {
			createReq.Antimicrobial.GeneralUsage = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when drug monitoring has zero len or is nil", func() {
			createReq.Antimicrobial.DrugMonitoring = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when adverse effects has zero len or is nil", func() {
			createReq.Antimicrobial.AdverseEffects = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when major interactions has zero len or is nil", func() {
			createReq.Antimicrobial.MajorInteractions = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when pharmacology has zero len or is nil", func() {
			createReq.Antimicrobial.Pharmacology = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when additional information has zero len or is nil", func() {
			createReq.Antimicrobial.AdditionalInformation = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when activity spectrum has zero len or is nil", func() {
			createReq.Antimicrobial.ActivitySpectrum = nil
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
	})

	Context("Creating antimicrobial with a name that already exists", func() {
		antimicrobialName := randomdata.SillyName()
		It("should succeed if antimicrobial name does not exist", func() {
			createReq.Antimicrobial.AntimicrobialName = antimicrobialName
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(createRes).ToNot(BeNil())
		})
		It("should fail with already exists error", func() {
			createReq.Antimicrobial.AntimicrobialName = antimicrobialName
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.AlreadyExists))
			Expect(createRes).To(BeNil())
		})
	})

	When("Creating antimicrobial with valid request", func() {
		It("should create antimicrobial in database without error", func() {
			createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(createRes).ToNot(BeNil())
		})
	})
})

func newAntimicrobial() *antimicrobial.Antimicrobial {
	cDiff := randomdata.Paragraph()
	if len(cDiff) > 50 {
		cDiff = cDiff[:50]
	}
	return &antimicrobial.Antimicrobial{
		AntimicrobialName:   randomdata.SillyName() + " " + randomdata.SillyName(),
		CDiff:               cDiff,
		OralBioavailability: randomdata.Paragraph()[:30],
		ApproximateCost:     "Ksh. 500",
		GeneralUsage: &antimicrobial.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		DrugMonitoring: &antimicrobial.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		AdverseEffects: &antimicrobial.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		MajorInteractions: &antimicrobial.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		Editors: &antimicrobial.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		AdditionalInformation: &antimicrobial.RepeatedString{
			Values: []string{randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph()},
		},
		Pharmacology: &antimicrobial.Pharmacology{
			PharmacologyInfos: []*antimicrobial.PharmacologyInfo{
				&antimicrobial.PharmacologyInfo{
					Key: "-", Value: randomdata.RandStringRunes(10),
				},
				&antimicrobial.PharmacologyInfo{
					Key: "+", Value: randomdata.RandStringRunes(10),
				},
				&antimicrobial.PharmacologyInfo{
					Key: "=", Value: randomdata.RandStringRunes(10),
				},
			},
		},
		ActivitySpectrum: &antimicrobial.SpectrumOfActivity{
			Spectrum: []*antimicrobial.Spectrum{
				&antimicrobial.Spectrum{
					Group: "Gram+",
					Microbes: []*antimicrobial.MicrobesInfo{
						&antimicrobial.MicrobesInfo{
							Name: randomdata.SillyName(), Id: randomdata.RandStringRunes(32),
						},
						&antimicrobial.MicrobesInfo{
							Name: randomdata.SillyName(), Id: randomdata.RandStringRunes(32),
						},
					},
				},
				&antimicrobial.Spectrum{
					Group: "Gram-",
					Microbes: []*antimicrobial.MicrobesInfo{
						&antimicrobial.MicrobesInfo{
							Name: randomdata.SillyName(), Id: randomdata.RandStringRunes(32),
						},
						&antimicrobial.MicrobesInfo{
							Name: randomdata.SillyName(), Id: randomdata.RandStringRunes(32),
						},
					},
				},
			},
		},
		UpdateTimeSec: time.Now().Unix(),
	}
}
