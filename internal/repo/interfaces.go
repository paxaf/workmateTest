package repo

import "github.com/paxaf/workmateTest/internal/entity"

type Repository interface {
	Set(key string, value entity.Task)
	Del(key string)
	Get(key string) (entity.Task, bool)
	GetAll() []entity.Task
}
