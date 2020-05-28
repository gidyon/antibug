package facility

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/facility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Gettin A Facility #get", func() {
	var (
		getReq *facility.GetFacilityRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &facility.GetFacilityRequest{
			FacilityId: "4af01c44-a7af-4517-b65e-94d5a88d4c34",
		}
		ctx = context.Background()
	})

	Describe("Retrieving facility with nil request", func() {
		It("should fail when request is nil", func() {
			getReq = nil
			getRes, err := FacilityAPI.GetFacility(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
	})

	When("Retrieving facility with missing/incorrect facility id", func() {
		It("should fail when facility id is missing", func() {
			getReq.FacilityId = ""
			getRes, err := FacilityAPI.GetFacility(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
		It("should fail when facility id doesn't exist", func() {
			getReq.FacilityId = "djwfb3fninfioef"
			getRes, err := FacilityAPI.GetFacility(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
			Expect(getRes).To(BeNil())
		})
	})

	When("Retrieving facility with valid request", func() {
		var facilityID string
		Context("Lets create a facility first", func() {
			It("should succeed", func() {
				createReq := &facility.AddFacilityRequest{
					Facility: newFacility(),
				}
				createRes, err := FacilityAPI.AddFacility(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				facilityID = createRes.FacilityId
			})
		})
		It("should get the facility", func() {
			getReq.FacilityId = facilityID
			getRes, err := FacilityAPI.GetFacility(ctx, getReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(getRes).ToNot(BeNil())
		})
	})
})
