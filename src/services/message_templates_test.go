package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMessageTemplates_Success(t *testing.T) {
	createTestTemplateFile(t)
	defer removeTestTemplateFile()

	templates, err := LoadMessageTemplates()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	if templates.Excellent.PT == "" {
		t.Error("Expected excellent PT template to be loaded")
	}

	if templates.Excellent.EN == "" {
		t.Error("Expected excellent EN template to be loaded")
	}

	if templates.Good.MinThreshold != 80 {
		t.Errorf("Expected good min threshold 80, got %d", templates.Good.MinThreshold)
	}

	if templates.Warning.MaxThreshold != 79 {
		t.Errorf("Expected warning max threshold 79, got %d", templates.Warning.MaxThreshold)
	}
}

func TestLoadMessageTemplates_FileNotFound(t *testing.T) {
	removeTestTemplateFile()

	_, err := LoadMessageTemplates()
	if err == nil {
		t.Error("Expected error when template file is missing")
	}
}

func TestLoadMessageTemplates_InvalidJSON(t *testing.T) {
	os.MkdirAll("config", 0755)
	invalidContent := `{"excellent": "invalid json`
	os.WriteFile(filepath.Join("config", "message_templates.json"), []byte(invalidContent), 0644)
	defer removeTestTemplateFile()

	_, err := LoadMessageTemplates()
	if err == nil {
		t.Error("Expected error when JSON is invalid")
	}
}

func createTestTemplateFile(t *testing.T) {
	templates := MessageTemplates{
		Excellent: TemplateCategory{
			Threshold: 100,
			PT:        "[Nome]: Assiduidade de 100% em [Curso]. Excelente!",
			EN:        "[Name]: 100% attendance in [Course]. Excellent!",
		},
		Good: TemplateCategory{
			MinThreshold: 80,
			MaxThreshold: 99,
			PT:           "[Nome]: Assiduidade de [X]% em [Curso]. Bom trabalho!",
			EN:           "[Name]: [X]% attendance in [Course]. Good job!",
		},
		Warning: TemplateCategory{
			MinThreshold: 50,
			MaxThreshold: 79,
			PT:           "[Nome]: Assiduidade de [X]% em [Curso]. Atenção!",
			EN:           "[Name]: [X]% attendance in [Course]. Warning!",
		},
		Critical: TemplateCategory{
			MaxThreshold: 49,
			PT:           "[Nome]: ALERTA! Assiduidade crítica de [X]% em [Curso].",
			EN:           "[Name]: ALERT! Critical attendance of [X]% in [Course].",
		},
	}

	os.MkdirAll("config", 0755)
	data, _ := json.MarshalIndent(templates, "", "  ")
	os.WriteFile(filepath.Join("config", "message_templates.json"), data, 0644)
}

func removeTestTemplateFile() {
	os.Remove(filepath.Join("config", "message_templates.json"))
}
