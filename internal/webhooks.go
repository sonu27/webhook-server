package internal

import "sync"

var webhooks = Webhooks{
	mutex: &sync.Mutex{},
	data:  []Webhook{},
}

type Webhooks struct {
	mutex *sync.Mutex
	data  []Webhook
}

func (w *Webhooks) Add(url, token string) {
	w.mutex.Lock()
	w.data = append(w.data, Webhook{
		Url:   url,
		Token: token,
	})
	w.mutex.Unlock()
}

func (w *Webhooks) Clone() []Webhook {
	w.mutex.Lock()
	clone := make([]Webhook, len(w.data))
	copy(clone, w.data)
	w.mutex.Unlock()
	return clone
}

type Webhook struct {
	Url   string
	Token string
}