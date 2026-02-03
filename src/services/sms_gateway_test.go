package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewSMSGatewayClient_DefaultURL(t *testing.T) {
	os.Unsetenv("SMS_GATEWAY_URL")

	client := NewSMSGatewayClient()

	if client.BaseURL != "https://smsgateway.omniscience.co.mz" {
		t.Errorf("Expected default URL, got %s", client.BaseURL)
	}
}

func TestNewSMSGatewayClient_CustomURL(t *testing.T) {
	customURL := "https://custom.sms.com"
	os.Setenv("SMS_GATEWAY_URL", customURL)
	defer os.Unsetenv("SMS_GATEWAY_URL")

	client := NewSMSGatewayClient()

	if client.BaseURL != customURL {
		t.Errorf("Expected %s, got %s", customURL, client.BaseURL)
	}
}

func TestSendSMS_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/messages" {
			t.Errorf("Expected /messages path, got %s", r.URL.Path)
		}

		var req SMSRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.ToNumber != "+258840000000" {
			t.Errorf("Expected phone +258840000000, got %s", req.ToNumber)
		}

		if req.Topic == "" {
			t.Error("Expected topic to be set")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SMSResponse{
			Message: "Message queued successfully",
			ID:      "msg_123",
		})
	}))
	defer server.Close()

	client := &SMSGatewayClient{
		BaseURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	err := client.SendSMS("+258840000000", "Test message")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendSMS_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer server.Close()

	client := &SMSGatewayClient{
		BaseURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	err := client.SendSMS("+258840000000", "Test message")
	if err == nil {
		t.Error("Expected error for server error response")
	}
}

func TestSendSMS_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request"}`))
	}))
	defer server.Close()

	client := &SMSGatewayClient{
		BaseURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	err := client.SendSMS("invalid-phone", "Test message")
	if err == nil {
		t.Error("Expected error for bad request")
	}
}

func TestSendSMS_CustomTopic(t *testing.T) {
	customTopic := "custom-topic"
	os.Setenv("SMS_TOPIC", customTopic)
	defer os.Unsetenv("SMS_TOPIC")

	receivedTopic := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SMSRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedTopic = req.Topic

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SMSResponse{
			Message: "Message queued successfully",
			ID:      "msg_123",
		})
	}))
	defer server.Close()

	client := &SMSGatewayClient{
		BaseURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	client.SendSMS("+258840000000", "Test message")

	if receivedTopic != customTopic {
		t.Errorf("Expected topic %s, got %s", customTopic, receivedTopic)
	}
}

func TestSendSMS_DefaultTopic(t *testing.T) {
	os.Unsetenv("SMS_TOPIC")

	receivedTopic := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SMSRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedTopic = req.Topic

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SMSResponse{
			Message: "Message queued successfully",
			ID:      "msg_123",
		})
	}))
	defer server.Close()

	client := &SMSGatewayClient{
		BaseURL:    server.URL,
		HTTPClient: &http.Client{},
	}

	client.SendSMS("+258840000000", "Test message")

	if receivedTopic != "attendance" {
		t.Errorf("Expected default topic 'attendance', got %s", receivedTopic)
	}
}

func TestSendAttendanceReports_EmptyList(t *testing.T) {
	students := []StudentAttendanceInfo{}

	err := SendAttendanceReports(students)
	if err != nil {
		t.Errorf("Expected no error for empty list, got %v", err)
	}
}

func TestSendAttendanceReports_WithStudents(t *testing.T) {
	createTestTemplateFile(t)
	defer removeTestTemplateFile()

	messageCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messageCount++
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SMSResponse{
			Message: "Message queued successfully",
			ID:      "msg_123",
		})
	}))
	defer server.Close()

	os.Setenv("SMS_GATEWAY_URL", server.URL)
	defer os.Unsetenv("SMS_GATEWAY_URL")

	students := []StudentAttendanceInfo{
		{
			StudentID:   1,
			StudentName: "João Silva",
			PhoneNumber: "+258840000001",
			CourseName:  "Matemática",
			Percentage:  100,
		},
		{
			StudentID:   2,
			StudentName: "Maria Santos",
			PhoneNumber: "+258840000002",
			CourseName:  "Essential English",
			Percentage:  85,
		},
	}

	err := SendAttendanceReports(students)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if messageCount != 2 {
		t.Errorf("Expected 2 messages sent, got %d", messageCount)
	}
}
