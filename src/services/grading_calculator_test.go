package services

import (
	"skulla-api/db"
	"testing"
)

func TestCalculateInformaticaFormula(t *testing.T) {
	gc := NewGradingCalculator()

	informaticaFormula := db.FormulaConfig{
		"type": "multi_stage",
		"stages": map[string]interface{}{
			"MDF": map[string]interface{}{
				"type": "weighted_average",
				"categories": map[string]interface{}{
					"Windows":    0.25,
					"Excel":      0.25,
					"Word":       0.25,
					"PowerPoint": 0.25,
				},
			},
			"MF": map[string]interface{}{
				"type":       "simple_average",
				"components": []interface{}{"MDF", "Exame"},
			},
		},
	}

	categoryScores := db.CategoryScores{
		"Windows": db.CategoryScore{
			CategoryID:   1,
			CategoryName: "Windows",
			Score:        85.0,
			Percentage:   85.0,
		},
		"Excel": db.CategoryScore{
			CategoryID:   2,
			CategoryName: "Excel",
			Score:        90.0,
			Percentage:   90.0,
		},
		"Word": db.CategoryScore{
			CategoryID:   3,
			CategoryName: "Word",
			Score:        88.0,
			Percentage:   88.0,
		},
		"PowerPoint": db.CategoryScore{
			CategoryID:   4,
			CategoryName: "PowerPoint",
			Score:        92.0,
			Percentage:   92.0,
		},
		"Exame": db.CategoryScore{
			CategoryID:   5,
			CategoryName: "Exame",
			Score:        80.0,
			Percentage:   80.0,
		},
	}

	result, err := gc.calculateCustomFormula(informaticaFormula, categoryScores, 1)
	if err != nil {
		t.Fatalf("Failed to calculate Informática formula: %v", err)
	}

	expectedMDF := (85.0 + 90.0 + 88.0 + 92.0) / 4.0
	expectedMF := (expectedMDF + 80.0) / 2.0

	if result != expectedMF {
		t.Errorf("Expected MF = %.2f, got %.2f", expectedMF, result)
		t.Logf("Expected MDF = %.2f", expectedMDF)
		t.Logf("Expected MF = (%.2f + 80.0) / 2 = %.2f", expectedMDF, expectedMF)
	} else {
		t.Logf("✓ Informática formula correct: MDF=%.2f, MF=%.2f", expectedMDF, result)
	}
}

func TestCalculateInglesLevel1Formula(t *testing.T) {
	gc := NewGradingCalculator()

	inglesFormula := db.FormulaConfig{
		"type": "multi_stage_level_based",
		"levels": map[string]interface{}{
			"1": map[string]interface{}{
				"stages": map[string]interface{}{
					"MDF": map[string]interface{}{
						"type": "weighted_average",
						"categories": map[string]interface{}{
							"Avaliação escrita": map[string]interface{}{"weight": 0.60, "count": 4},
							"Ditados":           map[string]interface{}{"weight": 0.20, "count": 8},
							"Composições":       map[string]interface{}{"weight": 0.10, "count": 2},
							"Karaoke":           map[string]interface{}{"weight": 0.10, "count": 2},
						},
					},
					"MF": map[string]interface{}{
						"type":       "simple_average",
						"components": []interface{}{"MDF", "Exame"},
					},
				},
			},
		},
	}

	categoryScores := db.CategoryScores{
		"Avaliação escrita": db.CategoryScore{
			CategoryID:   1,
			CategoryName: "Avaliação escrita",
			Score:        85.0,
			Percentage:   85.0,
		},
		"Ditados": db.CategoryScore{
			CategoryID:   2,
			CategoryName: "Ditados",
			Score:        90.0,
			Percentage:   90.0,
		},
		"Composições": db.CategoryScore{
			CategoryID:   3,
			CategoryName: "Composições",
			Score:        75.0,
			Percentage:   75.0,
		},
		"Karaoke": db.CategoryScore{
			CategoryID:   4,
			CategoryName: "Karaoke",
			Score:        80.0,
			Percentage:   80.0,
		},
		"Exame": db.CategoryScore{
			CategoryID:   5,
			CategoryName: "Exame",
			Score:        88.0,
			Percentage:   88.0,
		},
	}

	result, err := gc.calculateMultiStageLevelBased(inglesFormula, categoryScores, 1)
	if err != nil {
		t.Fatalf("Failed to calculate Inglês Level 1 formula: %v", err)
	}

	expectedMDF := (0.60 * 85.0) + (0.20 * 90.0) + (0.10 * 75.0) + (0.10 * 80.0)
	expectedMF := (expectedMDF + 88.0) / 2.0

	if result != expectedMF {
		t.Errorf("Expected MF = %.2f, got %.2f", expectedMF, result)
		t.Logf("Expected MDF = %.2f", expectedMDF)
		t.Logf("Expected MF = (%.2f + 88.0) / 2 = %.2f", expectedMDF, expectedMF)
	} else {
		t.Logf("✓ Inglês Level 1 formula correct: MDF=%.2f, MF=%.2f", expectedMDF, result)
	}
}

