package controllers

import (
	"broker/cmd/api/helpers"
	"broker/logs"
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPC interface {
	LogViaGRPC(w http.ResponseWriter, r *http.Request)
}

type gRPC struct {
	json *helpers.JsonResponse
}

func NewGRPCController() GRPC {
	return &gRPC{
		json: &helpers.JsonResponse{},
	}
}

func (g *gRPC) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LogViaGRPC")
	var requestPayload RequestPayload

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = g.json.WriteJSONError(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		_ = g.json.WriteJSONError(w, err)
		return
	}
	defer conn.Close()

	client := logs.NewLogServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		_ = g.json.WriteJSONError(w, err)
		return
	}

	resp := &helpers.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = resp.WriteJSON(w, http.StatusOK, nil)
}
