package pathogen

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("List Pathogens #list", func() {
	var (
		listReq *pathogen.ListPathogensRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &pathogen.ListPathogensRequest{
			View:      pathogen.PathogenView_LIST,
			PageToken: 0,
		}
		ctx = context.Background()
	})

	Describe("Listing pathogens with nil request", func() {
		It("should fail when request is nil", func() {
			listReq = nil
			listRes, err := PathogenAPI.ListPathogens(ctx, listReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(listRes).To(BeNil())
		})
	})

	When("Listing pathogens with weird request payload", func() {
		It("should succeed when page token is weird", func() {
			listReq.PageToken = int32(-45)
			listRes, err := PathogenAPI.ListPathogens(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
		})
	})

	When("Listing pathogens with valid request", func() {
		Context("Lets create at least one pathogen", func() {
			It("should succeed", func() {
				createReq := &pathogen.CreatePathogenRequest{
					Pathogen: newPathogen(),
				}
				createRes, err := PathogenAPI.CreatePathogen(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
			})
		})

		It("should list pathogens", func() {
			listRes, err := PathogenAPI.ListPathogens(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
			Expect(len(listRes.Pathogens)).ShouldNot(BeZero())
		})

		It("should list pathogens even when page token is large", func() {
			listReq.PageToken = 3000
			listRes, err := PathogenAPI.ListPathogens(ctx, listReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(listRes).ToNot(BeNil())
			Expect(len(listRes.Pathogens)).Should(BeZero())
		})
	})
})
