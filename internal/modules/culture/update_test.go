package culture

import (
	"context"
	"fmt"
	"github.com/gidyon/antibug/pkg/api/culture"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

var _ = Describe("Updating Culture Resource #update", func() {
	var (
		updateReq *culture.UpdateCultureRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &culture.UpdateCultureRequest{
			EditorId:  fmt.Sprint(rand.Int()),
			CultureId: fmt.Sprint(rand.Int()),
			Culture:   fakeCulture(),
		}
		ctx = context.Background()
	})

	When("Updating culture with malformed request", func() {
		It("should fail when request is nil", func() {
			updateReq.Culture = nil
			updateRes, err := CultureAPI.UpdateCulture(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when culture is nil", func() {
			updateReq.Culture = nil
			updateRes, err := CultureAPI.UpdateCulture(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when editor id is missing", func() {
			updateReq.EditorId = ""
			updateRes, err := CultureAPI.UpdateCulture(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when culture id is missing", func() {
			updateReq.CultureId = ""
			updateRes, err := CultureAPI.UpdateCulture(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when culture id is incorrect", func() {
			updateReq.CultureId = "bad"
			updateRes, err := CultureAPI.UpdateCulture(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
			Expect(updateRes).To(BeNil())
		})
	})

	When("Updating culture with well-formed request", func() {
		var (
			cultureID string
			source    string
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
				source = createReq.Culture.CultureSource
			})
		})

		Describe("Updating the culture", func() {
			It("should update the culture in database", func() {
				updateReq.CultureId = cultureID
				updateRes, err := CultureAPI.UpdateCulture(ctx, updateReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(updateRes).ShouldNot(BeNil())
			})
		})

		Describe("Getting the updated culture", func() {
			It("should succeed and updates shown", func() {
				getReq := &culture.GetCultureRequest{
					CultureId: cultureID,
				}
				getRes, err := CultureAPI.GetCulture(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())

				// Expect the source to have changed
				Expect(getRes.CultureSource).ShouldNot(Equal(source))
			})
		})
	})
})
