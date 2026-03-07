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
		"Avaliação escrita": db.CategoryScore{Score: 85.0, Percentage: 85.0},
		"Ditados":           db.CategoryScore{Score: 90.0, Percentage: 90.0},
		"Composições":       db.CategoryScore{Score: 75.0, Percentage: 75.0},
		"Speeches":          db.CategoryScore{Score: 78.0, Percentage: 78.0},
		"Listening":         db.CategoryScore{Score: 82.0, Percentage: 82.0},
		"Exame":             db.CategoryScore{Score: 88.0, Percentage: 88.0},
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
		"Avaliação escrita": db.CategoryScore{Score: 85.0, Percentage: 85.0},
		"Ditados":           db.CategoryScore{Score: 90.0, Percentage: 90.0},
		"Composições":       db.CategoryScore{Score: 75.0, Percentage: 75.0},
		"Speeches":          db.CategoryScore{Score: 78.0, Percentage: 78.0},
		"Listening":         db.CategoryScore{Score: 82.0, Percentage: 82.0},
		"Exame":             db.CategoryScore{Score: 88.0, Percentage: 88.0},
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
		"Assignments": db.CategoryScore{Score: 85.0, Percentage: 85.0},
		"Midterm":     db.CategoryScore{Score: 78.0, Percentage: 78.0},
		"Final":       db.CategoryScore{Score: 92.0, Percentage: 92.0},
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

func TestCalculateCategoryScore_AllGradesPresent(t *testing.T) {
	gc := NewGradingCalculator()

	category := db.EvaluationCategory{ID: 1, Name: "Quiz", MaxScore: 20.0, Cardinality: 3}

	s1, s2, s3 := 15.0, 18.0, 12.0
	grades := []db.Grade{
		{Score: &s1},
		{Score: &s2},
		{Score: &s3},
	}

	rawScore, percentage := gc.calculateCategoryScore(grades, category)

	expectedAvg := (15.0 + 18.0 + 12.0) / 3.0
	expectedPct := (expectedAvg / 20.0) * 100

	if rawScore != expectedAvg {
		t.Errorf("Expected rawScore = %.2f, got %.2f", expectedAvg, rawScore)
	}
	if percentage != expectedPct {
		t.Errorf("Expected percentage = %.2f, got %.2f", expectedPct, percentage)
	} else {
		t.Logf("✓ Category score correct: avg=%.2f, pct=%.2f%%", rawScore, percentage)
	}
}

func TestCalculateCategoryScore_MissingGradesPaddedWithZero(t *testing.T) {
	gc := NewGradingCalculator()

	category := db.EvaluationCategory{ID: 1, Name: "Quiz", MaxScore: 20.0, Cardinality: 4}

	s1, s2 := 16.0, 12.0
	grades := []db.Grade{
		{Score: &s1},
		{Score: &s2},
	}

	rawScore, percentage := gc.calculateCategoryScore(grades, category)

	expectedAvg := (16.0 + 12.0 + 0.0 + 0.0) / 4.0
	expectedPct := (expectedAvg / 20.0) * 100

	if rawScore != expectedAvg {
		t.Errorf("Expected rawScore = %.2f, got %.2f", expectedAvg, rawScore)
	}
	if percentage != expectedPct {
		t.Errorf("Expected percentage = %.2f, got %.2f", expectedPct, percentage)
	} else {
		t.Logf("✓ Missing grades padded correctly: avg=%.2f, pct=%.2f%%", rawScore, percentage)
	}
}

