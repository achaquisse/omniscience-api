package services

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type TemplateCategory struct {
	Threshold    int    `json:"threshold,omitempty"`
	MinThreshold int    `json:"min_threshold,omitempty"`
	MaxThreshold int    `json:"max_threshold,omitempty"`
	PT           string `json:"pt"`
	EN           string `json:"en"`
}

type MessageTemplates struct {
	Excellent TemplateCategory `json:"excellent"`
	Good      TemplateCategory `json:"good"`
	Warning   TemplateCategory `json:"warning"`
	Critical  TemplateCategory `json:"critical"`
}

func LoadMessageTemplates() (*MessageTemplates, error) {
	possiblePaths := []string{
		filepath.Join("config", "message_templates.json"),
		filepath.Join("..", "config", "message_templates.json"),
		"config/message_templates.json",
		"../config/message_templates.json",
	}

	var data []byte
	var err error

	for _, configPath := range possiblePaths {
		data, err = os.ReadFile(configPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var templates MessageTemplates
	err = json.Unmarshal(data, &templates)
	if err != nil {
		return nil, err
	}

	return &templates, nil
}
