package main

import (
	"net/http"
	"testing"
)

type ServiceResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func TestTestService_Call(t *testing.T) {

	testCases := []struct {
		method         string
		api            string
		expectedCode   int
		expectedStatus string
	}{
		{"GET", "/", 404, "ERROR"},
		{"GET", "/orders", 200, "SUCCESS"},
		{"GET", "/orders/955", 404, "ERROR"},
		{"GET", "/orders/956", 200, "SUCCESS"},
		{"POST", "/orders/", 200, "SUCCESS"},
	}

	for _, tc := range testCases {
		s := TestService{}
		response := s.Call(tc.api, tc.method, &http.Request{})
		if response.Code != tc.expectedCode {
			t.Errorf("Error in code. Expected: %v. Got: %v", tc.expectedCode, response.Code)
		}
		if response.Status != tc.expectedStatus {
			t.Errorf("Error in code. Expected: %s. Got: %s", tc.expectedStatus, response.Status)
		}
	}
}

type TestService struct{}

func (t TestService) Call(api, method string, r *http.Request) ServiceResponse {

	switch method {
	case "GET":
		switch api {
		case "/orders":
			return ServiceResponse{200, "SUCCESS", "OK", nil}
		case "/orders/955":
			return ServiceResponse{404, "ERROR", "NOT OK", nil}
		case "/orders/956":
			return ServiceResponse{200, "SUCCESS", "OK", nil}
		case "/orders/":
			return ServiceResponse{200, "SUCCESS", "OK", nil}
		default:
			return ServiceResponse{404, "ERROR", "This is an invalid API.", nil}
		}
	case "POST":
		return ServiceResponse{200, "SUCCESS", "OK", nil}
	default:
		return ServiceResponse{404, "ERROR", "Invalid Method.", nil}
	}
}
