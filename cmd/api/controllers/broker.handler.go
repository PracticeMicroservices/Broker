package controllers

import (
	"broker/cmd/api/helpers"
	"net/http"
)

type Broker interface {
	Broker(w http.ResponseWriter, r *http.Request)
}

type brokerController struct{}

func NewBrokerController() Broker {
	return &brokerController{}
}

func (b *brokerController) Broker(w http.ResponseWriter, r *http.Request) {
	res := &helpers.JsonResponse{
		Error:   false,
		Message: "Hello from the broker",
		Data:    nil,
	}
	_ = res.WriteJSON(w, http.StatusOK, nil)
}
