package services

import (
	"testing"
)

func TestGetMessageTemplate_Excellent_Portuguese(t *testing.T) {
	createTestTemplateFile(t)
	defer removeTestTemplateFile()

	template := GetMessageTemplate(100, "Informática Básica")

	if template == "" {
		t.Fatal("Expected non-empty template")
	}

	if template != "[Nome]: Assiduidade de 100% em [Curso]. Excelente!" {
		t.Errorf("Unexpected template: %s", template)
	}
}

func TestGetMessageTemplate_Good_Portuguese(t *testing.T) {
	createTestTemplateFile(t)
	defer removeTestTemplateFile()

	testCases := []float64{80, 85, 90, 95, 99}

	for _, percentage := range testCases {
		template := GetMessageTemplate(percentage, "Matemática")

		if template == "" {
			t.Fatalf("Expected non-empty template for %.0f%%", percentage)
		}

		if template != "[Nome]: Assiduidade de [X]% em [Curso]. Bom trabalho!" {
			t.Errorf("Unexpected template for %.0f%%: %s", percentage, template)
		}
	}
}

func TestGetMessageTemplate_Critical_Portuguese(t *testing.T) {
	createTestTemplateFile(t)
	defer removeTestTemplateFile()

	testCases := []float64{0, 10, 25, 49}

	for _, percentage := range testCases {
		template := GetMessageTemplate(percentage, "Português")

		if template == "" {
			t.Fatalf("Expected non-empty template for %.0f%%", percentage)
		}

		if template != "[Nome]: ALERTA! Assiduidade crítica de [X]% em [Curso]." {
			t.Errorf("Unexpected template for %.0f%%: %s", percentage, template)
		}
	}
}

func TestGetMessageTemplate_PortugueseDefault(t *testing.T) {
	createTestTemplateFile(t)
	defer removeTestTemplateFile()

	nonEnglishCourses := []string{
		"Matemática",
		"Física",
		"História",
		"Inglês",
	}

	for _, course := range nonEnglishCourses {
		template := GetMessageTemplate(100, course)

		if template != "[Nome]: Assiduidade de 100% em [Curso]. Excelente!" {
			t.Errorf("Expected Portuguese template for course '%s', got: %s", course, template)
		}
	}
}

func TestFormatMessage_AllPlaceholders(t *testing.T) {
	template := "[Nome]: Assiduidade de [X]% em [Curso]."
	studentName := "João Silva"
	courseName := "Matemática Avançada"
	percentage := 85.0

	result := FormatMessage(template, studentName, courseName, percentage)

	expected := "João Silva: Assiduidade de 85% em Matemática Avançada."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFormatMessage_EnglishTemplate(t *testing.T) {
	template := "[Name]: [X]% attendance in [Course]."
	studentName := "Mary Johnson"
	courseName := "Essential English"
	percentage := 100.0

	result := FormatMessage(template, studentName, courseName, percentage)

	expected := "Mary Johnson: 100% attendance in Essential English."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFormatMessage_ZeroPercentage(t *testing.T) {
	template := "[Nome]: Assiduidade de [X]%"
	studentName := "Test Student"
	courseName := "Test Course"
	percentage := 0.0

	result := FormatMessage(template, studentName, courseName, percentage)

	expected := "Test Student: Assiduidade de 0%"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFormatMessage_SpecialCharacters(t *testing.T) {
	template := "[Nome] - [Curso]: [X]%"
	studentName := "José António"
	courseName := "Informática & Redes"
	percentage := 95.0

	result := FormatMessage(template, studentName, courseName, percentage)

	expected := "José António - Informática & Redes: 95%"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFormatMessage_NoPlaceholders(t *testing.T) {
	template := "This is a static message"
	studentName := "Any Name"
	courseName := "Any Course"
	percentage := 50.0

	result := FormatMessage(template, studentName, courseName, percentage)

	if result != template {
		t.Errorf("Expected '%s', got '%s'", template, result)
	}
}
