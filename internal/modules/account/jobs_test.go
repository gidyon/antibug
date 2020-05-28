package account

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fakeJob() *account.Job {
	return &account.Job{
		FacilityName: randomdata.City() + " " + randomdata.RandStringRunes(5) + " Hospital",
		FacilityId:   randomdata.RandStringRunes(20),
		Role:         "Clinician and Lab Technician",
		JobId:        randomdata.RandStringRunes(10),
		Description:  randomdata.Paragraph(),
	}
}

func fakeJobs() *account.Jobs {
	return &account.Jobs{
		Jobs: []*account.Job{
			fakeJob(), fakeJob(),
		},
	}
}

var _ = Describe("Updating account jobs #jobs", func() {
	var (
		updateReq *account.UpdateJobsRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &account.UpdateJobsRequest{
			Jobs:      fakeJobs().Jobs,
			AccountId: fmt.Sprint(randomdata.Number(0, 10000)),
		}
		ctx = context.Background()
	})

	Describe("Updating account jobs with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := AccountAPI.UpdateJobs(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when jobs is nil", func() {
			updateReq.Jobs = nil
			updateRes, err := AccountAPI.UpdateJobs(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when id is missing", func() {
			updateReq.AccountId = ""
			updateRes, err := AccountAPI.UpdateJobs(ctx, updateReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	Describe("Updating account jobs with well-formed request", func() {
		var (
			accountID string
			jobs      []*account.Job
		)

		Context("Lets create an account first", func() {
			It("should succeed in creating the account", func() {
				createReq := &account.CreateAccountRequest{
					Account: fakeAccount(),
				}
				createRes, err := AccountAPI.CreateAccount(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				accountID = createRes.AccountId
			})
		})

		Describe("Lets update the account jobs", func() {
			It("should succeed because the request is good", func() {
				updateReq.AccountId = accountID
				updateRes, err := AccountAPI.UpdateJobs(ctx, updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
				jobs = updateReq.Jobs
			})
		})

		Describe("Getting the jobs", func() {
			It("should get the updated jobs", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetJobs(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())

				Expect(getRes.Jobs).ShouldNot(BeIdenticalTo(jobs))
			})
		})
	})
})

var _ = Describe("Getting jobs #jobs", func() {
	var (
		getReq *account.GetRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &account.GetRequest{
			AccountId: fmt.Sprint(randomdata.Number(10, 1000)),
		}
		ctx = context.Background()
	})

	Describe("Getting jobs with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			updateRes, err := AccountAPI.GetJobs(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
		It("should fail when id is missing", func() {
			getReq.AccountId = ""
			updateRes, err := AccountAPI.GetJobs(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
			Expect(updateRes).To(BeNil())
		})
	})

	Describe("Getting jobs with well-formed request", func() {
		var accountID string

		Context("Lets create an account first", func() {
			It("should succeed in creating the account", func() {
				createReq := &account.CreateAccountRequest{
					Account: fakeAccount(),
				}
				createRes, err := AccountAPI.CreateAccount(ctx, createReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(createRes).ToNot(BeNil())
				accountID = createRes.AccountId
			})
		})

		Describe("Lets update the account jobs", func() {
			It("should succeed because the request is good", func() {
				updateReq := &account.UpdateJobsRequest{
					Jobs:      fakeJobs().Jobs,
					AccountId: accountID,
				}
				updateReq.AccountId = accountID
				updateRes, err := AccountAPI.UpdateJobs(ctx, updateReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(updateRes).ToNot(BeNil())
			})
		})

		Describe("Getting the jobs", func() {
			It("should get the updated jobs", func() {
				getReq := &account.GetRequest{
					AccountId: accountID,
				}
				getRes, err := AccountAPI.GetJobs(ctx, getReq)
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Code(err)).To(Equal(codes.OK))
				Expect(getRes).ToNot(BeNil())
			})
		})
	})
})