func TestCalculateCategoryScore_ExcusedGradesIgnored(t *testing.T) {
	gc := NewGradingCalculator()

	category := db.EvaluationCategory{ID: 1, Name: "Quiz", MaxScore: 20.0, Cardinality: 3}

	s1, s2 := 18.0, 14.0
	grades := []db.Grade{
		{Score: &s1},
		{Score: &s2, IsExcused: true},
	}

	rawScore, percentage := gc.calculateCategoryScore(grades, category)

	expectedAvg := (18.0 + 0.0 + 0.0) / 3.0
	expectedPct := (expectedAvg / 20.0) * 100

	if rawScore != expectedAvg {
		t.Errorf("Expected rawScore = %.2f, got %.2f", expectedAvg, rawScore)
	}
	if percentage != expectedPct {
		t.Errorf("Expected percentage = %.2f, got %.2f", expectedPct, percentage)
	} else {
		t.Logf("✓ Excused grades ignored correctly: avg=%.2f, pct=%.2f%%", rawScore, percentage)
	}
}

func TestCalculateCategoryScore_NoGrades(t *testing.T) {
	gc := NewGradingCalculator()

	category := db.EvaluationCategory{ID: 1, Name: "Quiz", MaxScore: 20.0, Cardinality: 2}

	rawScore, percentage := gc.calculateCategoryScore([]db.Grade{}, category)

	if rawScore != 0.0 {
		t.Errorf("Expected rawScore = 0.0, got %.2f", rawScore)
	}
	if percentage != 0.0 {
		t.Errorf("Expected percentage = 0.0, got %.2f", percentage)
	} else {
		t.Logf("✓ No grades yields zero score correctly")
	}
}

func TestFrequencyGradeCalculation_MultiStage(t *testing.T) {
	gc := NewGradingCalculator()

	formula := db.FormulaConfig{
		"type": "multi_stage",
		"stages": map[string]interface{}{
			"MDF": map[string]interface{}{
				"type": "weighted_average",
				"categories": map[string]interface{}{
					"Quiz":       0.50,
					"Assignment": 0.50,
				},
			},
			"MF": map[string]interface{}{
				"type":       "simple_average",
				"components": []interface{}{"MDF", "Exame"},
			},
		},
	}

	quizCategory := db.EvaluationCategory{ID: 1, Name: "Quiz", MaxScore: 20.0, Cardinality: 2}
	assignCategory := db.EvaluationCategory{ID: 2, Name: "Assignment", MaxScore: 20.0, Cardinality: 1}
	examCategory := db.EvaluationCategory{ID: 3, Name: "Exame", MaxScore: 20.0, Cardinality: 1}

	q1, q2 := 14.0, 16.0
	quizGrades := []db.Grade{{Score: &q1}, {Score: &q2}}

	a1 := 18.0
	assignGrades := []db.Grade{{Score: &a1}}

	e1 := 12.0
	examGrades := []db.Grade{{Score: &e1}}

	quizRaw, quizPct := gc.calculateCategoryScore(quizGrades, quizCategory)
	assignRaw, assignPct := gc.calculateCategoryScore(assignGrades, assignCategory)
	examRaw, examPct := gc.calculateCategoryScore(examGrades, examCategory)

	categoryScores := db.CategoryScores{
		"Quiz":       {CategoryID: 1, CategoryName: "Quiz", Score: quizRaw, Percentage: quizPct},
		"Assignment": {CategoryID: 2, CategoryName: "Assignment", Score: assignRaw, Percentage: assignPct},
		"Exame":      {CategoryID: 3, CategoryName: "Exame", Score: examRaw, Percentage: examPct},
	}

	result, err := gc.calculateCustomFormula(formula, categoryScores, 1)
	if err != nil {
		t.Fatalf("Failed to calculate formula: %v", err)
	}

	expectedQuizPct := ((14.0 + 16.0) / 2.0 / 20.0) * 100
	expectedAssignPct := (18.0 / 20.0) * 100
	expectedExamPct := (12.0 / 20.0) * 100
	expectedMDF := (0.50*expectedQuizPct + 0.50*expectedAssignPct)
	expectedMF := (expectedMDF + expectedExamPct) / 2.0

	tolerance := 0.001
	if result < expectedMF-tolerance || result > expectedMF+tolerance {
		t.Errorf("Expected MF = %.4f, got %.4f", expectedMF, result)
	} else {
		t.Logf("✓ Full grading pipeline correct: QuizPct=%.2f%%, AssignPct=%.2f%%, ExamPct=%.2f%%, MDF=%.2f%%, MF=%.2f%%",
			expectedQuizPct, expectedAssignPct, expectedExamPct, expectedMDF, result)
	}
}

