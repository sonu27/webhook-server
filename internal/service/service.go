package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"webhook-server-test/internal/webhooks"
)

func NewService(logger *zap.Logger, client *http.Client) *Service {
	return &Service{
		l: logger,
		c: client,
		w: new(webhooks.Webhooks),
	}
}

type Service struct {
	l *zap.Logger
	c *http.Client
	w *webhooks.Webhooks
}

func (s *Service) CreateWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := fmt.Errorf("json decode: %w", err)
		s.l.Error(e.Error())
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	s.w.Add(req.Url, req.Token)

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) FireWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req FireWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := fmt.Errorf("json decode: %w", err)
		s.l.Error(e.Error())
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	webhooksClone := s.w.Get()

	w.WriteHeader(http.StatusAccepted)

	for _, v := range webhooksClone {
		fmt.Println(v)
		body := map[string]interface{}{
			"token":   v.Token,
			"payload": req.Payload,
		}
		fmt.Println(v.Url, body)

		go makeReq(s.c, v.Url, body)
	}
}

func makeReq(client *http.Client, url string, body map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovering from panic in makeReq: %v \n", r)
		}
	}()

	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(resp.StatusCode, bodyString)
}
