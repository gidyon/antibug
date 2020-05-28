package facility

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/facility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("List Facilities #list", func() {
	var (
		listReq *facility.ListFacilitiesRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &facility.ListFacilitiesRequest{
			PageToken: 0,
		}
		ctx = context.Background()
	})

	Describe("Listing facilities with nil request", func() {
		It("should fail when request is nil", func() {
			listReq = nil
			listRes, err := FacilityAPI.ListFacilities(ctx, listReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(listRes).To(BeNil())
		})
	})

	When("Listing facilities with weird request payload", func() {
		It("should succeed even when page token is weird", func() {
			listReq.PageToken = int32(-45)
			listRes, err := FacilityAPI.ListFacilities(context.Background(), listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
		})
	})

	When("Listing facilities with valid request", func() {
		Context("Lets create at least one facility", func() {
			It("should succeed", func() {
				createReq := &facility.AddFacilityRequest{
					Facility: newFacility(),
				}
				createRes, err := FacilityAPI.AddFacility(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
			})
		})

		It("should list facilities", func() {
			listRes, err := FacilityAPI.ListFacilities(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
			Expect(len(listRes.Facilities)).ShouldNot(BeZero())
		})

		It("should list facilities even when page token is large", func() {
			listReq.PageToken = 3000
			listRes, err := FacilityAPI.ListFacilities(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
			Expect(len(listRes.Facilities)).Should(BeZero())
		})
	})
})
