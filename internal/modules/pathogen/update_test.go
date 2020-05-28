package pathogen

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/pathogen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Updating Pathogen #update", func() {
	var (
		updateReq *pathogen.UpdatePathogenRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &pathogen.UpdatePathogenRequest{
			PathogenId: "8669bfd8-768d-4b7c-bb90-b8540036cb1b",
			Pathogen:   newPathogen(),
		}
		ctx = context.Background()
	})

	Describe("Updating pathogen with nil request", func() {
		It("should fail when request is nil", func() {
			updateReq = nil
			updateRes, err := PathogenAPI.UpdatePathogen(context.Background(), updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	When("Updating pathogen with incorrect/missing pathogen values", func() {
		It("should fail when pathogen is nil", func() {
			updateReq.Pathogen = nil
			updateRes, err := PathogenAPI.UpdatePathogen(context.Background(), updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	When("Updating pathogen with valid request", func() {
		var (
			pathogenID   string
			pathogenName string
		)

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
				pathogenName = createReq.Pathogen.PathogenName
			})
		})

		When("Updating the pathogen name", func() {
			It("should succeed", func() {
				updateReq.Pathogen.PathogenName = randomdata.SillyName()
				updateReq.PathogenId = pathogenID
				updateRes, err := PathogenAPI.UpdatePathogen(context.Background(), updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		When("Getting the pathogen", func() {
			It("should get the updated pathogen", func() {
				getReq := &pathogen.GetPathogenRequest{
					View:       pathogen.PathogenView_FULL,
					PathogenId: pathogenID,
				}
				getRes, err := PathogenAPI.GetPathogen(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				// Expect the name of the pathogen to be different than original created name
				Expect(pathogenName).ShouldNot(Equal(getRes.PathogenName))
			})
		})
	})
})
