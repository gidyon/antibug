package pathogen

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Deleting Pathogen #delete", func() {
	var (
		delReq *pathogen.DeletePathogenRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &pathogen.DeletePathogenRequest{
			PathogenId: "8669bfd8-768d-4b7c-bb90-b8540036cb1b",
		}
		ctx = context.Background()
	})

	Describe("Deleting pathogen with nil request", func() {
		It("should fail when request is nil", func() {
			delReq = nil
			delRes, err := PathogenAPI.DeletePathogen(context.Background(), delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
	})

	When("Deleting pathogen with missing/incorrect pathogen id", func() {
		It("should fail when pathogen id is missing", func() {
			delReq.PathogenId = ""
			delRes, err := PathogenAPI.DeletePathogen(context.Background(), delReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(delRes).To(BeNil())
		})
	})

	When("Deleting pathogen with correct pathogen id", func() {
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

		When("Deleting the pathogen name", func() {
			It("should succeed", func() {
				delReq.PathogenId = pathogenID
				updateRes, err := PathogenAPI.DeletePathogen(ctx, delReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		When("Getting the pathogen", func() {
			It("should fail because its deleted", func() {
				getReq := &pathogen.GetPathogenRequest{
					View:       pathogen.PathogenView_FULL,
					PathogenId: pathogenID,
				}
				getRes, err := PathogenAPI.GetPathogen(ctx, getReq)
				Expect(err).To(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.NotFound))
				Expect(getRes).To(BeNil())
			})
		})
	})
})
