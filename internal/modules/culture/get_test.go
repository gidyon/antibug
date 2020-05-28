package culture

import (
	"context"
	"fmt"
	"github.com/gidyon/antibug/pkg/api/culture"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

var _ = Describe("Getting Culture Resource #get", func() {
	var (
		getReq *culture.GetCultureRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &culture.GetCultureRequest{
			CultureId: fmt.Sprint(rand.Int()),
		}
		ctx = context.Background()
	})

	When("Get request has missing/incorrect information", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := CultureAPI.GetCulture(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
		It("should fail when culture id is missing", func() {
			getReq.CultureId = ""
			getRes, err := CultureAPI.GetCulture(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
		It("should fail when culture id is incorrect", func() {
			getReq.CultureId = "bad"
			getRes, err := CultureAPI.GetCulture(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
			Expect(getRes).To(BeNil())
		})
	})

	Describe("Get culture call success", func() {
		var (
			cultureID string
		)

		Context("Lets create culture first", func() {
			It("should succeed", func() {
				createReq := &culture.CreateCultureRequest{
					Culture: fakeCulture(),
				}
				createRes, err := CultureAPI.CreateCulture(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				cultureID = createRes.CultureId
			})
		})

		Describe("Getting the culture", func() {
			It("should succeed", func() {
				getReq := &culture.GetCultureRequest{
					CultureId: cultureID,
				}
				getRes, err := CultureAPI.GetCulture(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
			})
		})
	})
})
