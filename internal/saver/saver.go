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

// Saver is the interface for saving teams.
// Actual saving to db is done either on ticker event or
// on when Close() method is called.
type Saver interface {
	Save(team models.Team) error
	Close()
}

// saver is the struct that implements Saver interface.
type saver struct {
	flusher  flusher.Flusher
	teams    []models.Team
	teamsCh  chan models.Team
	doneCh   chan int
	ticker   *time.Ticker
	state    state
	capacity uint
}

// NewSaver is the constructor method for saver struct.
// In addition, it constructs ticker.
func NewSaver(capacity uint, flusher flusher.Flusher, interval time.Duration) *saver {
	if capacity == 0 || interval <= 0 {
		return nil
	}

	s := &saver{
		flusher:  flusher,
		teams:    make([]models.Team, 0, capacity),
		teamsCh:  make(chan models.Team),
		doneCh:   make(chan int),
		ticker:   time.NewTicker(interval),
		state:    initialized,
		capacity: capacity,
	}

	go func() {
		defer s.ticker.Stop()

		for {
			select {
			case team := <-s.teamsCh:
				s.teams = append(s.teams, team)
				if uint(len(s.teams)) >= s.capacity {
					s.flush()
				}
			case <-s.ticker.C:
				s.flush()
			case <-s.doneCh:
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

// Save is the method for adding new team to the save channel.
func (s *saver) Save(team models.Team) error {
	if s.state == closed {
		return errors.New("cannot save to the closed saver")
	}

	s.teamsCh <- team

	return nil
}

// Close is the method for closing the saver.
// It sends value to the done channel to close it.
func (s *saver) Close() {
	if s.state == closed {
		return
	}

	s.doneCh <- 1
}
