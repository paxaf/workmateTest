package usecase

import (
	"sync/atomic"

	"github.com/paxaf/BrandScoutTest/internal/entity"
	"github.com/paxaf/BrandScoutTest/internal/repo"
)

type Usecase interface {
	Delete(key string) error
	GetAll() []entity.Task
	Set(value entity.Task)
	Get(key string) (entity.Task, bool)
	Update(key string, value entity.Task)
}

type usecase struct {
	repo       repo.Repository
	keyCounter atomic.Int64
}

func New(repo repo.Repository) *usecase {
	return &usecase{
		repo:       repo,
		keyCounter: atomic.Int64{},
	}
}
