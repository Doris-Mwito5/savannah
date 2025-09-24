// internal/processor/email.go
package processor

import (
    "crypto/tls"
    "fmt"
    "github/Doris-Mwito5/savannah-pos/internal/models"
    "net/smtp"
    "strings"
)

type EmailClient interface {
    SendEmail(req *models.EmailRequest) (*models.EmailResponse, error)
    SendOrderNotification(order *models.Order) error
}

type emailClient struct {
    emailService *models.EmailService
}

func NewEmailClient(emailService *models.EmailService) EmailClient {
    return &emailClient{
        emailService: emailService,
    }
}

// SendEmail sends a generic email using SMTP
func (e *emailClient) SendEmail(req *models.EmailRequest) (*models.EmailResponse, error) {
    // Validate email request
    if err := e.validateEmailRequest(req); err != nil {
        return nil, fmt.Errorf("invalid email request: %w", err)
    }

    // Set default content type if not specified
    if req.ContentType == "" {
        req.ContentType = "text/html"
    }

    // Prepare authentication
    auth := smtp.PlainAuth("", 
        e.emailService.Username, 
        e.emailService.Password, 
        e.emailService.SMTPHost,
    )

    // Prepare message
    message := e.buildMessage(req)

    // Send email
    err := e.sendSMTP(req.To, message, auth)
    if err != nil {
        return &models.EmailResponse{
            Success: false,
            Message: fmt.Sprintf("Failed to send email: %v", err),
        }, err
    }

    return &models.EmailResponse{
        Success:    true,
        Message:    "Email sent successfully",
        Recipients: req.To,
    }, nil
}

// SendOrderNotification sends an order notification email to administrators
func (e *emailClient) SendOrderNotification(order *models.Order) error {
    // Create email subject
    subject := fmt.Sprintf("📦 New Order Received - #%s", order.ReferenceNumber)

    // Create email body
    body := e.createOrderEmailBody(order)

    // Prepare email request
    emailReq := &models.EmailRequest{
        To:          []string{e.emailService.AdminEmail},
        Subject:     subject,
        Body:        body,
        ContentType: "text/html",
    }

    // Send the email
    response, err := e.SendEmail(emailReq)
    if err != nil {
        return fmt.Errorf("failed to send order notification email: %w", err)
    }

    if !response.Success {
        return fmt.Errorf("email sending failed: %s", response.Message)
    }

    return nil
}

// validateEmailRequest validates the email request parameters
func (e *emailClient) validateEmailRequest(req *models.EmailRequest) error {
    if len(req.To) == 0 {
        return fmt.Errorf("recipient list cannot be empty")
    }

    for _, email := range req.To {
        if !e.isValidEmail(email) {
            return fmt.Errorf("invalid email address: %s", email)
        }
    }

    if strings.TrimSpace(req.Subject) == "" {
        return fmt.Errorf("email subject cannot be empty")
    }

    if strings.TrimSpace(req.Body) == "" {
        return fmt.Errorf("email body cannot be empty")
    }

    return nil
}

// isValidEmail performs basic email validation
func (e *emailClient) isValidEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// buildMessage constructs the email message in proper format
func (e *emailClient) buildMessage(req *models.EmailRequest) []byte {
    fromHeader := fmt.Sprintf("From: %s <%s>\r\n", e.emailService.FromName, e.emailService.FromEmail)
    toHeader := fmt.Sprintf("To: %s\r\n", strings.Join(req.To, ","))
    subjectHeader := fmt.Sprintf("Subject: %s\r\n", req.Subject)
    contentTypeHeader := fmt.Sprintf("Content-Type: %s; charset=UTF-8\r\n", req.ContentType)
    
    message := []byte(fromHeader + 
        toHeader + 
        subjectHeader + 
        contentTypeHeader + 
        "\r\n" + 
        req.Body)
    
    return message
}

