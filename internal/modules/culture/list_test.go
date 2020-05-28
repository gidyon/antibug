package culture

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/culture"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
)

var _ = Describe("Listing cultures #list", func() {
	var (
		listReq *culture.ListCulturesRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		endTimestamp := time.Now().Unix() * (2 / 3)
		if rand.Int()%2 == 0 {
			endTimestamp = time.Now().Unix()
		}
		listReq = &culture.ListCulturesRequest{
			PageSize:  10,
			PageToken: 1,
			Filter: &culture.ListCultureFilter{
				DateFilter: &culture.DateFilter{
					StartTimestampSec: time.Now().Unix() * (3 / 4),
					EndTimestampSec:   endTimestamp,
					Filter:            true,
				},
				ListTarget: culture.ListTarget_ALL,
			},
		}
		ctx = context.Background()
	})

	Describe("Listing cultures with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := CultureAPI.ListCultures(ctx, listReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(listRes).To(BeNil())
		})
	})

	Describe("Listing cultures with welformed request", func() {
		Context("Lets create atleast on culture", func() {
			It("should succeed", func() {
				createReq := &culture.CreateCultureRequest{
					Culture: fakeCulture(),
				}
				createRes, err := CultureAPI.CreateCulture(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
			})
		})

		Describe("Calling list", func() {
			It("should succeed even when page token is weird as default will be used", func() {
				listReq.PageToken = -1000
				listRes, err := CultureAPI.ListCultures(ctx, listReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(listRes).ToNot(BeNil())
			})
			It("should succeed when page size is weird as default will be used", func() {
				listReq.PageSize = -100
				listRes, err := CultureAPI.ListCultures(ctx, listReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(listRes).ToNot(BeNil())
			})
		})

	})

	Describe("Listing cultures for a county ", func() {
		var countyCode string
		Context("Lets create atleast on culture", func() {
			It("should succeed", func() {
				createReq := &culture.CreateCultureRequest{
					Culture: fakeCulture(),
				}
				createRes, err := CultureAPI.CreateCulture(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				countyCode = createReq.Culture.CountyCode
			})
		})

		Describe("Calling list for a county", func() {
			It("should succeed when page size is weird as default will be used", func() {
				listReq.Filter = &culture.ListCultureFilter{
					ListTarget: culture.ListTarget_COUNTY,
					TargetIds:  []string{countyCode},
				}
				listRes, err := CultureAPI.ListCultures(ctx, listReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(listRes).ToNot(BeNil())

				for _, culturePB := range listRes.GetCultures() {
					Expect(culturePB.CountyCode).Should(BeElementOf(listReq.Filter.GetTargetIds()))
				}
			})
		})
	})

	Describe("Listing cultures for a subcounty ", func() {
		var subCountyCode string
		Context("Lets create atleast on culture", func() {
			It("should succeed", func() {
				createReq := &culture.CreateCultureRequest{
					Culture: fakeCulture(),
				}
				createRes, err := CultureAPI.CreateCulture(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				subCountyCode = createReq.Culture.SubCountyCode
			})
		})

		Describe("Calling list for a county", func() {
			It("should succeed when page size is weird as default will be used", func() {
				listReq.Filter = &culture.ListCultureFilter{
					ListTarget: culture.ListTarget_SUB_COUNTY,
					TargetIds:  []string{subCountyCode},
				}
				listRes, err := CultureAPI.ListCultures(ctx, listReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(listRes).ToNot(BeNil())

				for _, culturePB := range listRes.GetCultures() {
					Expect(culturePB.SubCountyCode).Should(BeElementOf(listReq.Filter.GetTargetIds()))
				}
			})
		})
	})

	Describe("Listing cultures for a hospital ", func() {
		var hospitalID string
		Context("Lets create atleast on culture", func() {
			It("should succeed", func() {
				createReq := &culture.CreateCultureRequest{
					Culture: fakeCulture(),
				}
				createRes, err := CultureAPI.CreateCulture(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				hospitalID = createReq.Culture.HospitalId
			})
		})

		Describe("Calling list for a county", func() {
			It("should succeed when page size is weird as default will be used", func() {
				listReq.Filter = &culture.ListCultureFilter{
					ListTarget: culture.ListTarget_HOSPITAL,
					TargetIds:  []string{hospitalID},
				}
				listRes, err := CultureAPI.ListCultures(ctx, listReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(listRes).ToNot(BeNil())

				for _, culturePB := range listRes.GetCultures() {
					Expect(culturePB.HospitalId).Should(BeElementOf(listReq.Filter.GetTargetIds()))
				}
			})
		})
	})

	Describe("Listing cultures for a patient ", func() {
		var patientID string
		Context("Lets create atleast on culture", func() {
			It("should succeed", func() {
				createReq := &culture.CreateCultureRequest{
					Culture: fakeCulture(),
				}
				createRes, err := CultureAPI.CreateCulture(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				patientID = createReq.Culture.PatientId
			})
		})

		Describe("Calling list for a county", func() {
			It("should succeed when page size is weird as default will be used", func() {
				listReq.Filter = &culture.ListCultureFilter{
					ListTarget: culture.ListTarget_PATIENT,
					TargetIds:  []string{patientID},
				}
				listRes, err := CultureAPI.ListCultures(ctx, listReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(listRes).ToNot(BeNil())

				for _, culturePB := range listRes.GetCultures() {
					Expect(culturePB.PatientId).Should(BeElementOf(listReq.Filter.GetTargetIds()))
				}
			})
		})
	})

	Describe("Listing cultures for a lab technician ", func() {
		var labTechID string
		Context("Lets create atleast on culture", func() {
			It("should succeed", func() {
				createReq := &culture.CreateCultureRequest{
					Culture: fakeCulture(),
				}
				createRes, err := CultureAPI.CreateCulture(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				labTechID = createReq.Culture.LabTechId
			})
		})

		Describe("Calling list for a county", func() {
			It("should succeed when page size is weird as default will be used", func() {
				listReq.Filter = &culture.ListCultureFilter{
					ListTarget: culture.ListTarget_LAB_TECHNICIAN,
					TargetIds:  []string{labTechID},
				}
				listRes, err := CultureAPI.ListCultures(ctx, listReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(listRes).ToNot(BeNil())

				for _, culturePB := range listRes.GetCultures() {
					Expect(culturePB.LabTechId).Should(BeElementOf(listReq.Filter.GetTargetIds()))
				}
			})
		})
	})
})
