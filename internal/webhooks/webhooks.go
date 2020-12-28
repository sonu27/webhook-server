package webhooks

type webhook struct {
	Url   string
	Token string
}

type Webhooks struct {
	data []webhook
}

func (w *Webhooks) Add(url, token string) {
	w.data = append(w.data, webhook{
		Url:   url,
		Token: token,
	})
}

func (w *Webhooks) Get() []webhook {
	clone := make([]webhook, len(w.data))
	copy(clone, w.data)
	return clone
}