func TestCalculateInglesLevel2Formula(t *testing.T) {
	gc := NewGradingCalculator()

	inglesFormula := db.FormulaConfig{
		"type": "multi_stage_level_based",
		"levels": map[string]interface{}{
			"2": map[string]interface{}{
				"stages": map[string]interface{}{
					"MDF": map[string]interface{}{
						"type": "weighted_average",
						"categories": map[string]interface{}{
							"Avaliação escrita": map[string]interface{}{"weight": 0.60, "count": 4},
							"Ditados":           map[string]interface{}{"weight": 0.10, "count": 4},
							"Composições":       map[string]interface{}{"weight": 0.10, "count": 2},
							"Karaoke":           map[string]interface{}{"weight": 0.10, "count": 2},
							"Listening":         map[string]interface{}{"weight": 0.10, "count": 1},
						},
					},
					"MF": map[string]interface{}{
						"type":       "simple_average",
						"components": []interface{}{"MDF", "Exame"},
					},
				},
			},
		},
	}

	categoryScores := db.CategoryScores{
		"Avaliação escrita": db.CategoryScore{
			CategoryID:   1,
			CategoryName: "Avaliação escrita",
			Score:        85.0,
			Percentage:   85.0,
		},
		"Ditados": db.CategoryScore{
			CategoryID:   2,
			CategoryName: "Ditados",
			Score:        90.0,
			Percentage:   90.0,
		},
		"Composições": db.CategoryScore{
			CategoryID:   3,
			CategoryName: "Composições",
			Score:        75.0,
			Percentage:   75.0,
		},
		"Karaoke": db.CategoryScore{
			CategoryID:   4,
			CategoryName: "Karaoke",
			Score:        80.0,
			Percentage:   80.0,
		},
		"Listening": db.CategoryScore{
			CategoryID:   6,
			CategoryName: "Listening",
			Score:        82.0,
			Percentage:   82.0,
		},
		"Exame": db.CategoryScore{
			CategoryID:   5,
			CategoryName: "Exame",
			Score:        88.0,
			Percentage:   88.0,
		},
	}

	result, err := gc.calculateMultiStageLevelBased(inglesFormula, categoryScores, 2)
	if err != nil {
		t.Fatalf("Failed to calculate Inglês Level 2 formula: %v", err)
	}

	expectedMDF := (0.60 * 85.0) + (0.10 * 90.0) + (0.10 * 75.0) + (0.10 * 80.0) + (0.10 * 82.0)
	expectedMF := (expectedMDF + 88.0) / 2.0

	tolerance := 0.01
	if result < (expectedMF-tolerance) || result > (expectedMF+tolerance) {
		t.Errorf("Expected MF = %.2f, got %.2f", expectedMF, result)
		t.Logf("Expected MDF = %.2f", expectedMDF)
		t.Logf("Expected MF = (%.2f + 88.0) / 2 = %.2f", expectedMDF, expectedMF)
	} else {
		t.Logf("✓ Inglês Level 2 formula correct: MDF=%.2f, MF=%.2f", expectedMDF, result)
	}
}

