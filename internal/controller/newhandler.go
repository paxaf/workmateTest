package controller

import (
	"github.com/paxaf/workmateTest/internal/usecase"
)

type UsecaseHandler struct {
	service usecase.Usecase
}

func New(s usecase.Usecase) *UsecaseHandler {
	return &UsecaseHandler{service: s}
}
