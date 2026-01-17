package rest

import (
	"encoding/json"
	"fmt"
	"skulla-api/db"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestListStudentClasses_Success(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/student-classes", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var classes []db.StudentClass
	if err := json.Unmarshal(resp.Body.Bytes(), &classes); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(classes) != 3 {
		t.Errorf("Expected 3 classes, got %d", len(classes))
	}
}

func TestListStudentClasses_WithDateFilters(t *testing.T) {
	app := setupTestApp(t)

	now := time.Now()
	startDate := now.AddDate(0, -3, 0).Format("2006-01-02")
	endDate := now.AddDate(0, 1, 0).Format("2006-01-02")

	resp, err := makeRequest(app, "GET", fmt.Sprintf("/student-classes?startDate=%s&endDate=%s", startDate, endDate), testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var classes []db.StudentClass
	if err := json.Unmarshal(resp.Body.Bytes(), &classes); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(classes) != 2 {
		t.Errorf("Expected 2 classes (filtered by date), got %d", len(classes))
	}
}

func TestListStudentClasses_InvalidDateFormat(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/student-classes?startDate=invalid-date", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}

	var errorResp map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &errorResp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if errorResp["error"] == "" {
		t.Error("Expected error message in response")
	}
}

func TestListStudentClasses_Unauthorized_MissingToken(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/student-classes", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.Code)
	}
}