func TestFinalGradeIncreasesAsGradesAreAdded(t *testing.T) {
	gc := NewGradingCalculator()

	formula := db.FormulaConfig{
		"type": "weighted_average",
		"categories": map[string]interface{}{
			"Quiz":  0.60,
			"Exame": 0.40,
		},
	}

	quizCategory := db.EvaluationCategory{ID: 1, Name: "Quiz", MaxScore: 20.0, Cardinality: 3}
	examCategory := db.EvaluationCategory{ID: 2, Name: "Exame", MaxScore: 20.0, Cardinality: 1}

	e1 := 15.0
	examGrades := []db.Grade{{Score: &e1}}
	_, examPct := gc.calculateCategoryScore(examGrades, examCategory)

	q1 := 10.0
	gradesRound1 := []db.Grade{{Score: &q1}}
	_, quizPct1 := gc.calculateCategoryScore(gradesRound1, quizCategory)

	result1, _ := gc.calculateCustomFormula(formula, db.CategoryScores{
		"Quiz":  {Score: q1, Percentage: quizPct1},
		"Exame": {Score: e1, Percentage: examPct},
	}, 1)

	q2 := 18.0
	gradesRound2 := []db.Grade{{Score: &q1}, {Score: &q2}}
	_, quizPct2 := gc.calculateCategoryScore(gradesRound2, quizCategory)

	result2, _ := gc.calculateCustomFormula(formula, db.CategoryScores{
		"Quiz":  {Score: (q1 + q2) / 2.0, Percentage: quizPct2},
		"Exame": {Score: e1, Percentage: examPct},
	}, 1)

	if result2 <= result1 {
		t.Errorf("Final grade should increase as higher grades are added: before=%.4f, after=%.4f", result1, result2)
	} else {
		t.Logf("✓ Final grade increases as grades registered: %.4f -> %.4f", result1, result2)
	}
}

func TestExamGradeImpactOnFinalGrade(t *testing.T) {
	gc := NewGradingCalculator()

	formula := db.FormulaConfig{
		"type": "weighted_average",
		"categories": map[string]interface{}{
			"Quiz":  0.60,
			"Exame": 0.40,
		},
	}

	quizCategory := db.EvaluationCategory{ID: 1, Name: "Quiz", MaxScore: 20.0, Cardinality: 2}
	examCategory := db.EvaluationCategory{ID: 2, Name: "Exame", MaxScore: 20.0, Cardinality: 1}

	q1, q2 := 16.0, 14.0
	_, quizPct := gc.calculateCategoryScore([]db.Grade{{Score: &q1}, {Score: &q2}}, quizCategory)

	lowExam := 8.0
	_, lowExamPct := gc.calculateCategoryScore([]db.Grade{{Score: &lowExam}}, examCategory)

	highExam := 20.0
	_, highExamPct := gc.calculateCategoryScore([]db.Grade{{Score: &highExam}}, examCategory)

	resultLow, _ := gc.calculateCustomFormula(formula, db.CategoryScores{
		"Quiz":  {Score: (q1 + q2) / 2.0, Percentage: quizPct},
		"Exame": {Score: lowExam, Percentage: lowExamPct},
	}, 1)

	resultHigh, _ := gc.calculateCustomFormula(formula, db.CategoryScores{
		"Quiz":  {Score: (q1 + q2) / 2.0, Percentage: quizPct},
		"Exame": {Score: highExam, Percentage: highExamPct},
	}, 1)

	if resultHigh <= resultLow {
		t.Errorf("Higher exam grade should yield higher final grade: low=%.4f, high=%.4f", resultLow, resultHigh)
	} else {
		t.Logf("✓ Exam grade impacts final grade correctly: lowExam=%.4f -> highExam=%.4f", resultLow, resultHigh)
	}
}
