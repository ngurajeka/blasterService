package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

type BlasterService interface {
	Send(string, string) string
	GetStatus(int) string
}

type blasterService struct{}

func (self blasterService) Send(target, message string) string {
	return fmt.Sprintf("Sending Message: %s, To: %s", message, target)
}

func (self blasterService) GetStatus(id int) string {
	return fmt.Sprintf("Getting Status queue id: %d", id)
}

type sendRequest struct {
	Target  string `json:"target"`
	Message string `json:"message"`
}

type sendResponse struct {
	Response string   `json:"response,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

type statusRequest struct {
	ID int `json:"id"`
}

type statusResponse struct {
	Response string `json:"response"`
}

func makeSendEndpoint(blaster BlasterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var errors []string
		req := request.(sendRequest)
		if req.Target == "" {
			errors = append(errors, "Target Empty")
		}
		if req.Message == "" {
			errors = append(errors, "Message Empty")
		}
		if len(errors) > 0 {
			return sendResponse{Errors: errors}, nil
		}
		result := blaster.Send(req.Target, req.Message)
		return sendResponse{Response: result}, nil
	}
}

func makeGetStatusEndpoint(blaster BlasterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(statusRequest)
		result := blaster.GetStatus(req.ID)
		return statusResponse{Response: result}, nil
	}
}

func main() {
	ctx := context.Background()
	svc := blasterService{}

	handler1 := httptransport.NewServer(
		ctx,
		makeSendEndpoint(svc),
		decodeSendRequest,
		encodeResponse,
	)

	http.Handle("/send", handler1)
	http.ListenAndServe(":8080", nil)
}

func decodeSendRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request sendRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
