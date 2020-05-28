package facility

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/facility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Deleting Facility #remove", func() {
	var (
		delReq *facility.RemoveFacilityRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &facility.RemoveFacilityRequest{
			FacilityId: "1234",
		}
		ctx = context.Background()
	})

	Describe("Deleting facility with nil request", func() {
		It("should fail when request is nil", func() {
			delReq = nil
			delRes, err := FacilityAPI.RemoveFacility(ctx, delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
	})

	When("Deleting facility with missing/incorrect facility id", func() {
		It("should fail when facility id is missing", func() {
			delReq.FacilityId = ""
			delRes, err := FacilityAPI.RemoveFacility(ctx, delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
		It("should delete nothing in database and exit without error", func() {
			delRes, err := FacilityAPI.RemoveFacility(ctx, delReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(delRes).ToNot(BeNil())
		})
	})

	When("Deleting facility with correct facility id", func() {
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
			getReq := &facility.GetFacilityRequest{
				FacilityId: facilityID,
			}
			getRes, err := FacilityAPI.GetFacility(ctx, getReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(getRes).ToNot(BeNil())
		})

		When("We delete the facility, tryig to get it, should returns 404", func() {
			It("should delete facility in database without error", func() {
				delReq.FacilityId = facilityID
				delRes, err := FacilityAPI.RemoveFacility(ctx, delReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(delRes).ToNot(BeNil())
			})

			It("should return not found since the facility is deleted", func() {
				getReq := &facility.GetFacilityRequest{
					FacilityId: facilityID,
				}
				getRes, err := FacilityAPI.GetFacility(ctx, getReq)
				Expect(err).To(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.NotFound))
				Expect(getRes).To(BeNil())
			})
		})
	})
})
