package antimicrobial

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Search Antimicrobials #search", func() {
	var (
		searchReq *antimicrobial.SearchAntimicrobialsRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		searchReq = &antimicrobial.SearchAntimicrobialsRequest{
			Query:     "Burnpaper",
			PageToken: 0,
		}
		ctx = context.Background()
	})

	Describe("Searching antimicrobials with nil request", func() {
		It("should fail when request is nil", func() {
			searchReq = nil
			searchRes, err := AntimicrobialAPI.SearchAntimicrobials(context.Background(), searchReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(searchRes).To(BeNil())
		})
	})

	When("Searching antimicrobials with weird request payload", func() {
		It("should succeed when page token is weird", func() {
			searchReq.PageToken = int32(-45)
			searchRes, err := AntimicrobialAPI.SearchAntimicrobials(context.Background(), searchReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(searchRes).ToNot(BeNil())
		})
	})

	When("Searching antimicrobials with valid request", func() {
		var antimicrobialName string
		Context("Lets create atleast one antmicrobial", func() {
			It("should succeed", func() {
				createReq := &antimicrobial.CreateAntimicrobialRequest{
					Antimicrobial: newAntimicrobial(),
				}
				createRes, err := AntimicrobialAPI.CreateAntimicrobial(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				antimicrobialName = createReq.Antimicrobial.AntimicrobialName
			})
		})

		Context("Lets update search query", func() {
			BeforeEach(func() {
				searchReq.Query = antimicrobialName
			})

			It("should search antimicrobials", func() {
				searchRes, err := AntimicrobialAPI.SearchAntimicrobials(ctx, searchReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(searchRes).ToNot(BeNil())
				Expect(len(searchRes.Antimicrobials)).ShouldNot(BeZero())
			})

			It("should search antimicrobials even when page token is large", func() {
				searchReq.PageToken = 300
				searchRes, err := AntimicrobialAPI.SearchAntimicrobials(ctx, searchReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(searchRes).ToNot(BeNil())
				Expect(len(searchRes.Antimicrobials)).Should(BeZero())
			})
		})

		Context("Searching antimicrobial with missing query", func() {
			It("should succeed but return 0 results", func() {
				searchReq.Query = ""
				searchRes, err := AntimicrobialAPI.SearchAntimicrobials(ctx, searchReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(searchRes).ToNot(BeNil())
				Expect(len(searchRes.Antimicrobials)).Should(BeZero())
			})
		})
	})
})
