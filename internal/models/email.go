package models

// EmailService holds email service configuration
type EmailService struct {
    SMTPHost     string `json:"smtp_host"`
    SMTPPort     int    `json:"smtp_port"`
    Username     string `json:"username"`
    Password     string `json:"password"`
    FromEmail    string `json:"from_email"`
    FromName     string `json:"from_name"`
    AdminEmail   string `json:"admin_email"` // Administrator email address
}

// EmailRequest represents an email sending request
type EmailRequest struct {
    To          []string `json:"to"`
    Subject     string   `json:"subject"`
    Body        string   `json:"body"`
    ContentType string   `json:"content_type"` // "text/plain" or "text/html"
}

// EmailResponse represents the email sending response
type EmailResponse struct {
    Success    bool     `json:"success"`
    Message    string   `json:"message"`
    MessageID  string   `json:"message_id,omitempty"`
    Recipients []string `json:"recipients,omitempty"`
}