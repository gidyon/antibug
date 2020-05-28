package culture

import (
	"context"
	"fmt"
	"github.com/gidyon/antibug/pkg/api/culture"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

var _ = Describe("Deleting culture resource #delete", func() {
	var (
		delReq *culture.DeleteCultureRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &culture.DeleteCultureRequest{
			CultureId: fmt.Sprint(rand.Int()),
		}
		ctx = context.Background()
	})

	When("Deleting culture resource with malformed request", func() {
		It("should fail when the request is nil", func() {
			delReq = nil
			delRes, err := CultureAPI.DeleteCulture(ctx, delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
		It("should fail when culture id is missing", func() {
			delReq.CultureId = ""
			delRes, err := CultureAPI.DeleteCulture(ctx, delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
	})

	Describe("Deleting culture with well-formed request", func() {
		var cultureID string

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

		Describe("Deleting the culture", func() {
			It("should delete the culture in database", func() {
				delReq.CultureId = cultureID
				updateRes, err := CultureAPI.DeleteCulture(ctx, delReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(updateRes).ShouldNot(BeNil())
			})
		})

		Describe("Getting the updated culture", func() {
			It("should fail because the culture is deleted", func() {
				getReq := &culture.GetCultureRequest{
					CultureId: cultureID,
				}
				getRes, err := CultureAPI.GetCulture(ctx, getReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.NotFound))
				Expect(getRes).Should(BeNil())
			})
		})
	})
})
