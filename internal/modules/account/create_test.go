package account

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Creating an account #create", func() {
	var (
		createReq *account.CreateAccountRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &account.CreateAccountRequest{
			Account:         fakeAccount(),
			Password:        "hakty11",
			ConfirmPassword: "hakty11",
		}
		ctx = context.Background()
	})

	Describe("Creating an account with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when the account is nil", func() {
			createReq.Account = nil
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when the first name is missing", func() {
			createReq.Account.FirstName = ""
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when the last name is missing", func() {
			createReq.Account.LastName = ""
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when the phone is missing", func() {
			createReq.Account.Phone = ""
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when the email is missing", func() {
			createReq.Account.Email = ""
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
		It("should fail when the gender is missing", func() {
			createReq.Account.Gender = ""
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(createRes).To(BeNil())
		})
	})

	Describe("Creating account with a valid request", func() {
		It("should succeed", func() {
			createRes, err := AccountAPI.CreateAccount(ctx, createReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(createRes).ToNot(BeNil())
		})
	})
})

func fakeAccount() *account.Account {
	return &account.Account{
		FirstName:  randomdata.FirstName(randomdata.Male),
		LastName:   randomdata.LastName(),
		Email:      randomdata.Email(),
		Phone:      randomdata.PhoneNumber(),
		ProfileUrl: randomdata.UserAgentString(),
		Gender:     "female",
		Group:      auth.Admin,
	}
}
