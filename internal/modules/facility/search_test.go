package facility

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/facility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Search Facilities #search", func() {
	var (
		searchReq *facility.SearchFacilitiesRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		searchReq = &facility.SearchFacilitiesRequest{
			Query:     "Burnpaper",
			PageToken: 0,
		}
		ctx = context.Background()
	})

	Describe("Searching facilities with nil request", func() {
		It("should fail when request is nil", func() {
			searchReq = nil
			searchRes, err := FacilityAPI.SearchFacilities(context.Background(), searchReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(searchRes).To(BeNil())
		})
	})

	When("Searching facilities with weird request payload", func() {
		It("should succeed when page token is weird", func() {
			searchReq.PageToken = int32(-45)
			searchRes, err := FacilityAPI.SearchFacilities(context.Background(), searchReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(searchRes).ToNot(BeNil())
		})
		It("should return empty results when query is empty", func() {
			searchReq.Query = ""
			searchRes, err := FacilityAPI.SearchFacilities(context.Background(), searchReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(searchRes).ToNot(BeNil())
			Expect(len(searchRes.Facilities)).Should(BeZero())
		})
	})

	When("Searching facilities with valid request", func() {
		var facilityName string
		Context("Lets add at least one facility", func() {
			It("should succeed", func() {
				addReq := &facility.AddFacilityRequest{
					Facility: newFacility(),
				}
				addRes, err := FacilityAPI.AddFacility(ctx, addReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(addRes).ToNot(BeNil())
				facilityName = addReq.Facility.FacilityName
			})
		})

		Context("Lets update search query", func() {
			BeforeEach(func() {
				searchReq.Query = facilityName
			})

			It("should search facilities", func() {
				searchRes, err := FacilityAPI.SearchFacilities(ctx, searchReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(searchRes).ToNot(BeNil())
				Expect(len(searchRes.Facilities)).ShouldNot(BeZero())
			})

			It("should search facilities even when page token is large", func() {
				searchReq.PageToken = 300
				searchRes, err := FacilityAPI.SearchFacilities(ctx, searchReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(searchRes).ToNot(BeNil())
				Expect(len(searchRes.Facilities)).Should(BeZero())
			})
		})
	})
})
