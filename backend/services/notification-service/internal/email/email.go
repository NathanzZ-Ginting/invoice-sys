package email

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// SimpleEmailService represents a simple email service
type SimpleEmailService struct {
	APIKey  string
	BaseURL string
}

// SendEmail sends email using available service (console fallback)
func SendEmail(to, subject, body string, attachmentName string, attachment []byte) error {
	// Try Resend first
	if apiKey := os.Getenv("RESEND_API_KEY"); apiKey != "" {
		return sendWithResend(to, subject, body, attachmentName, attachment)
	}

	// Try Mailgun
	if apiKey := os.Getenv("MAILGUN_API_KEY"); apiKey != "" {
		return sendWithMailgun(to, subject, body, attachmentName, attachment)
	}

	// Fallback to console (development mode)
	fmt.Printf("=== EMAIL SIMULATION (No service configured) ===\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Body: %s\n", body)
	if attachmentName != "" {
		fmt.Printf("Attachment: %s (%d bytes)\n", attachmentName, len(attachment))
	}
	fmt.Printf("=================================================\n")

	return nil
}

// sendWithResend sends email using Resend API
func sendWithResend(to, subject, body, attachmentName string, attachment []byte) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY not set")
	}

	payload := map[string]interface{}{
		"from":    "Invoice Generator <noreply@resend.dev>",
		"to":      []string{to},
		"subject": subject,
		"text":    body,
	}

	if attachmentName != "" && len(attachment) > 0 {
		payload["attachments"] = []map[string]interface{}{
			{
				"filename": attachmentName,
				"content":  base64.StdEncoding.EncodeToString(attachment),
			},
		}
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("resend send failed with status: %d", resp.StatusCode)
	}

	return nil
}

// sendWithMailgun sends email using Mailgun API
func sendWithMailgun(to, subject, body, attachmentName string, attachment []byte) error {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	domain := os.Getenv("MAILGUN_DOMAIN")
	if apiKey == "" || domain == "" {
		return fmt.Errorf("MAILGUN_API_KEY and MAILGUN_DOMAIN not set")
	}

	// For simplicity, send without attachment first
	payload := fmt.Sprintf("from=Invoice Generator <noreply@%s>&to=%s&subject=%s&text=%s",
		domain, to, subject, body)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", domain), bytes.NewBufferString(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("api:"+apiKey)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("mailgun send failed with status: %d", resp.StatusCode)
	}

	return nil
}