package antibiogram

import (
	"context"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting pathogens antibiogram #listpathogens", func() {
	var (
		filter *antibiogram.Filter
		ctx    context.Context
	)

	BeforeEach(func() {
		filter = fakeFilter(subjectPathogen)
		ctx = context.Background()
	})

	Describe("Getting pathogens antibiogram with malformed request", func() {
		It("should fail when the request is nil", func() {
			filter = nil
			pathogensAntibiogram, err := AntibiogramAPI.GenPathogensAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(pathogensAntibiogram).Should(BeNil())
		})
		It("should fail input values is missing", func() {
			filter.InputValues = nil
			pathogensAntibiogram, err := AntibiogramAPI.GenPathogensAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(pathogensAntibiogram).Should(BeNil())
		})
		It("should fail scope values is missing and region scope is not country", func() {
			filter.ScopeValues = nil
			filter.RegionScope = antibiogram.RegionScope_COUNTY
			pathogensAntibiogram, err := AntibiogramAPI.GenPathogensAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(pathogensAntibiogram).Should(BeNil())
		})
	})

	Describe("Getting antibiogram with well-formed request", func() {
		It("succeed because filter is well formed", func() {
			pathogensAntibiogram, err := AntibiogramAPI.GenPathogensAntibiogram(ctx, filter)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(pathogensAntibiogram).ShouldNot(BeNil())
			Expect(len(pathogensAntibiogram.Antibiograms)).ShouldNot(BeZero())

			// Loop over antibiograms
			for _, antibiogramData := range pathogensAntibiogram.Antibiograms {
				Expect(len(antibiogramData.Susceptibilities)).ShouldNot(BeZero())
			}
		})
	})
})
