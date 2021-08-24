package saver_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozoncp/ocp-team-api/internal/mocks"
	"github.com/ozoncp/ocp-team-api/internal/models"
	"github.com/ozoncp/ocp-team-api/internal/saver"
	"time"
)

var _ = Describe("Saver", func() {
	var (
		ctrl	*gomock.Controller
		mockFlusher *mocks.MockFlusher
		s saver.Saver
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockFlusher = mocks.NewMockFlusher(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when saver has no capacity", func() {
		It("returns nil on saver creation", func() {
			mockFlusher.EXPECT().Flush(gomock.Any(), gomock.Any()).Times(0)

			s = saver.NewSaver(0, mockFlusher, time.Second)

			gomega.Expect(s).Should(gomega.BeNil())

			gomega.Expect(func() {
				s.Close()
			}).Should(gomega.Panic())
		})
	})

	Context("when saver's flush interval is negative or zero", func() {
		It("returns nil on saver creation", func() {
			mockFlusher.EXPECT().Flush(gomock.Any(), gomock.Any()).Times(0)

			s = saver.NewSaver(1, mockFlusher, -1 * time.Second)

			gomega.Expect(s).Should(gomega.BeNil())

			gomega.Expect(func() {
				s.Close()
			}).Should(gomega.Panic())
		})
	})

	Context("when saver's capacity overloaded", func() {
		It("flushes elements", func() {
			mockFlusher.EXPECT().Flush(gomock.Any(), gomock.Any()).Return([]models.Team{}).AnyTimes()
			s = saver.NewSaver(1, mockFlusher, 10 * time.Second)

			for i := 0; i < 5; i++ {
				_ = s.Save(models.Team{Id: uint64(i), Name: "Name", Description: "Desc"})
			}

			s.Close()
		})
	})

	Context("when try to close saver multiple times", func() {
		It("does not panic", func() {
			mockFlusher.EXPECT().Flush(gomock.Any(), gomock.Any()).Return(nil).Times(1)

			s = saver.NewSaver(10, mockFlusher, 10 * time.Second)

			s.Close()
			gomega.Expect(func(){
				s.Close()
			}).ShouldNot(gomega.Panic())
		})
	})

	Context("when try to Save() on invalid closed state", func() {
		It("returns error", func() {
			mockFlusher.EXPECT().Flush(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			s = saver.NewSaver(10, mockFlusher, 10 * time.Second)

			s.Close()
			err := s.Save(models.Team{Id: 0, Name: "Name", Description: "Desc"})
			gomega.Expect(err).ShouldNot(gomega.BeNil())
		})
	})
})
