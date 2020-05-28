package account

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fakeStarredFacility() *account.Facility {
	return &account.Facility{
		FacilityName: randomdata.City() + " " + randomdata.RandStringRunes(5) + " Facility",
		FacilityId:   randomdata.RandStringRunes(20),
	}
}

func fakeStarredFacilities() *account.StarredFacilities {
	return &account.StarredFacilities{
		Facilities: []*account.Facility{
			fakeStarredFacility(), fakeStarredFacility(),
		},
	}
}

var _ = Describe("Updating account starred hospitals #starredhospitals", func() {
	var (
		updateReq *account.UpdateStarredFacilitiesRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &account.UpdateStarredFacilitiesRequest{
			Facilities: fakeStarredFacilities().Facilities,
			AccountId:  fmt.Sprint(randomdata.Number(0, 10000)),
		}
		ctx = context.Background()
	})

	Describe("Updating account starred hospitals with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := AccountAPI.UpdateStarredFacilities(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when starred hospitals is nil", func() {
			updateReq.Facilities = nil
			updateRes, err := AccountAPI.UpdateStarredFacilities(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when id is missing", func() {
			updateReq.AccountId = ""
			updateRes, err := AccountAPI.UpdateStarredFacilities(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	Describe("Updating account starred hospitals with well-formed request", func() {
		var (
			accountID         string
			starredFacilities []*account.Facility
		)

		Context("Lets create an account first", func() {
			It("should succeed in creating the account", func() {
				createReq := &account.CreateAccountRequest{
					Account: fakeAccount(),
				}
				createRes, err := AccountAPI.CreateAccount(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				accountID = createRes.AccountId
			})
		})

		Describe("Lets update the account starred hospitals", func() {
			It("should succeed because the request is good", func() {
				updateReq.AccountId = accountID
				updateRes, err := AccountAPI.UpdateStarredFacilities(ctx, updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
				starredFacilities = updateReq.Facilities
			})
		})

		Describe("Getting the starred hospitals", func() {
			It("should get the updated starred hospitals", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetStarredFacilities(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				Expect(getRes.Facilities).ShouldNot(BeIdenticalTo(starredFacilities))
			})
		})
	})
})

var _ = Describe("Getting starred hospitals #starredhospitals", func() {
	var (
		getReq *account.GetRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &account.GetRequest{
			AccountId: fmt.Sprint(randomdata.Number(10, 1000)),
		}
		ctx = context.Background()
	})

	Describe("Getting starred hospitals with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			updateRes, err := AccountAPI.GetStarredFacilities(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when id is missing", func() {
			getReq.AccountId = ""
			updateRes, err := AccountAPI.GetStarredFacilities(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	Describe("Getting starred hospitals with well-formed request", func() {
		var accountID string

		Context("Lets create an account first", func() {
			It("should succeed in creating the account", func() {
				createReq := &account.CreateAccountRequest{
					Account: fakeAccount(),
				}
				createRes, err := AccountAPI.CreateAccount(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				accountID = createRes.AccountId
			})
		})

		Describe("Lets update the account starred hospitals", func() {
			It("should succeed because the request is good", func() {
				updateReq := &account.UpdateStarredFacilitiesRequest{
					Facilities: fakeStarredFacilities().Facilities,
					AccountId:  accountID,
				}
				updateReq.AccountId = accountID
				updateRes, err := AccountAPI.UpdateStarredFacilities(ctx, updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		Describe("Getting the starred hospitals", func() {
			It("should get the updated starred hospitals", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetStarredFacilities(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())
			})
		})
	})
})
