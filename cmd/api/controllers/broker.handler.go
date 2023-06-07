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
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	case "log":
		b.logItem(w, requestPayload.Log)
	case "mail":
		b.sendMail(w, requestPayload.Mail)
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

	payload := &helpers.JsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}
	_ = payload.WriteJSON(w, http.StatusOK, nil)
}

func (b *brokerController) logItem(w http.ResponseWriter, l LogPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://logger-service/logger", bytes.NewBuffer(jsonData))
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
	if response.StatusCode != http.StatusOK {
		_ = b.json.WriteJSONError(w, errors.New("error calling logger service"), http.StatusInternalServerError)
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

	payload := &helpers.JsonResponse{
		Error:   false,
		Message: "Logged",
		Data:    jsonFromService.Data,
	}
	_ = payload.WriteJSON(w, http.StatusOK, nil)
}

func (b *brokerController) sendMail(w http.ResponseWriter, m MailPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonData))
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
	if response.StatusCode != http.StatusAccepted {
		_ = b.json.WriteJSONError(w, errors.New("error calling mail service"), http.StatusInternalServerError)
		return
	}

	//create a variable we'll read response.Body into

	payload := &helpers.JsonResponse{
		Error:   false,
		Message: "Mail sent to " + m.To,
	}
	_ = payload.WriteJSON(w, http.StatusOK, nil)
}
