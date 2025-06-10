package worker

import (
	"log"
	"math/rand"
	"time"

	"github.com/paxaf/BrandScoutTest/internal/entity"
	"github.com/paxaf/BrandScoutTest/internal/usecase"
)

const (
	timeProgress = 15 * time.Second
)

type Scheduler struct {
	service     usecase.Usecase
	Tickers     tickers
	StopChannel chan struct{}
}

type tickers struct {
	SetInProgress *time.Ticker
}

func NewScheduler(service usecase.Usecase) *Scheduler {
	return &Scheduler{
		service:     service,
		StopChannel: make(chan struct{}),
		Tickers: tickers{
			SetInProgress: time.NewTicker(timeProgress),
		},
	}
}

func (s *Scheduler) Start() {
	go s.run()
	log.Println("сервис обработки задач запущен")
}

func (s *Scheduler) Stop() {
	if s.Tickers.SetInProgress != nil {
		s.Tickers.SetInProgress.Stop()
	}
	close(s.StopChannel)
	log.Println("Сервис планировщика остановлен")
}

func (s *Scheduler) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(nil, "Паника в планировщике", map[string]interface{}{
				"panic": r,
			})
			go s.run()
		}
	}()

	for {
		select {
		case <-s.Tickers.SetInProgress.C:
			go s.setupInProgress()
		case <-s.StopChannel:
			return
		}
	}
}

func (s *Scheduler) setupInProgress() {
	tasks := s.service.GetAll()
	for _, val := range tasks {
		if val.Status == entity.StatusNew {
			val.Status = entity.StatusInProgress
			s.service.Update(val.Id, val)
			go func(val entity.Task) {
				time.Sleep(time.Duration((rand.Intn(3) + 3)) * time.Minute)
				val.Status = entity.StatusDone
				s.service.Update(val.Id, val)
			}(val)
		}
	}
}
