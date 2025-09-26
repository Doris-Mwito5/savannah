package notification

import (
	"fmt"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/processor"
	"log"
	"strings"
)

type OrderNotification interface {
	SendOrderSMS(order *models.Order) error
	SendOrderEmail(order *models.Order) error
	SendOrderNotifications(order *models.Order) error
}

type orderNotification struct {
	smsProcessor   processor.SMSClient
	emailProcessor processor.EmailClient
}

func NewOrderNotification(
	smsProcessor processor.SMSClient,
	emailProcessor processor.EmailClient,
) OrderNotification {
	return &orderNotification{
		smsProcessor: smsProcessor,
		emailProcessor: emailProcessor,
	}
}

func (n *orderNotification) SendOrderSMS(order *models.Order) error {
	message := fmt.Sprintf(
		"🎉 Order Confirmed!\n"+
			"Order ID: %d\n"+
			"Total: $%.2f\n"+
			"Items: %d\n"+
			"We'll notify you when it's ready. Thank you!",
		order.ID,
		order.TotalAmount,
		order.TotalItems,
	)
	log.Printf("message: %s\n", message)
	return n.sendSMSWithValidation(order.PhoneNumber, message)
}

func (n *orderNotification) sendSMSWithValidation(phone, message string) error {
	if !n.isValidPhoneNumber(phone) {
		return fmt.Errorf("invalid phone number format: %s", phone)
	}

	response, err := n.smsProcessor.SendSMS(phone, message)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	log.Printf("response: %v\n", response)

	if err := n.smsProcessor.ValidateResponse(response); err != nil {
		return fmt.Errorf("SMS delivery failed: %w", err)
	}

	if len(response.SMSMessageData.Recipients) > 0 {
		recipient := response.SMSMessageData.Recipients[0]
		fmt.Printf("✅ SMS sent to %s | ID: %s | Cost: %s\n",
			recipient.Number, recipient.MessageID, recipient.Cost)
	}

	return nil
}

func (n *orderNotification) isValidPhoneNumber(phone string) bool {
	if !strings.HasPrefix(phone, "+") {
		return false
	}

	digits := strings.ReplaceAll(phone[1:], " ", "")
	if len(digits) < 10 || len(digits) > 15 {
		return false
	}

	for _, char := range digits {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}



// SendOrderNotifications sends both SMS and email notifications
func (n *orderNotification) SendOrderNotifications(order *models.Order) error {
    // Send SMS notification (async - don't block on failure)
    go func() {
        if err := n.SendOrderSMS(order); err != nil {
            log.Printf("Failed to send order SMS: %v", err)
        }
    }()

    // Send email notification 
    go func() {
        if err := n.SendOrderEmail(order); err != nil {
            log.Printf("Failed to send order email: %v", err)
        }
    }()

    return nil
}

// SendOrderEmail sends email notification to administrator
func (n *orderNotification) SendOrderEmail(order *models.Order) error {
    if n.emailProcessor == nil {
        return fmt.Errorf("email processor not configured")
    }

    err := n.emailProcessor.SendOrderNotification(order)
    if err != nil {
        return fmt.Errorf("failed to send order email: %w", err)
    }

    log.Printf("✅ Order notification email sent for order #%s", order.ReferenceNumber)
    return nil
}