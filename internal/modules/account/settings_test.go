package account

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fakeSettings() *account.Settings {
	return &account.Settings{
		Settings: map[string]bool{
			"notifications": true,
			"account_watch": true,
		},
	}
}

var _ = Describe("Updating account settings #settings", func() {
	var (
		updateReq *account.UpdateSettingsRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &account.UpdateSettingsRequest{
			Settings:  fakeSettings(),
			AccountId: fmt.Sprint(randomdata.Number(0, 10000)),
		}
		ctx = context.Background()
	})

	Describe("Updating account settings with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := AccountAPI.UpdateSettings(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when settings is nil", func() {
			updateReq.Settings = nil
			updateRes, err := AccountAPI.UpdateSettings(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when id is missing", func() {
			updateReq.AccountId = ""
			updateRes, err := AccountAPI.UpdateSettings(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when id is incorrect", func() {
			updateRes, err := AccountAPI.UpdateSettings(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	Describe("Updating account settings with well-formed request", func() {
		var (
			accountID string
			settings  map[string]bool
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

		Describe("Lets update the account settings", func() {
			It("should succeed because the request is good", func() {
				updateReq.AccountId = accountID
				updateRes, err := AccountAPI.UpdateSettings(ctx, updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
				settings = updateReq.Settings.Settings
			})
		})

		Describe("Getting the settings", func() {
			It("should get the updated settings", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetSettings(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				Expect(getRes.Settings).ShouldNot(BeIdenticalTo(settings))
			})
		})
	})
})

var _ = Describe("Getting settings #settings", func() {
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

	Describe("Getting settings with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			updateRes, err := AccountAPI.GetSettings(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when id is missing", func() {
			getReq.AccountId = ""
			updateRes, err := AccountAPI.GetSettings(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	Describe("Getting settings with well-formed request", func() {
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

		Describe("Lets update the account settings", func() {
			It("should succeed because the request is good", func() {
				updateReq := &account.UpdateSettingsRequest{
					Settings:  fakeSettings(),
					AccountId: accountID,
				}
				updateReq.AccountId = accountID
				updateRes, err := AccountAPI.UpdateSettings(ctx, updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		Describe("Getting the settings", func() {
			It("should get the updated settings", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetSettings(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())
			})
		})
	})
})
