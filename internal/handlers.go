package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 5*time.Second,
}

type CreateWebhookRequest struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

func createWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := fmt.Errorf("json decode: %w", err)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	webhooks.Add(req.Url, req.Token)

	w.WriteHeader(http.StatusCreated)
}

type TestWebhookRequest struct {
	Payload interface{} `json:"payload"`
}

func testWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req TestWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := fmt.Errorf("json decode: %w", err)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	clone := webhooks.Clone()

	w.WriteHeader(http.StatusAccepted)

	for _, v := range clone {
		fmt.Println(v)
		body := map[string]interface{} {
			"token": v.Token,
			"payload": req.Payload,
		}
		fmt.Println(v.Url, body)

		go makeReq(v.Url, body)
	}
}

func makeReq(url string, body map[string]interface{}) {
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

	resp, err := Client.Do(req)
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
