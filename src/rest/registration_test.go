package rest

import (
	"encoding/json"
	"skulla-api/db"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestListRegistrations_Success(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/registrations?studentClassId=1", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var registrations []db.Registration
	if err := json.Unmarshal(resp.Body.Bytes(), &registrations); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(registrations) != 3 {
		t.Errorf("Expected 3 registrations, got %d", len(registrations))
	}
}

func TestListRegistrations_MissingStudentClassId(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/registrations", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestListRegistrations_InvalidStudentClassId(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/registrations?studentClassId=invalid", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestListRegistrations_Unauthorized_WrongTeacher(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/registrations?studentClassId=1", testTeacherEmail2, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}