func TestCalculateInglesLevel3Formula(t *testing.T) {
	gc := NewGradingCalculator()

	inglesFormula := db.FormulaConfig{
		"type": "multi_stage_level_based",
		"levels": map[string]interface{}{
			"3": map[string]interface{}{
				"stages": map[string]interface{}{
					"MDF": map[string]interface{}{
						"type": "weighted_average",
						"categories": map[string]interface{}{
							"Avaliação escrita": map[string]interface{}{"weight": 0.60, "count": 4},
							"Ditados":           map[string]interface{}{"weight": 0.10, "count": 4},
							"Composições":       map[string]interface{}{"weight": 0.10, "count": 2},
							"Speeches":          map[string]interface{}{"weight": 0.10, "count": 2},
							"Listening":         map[string]interface{}{"weight": 0.10, "count": 1},
						},
					},
					"MF": map[string]interface{}{
						"type":       "simple_average",
						"components": []interface{}{"MDF", "Exame"},
					},
				},
			},
		},
	}

	categoryScores := db.CategoryScores{
		"Avaliação escrita": db.CategoryScore{Score: 85.0},
		"Ditados":           db.CategoryScore{Score: 90.0},
		"Composições":       db.CategoryScore{Score: 75.0},
		"Speeches":          db.CategoryScore{Score: 78.0},
		"Listening":         db.CategoryScore{Score: 82.0},
		"Exame":             db.CategoryScore{Score: 88.0},
	}

	result, err := gc.calculateMultiStageLevelBased(inglesFormula, categoryScores, 3)
	if err != nil {
		t.Fatalf("Failed to calculate Inglês Level 3 formula: %v", err)
	}

	expectedMDF := (0.60 * 85.0) + (0.10 * 90.0) + (0.10 * 75.0) + (0.10 * 78.0) + (0.10 * 82.0)
	expectedMF := (expectedMDF + 88.0) / 2.0

	if result != expectedMF {
		t.Errorf("Expected MF = %.2f, got %.2f", expectedMF, result)
	} else {
		t.Logf("✓ Inglês Level 3 formula correct: MDF=%.2f, MF=%.2f", expectedMDF, result)
	}
}

func TestCalculateInglesLevel4Formula(t *testing.T) {
	gc := NewGradingCalculator()

	inglesFormula := db.FormulaConfig{
		"type": "multi_stage_level_based",
		"levels": map[string]interface{}{
			"4": map[string]interface{}{
				"stages": map[string]interface{}{
					"MDF": map[string]interface{}{
						"type": "weighted_average",
						"categories": map[string]interface{}{
							"Avaliação escrita": map[string]interface{}{"weight": 0.60, "count": 4},
							"Ditados":           map[string]interface{}{"weight": 0.10, "count": 4},
							"Composições":       map[string]interface{}{"weight": 0.10, "count": 2},
							"Speeches":          map[string]interface{}{"weight": 0.10, "count": 2},
							"Listening":         map[string]interface{}{"weight": 0.10, "count": 1},
						},
					},
					"MF": map[string]interface{}{
						"type":       "simple_average",
						"components": []interface{}{"MDF", "Exame"},
					},
				},
			},
		},
	}

	categoryScores := db.CategoryScores{
		"Avaliação escrita": db.CategoryScore{Score: 85.0},
		"Ditados":           db.CategoryScore{Score: 90.0},
		"Composições":       db.CategoryScore{Score: 75.0},
		"Speeches":          db.CategoryScore{Score: 78.0},
		"Listening":         db.CategoryScore{Score: 82.0},
		"Exame":             db.CategoryScore{Score: 88.0},
	}

	result, err := gc.calculateMultiStageLevelBased(inglesFormula, categoryScores, 4)
	if err != nil {
		t.Fatalf("Failed to calculate Inglês Level 4 formula: %v", err)
	}

	expectedMDF := (0.60 * 85.0) + (0.10 * 90.0) + (0.10 * 75.0) + (0.10 * 78.0) + (0.10 * 82.0)
	expectedMF := (expectedMDF + 88.0) / 2.0

	if result != expectedMF {
		t.Errorf("Expected MF = %.2f, got %.2f", expectedMF, result)
	} else {
		t.Logf("✓ Inglês Level 4 formula correct: MDF=%.2f, MF=%.2f", expectedMDF, result)
	}
}

func TestWeightedAverageFormula(t *testing.T) {
	gc := NewGradingCalculator()

	formula := db.FormulaConfig{
		"type": "weighted_average",
		"categories": map[string]interface{}{
			"Assignments": 0.40,
			"Midterm":     0.30,
			"Final":       0.30,
		},
	}

	categoryScores := db.CategoryScores{
		"Assignments": db.CategoryScore{Score: 85.0},
		"Midterm":     db.CategoryScore{Score: 78.0},
		"Final":       db.CategoryScore{Score: 92.0},
	}

	result, err := gc.calculateCustomFormula(formula, categoryScores, 1)
	if err != nil {
		t.Fatalf("Failed to calculate weighted average: %v", err)
	}

	expected := (0.40 * 85.0) + (0.30 * 78.0) + (0.30 * 92.0)
	if result != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, result)
	} else {
		t.Logf("✓ Weighted average correct: %.2f", result)
	}
}
