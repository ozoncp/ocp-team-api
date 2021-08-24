package api_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-team-api/internal/api"
	"github.com/ozoncp/ocp-team-api/internal/mocks"
	"github.com/ozoncp/ocp-team-api/internal/models"
	desc "github.com/ozoncp/ocp-team-api/pkg/ocp-team-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Api", func() {

	var (
		ctrl *gomock.Controller

		s        desc.OcpTeamApiServer
		mockRepo *mocks.MockRepo
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockRepo = mocks.NewMockRepo(ctrl)
		s = api.NewOcpTeamApi(mockRepo)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CreateTeamV1()", func() {
		It("returns response", func() {
			expectedResponse := &desc.CreateTeamV1Response{Id: uint64(1)}

			mockRepo.EXPECT().AddTeam(gomock.Any(), gomock.Any()).Return(uint64(1), nil)

			req := &desc.CreateTeamV1Request{Name: "Name", Description: "Description"}

			actualResponse, err := s.CreateTeamV1(context.Background(), req)
			Expect(err).Should(BeNil())
			Expect(actualResponse.Id).Should(Equal(expectedResponse.Id))
		})
	})

	Context("GetTeamV1()", func() {
		It("get existing team by id", func() {
			expectedResponse := &desc.GetTeamV1Response{Team: &desc.Team{
				Id:          uint64(1),
				Name:        "Name",
				Description: "Description",
			}}

			mockRepo.EXPECT().GetTeam(gomock.Any(), gomock.Any()).Return(
				&models.Team{
					Id:          uint64(1),
					Name:        "Name",
					Description: "Description",
				}, nil)

			req := &desc.GetTeamV1Request{Id: uint64(1)}

			actualResponse, err := s.GetTeamV1(context.Background(), req)
			Expect(err).Should(BeNil())
			Expect(actualResponse.Team.Id).Should(Equal(expectedResponse.Team.Id))
			Expect(actualResponse.Team.Name).Should(Equal(expectedResponse.Team.Name))
			Expect(actualResponse.Team.Description).Should(Equal(expectedResponse.Team.Description))
		})

		It("get non-existing team by id", func() {
			mockRepo.EXPECT().GetTeam(gomock.Any(), gomock.Any()).Return(
				nil, status.Error(codes.NotFound, ""))

			req := &desc.GetTeamV1Request{Id: uint64(1)}

			actualResponse, err := s.GetTeamV1(context.Background(), req)
			Expect(actualResponse).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("RemoveTeamV1()", func() {
		It("removes existing element", func() {
			mockRepo.EXPECT().RemoveTeam(gomock.Any(), gomock.Any()).Return(nil)

			req := &desc.RemoveTeamV1Request{Id: uint64(1)}
			expectedResponse := &desc.RemoveTeamV1Response{}

			actualResponse, err := s.RemoveTeamV1(context.Background(), req)
			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})
	})

	Context("ListTeamsV1()", func() {
		It("return nothing when limit and offset are default", func() {
			mockRepo.EXPECT().ListTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Team{}, nil)

			req := &desc.ListTeamsV1Request{}
			expectedResponse := &desc.ListTeamsV1Response{Teams: []*desc.Team{}}

			actualResponse, err := s.ListTeamsV1(context.Background(), req)

			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})
	})

	Context("ListTeamsV1()", func() {
		It("return nothing when limit and offset are default", func() {
			mockRepo.EXPECT().ListTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Team{}, nil)

			req := &desc.ListTeamsV1Request{}
			expectedResponse := &desc.ListTeamsV1Response{Teams: []*desc.Team{}}

			actualResponse, err := s.ListTeamsV1(context.Background(), req)

			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})

		It("return teams when limit and offset are set", func() {
			mockRepo.EXPECT().ListTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				[]models.Team{
					{Id: uint64(1), Name: "Name", Description: "Description"},
					{Id: uint64(2), Name: "Name", Description: "Description"},
				}, nil)

			req := &desc.ListTeamsV1Request{Limit: 2, Offset: 2}
			expectedResponse := &desc.ListTeamsV1Response{Teams: []*desc.Team{
				{Id: uint64(1), Name: "Name", Description: "Description"},
				{Id: uint64(2), Name: "Name", Description: "Description"},
			}}

			actualResponse, err := s.ListTeamsV1(context.Background(), req)

			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})
	})
})