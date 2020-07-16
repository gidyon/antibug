package antibiogram

import (
	"context"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting antimicrobials antibiogram #antimicrobials", func() {
	var (
		filter *antibiogram.Filter
		ctx    context.Context
	)

	BeforeEach(func() {
		filter = fakeFilter(subjectAntimicrobial)
		ctx = context.Background()
	})

	Describe("Getting antimicrobials antibiogram with malformed request", func() {
		It("should fail when the request is nil", func() {
			filter = nil
			antimicrobialsAntibiogram, err := AntibiogramAPI.GenAntimicrobialsAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(antimicrobialsAntibiogram).Should(BeNil())
		})
		It("should fail input values is missing", func() {
			filter.InputValues = nil
			antimicrobialsAntibiogram, err := AntibiogramAPI.GenAntimicrobialsAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(antimicrobialsAntibiogram).Should(BeNil())
		})
		It("should fail scope values is missing and region scope is not country", func() {
			filter.ScopeValues = nil
			filter.RegionScope = antibiogram.RegionScope_COUNTY
			antimicrobialsAntibiogram, err := AntibiogramAPI.GenAntimicrobialsAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(antimicrobialsAntibiogram).Should(BeNil())
		})
	})

	Describe("Getting antibiogram with well-formed request", func() {
		It("succeed because filter is well formed", func() {
			antimicrobialsAntibiogram, err := AntibiogramAPI.GenAntimicrobialsAntibiogram(ctx, filter)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(antimicrobialsAntibiogram).ShouldNot(BeNil())
			Expect(len(antimicrobialsAntibiogram.Antibiograms)).ShouldNot(BeZero())

			// Loop over antibiograms
			for _, antibiogramData := range antimicrobialsAntibiogram.Antibiograms {
				Expect(len(antibiogramData.Susceptibilities)).ShouldNot(BeZero())
			}
		})
	})
})
