package account

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/account"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting user account #get", func() {
	var (
		getReq *account.GetRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &account.GetRequest{
			AccountId: uuid.New().String(),
		}
		ctx = context.Background()
	})

	Describe("Getting user account with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := AccountAPI.GetAccount(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
		It("should fail when the account id is missing", func() {
			getReq.AccountId = ""
			getRes, err := AccountAPI.GetAccount(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
		It("should fail when the account id is incorrect", func() {
			getRes, err := AccountAPI.GetAccount(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
			Expect(getRes).To(BeNil())
		})
	})

	Describe("Getting an account with well-formed request", func() {
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

		Describe("Getting the account", func() {
			It("should get the created account", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetAccount(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())
			})
		})
	})
})
