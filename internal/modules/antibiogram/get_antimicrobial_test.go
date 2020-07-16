package antibiogram

import (
	"context"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting antimicrobial antibiogram #antimicrobial", func() {
	var (
		filter *antibiogram.Filter
		ctx    context.Context
	)

	BeforeEach(func() {
		filter = fakeFilter(subjectAntimicrobial)
		ctx = context.Background()
	})

	Describe("Getting antimicrobial antibiogram with malformed request", func() {
		It("should fail when the request is nil", func() {
			filter = nil
			antimicrobialAntibiogram, err := AntibiogramAPI.GenAntimicrobialAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(antimicrobialAntibiogram).Should(BeNil())
		})
		It("should fail input values is missing", func() {
			filter.InputValues = nil
			antimicrobialAntibiogram, err := AntibiogramAPI.GenAntimicrobialAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(antimicrobialAntibiogram).Should(BeNil())
		})
		It("should fail scope values is missing and region scope is not country", func() {
			filter.ScopeValues = nil
			filter.RegionScope = antibiogram.RegionScope_COUNTY
			antimicrobialAntibiogram, err := AntibiogramAPI.GenAntimicrobialAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(antimicrobialAntibiogram).Should(BeNil())
		})
	})

	Describe("Getting antibiogram with well-formed request", func() {
		It("succeed because filter is well formed", func() {
			antimicrobialAntibiogram, err := AntibiogramAPI.GenAntimicrobialAntibiogram(ctx, filter)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(antimicrobialAntibiogram).ShouldNot(BeNil())
			Expect(antimicrobialAntibiogram.AntimicrobialName).Should(Equal(filter.InputValues[0].Name))
			Expect(antimicrobialAntibiogram.AntimicrobialId).Should(Equal(filter.InputValues[0].Id))
			Expect(len(antimicrobialAntibiogram.Susceptibilities)).ShouldNot(BeZero())
		})
	})
})
