package antibiogram

import (
	"context"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting pathogen antibiogram #getpathogen", func() {
	var (
		filter *antibiogram.Filter
		ctx    context.Context
	)

	BeforeEach(func() {
		filter = fakeFilter(subjectPathogen)
		ctx = context.Background()
	})

	Describe("Getting pathogen antibiogram with malformed request", func() {
		It("should fail when the request is nil", func() {
			filter = nil
			pathogenAntibiogram, err := AntibiogramAPI.GenPathogenAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(pathogenAntibiogram).Should(BeNil())
		})
		It("should fail input values is missing", func() {
			filter.InputValues = nil
			pathogenAntibiogram, err := AntibiogramAPI.GenPathogenAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(pathogenAntibiogram).Should(BeNil())
		})
		It("should fail scope values is missing and region scope is not country", func() {
			filter.ScopeValues = nil
			filter.RegionScope = antibiogram.RegionScope_COUNTY
			pathogenAntibiogram, err := AntibiogramAPI.GenPathogenAntibiogram(ctx, filter)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(pathogenAntibiogram).Should(BeNil())
		})
	})

	Describe("Getting antibiogram with well-formed request", func() {
		It("succeed because filter is well formed", func() {
			pathogenAntibiogram, err := AntibiogramAPI.GenPathogenAntibiogram(ctx, filter)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(pathogenAntibiogram).ShouldNot(BeNil())
			Expect(pathogenAntibiogram.PathogenName).Should(Equal(filter.InputValues[0].Name))
			Expect(pathogenAntibiogram.PathogenId).Should(Equal(filter.InputValues[0].Id))
			Expect(len(pathogenAntibiogram.Susceptibilities)).ShouldNot(BeZero())
		})
	})
})
