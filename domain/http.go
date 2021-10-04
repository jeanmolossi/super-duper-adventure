package domain

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type HttpRequest struct {
	Client *http.Client
	Method string
	Url    string
	Body   io.Reader
}

type HttpResponse struct {
	Status     string
	StatusCode int
	Body       []byte
}

func NewRequest() *HttpRequest {
	return &HttpRequest{
		Client: &http.Client{
			Timeout: time.Second * 60,
		},
		Method: "GET",
		Body:   strings.NewReader(""),
	}
}

func (h *HttpRequest) Request() (*http.Response, error) {
	r := regexp.MustCompile("localhost")
	url := r.ReplaceAllString(h.Url, "gsr_mock_api")
	request, err := http.NewRequest(h.Method, url, h.Body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Origin-App", "true")

	response, err := h.Client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h *HttpRequest) GetResponse(response *http.Response) (*HttpResponse, error) {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Body:       body,
	}, nil
}
