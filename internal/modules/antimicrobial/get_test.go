package antimicrobial

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting Antimicrobial #get", func() {
	var (
		getReq *antimicrobial.GetAntimicrobialRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &antimicrobial.GetAntimicrobialRequest{
			View:            antimicrobial.AntimicrobialView_FULL,
			AntimicrobialId: "8f344acb-d6ce-4db0-bda2-c732f2464fc8",
		}
		ctx = context.Background()
	})

	Describe("Getting antimicrobial with nil request", func() {
		It("should fail when request is nil", func() {
			getReq = nil
			getRes, err := AntimicrobialAPI.GetAntimicrobial(context.Background(), getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
	})

	When("Getting antimicrobial with missing/incorrect antimicrobial id", func() {
		It("should fail when antimicrobial id is missing", func() {
			getReq.AntimicrobialId = ""
			getRes, err := AntimicrobialAPI.GetAntimicrobial(context.Background(), getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})

		It("should fail when antimicrobial id doesn't exist", func() {
			getReq.AntimicrobialId = "bad-val"
			getRes, err := AntimicrobialAPI.GetAntimicrobial(context.Background(), getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(getRes).Should(BeNil())
		})
	})

	When("Getting antimicrobial with valid request", func() {
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

		It("should get the antimicrobial", func() {
			getReq.AntimicrobialId = antimicrobialID
			getRes, err := AntimicrobialAPI.GetAntimicrobial(ctx, getReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(getRes).ToNot(BeNil())
		})
	})
})
