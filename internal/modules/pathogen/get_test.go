package pathogen

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting Pathogen #get", func() {
	var (
		getReq *pathogen.GetPathogenRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &pathogen.GetPathogenRequest{
			View:       pathogen.PathogenView_FULL,
			PathogenId: "8f344acb-d6ce-4db0-bda2-c732f2464fc8",
		}
		ctx = context.Background()
	})

	Describe("Getting pathogen with nil request", func() {
		It("should fail when request is nil", func() {
			getReq = nil
			getRes, err := PathogenAPI.GetPathogen(context.Background(), getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
	})

	When("Getting pathogen with missing/incorrect pathogen id", func() {
		It("should fail when pathogen id is missing", func() {
			getReq.PathogenId = ""
			getRes, err := PathogenAPI.GetPathogen(context.Background(), getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(getRes).To(BeNil())
		})
		It("should fail when pathogen id doesn't exist", func() {
			getReq.PathogenId = "bad-val"
			getRes, err := PathogenAPI.GetPathogen(context.Background(), getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(getRes).Should(BeNil())
		})
	})

	When("Getting pathogen with valid request", func() {
		var pathogenID string
		Context("Lets create pathogen first", func() {
			It("should succeed", func() {
				createReq := &pathogen.CreatePathogenRequest{
					Pathogen: newPathogen(),
				}
				createRes, err := PathogenAPI.CreatePathogen(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				pathogenID = createRes.PathogenId
			})
		})

		It("should get the pathogen", func() {
			getReq.PathogenId = pathogenID
			getRes, err := PathogenAPI.GetPathogen(ctx, getReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.OK))
			Expect(getRes).ToNot(BeNil())
		})
	})
})
