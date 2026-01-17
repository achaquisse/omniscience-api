package rest

import (
	"encoding/json"
	"skulla-api"
	"skulla-api/db"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRecordAttendance_Create_Success(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := map[string]interface{}{
		"registration_id": 1,
		"date":            "2024-01-20",
		"status":          "PRESENT",
		"remarks":         "Test attendance",
	}

	resp, err := main.makeRequest(app, "POST", "/attendance", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "Attendance recorded successfully" {
		t.Errorf("Unexpected message: %v", response["message"])
	}
}

func TestRecordAttendance_Update_Success(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := map[string]interface{}{
		"registration_id": 1,
		"date":            "2024-01-15",
		"status":          "LATE",
		"remarks":         "Updated status",
	}

	resp, err := main.makeRequest(app, "POST", "/attendance", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}

func TestRecordAttendance_InvalidRequestBody(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := map[string]interface{}{
		"registration_id": "invalid",
	}

	resp, err := main.makeRequest(app, "POST", "/attendance", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestRecordAttendance_MissingRequiredFields(t *testing.T) {
	app := main.setupTestApp(t)

	testCases := []map[string]interface{}{
		{"date": "2024-01-20", "status": "PRESENT"},
		{"registration_id": 1, "status": "PRESENT"},
		{"registration_id": 1, "date": "2024-01-20"},
	}

	for i, reqBody := range testCases {
		resp, err := main.makeRequest(app, "POST", "/attendance", main.testTeacherEmail, reqBody)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}

		if resp.Code != fiber.StatusBadRequest {
			t.Errorf("Test case %d: Expected status 400, got %d", i, resp.Code)
		}
	}
}

func TestRecordAttendance_InvalidStatus(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := map[string]interface{}{
		"registration_id": 1,
		"date":            "2024-01-20",
		"status":          "INVALID_STATUS",
		"remarks":         "",
	}

	resp, err := main.makeRequest(app, "POST", "/attendance", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestRecordAttendance_RegistrationNotFound(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := map[string]interface{}{
		"registration_id": 9999,
		"date":            "2024-01-20",
		"status":          "PRESENT",
		"remarks":         "",
	}

	resp, err := main.makeRequest(app, "POST", "/attendance", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestRecordBulkAttendance_Success(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := []map[string]interface{}{
		{
			"registration_id": 1,
			"date":            "2024-01-25",
			"status":          "PRESENT",
			"remarks":         "Bulk test 1",
		},
		{
			"registration_id": 2,
			"date":            "2024-01-25",
			"status":          "ABSENT",
			"remarks":         "Bulk test 2",
		},
		{
			"registration_id": 3,
			"date":            "2024-01-25",
			"status":          "LATE",
			"remarks":         "Bulk test 3",
		},
	}

	resp, err := main.makeRequest(app, "POST", "/attendance/bulk", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["records_processed"].(float64) != 3 {
		t.Errorf("Expected 3 records processed, got %v", response["records_processed"])
	}
}

func TestRecordBulkAttendance_EmptyArray(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := []map[string]interface{}{}

	resp, err := main.makeRequest(app, "POST", "/attendance/bulk", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestRecordBulkAttendance_InvalidRecord(t *testing.T) {
	app := main.setupTestApp(t)

	reqBody := []map[string]interface{}{
		{
			"registration_id": 1,
			"date":            "2024-01-25",
			"status":          "PRESENT",
		},
		{
			"registration_id": 9999,
			"date":            "2024-01-25",
			"status":          "PRESENT",
		},
	}

	resp, err := main.makeRequest(app, "POST", "/attendance/bulk", main.testTeacherEmail, reqBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}

func TestGetStudentAttendanceReport_WithClassId(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/report?student_id=1&student_class_id=1&start_date=2024-01-01&end_date=2024-01-31", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var report db.DetailedAttendanceReport
	if err := json.Unmarshal(resp.Body.Bytes(), &report); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if report.Summary.TotalDays != 3 {
		t.Errorf("Expected 3 total days, got %d", report.Summary.TotalDays)
	}

	if report.Summary.PresentCount != 2 {
		t.Errorf("Expected 2 present days, got %d", report.Summary.PresentCount)
	}
}

func TestGetStudentAttendanceReport_WithoutClassId(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/report?student_id=1&start_date=2024-01-01&end_date=2024-01-31", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var report db.AggregatedStudentAttendanceReport
	if err := json.Unmarshal(resp.Body.Bytes(), &report); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if report.OverallSummary.TotalDays != 5 {
		t.Errorf("Expected 5 total days (across all classes), got %d", report.OverallSummary.TotalDays)
	}

	if len(report.ByClass) != 2 {
		t.Errorf("Expected 2 classes, got %d", len(report.ByClass))
	}
}

func TestGetStudentAttendanceReport_MissingStudentId(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/report", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGetStudentAttendanceReport_InvalidStudentId(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/report?student_id=invalid", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGetStudentAttendanceReport_InvalidDateFormat(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/report?student_id=1&start_date=invalid-date", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGetClassAttendanceReport_Success(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/class-report?student_class_id=1&start_date=2024-01-01&end_date=2024-01-31", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var report db.ClassAttendanceReport
	if err := json.Unmarshal(resp.Body.Bytes(), &report); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if report.TotalStudents != 3 {
		t.Errorf("Expected 3 total students, got %d", report.TotalStudents)
	}

	if report.Period != "all" {
		t.Errorf("Expected period 'all', got %s", report.Period)
	}
}

func TestGetClassAttendanceReport_WithPeriodDay(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/class-report?student_class_id=1&start_date=2024-01-01&end_date=2024-01-31&period=day", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var report db.ClassAttendanceReport
	if err := json.Unmarshal(resp.Body.Bytes(), &report); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if report.Period != "day" {
		t.Errorf("Expected period 'day', got %s", report.Period)
	}

	if len(report.DailyData) == 0 {
		t.Error("Expected daily data to be populated")
	}
}

func TestGetClassAttendanceReport_WithPeriodWeek(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/class-report?student_class_id=1&period=week", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var report db.ClassAttendanceReport
	if err := json.Unmarshal(resp.Body.Bytes(), &report); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if report.Period != "week" {
		t.Errorf("Expected period 'week', got %s", report.Period)
	}
}

func TestGetClassAttendanceReport_WithPeriodMonth(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/class-report?student_class_id=1&period=month", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var report db.ClassAttendanceReport
	if err := json.Unmarshal(resp.Body.Bytes(), &report); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if report.Period != "month" {
		t.Errorf("Expected period 'month', got %s", report.Period)
	}
}

func TestGetClassAttendanceReport_InvalidPeriod(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/class-report?student_class_id=1&period=invalid", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGetClassAttendanceReport_MissingStudentClassId(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/class-report", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGetClassAttendanceReport_InvalidStudentClassId(t *testing.T) {
	app := main.setupTestApp(t)

	resp, err := main.makeRequest(app, "GET", "/attendance/class-report?student_class_id=invalid", main.testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}
