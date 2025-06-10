package storage

import (
	"log"

	"github.com/paxaf/workmateTest/internal/entity"
)

type Engine struct {
	partition *HashTable
}

func NewEngine() (*Engine, error) {
	engine := &Engine{
		partition: NewHashTable(),
	}
	engine.partition = NewHashTable()
	return engine, nil
}

func (e *Engine) Set(key string, value entity.Task) {
	e.partition.Set(key, value)
	log.Println("succeseful set query")
}

func (e *Engine) Get(key string) (entity.Task, bool) {
	value, found := e.partition.Get(key)
	log.Println("succesefull get query")
	return value, found
}

func (e *Engine) Del(key string) {
	e.partition.Del(key)
	log.Println("succesefull delete query")
}

func (e *Engine) GetAll() []entity.Task {
	var res []entity.Task
	for _, val := range e.partition.data {
		res = append(res, val)
	}
	return res
}
