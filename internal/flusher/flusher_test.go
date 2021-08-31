package flusher_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozoncp/ocp-team-api/internal/flusher"
	"github.com/ozoncp/ocp-team-api/internal/mocks"
	"github.com/ozoncp/ocp-team-api/internal/models"
)

var _ = Describe("Flusher", func() {

	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockRepo
		f        flusher.Flusher
		teams    []models.Team
	)

	mockError := errors.New("error")

	emptyTeams := make([]models.Team, 0)
	nonEmptyTeams := []models.Team{
		{1, "Team1", "Desc1", false},
		{2, "Team2", "Desc2", false},
		{3, "Team3", "Desc3", false},
		{4, "Team4", "Desc4", false},
		{5, "Team5", "Desc5", false},
	}

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockRepo = mocks.NewMockRepo(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("invalid flusher ", func() {
		BeforeEach(func() {
			mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, nil).Times(0)

			teams = nonEmptyTeams
		})

		It("chunk size is negative", func() {
			f = flusher.NewFlusher(-1, mockRepo)
		})

		It("chunk size is equal to 0", func() {
			f = flusher.NewFlusher(0, mockRepo)
		})

		AfterEach(func() {
			gomega.Expect(f.Flush(context.TODO(), teams)).Should(gomega.BeEmpty())
		})
	})

	Context("valid flusher", func() {
		JustBeforeEach(func() {
			f = flusher.NewFlusher(2, mockRepo)
		})

		Context("when there are no teams to be flushed", func() {
			BeforeEach(func() {
				teams = emptyTeams
			})

			It("returns empty slice", func() {
				mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, nil).Times(0)

				gomega.Expect(f.Flush(context.TODO(), teams)).Should(gomega.BeEmpty())
			})
		})

		Context("when all teams successfully flushed", func() {
			BeforeEach(func() {
				teams = nonEmptyTeams
			})

			It("returns empty slice", func() {
				mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, nil).Times(3)

				gomega.Expect(f.Flush(context.TODO(), teams)).Should(gomega.BeEmpty())
			})
		})

		Context("when teams failed to be flushed", func() {
			BeforeEach(func() {
				teams = nonEmptyTeams
			})

			It("cannot flush all teams", func() {
				mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, mockError).Times(3)

				gomega.Expect(f.Flush(context.TODO(), teams)).Should(gomega.Equal(teams))
			})

			It("cannot flush last 3 teams", func() {
				mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, mockError).Times(2)

				gomega.Expect(f.Flush(context.TODO(), teams)).Should(gomega.Equal(teams[2:]))
			})

			It("cannot flush first 2 teams", func() {
				mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, mockError).Times(1)
				mockRepo.EXPECT().CreateTeams(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)

				gomega.Expect(f.Flush(context.TODO(), teams)).Should(gomega.Equal(teams[:2]))
			})
		})
	})
})
