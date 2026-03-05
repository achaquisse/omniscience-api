package rest

import (
	"bytes"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestGenerateCertificate_Success_English(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=english&studentName=John%20Smith&certDescription=Successfully%20completed%20the%20course", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	pdfContent := resp.Body.Bytes()

	if len(pdfContent) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfContent, []byte("%PDF-")) {
		t.Error("Response does not contain valid PDF magic bytes")
	}

	if !bytes.Contains(pdfContent, []byte("%%EOF")) {
		t.Error("PDF does not contain EOF marker")
	}
}

func TestGenerateCertificate_Success_Portuguese(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=portuguese&studentName=Maria%20Silva&certDescription=Concluiu%20com%20sucesso%20o%20curso%20de%20exemplo", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	pdfContent := resp.Body.Bytes()

	if len(pdfContent) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfContent, []byte("%PDF-")) {
		t.Error("Response does not contain valid PDF magic bytes")
	}

	if !bytes.Contains(pdfContent, []byte("%%EOF")) {
		t.Error("PDF does not contain EOF marker")
	}
}

func TestGenerateCertificate_InvalidTemplate(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=spanish&studentName=John%20Smith&certDescription=Successfully%20completed", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGenerateCertificate_ShortStudentName(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=english&studentName=John&certDescription=Successfully%20completed", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGenerateCertificate_ShortDescription(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=english&studentName=John%20Smith&certDescription=Short", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGenerateCertificate_MissingTemplate(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?studentName=John%20Smith&certDescription=Successfully%20completed", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGenerateCertificate_MissingStudentName(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=english&certDescription=Successfully%20completed", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGenerateCertificate_MissingDescription(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=english&studentName=John%20Smith", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestGenerateCertificate_SpecialCharacters(t *testing.T) {
	app := setupTestApp(t)

	resp, err := makeRequest(app, "GET", "/certificates?template=portuguese&studentName=José%20António&certDescription=Concluiu%20com%20êxito%20o%20curso%20avançado", testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	pdfContent := resp.Body.Bytes()

	if !bytes.HasPrefix(pdfContent, []byte("%PDF-")) {
		t.Error("Response does not contain valid PDF magic bytes")
	}
}

func TestGenerateCertificate_LongDescription(t *testing.T) {
	app := setupTestApp(t)

	longDesc := "This%20is%20a%20very%20long%20certificate%20description%20that%20tests%20the%20word%20wrap%20functionality"
	resp, err := makeRequest(app, "GET", "/certificates?template=english&studentName=Elizabeth%20Johnson&certDescription="+longDesc, testTeacherEmail, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Code != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	pdfContent := resp.Body.Bytes()

	if len(pdfContent) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfContent, []byte("%PDF-")) {
		t.Error("Response does not contain valid PDF magic bytes")
	}
}
