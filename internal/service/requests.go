package service

type CreateWebhookRequest struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

type FireWebhookRequest struct {
	Payload interface{} `json:"payload"`
}
