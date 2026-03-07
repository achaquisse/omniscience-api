package rest

import (
	"encoding/json"
	"fmt"
	"skulla-api/db"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func setupGradeTestData(t *testing.T) (*fiber.App, uint, uint) {
	app := setupTestApp(t)

	maxScore := 20.0
	category := db.EvaluationCategory{
		ID:          1,
		CourseID:    1,
		Name:        "Quiz",
		MaxScore:    maxScore,
		Cardinality: 3,
		IsActive:    true,
	}
	if err := db.DB().Create(&category).Error; err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	score := 15.0
	pct := 75.0
	gradeDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	grade := db.Grade{
		ID:             1,
		CategoryID:     1,
		StudentClassID: 1,
		RegistrationID: 1,
		Name:           "Quiz 1",
		Date:           &gradeDate,
		MaxScore:       maxScore,
		Score:          &score,
		Percentage:     &pct,
	}
	if err := db.DB().Create(&grade).Error; err != nil {
		t.Fatalf("Failed to create grade: %v", err)
	}

	return app, 1, 1
}

func TestUpdateGrade_Success_MarksAndDate(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{
		"marks": 18.5,
		"date":  "2024-03-01",
	}

	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result["score"] != 18.5 {
		t.Errorf("Expected score 18.5, got %v", result["score"])
	}

	if result["updated_by"] != testTeacherEmail {
		t.Errorf("Expected updated_by %s, got %v", testTeacherEmail, result["updated_by"])
	}
}

func TestUpdateGrade_Success_MarksOnly(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{
		"marks": 10.0,
	}

	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}

func TestUpdateGrade_Success_DateOnly(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{
		"date": "2024-05-20",
	}

	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}

func TestUpdateGrade_InvalidRegistrationId(t *testing.T) {
	app, _, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{"marks": 10.0}
	url := fmt.Sprintf("/grades/abc/%d", gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestUpdateGrade_InvalidGradeId(t *testing.T) {
	app, registrationId, _ := setupGradeTestData(t)

	body := map[string]interface{}{"marks": 10.0}
	url := fmt.Sprintf("/grades/%d/abc", registrationId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestUpdateGrade_MarksAbove20(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{"marks": 21.0}
	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}

func TestUpdateGrade_MarksNegative(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{"marks": -1.0}
	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestUpdateGrade_InvalidDateFormat(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{"date": "01/03/2024"}
	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestUpdateGrade_EmptyBody(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestUpdateGrade_GradeNotFound(t *testing.T) {
	app, registrationId, _ := setupGradeTestData(t)

	body := map[string]interface{}{"marks": 10.0}
	url := fmt.Sprintf("/grades/%d/9999", registrationId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestUpdateGrade_GradeBelongsToAnotherRegistration(t *testing.T) {
	app, _, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{"marks": 10.0}
	url := fmt.Sprintf("/grades/2/%d", gradeId)
	resp, err := makeRequest(app, "PATCH", url, testTeacherEmail, body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}

func TestUpdateGrade_Unauthorized(t *testing.T) {
	app, registrationId, gradeId := setupGradeTestData(t)

	body := map[string]interface{}{"marks": 10.0}
	url := fmt.Sprintf("/grades/%d/%d", registrationId, gradeId)
	resp, err := makeRequest(app, "PATCH", url, "", body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.Code)
	}
}