// sendSMTP handles the actual SMTP connection and email sending
func (e *emailClient) sendSMTP(to []string, message []byte, auth smtp.Auth) error {
    // SMTP server address
    smtpAddr := fmt.Sprintf("%s:%d", e.emailService.SMTPHost, e.emailService.SMTPPort)

    // For TLS connection (most modern SMTP servers require this)
    tlsConfig := &tls.Config{
        ServerName: e.emailService.SMTPHost,
    }

    // Connect to SMTP server
    conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
    if err != nil {
        return fmt.Errorf("TLS connection failed: %w", err)
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, e.emailService.SMTPHost)
    if err != nil {
        return fmt.Errorf("SMTP client creation failed: %w", err)
    }
    defer client.Close()

    // Authenticate
    if err := client.Auth(auth); err != nil {
        return fmt.Errorf("SMTP authentication failed: %w", err)
    }

    // Set sender
    if err := client.Mail(e.emailService.FromEmail); err != nil {
        return fmt.Errorf("setting sender failed: %w", err)
    }

    // Set recipients
    for _, recipient := range to {
        if err := client.Rcpt(recipient); err != nil {
            return fmt.Errorf("adding recipient failed: %w", err)
        }
    }

    // Send email data
    w, err := client.Data()
    if err != nil {
        return fmt.Errorf("getting data writer failed: %w", err)
    }

    _, err = w.Write(message)
    if err != nil {
        return fmt.Errorf("writing message failed: %w", err)
    }

    err = w.Close()
    if err != nil {
        return fmt.Errorf("closing data writer failed: %w", err)
    }

    return client.Quit()
}

// createOrderEmailBody creates a nicely formatted HTML email for order notifications
func (e *emailClient) createOrderEmailBody(order *models.Order) string {
    return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f8f9fa; padding: 20px; text-align: center; border-radius: 5px; }
        .order-info { background: #fff; padding: 20px; margin: 20px 0; border: 1px solid #ddd; border-radius: 5px; }
        .order-detail { margin: 10px 0; }
        .status-pending { color: #856404; background: #fff3cd; padding: 5px 10px; border-radius: 3px; }
        .status-confirmed { color: #155724; background: #d4edda; padding: 5px 10px; border-radius: 3px; }
        .table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        .table th, .table td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        .table th { background: #f8f9fa; }
        .total-row { font-weight: bold; background: #f8f9fa; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🛍️ New Order Received</h1>
            <p>Savannah POS System Notification</p>
        </div>
        
        <div class="order-info">
            <h2>Order Details</h2>
            
            <div class="order-detail">
                <strong>Order Reference:</strong> %s
            </div>
            <div class="order-detail">
                <strong>Order ID:</strong> %d
            </div>
            <div class="order-detail">
                <strong>Order Date:</strong> %s
            </div>
            <div class="order-detail">
                <strong>Customer Phone:</strong> %s
            </div>
            <div class="order-detail">
                <strong>Order Status:</strong> <span class="status-%s">%s</span>
            </div>
            <div class="order-detail">
                <strong>Payment Method:</strong> %s
            </div>
            <div class="order-detail">
                <strong>Order Medium:</strong> %s
            </div>
        </div>

        <div class="order-info">
            <h2>Order Items</h2>
            <table class="table">
                <thead>
                    <tr>
                        <th>Product ID</th>
                        <th>Unit Price</th>
                        <th>Quantity</th>
                        <th>Total</th>
                    </tr>
                </thead>
                <tbody>
                    %s
                </tbody>
                <tfoot>
                    <tr class="total-row">
                        <td colspan="3"><strong>Subtotal:</strong></td>
                        <td>$%.2f</td>
                    </tr>
                    <tr class="total-row">
                        <td colspan="3"><strong>Discount:</strong></td>
                        <td>$%.2f</td>
                    </tr>
                    <tr class="total-row">
                        <td colspan="3"><strong>Total Amount:</strong></td>
                        <td>$%.2f</td>
                    </tr>
                </tfoot>
            </table>
        </div>

        <div class="order-info">
            <p><strong>Action Required:</strong> Please review this order in the admin dashboard.</p>
            <p>This is an automated notification from Savannah POS System.</p>
        </div>
    </div>
</body>
</html>`,
        order.ReferenceNumber,
        order.ID,
        order.CreatedAt.Format("2006-01-02 15:04:05"),
        order.PhoneNumber,
        strings.ToLower(string(order.OrderStatus)),
        order.OrderStatus,
        order.PaymentMethod,
        order.OrderMedium,
        e.generateOrderItemsHTML(order),
        order.TotalAmount,
        e.getDiscountAmount(order),
        order.TotalAmount,
    )
}

// generateOrderItemsHTML generates HTML for order items table
func (e *emailClient) generateOrderItemsHTML(order *models.Order) string {
    return `
        <tr>
            <td>1</td>
            <td>$999.99</td>
            <td>2</td>
            <td>$1999.98</td>
        </tr>
        <tr>
            <td>2</td>
            <td>$9999.99</td>
            <td>1</td>
            <td>$9999.99</td>
        </tr>`
}

// getDiscountAmount returns the discount amount (handles nil discount)
func (e *emailClient) getDiscountAmount(order *models.Order) float64 {
    if order.Discount == nil {
        return 0.0
    }
    return *order.Discount
}