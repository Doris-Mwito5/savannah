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

// SMSClient defines the contract for SMS operations
// This interface allows for different implementations and makes testing easier
type SMSClient interface {
	// SendSMS sends an SMS message to the specified phone number
	SendSMS(to, message string) (*models.SMSResponse, error)
	
	// ValidateResponse checks if the SMS was delivered successfully
	ValidateResponse(response *models.SMSResponse) error
}

// smsClient is the concrete implementation of SMSClient using Africa's Talking API
type smsClient struct {
	SMSService *models.SMSService  // Configuration for the SMS service
	client     *http.Client        // HTTP client for making API requests
}

// NewSMSClient creates a new instance of SMSClient with dependency injection
// This constructor pattern allows for flexible testing and configuration
func NewSMSClient(
	SMSService *models.SMSService, // SMS service configuration (credentials, URLs)
	client     *http.Client,       // HTTP client (can be customized for timeouts, etc.)
) SMSClient {
	return &smsClient{
		SMSService: SMSService,
		client:     client,
	}
}

// SendSMS sends an SMS message using Africa's Talking API
// It handles phone number formatting, API communication, and response parsing
func (p *smsClient) SendSMS(to, message string) (*models.SMSResponse, error) {
	// Validate input parameters to prevent unnecessary API calls
	if to == "" || message == "" {
		return nil, fmt.Errorf("phone number and message cannot be empty")
	}

	// Format phone number to E.164 international format (+254712345678)
	// This ensures consistency and compatibility with the SMS gateway
	formattedTo, err := p.formatPhoneNumber(to)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	// Prepare the request payload for Africa's Talking API
	// Using application/x-www-form-urlencoded as required by their API
	data := url.Values{}
	data.Set("username", p.SMSService.Username)  // Africa's Talking username
	data.Set("to", formattedTo)                  // Recipient phone number
	data.Set("message", message)                 // SMS content

	// Create HTTP POST request to Africa's Talking messaging endpoint
	req, err := http.NewRequest("POST", p.SMSService.BaseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required headers for Africa's Talking API authentication
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Required content type
	req.Header.Set("ApiKey", p.SMSService.ApiKey)                       // API key for authentication
	req.Header.Set("Accept", "application/json")                        // Expect JSON response
	req.Header.Set("User-Agent", "Savannah-POS/1.0")                    // Identify our application

	// Execute the HTTP request with the configured client
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	// Ensure response body is closed to prevent resource leaks
	defer resp.Body.Close()

	// Read the entire response body for processing
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle specific HTTP status codes with meaningful error messages
	if resp.StatusCode == 401 {
		// 401 Unauthorized - typically invalid API key or username
		return nil, fmt.Errorf("authentication failed - check your API key and username. Status: %d, Body: %s", 
			resp.StatusCode, string(body))
	}

	// Check for any non-success HTTP status codes (outside 200-299 range)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response into our structured SMSResponse model
	var smsResponse models.SMSResponse
	if err := json.Unmarshal(body, &smsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &smsResponse, nil
}

// formatPhoneNumber normalizes and validates phone numbers to E.164 format
// This ensures all numbers are in a consistent format for the SMS gateway
func (p *smsClient) formatPhoneNumber(phone string) (string, error) {
	// Remove common formatting characters to get raw digits
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Ensure the number starts with country code for international format
	if !strings.HasPrefix(cleaned, "+") {
		// Handle Kenyan phone numbers specifically:
		// - Numbers starting with 0 (like 0712345678) become +254712345678
		// - 9-digit numbers (like 712345678) become +254712345678
		if strings.HasPrefix(cleaned, "0") && len(cleaned) == 10 {
			// Convert local format (0712345678) to international (+254712345678)
			cleaned = "+254" + cleaned[1:]
		} else if len(cleaned) == 9 {
			// Convert 9-digit numbers to international format
			cleaned = "+254" + cleaned
		} else {
			return "", fmt.Errorf("invalid phone number format: %s", phone)
		}
	}

	// Extract digits only (without + prefix) for validation
	digits := strings.TrimPrefix(cleaned, "+")
	
	// Validate phone number length (typical international numbers are 9-12 digits)
	if len(digits) < 9 || len(digits) > 12 {
		return "", fmt.Errorf("invalid phone number length: %s", phone)
	}

	// Ensure all characters after + are numeric digits
	for _, char := range digits {
		if char < '0' || char > '9' {
			return "", fmt.Errorf("invalid characters in phone number: %s", phone)
		}
	}

	return cleaned, nil
}

// ValidateResponse checks if the SMS was successfully delivered to all recipients
// It examines the Africa's Talking API response for delivery status and error codes
func (p *smsClient) ValidateResponse(response *models.SMSResponse) error {
	// Check if there are any recipients in the response
	// Note: No nil check needed since len() returns 0 for nil slices
	if len(response.SMSMessageData.Recipients) == 0 {
		// If no recipients, check if the API provided an error message
		if response.SMSMessageData.Message != "" && response.SMSMessageData.Message != "Sent" {
			return fmt.Errorf("SMS API error: %s", response.SMSMessageData.Message)
		}
		return fmt.Errorf("no recipients in response")
	}

	// Check each recipient's delivery status
	for _, recipient := range response.SMSMessageData.Recipients {
		// Africa's Talking success status codes:
		// - 100: Processed (message is being processed)
		// - 101: Sent (message sent to mobile network)
		// - 102: Queued (message queued for delivery)
		// Any other code indicates a failure
		if recipient.StatusCode != 100 && recipient.StatusCode != 101 && recipient.StatusCode != 102 {
			return fmt.Errorf("SMS failed for %s: %s (code: %d)",
				recipient.Number, recipient.Status, recipient.StatusCode)
		}
	}

	// All recipients received the message successfully
	return nil
}