package account

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Login a user #login", func() {
	var (
		loginReq *account.LoginRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		loginReq = &account.LoginRequest{
			Username: randomdata.Email(),
			Password: randomdata.RandStringRunes(10),
		}
		ctx = context.Background()
	})

	Describe("Login with malformed request", func() {
		It("should fail when the request is nil", func() {
			loginReq = nil
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(loginRes).To(BeNil())
		})
		It("should fail when username is missing", func() {
			loginReq.Username = ""
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(loginRes).To(BeNil())
		})
		It("should fail when password is missing", func() {
			loginReq.Password = ""
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(loginRes).To(BeNil())
		})
	})

	Describe("Login with wellformed request", func() {
		var (
			userName  string
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
				userName = createReq.Account.Email
				accountID = createRes.AccountId
			})
		})

		Describe("Login to account", func() {
			It("should fail because account is not active", func() {
				loginReq.Username = userName
				loginRes, err := AccountAPI.Login(ctx, loginReq)
				Expect(err).To(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.PermissionDenied))
				Expect(loginRes).To(BeNil())
			})
		})

		Describe("Lets activate the account", func() {
			It("should succeed", func() {
				activateReq := &account.ActivateAccountRequest{
					AccountId: accountID,
				}
				loginRes, err := AccountAPI.ActivateAccount(ctx, activateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(loginRes).ToNot(BeNil())
			})
		})

		Describe("Login to account", func() {
			It("should succeed because account is now active", func() {
				loginReq.Username = userName
				loginReq.Password = "hakty11"
				loginRes, err := AccountAPI.Login(ctx, loginReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(loginRes).ToNot(BeNil())

				// Fields should be populated
				Expect(loginRes.Token).ShouldNot(BeZero())
				Expect(loginRes.AccountId).ShouldNot(BeZero())
				Expect(loginRes.AccountGroup).ShouldNot(BeZero())
			})

			It("should fail if password is incorrect", func() {
				loginReq.Username = userName
				loginRes, err := AccountAPI.Login(ctx, loginReq)
				Expect(err).To(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.Unauthenticated))
				Expect(loginRes).To(BeNil())
			})
		})
	})
})
