package antimicrobial

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/antimicrobial"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Updating Antimicrobial #update", func() {
	var (
		updateReq *antimicrobial.UpdateAntimicrobialRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &antimicrobial.UpdateAntimicrobialRequest{
			AntimicrobialId: "8669bfd8-768d-4b7c-bb90-b8540036cb1b",
			Antimicrobial:   newAntimicrobial(),
		}
		ctx = context.Background()
	})

	Describe("Updating antimicrobial with nil request", func() {
		It("should fail when request is nil", func() {
			updateReq = nil
			updateRes, err := AntimicrobialAPI.UpdateAntimicrobial(context.Background(), updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	When("Updating antimicrobial with incorrect/missing antimicrobial values", func() {
		It("should fail when antimicrobial is nil", func() {
			updateReq.Antimicrobial = nil
			updateRes, err := AntimicrobialAPI.UpdateAntimicrobial(context.Background(), updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	When("Updating antimicrobial with valid request", func() {
		var (
			antimicrobialID   string
			antimicrobialName string
		)

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
				antimicrobialName = createReq.Antimicrobial.AntimicrobialName
			})
		})

		When("Updating the antimicrobial name", func() {
			It("should succeed", func() {
				updateReq.Antimicrobial.AntimicrobialName = randomdata.SillyName()
				updateReq.AntimicrobialId = antimicrobialID
				updateRes, err := AntimicrobialAPI.UpdateAntimicrobial(context.Background(), updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		When("Getting the antimicrobial", func() {
			It("should get the updated antimicrobial", func() {
				getReq := &antimicrobial.GetAntimicrobialRequest{
					View:            antimicrobial.AntimicrobialView_FULL,
					AntimicrobialId: antimicrobialID,
				}
				getRes, err := AntimicrobialAPI.GetAntimicrobial(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				// Expect the name of the antimicrobial to be different than original created name
				Expect(antimicrobialName).ShouldNot(Equal(getRes.AntimicrobialName))
			})
		})
	})
})
