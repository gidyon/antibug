package account

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Updating user account #update", func() {
	var (
		updateReq *account.UpdateAccountRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &account.UpdateAccountRequest{
			Account:   fakeAccount(),
			AccountId: fmt.Sprint(randomdata.Number(1, 100000)),
		}
		ctx = context.Background()
	})

	Describe("Updating account with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := AccountAPI.UpdateAccount(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when account is nil", func() {
			updateReq.Account = nil
			updateRes, err := AccountAPI.UpdateAccount(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when account id is missing", func() {
			updateReq.AccountId = ""
			updateRes, err := AccountAPI.UpdateAccount(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	Describe("Updating account with well-formed request", func() {
		var (
			accountID string
			firstName string
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
				firstName = createReq.Account.FirstName
			})
		})

		Describe("Lets update the account", func() {
			It("should succeed because the request is good", func() {
				updateReq.AccountId = accountID
				updateRes, err := AccountAPI.UpdateAccount(ctx, updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		Describe("Getting the account", func() {
			It("should get the most up to date account", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetAccount(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				// The account should have changed now
				Expect(getRes.FirstName).ShouldNot(Equal(firstName))
			})
		})
	})
})
