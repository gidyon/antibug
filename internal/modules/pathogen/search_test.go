package pathogen

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Search Pathogens #search", func() {
	var (
		searchReq *pathogen.SearchPathogensRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		searchReq = &pathogen.SearchPathogensRequest{
			Query:     "Burnpaper",
			PageToken: 0,
		}
		ctx = context.Background()
	})

	Describe("Searching pathogens with nil request", func() {
		It("should fail when request is nil", func() {
			searchReq = nil
			searchRes, err := PathogenAPI.SearchPathogens(context.Background(), searchReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(searchRes).To(BeNil())
		})
	})

	When("Searching pathogens with weird request payload", func() {
		It("should succeed when page token is weird", func() {
			searchReq.PageToken = int32(-45)
			searchRes, err := PathogenAPI.SearchPathogens(context.Background(), searchReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(searchRes).ToNot(BeNil())
		})
	})

	When("Searching pathogens with valid request", func() {
		var pathogenName string
		Context("Lets create atleast one antmicrobial", func() {
			It("should succeed", func() {
				createReq := &pathogen.CreatePathogenRequest{
					Pathogen: newPathogen(),
				}
				createRes, err := PathogenAPI.CreatePathogen(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				pathogenName = createReq.Pathogen.PathogenName
			})
		})

		Context("Lets update search query", func() {
			BeforeEach(func() {
				searchReq.Query = pathogenName
			})

			It("should search pathogens", func() {
				searchRes, err := PathogenAPI.SearchPathogens(ctx, searchReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(searchRes).ToNot(BeNil())
				Expect(len(searchRes.Pathogens)).ShouldNot(BeZero())
			})

			It("should search pathogens even when page token is large", func() {
				searchReq.PageToken = 300
				searchRes, err := PathogenAPI.SearchPathogens(ctx, searchReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(searchRes).ToNot(BeNil())
				Expect(len(searchRes.Pathogens)).Should(BeZero())
			})
		})
	})
})
