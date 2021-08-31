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

		s                 desc.OcpTeamApiServer
		mockRepo          *mocks.MockIRepo
		mockKafkaProducer *mocks.MockIProducer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockRepo = mocks.NewMockIRepo(ctrl)
		mockKafkaProducer = mocks.NewMockIProducer(ctrl)
		s = api.NewOcpTeamApi(mockRepo, mockKafkaProducer)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CreateTeamV1()", func() {
		It("returns response", func() {
			mockRepo.EXPECT().CreateTeam(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			mockKafkaProducer.EXPECT().Send(gomock.Any()).Return(nil)

			req := &desc.CreateTeamV1Request{Name: "Name", Description: "Description"}

			_, err := s.CreateTeamV1(context.Background(), req)
			Expect(err).Should(BeNil())
		})
	})

	Context("GetTeamV1()", func() {
		It("get existing team by id", func() {
			mockKafkaProducer.EXPECT().Send(gomock.Any()).Return(nil).Times(0)

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
			mockKafkaProducer.EXPECT().Send(gomock.Any()).Return(nil).Times(0)

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
			mockKafkaProducer.EXPECT().Send(gomock.Any()).Return(nil).Times(1)

			mockRepo.EXPECT().RemoveTeam(gomock.Any(), gomock.Any()).Return(nil)

			req := &desc.RemoveTeamV1Request{Id: uint64(1)}
			expectedResponse := &desc.RemoveTeamV1Response{}

			actualResponse, err := s.RemoveTeamV1(context.Background(), req)
			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})
	})

	Context("UpdateTeamV1()", func() {
		It("updates existing element", func() {
			mockKafkaProducer.EXPECT().Send(gomock.Any()).Return(nil).Times(1)

			mockRepo.EXPECT().UpdateTeam(gomock.Any(), gomock.Any()).Return(nil)

			req := &desc.UpdateTeamV1Request{Team: &desc.Team{Id: uint64(1), Name: "Name1", Description: "Descr1"}}
			expectedResponse := &desc.UpdateTeamV1Response{}

			actualResponse, err := s.UpdateTeamV1(context.Background(), req)
			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})
	})

	Context("ListTeamsV1()", func() {
		It("return nothing when limit and offset are minimal", func() {
			mockKafkaProducer.EXPECT().Send(gomock.Any()).Return(nil).Times(0)

			mockRepo.EXPECT().ListTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Team{}, uint64(0), nil)

			req := &desc.ListTeamsV1Request{Limit: uint64(1)}
			expectedResponse := &desc.ListTeamsV1Response{Total: uint64(0), Teams: []*desc.Team{}}

			actualResponse, err := s.ListTeamsV1(context.Background(), req)

			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})

		It("return teams when limit and offset are set", func() {
			mockKafkaProducer.EXPECT().Send(gomock.Any()).Return(nil).Times(0)

			mockRepo.EXPECT().ListTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				[]models.Team{
					{Id: uint64(1), Name: "Name", Description: "Description"},
					{Id: uint64(2), Name: "Name", Description: "Description"},
				}, uint64(2), nil)

			req := &desc.ListTeamsV1Request{Limit: 2, Offset: 2}
			expectedResponse := &desc.ListTeamsV1Response{Total: uint64(2), Teams: []*desc.Team{
				{Id: uint64(1), Name: "Name", Description: "Description"},
				{Id: uint64(2), Name: "Name", Description: "Description"},
			}}

			actualResponse, err := s.ListTeamsV1(context.Background(), req)

			Expect(err).Should(BeNil())
			Expect(actualResponse).Should(Equal(expectedResponse))
		})
	})
})
