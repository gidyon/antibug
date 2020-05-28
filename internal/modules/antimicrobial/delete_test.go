package antimicrobial

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Deleting Antimicrobial #delete", func() {
	var (
		delReq *antimicrobial.DeleteAntimicrobialRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &antimicrobial.DeleteAntimicrobialRequest{
			AntimicrobialId: "8669bfd8-768d-4b7c-bb90-b8540036cb1b",
		}
		ctx = context.Background()
	})

	Describe("Deleting antimicrobial with nil request", func() {
		It("should fail when request is nil", func() {
			delReq = nil
			delRes, err := AntimicrobialAPI.DeleteAntimicrobial(context.Background(), delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
	})

	When("Deleting antimicrobial with missing/incorrect antimicrobial id", func() {
		It("should fail when antimicrobial id is missing", func() {
			delReq.AntimicrobialId = ""
			delRes, err := AntimicrobialAPI.DeleteAntimicrobial(context.Background(), delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
	})

	When("Deleting antimicrobial with correct antimicrobial id", func() {
		var antimicrobialID string
		Context("Lets create antimicrobial first", func() {
			It("should succeed", func() {
				createReq := &antimicrobial.CreateAntimicrobialRequest{
					Antimicrobial: newAntimicrobial(),
				}
				createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				antimicrobialID = createRes.AntimicrobialId
			})
		})

		When("Deleting the antimicrobial name", func() {
			It("should succeed", func() {
				delReq.AntimicrobialId = antimicrobialID
				updateRes, err := AntimicrobialAPI.DeleteAntimicrobial(ctx, delReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		When("Getting the antimicrobial", func() {
			It("should fail because its deleted", func() {
				getReq := &antimicrobial.GetAntimicrobialRequest{
					View:            antimicrobial.AntimicrobialView_FULL,
					AntimicrobialId: antimicrobialID,
				}
				getRes, err := AntimicrobialAPI.GetAntimicrobial(ctx, getReq)
				Expect(err).To(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.NotFound))
				Expect(getRes).To(BeNil())
			})
		})
	})
})
