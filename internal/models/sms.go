package models

type SMSService struct {
	ApiKey   string `json:"api_key"`
	Username string `json:"username"`
	BaseURL  string `json:"base_url"`
	Env      string `json:"env"`
}

type SMSRequest struct {
	Username string `json:"username"`
	To       string `json:"to"`
	Message  string `json:"message"`
	From     string `json:"from,omitempty"`
}

type SMSResponse struct {
	SMSMessageData SMSMessageData `json:"SMSMessageData"`
}

type SMSMessageData struct {
	Message    string      `json:"Message"`
	Recipients []Recipient `json:"Recipients"`
}

type Recipient struct {
	StatusCode   int64  `json:"statusCode"`
	Number       string `json:"number"`
	Status       string `json:"status"`
	Cost         string `json:"cost"`
	MessageID    string `json:"messageId"`
	MessageParts int64  `json:"messageParts"`
}
