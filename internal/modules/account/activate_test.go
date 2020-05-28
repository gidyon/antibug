package account

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Activate account for a user #activate", func() {
	var (
		activateReq *account.ActivateAccountRequest
		ctx         context.Context
	)

	BeforeEach(func() {
		activateReq = &account.ActivateAccountRequest{
			AccountId: randomdata.RandStringRunes(10),
		}
		ctx = context.Background()
	})

	Describe("ActivateAccount with malformed request", func() {
		It("should fail when the request is nil", func() {
			activateReq = nil
			activateRes, err := AccountAPI.ActivateAccount(ctx, activateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(activateRes).To(BeNil())
		})
		It("should fail when account id is missing", func() {
			activateReq.AccountId = ""
			activateRes, err := AccountAPI.ActivateAccount(ctx, activateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(activateRes).To(BeNil())
		})
	})

	Describe("ActivateAccount with wellformed request", func() {
		var (
			accountID string
		)
		Context("Lets create an account first", func() {
			It("should succeed in creating the account", func() {
				createReq := &account.CreateAccountRequest{
					Account:         fakeAccount(),
					Password:        "hakty11",
					ConfirmPassword: "hakty11",
				}
				createRes, err := AccountAPI.CreateAccount(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				accountID = createRes.AccountId
			})
		})

		Describe("Lets get the account", func() {
			It("should be inactive because account is not active", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetAccount(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				// Account active is false
				Expect(getRes.Active).Should(BeFalse())
			})
		})

		Describe("Lets activate the account", func() {
			It("should succeed", func() {
				activateReq := &account.ActivateAccountRequest{
					AccountId: accountID,
				}
				activateRes, err := AccountAPI.ActivateAccount(ctx, activateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(activateRes).ToNot(BeNil())
			})
		})

		Describe("Lets get the account", func() {
			It("should be active because we've made it active", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetAccount(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				// Account active is false
				Expect(getRes.Active).Should(BeTrue())
			})
		})
	})
})
