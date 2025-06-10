package usecase

import (
	"errors"
	"log"
	"strconv"

	"github.com/paxaf/BrandScoutTest/internal/entity"
)

func (uc *usecase) Delete(key string) error {
	_, ok := uc.repo.Get(key)
	if !ok {
		return errors.New("no rows affected")
	}
	uc.repo.Del(key)
	return nil
}

func (uc *usecase) GetAll() []entity.Task {
	return uc.repo.GetAll()
}

func (uc *usecase) Get(key string) (entity.Task, bool) {
	return uc.repo.Get(key)
}

func (uc *usecase) Set(value entity.Task) {
	key := uc.keyCounter.Add(1)
	keyStr := strconv.Itoa(int(key))
	value.Id = keyStr
	value.Status = "new"
	uc.repo.Set(keyStr, value)
	log.Println("successeful set value")
}
