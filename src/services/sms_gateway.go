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

type MessageResult struct {
	StudentID    uint    `json:"student_id"`
	StudentName  string  `json:"student_name"`
	PhoneNumber  string  `json:"phone_number"`
	CourseName   string  `json:"course_name"`
	Percentage   float64 `json:"percentage"`
	Message      string  `json:"message"`
	Success      bool    `json:"success"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

type SendReportResult struct {
	TotalMessages   int             `json:"total_messages"`
	SuccessCount    int             `json:"success_count"`
	FailureCount    int             `json:"failure_count"`
	MessagesSent    []MessageResult `json:"messages_sent"`
	MessagesFailed  []MessageResult `json:"messages_failed"`
	MessagesPlanned []MessageResult `json:"messages_planned,omitempty"`
}

func SendAttendanceReportsWithResult(students []StudentAttendanceInfo, dryRun bool) (*SendReportResult, error) {
	result := &SendReportResult{
		TotalMessages:   len(students),
		MessagesSent:    []MessageResult{},
		MessagesFailed:  []MessageResult{},
		MessagesPlanned: []MessageResult{},
	}

	client := NewSMSGatewayClient()

	for _, student := range students {
		template := GetMessageTemplate(student.Percentage, student.CourseName)
		if template == "" {
			msgResult := MessageResult{
				StudentID:    student.StudentID,
				StudentName:  student.StudentName,
				PhoneNumber:  student.PhoneNumber,
				CourseName:   student.CourseName,
				Percentage:   student.Percentage,
				Success:      false,
				ErrorMessage: "Failed to get message template",
			}
			if dryRun {
				result.MessagesPlanned = append(result.MessagesPlanned, msgResult)
			} else {
				result.MessagesFailed = append(result.MessagesFailed, msgResult)
				result.FailureCount++
			}
			continue
		}

		message := FormatMessage(template, student.StudentName, student.CourseName, student.Percentage)

		msgResult := MessageResult{
			StudentID:   student.StudentID,
			StudentName: student.StudentName,
			PhoneNumber: student.PhoneNumber,
			CourseName:  student.CourseName,
			Percentage:  student.Percentage,
			Message:     message,
		}

		if dryRun {
			msgResult.Success = true
			result.MessagesPlanned = append(result.MessagesPlanned, msgResult)
		} else {
			err := client.SendSMS(student.PhoneNumber, message)
			if err != nil {
				msgResult.Success = false
				msgResult.ErrorMessage = err.Error()
				result.MessagesFailed = append(result.MessagesFailed, msgResult)
				result.FailureCount++
				fmt.Printf("Failed to send SMS to %s (%s): %v\n", student.StudentName, student.PhoneNumber, err)
			} else {
				msgResult.Success = true
				result.MessagesSent = append(result.MessagesSent, msgResult)
				result.SuccessCount++
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	if dryRun {
		fmt.Printf("\nDry run: %d messages would be sent\n", len(result.MessagesPlanned))
	} else {
		fmt.Printf("\nAttendance reports sent: %d successful, %d failed\n", result.SuccessCount, result.FailureCount)
	}

	return result, nil
}

func SendAttendanceReports(students []StudentAttendanceInfo) error {
	_, err := SendAttendanceReportsWithResult(students, false)
	return err
}
