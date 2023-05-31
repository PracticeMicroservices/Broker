package controllers

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Broker interface {
	Broker(w http.ResponseWriter, r *http.Request)
}

type brokerController struct{}

func NewBrokerController() Broker {
	return &brokerController{}
}

func (b *brokerController) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Ping the broker",
		Data:    nil,
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}
