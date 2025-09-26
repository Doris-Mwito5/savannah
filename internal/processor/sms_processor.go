package processor

import (
	"encoding/json"
	"fmt"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SMSClient interface {
	SendSMS(to, message string) (*models.SMSResponse, error)
	
	ValidateResponse(response *models.SMSResponse) error
}

type smsClient struct {
	SMSService *models.SMSService  
	client     *http.Client        
}


func NewSMSClient(
	SMSService *models.SMSService, 
	client     *http.Client,       
) SMSClient {
	return &smsClient{
		SMSService: SMSService,
		client:     client,
	}
}


func (p *smsClient) SendSMS(to, message string) (*models.SMSResponse, error) {
	if to == "" || message == "" {
		return nil, fmt.Errorf("phone number and message cannot be empty")
	}

	formattedTo, err := p.formatPhoneNumber(to)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	data := url.Values{}
	data.Set("username", p.SMSService.Username)  
	data.Set("to", formattedTo)                  
	data.Set("message", message)                 

	req, err := http.NewRequest("POST", p.SMSService.BaseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") 
	req.Header.Set("ApiKey", p.SMSService.ApiKey)                       
	req.Header.Set("Accept", "application/json")                       
	req.Header.Set("User-Agent", "Savannah-POS/1.0")                    
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("authentication failed - check your API key and username. Status: %d, Body: %s", 
			resp.StatusCode, string(body))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var smsResponse models.SMSResponse
	if err := json.Unmarshal(body, &smsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &smsResponse, nil
}


func (p *smsClient) formatPhoneNumber(phone string) (string, error) {
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	if !strings.HasPrefix(cleaned, "+") {
		
		if strings.HasPrefix(cleaned, "0") && len(cleaned) == 10 {
			cleaned = "+254" + cleaned[1:]
		} else if len(cleaned) == 9 {
			cleaned = "+254" + cleaned
		} else {
			return "", fmt.Errorf("invalid phone number format: %s", phone)
		}
	}

	digits := strings.TrimPrefix(cleaned, "+")
	
	if len(digits) < 9 || len(digits) > 12 {
		return "", fmt.Errorf("invalid phone number length: %s", phone)
	}

	for _, char := range digits {
		if char < '0' || char > '9' {
			return "", fmt.Errorf("invalid characters in phone number: %s", phone)
		}
	}

	return cleaned, nil
}

func (p *smsClient) ValidateResponse(response *models.SMSResponse) error {
	if len(response.SMSMessageData.Recipients) == 0 {
		if response.SMSMessageData.Message != "" && response.SMSMessageData.Message != "Sent" {
			return fmt.Errorf("SMS API error: %s", response.SMSMessageData.Message)
		}
		return fmt.Errorf("no recipients in response")
	}

	for _, recipient := range response.SMSMessageData.Recipients {

		if recipient.StatusCode != 100 && recipient.StatusCode != 101 && recipient.StatusCode != 102 {
			return fmt.Errorf("SMS failed for %s: %s (code: %d)",
				recipient.Number, recipient.Status, recipient.StatusCode)
		}
	}

	return nil
}