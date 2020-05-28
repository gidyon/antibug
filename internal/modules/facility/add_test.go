package facility

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/facility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

var _ = Describe("Adding A Facility #add", func() {
	var (
		addReq *facility.AddFacilityRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		addReq = &facility.AddFacilityRequest{
			Facility: newFacility(),
		}
		ctx = context.Background()
	})

	Describe("Creating facility facility with nil request", func() {
		It("should fail when request is nil", func() {
			addReq = nil
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(addRes).To(BeNil())
		})
	})

	When("Creating facility with some missing fields", func() {
		It("should fail when facility name isn't provided", func() {
			addReq.Facility.FacilityName = ""
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(addRes).To(BeNil())
		})
		It("should fail when county name isn't provided", func() {
			addReq.Facility.County = ""
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(addRes).To(BeNil())
		})
		It("should fail when county code isn't provided", func() {
			addReq.Facility.CountyCode = 0
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(addRes).To(BeNil())
		})
		It("should fail when sub county name isn't provided", func() {
			addReq.Facility.SubCounty = ""
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(addRes).To(BeNil())
		})
		It("should fail when sub county code isn't provided", func() {
			addReq.Facility.SubCountyCode = 0
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(addRes).To(BeNil())
		})
	})

	When("Creating facility with valid request", func() {
		var faciltyName string
		It("should add facility in database without error", func() {
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(addRes).ToNot(BeNil())
			faciltyName = addReq.Facility.FacilityName
		})
		It("should fail if the name is already registered", func() {
			addReq.Facility.FacilityName = faciltyName
			addRes, err := FacilityAPI.AddFacility(ctx, addReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.AlreadyExists))
			Expect(addRes).To(BeNil())
		})
	})
})

// adds a new facility resource
func newFacility() *facility.Facility {
	return &facility.Facility{
		FacilityName:  randomdata.City() + " " + randomdata.RandStringRunes(5) + " facility",
		County:        randomdata.State(randomdata.Small),
		CountyCode:    rand.Int31(),
		SubCounty:     randomdata.ProvinceForCountry("US"),
		SubCountyCode: rand.Int31(),
	}
}
