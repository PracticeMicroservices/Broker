package controllers

import (
	"broker/cmd/api/helpers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Broker interface {
	Broker(w http.ResponseWriter, r *http.Request)
	HandleSubmission(w http.ResponseWriter, r *http.Request)
}

type brokerController struct {
	json *helpers.JsonResponse
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewBrokerController() Broker {
	return &brokerController{
		json: &helpers.JsonResponse{},
	}
}

func (b *brokerController) Broker(w http.ResponseWriter, r *http.Request) {
	res := &helpers.JsonResponse{
		Error:   false,
		Message: "Hello from the broker",
		Data:    nil,
	}
	_ = res.WriteJSON(w, http.StatusOK, nil)
}

func (b *brokerController) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	requestPayload := RequestPayload{}

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = b.json.WriteJSONError(w, err, http.StatusBadRequest)
		return
	}
	switch requestPayload.Action {
	case "auth":
		b.Authenticate(w, requestPayload.Auth)
	default:
		_ = b.json.WriteJSONError(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (b *brokerController) Authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice

	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service

	request, err := http.NewRequest("POST", "http://authentication-service/authentication", bytes.NewBuffer(jsonData))
	if err != nil {
		_ = b.json.WriteJSONError(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		_ = b.json.WriteJSONError(w, err)
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	fmt.Println("Response: ", response)
	if response.StatusCode == http.StatusBadRequest {
		_ = b.json.WriteJSONError(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	} else if response.StatusCode != http.StatusOK {
		_ = b.json.WriteJSONError(w, errors.New("error calling auth service"), http.StatusInternalServerError)
		return
	}

	//create a variable we'll read response.Body into

	jsonFromService := &helpers.JsonResponse{}
	err = json.NewDecoder(response.Body).Decode(jsonFromService)
	if err != nil {
		_ = b.json.WriteJSONError(w, err)
		return
	}
	if jsonFromService.Error {
		_ = b.json.WriteJSONError(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	payload := &helpers.JsonResponse{}
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	_ = payload.WriteJSON(w, http.StatusOK, nil)

}
