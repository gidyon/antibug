package antimicrobial

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("List Antimicrobials #list", func() {
	var (
		listReq *antimicrobial.ListAntimicrobialsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &antimicrobial.ListAntimicrobialsRequest{
			View:      antimicrobial.AntimicrobialView_LIST,
			PageToken: 0,
		}
		ctx = context.Background()
	})

	Describe("Listing antimicrobials with nil request", func() {
		It("should fail when request is nil", func() {
			listReq = nil
			listRes, err := AntimicrobialAPI.ListAntimicrobials(ctx, listReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(listRes).To(BeNil())
		})
	})

	When("Listing antimicrobials with weird request payload", func() {
		It("should succeed when page token is weird", func() {
			listReq.PageToken = int32(-45)
			listRes, err := AntimicrobialAPI.ListAntimicrobials(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
		})
	})

	When("Listing antimicrobials with valid request", func() {
		Context("Lets create at least one antimicrobial", func() {
			It("should succeed", func() {
				createReq := &antimicrobial.CreateAntimicrobialRequest{
					Antimicrobial: newAntimicrobial(),
				}
				createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
			})
		})

		It("should list antimicrobials", func() {
			listRes, err := AntimicrobialAPI.ListAntimicrobials(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
			Expect(len(listRes.Antimicrobials)).ShouldNot(BeZero())
		})

		It("should list antimicrobials even when page token is large", func() {
			listReq.PageToken = 3000
			listRes, err := AntimicrobialAPI.ListAntimicrobials(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
			Expect(len(listRes.Antimicrobials)).Should(BeZero())
		})
	})
})
