package services

import (
	"bytes"
	"testing"
)

func TestGenerateCertificate_EnglishTemplate(t *testing.T) {
	pdfBytes, err := GenerateCertificate("english", "John Smith", "Successfully completed the Advanced Programming Course")
	if err != nil {
		t.Fatalf("GenerateCertificate failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Error("Generated content does not contain valid PDF magic bytes")
	}

	if !bytes.Contains(pdfBytes, []byte("%%EOF")) {
		t.Error("PDF does not contain EOF marker")
	}
}

func TestGenerateCertificate_PortugueseTemplate(t *testing.T) {
	pdfBytes, err := GenerateCertificate("portuguese", "Maria Silva", "Concluiu com sucesso o curso de Programação Avançada")
	if err != nil {
		t.Fatalf("GenerateCertificate failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Error("Generated content does not contain valid PDF magic bytes")
	}

	if !bytes.Contains(pdfBytes, []byte("%%EOF")) {
		t.Error("PDF does not contain EOF marker")
	}
}

func TestGenerateCertificate_InvalidTemplate(t *testing.T) {
	_, err := GenerateCertificate("spanish", "John Smith", "Successfully completed")
	if err == nil {
		t.Error("Expected error for invalid template, got nil")
	}

	if err.Error() != "invalid certificate template: spanish" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGenerateCertificate_WithSpecialCharacters(t *testing.T) {
	pdfBytes, err := GenerateCertificate("portuguese", "José António", "Concluiu com êxito o curso avançado de informática")
	if err != nil {
		t.Fatalf("GenerateCertificate failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Error("Generated content does not contain valid PDF magic bytes")
	}
}

func TestGenerateCertificate_LongDescription(t *testing.T) {
	longDescription := "This is a very long certificate description that will test the word wrapping functionality. It should span multiple lines and still generate a valid PDF document without any errors. The system must handle text that extends beyond the normal single-line capacity."

	pdfBytes, err := GenerateCertificate("english", "Elizabeth Johnson", longDescription)
	if err != nil {
		t.Fatalf("GenerateCertificate failed with long description: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Error("Generated content does not contain valid PDF magic bytes")
	}
}

func TestGenerateCertificate_MinimumValidInputs(t *testing.T) {
	pdfBytes, err := GenerateCertificate("english", "Al Be", "Short text")
	if err != nil {
		t.Fatalf("GenerateCertificate failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("PDF content is empty")
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Error("Generated content does not contain valid PDF magic bytes")
	}
}

func TestGenerateCertificate_PDFStructure(t *testing.T) {
	pdfBytes, err := GenerateCertificate("english", "Test Student", "Test Description for Certificate")
	if err != nil {
		t.Fatalf("GenerateCertificate failed: %v", err)
	}

	requiredPDFElements := [][]byte{
		[]byte("%PDF-"),
		[]byte("%%EOF"),
		[]byte("/Type"),
		[]byte("/Page"),
	}

	for _, element := range requiredPDFElements {
		if !bytes.Contains(pdfBytes, element) {
			t.Errorf("PDF missing required element: %s", string(element))
		}
	}
}
