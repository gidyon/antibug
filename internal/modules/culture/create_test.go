package culture

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/culture"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Creating culture resource #create", func() {
	var (
		createReq *culture.CreateCultureRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &culture.CreateCultureRequest{
			Culture: FakeCulture(),
		}
		ctx = context.Background()
	})

	Describe("Creating a culture resource with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when lab tech id is missing", func() {
			createReq.Culture.LabTechId = ""
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when hospital id is missing", func() {
			createReq.Culture.HospitalId = ""
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when county code is missing", func() {
			createReq.Culture.CountyCode = ""
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when sub-county code is missing", func() {
			createReq.Culture.SubCountyCode = ""
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when patient age is missing", func() {
			createReq.Culture.PatientAge = -1
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when patient gender is missing", func() {
			createReq.Culture.PatientGender = ""
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when patient id is missing", func() {
			createReq.Culture.PatientId = ""
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when culture source is missing", func() {
			createReq.Culture.CultureSource = ""
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when culture createResults is nil or zero", func() {
			createReq.Culture.CultureResults = nil
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when pathogens found is nil or zero", func() {
			createReq.Culture.PathogensFound = nil
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when antimicrobials used is nil or zero", func() {
			createReq.Culture.AntimicrobialsUsed = nil
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
	})

	Describe("Creating culture with well-formed request", func() {
		It("should create culture and save it in database", func() {
			createRes, err := CultureAPI.CreateCulture(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(createRes.CultureId).ShouldNot(BeZero())
			Expect(status.Code(err)).To(Equal(codes.OK))
		})
	})
})
