package saver

import (
	"context"
	"errors"
	"github.com/ozoncp/ocp-team-api/internal/flusher"
	"github.com/ozoncp/ocp-team-api/internal/models"
	"time"
)

type state uint8

const (
	initialized state = iota
	closed
)

type Saver interface {
	Save(team models.Team) error
	Close()
}

type saver struct {
	flusher flusher.Flusher
	teams   []models.Team
	teamsCh chan models.Team
	doneCh  chan int
	ticker  *time.Ticker
	state state
	capacity uint
}

func NewSaver(capacity uint, flusher flusher.Flusher, interval time.Duration) Saver {
	if capacity == 0 || interval <= 0 {
		return nil
	}

	s := &saver{
		flusher: flusher,
		teams:   make([]models.Team, 0, capacity),
		teamsCh: make(chan models.Team),
		doneCh:  make(chan int),
		ticker:  time.NewTicker(interval),
		state: initialized,
		capacity: capacity,
	}

	go func(){
		defer s.ticker.Stop()

		for {
			select {
			case team := <- s.teamsCh:
				s.teams = append(s.teams, team)
				if uint(len(s.teams)) >= s.capacity {
					s.flush()
				}
			case <- s.ticker.C:
				s.flush()
			case <- s.doneCh:
				s.state = closed
				s.flusher.Flush(context.TODO(), s.teams)
				close(s.doneCh)
				close(s.teamsCh)
				return
			}
		}
	}()

	return s
}

func (s *saver) flush() {
	failed := s.flusher.Flush(context.TODO(), s.teams)
	s.teams = make([]models.Team, 0, s.capacity)
	s.teams = append(s.teams, failed...)
}


func (s *saver) Save(team models.Team) error {
	if s.state == closed {
		return errors.New("cannot save to the closed saver")
	}

	s.teamsCh <- team

	return nil
}

func (s *saver) Close() {
	if s.state == closed {
		return
	}

	s.doneCh <- 1
}