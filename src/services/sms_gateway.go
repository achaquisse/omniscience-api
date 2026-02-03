package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type SMSRequest struct {
	Topic    string `json:"topic"`
	ToNumber string `json:"to_number"`
	Body     string `json:"body"`
}

type SMSResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type SMSGatewayClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewSMSGatewayClient() *SMSGatewayClient {
	baseURL := os.Getenv("SMS_GATEWAY_URL")
	if baseURL == "" {
		baseURL = "https://smsgateway.omniscience.co.mz"
	}

	return &SMSGatewayClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *SMSGatewayClient) SendSMS(phoneNumber, message string) error {
	topic := os.Getenv("SMS_TOPIC")
	if topic == "" {
		topic = "attendance"
	}

	smsReq := SMSRequest{
		Topic:    topic,
		ToNumber: phoneNumber,
		Body:     message,
	}

	jsonData, err := json.Marshal(smsReq)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("SMS gateway returned status %d: %s", resp.StatusCode, string(body))
	}

	var smsResp SMSResponse
	if err := json.Unmarshal(body, &smsResp); err != nil {
		return fmt.Errorf("failed to parse SMS response: %w", err)
	}

	fmt.Printf("SMS sent successfully. ID: %s, Phone: %s\n", smsResp.ID, phoneNumber)
	return nil
}

func SendAttendanceReports(students []StudentAttendanceInfo) error {
	client := NewSMSGatewayClient()
	successCount := 0
	failCount := 0

	for _, student := range students {
		template := GetMessageTemplate(student.Percentage, student.CourseName)
		if template == "" {
			fmt.Printf("Failed to get template for student %s\n", student.StudentName)
			failCount++
			continue
		}

		message := FormatMessage(template, student.StudentName, student.CourseName, student.Percentage)

		err := client.SendSMS(student.PhoneNumber, message)
		if err != nil {
			fmt.Printf("Failed to send SMS to %s (%s): %v\n", student.StudentName, student.PhoneNumber, err)
			failCount++
		} else {
			successCount++
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\nAttendance reports sent: %d successful, %d failed\n", successCount, failCount)
	return nil
}
