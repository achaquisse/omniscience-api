package services

import (
	"errors"
	"fmt"
	"math"
	"skulla-api/db"
)

type GradingCalculator struct{}

func NewGradingCalculator() *GradingCalculator {
	return &GradingCalculator{}
}

func (gc *GradingCalculator) CalculatePercentage(score, maxScore float64) float64 {
	if maxScore == 0 {
		return 0
	}
	return math.Round((score/maxScore)*10000) / 100
}

type CategoryWeightConfig struct {
	Weight float64 `json:"weight"`
	Count  int     `json:"count"`
}

type StageConfig struct {
	Type       string                 `json:"type"`
	Categories map[string]interface{} `json:"categories,omitempty"`
	Components []string               `json:"components,omitempty"`
}

type LevelConfig struct {
	Stages map[string]StageConfig `json:"stages"`
}

func (gc *GradingCalculator) calculateCustomFormula(
	config db.FormulaConfig,
	categoryScores db.CategoryScores,
	registrationID uint,
) (float64, error) {
	formulaType, ok := config["type"].(string)
	if !ok {
		return 0, errors.New("invalid formula config: missing or invalid type")
	}

	switch formulaType {
	case "multi_stage":
		return gc.calculateMultiStage(config, categoryScores)
	case "multi_stage_level_based":
		var registration db.Registration
		err := db.DB().First(&registration, registrationID).Error
		if err != nil {
			return 0, err
		}

		level := registration.StudentClass.LevelId

		return gc.calculateMultiStageLevelBased(config, categoryScores, level)
	case "weighted_average":
		return gc.calculateWeightedAverage(config, categoryScores)
	default:
		return 0, fmt.Errorf("unsupported formula type: %s", formulaType)
	}
}

func (gc *GradingCalculator) calculateMultiStage(
	config db.FormulaConfig,
	categoryScores db.CategoryScores,
) (float64, error) {
	stagesRaw, ok := config["stages"]
	if !ok {
		return 0, errors.New("multi_stage formula missing 'stages'")
	}

	stages, ok := stagesRaw.(map[string]interface{})
	if !ok {
		return 0, errors.New("invalid stages format")
	}

	stageResults := make(map[string]float64)

	for stageName, stageDataRaw := range stages {
		stageData, ok := stageDataRaw.(map[string]interface{})
		if !ok {
			continue
		}

		stageType, _ := stageData["type"].(string)

		switch stageType {
		case "weighted_average":
			categoriesRaw, _ := stageData["categories"].(map[string]interface{})
			score := gc.calculateWeightedAverageFromCategories(categoriesRaw, categoryScores)
			stageResults[stageName] = score

		case "simple_average":
			componentsRaw, _ := stageData["components"].([]interface{})
			var sum float64
			var count int

			for _, compRaw := range componentsRaw {
				compName, _ := compRaw.(string)
				if val, exists := stageResults[compName]; exists {
					sum += val
					count++
				} else if catScore, exists := categoryScores[compName]; exists {
					sum += catScore.Percentage
					count++
				}
			}

			if count > 0 {
				stageResults[stageName] = sum / float64(count)
			}
		}
	}

	if mf, exists := stageResults["MF"]; exists {
		return mf, nil
	}

	return 0, errors.New("MF stage not found in results")
}

func (gc *GradingCalculator) calculateMultiStageLevelBased(
	config db.FormulaConfig,
	categoryScores db.CategoryScores,
	level int,
) (float64, error) {
	levelsRaw, ok := config["levels"]
	if !ok {
		return 0, errors.New("multi_stage_level_based formula missing 'levels'")
	}

	levels, ok := levelsRaw.(map[string]interface{})
	if !ok {
		return 0, errors.New("invalid levels format")
	}

	levelKey := fmt.Sprintf("%d", level)
	levelConfigRaw, ok := levels[levelKey]
	if !ok {
		return 0, fmt.Errorf("level %d not found in formula config", level)
	}

	levelConfig, ok := levelConfigRaw.(map[string]interface{})
	if !ok {
		return 0, errors.New("invalid level config format")
	}

	stagesRaw, ok := levelConfig["stages"]
	if !ok {
		return 0, errors.New("level config missing 'stages'")
	}

	return gc.calculateMultiStage(map[string]interface{}{
		"stages": stagesRaw,
	}, categoryScores)
}

func (gc *GradingCalculator) calculateWeightedAverage(
	config db.FormulaConfig,
	categoryScores db.CategoryScores,
) (float64, error) {
	categoriesRaw, ok := config["categories"]
	if !ok {
		return 0, errors.New("weighted_average formula missing 'categories'")
	}

	categories, ok := categoriesRaw.(map[string]interface{})
	if !ok {
		return 0, errors.New("invalid categories format")
	}

	return gc.calculateWeightedAverageFromCategories(categories, categoryScores), nil
}

func (gc *GradingCalculator) calculateWeightedAverageFromCategories(
	weightConfig map[string]interface{},
	categoryScores db.CategoryScores,
) float64 {
	var weightedSum float64
	var totalWeight float64

	for catName, weightRaw := range weightConfig {
		catScore, exists := categoryScores[catName]
		if !exists {
			continue
		}

		switch w := weightRaw.(type) {
		case float64:
			weightedSum += catScore.Percentage * w
			totalWeight += w
		case map[string]interface{}:
			weight, _ := w["weight"].(float64)
			weightedSum += catScore.Percentage * weight
			totalWeight += weight
		}
	}

	if totalWeight > 0 {
		return weightedSum / totalWeight
	}
	return 0
}

func (gc *GradingCalculator) calculateCategoryScore(grades []db.Grade, category db.EvaluationCategory) (float64, float64) {
	if category.Cardinality <= 0 {
		category.Cardinality = 1
	}

	rawScores := make([]float64, 0, category.Cardinality)

	for _, grade := range grades {
		if grade.Score != nil && !grade.IsExcused {
			rawScores = append(rawScores, *grade.Score)
		}
	}

	missingGrades := category.Cardinality - len(rawScores)
	if missingGrades > 0 {
		for i := 0; i < missingGrades; i++ {
			rawScores = append(rawScores, 0)
		}
	}

	if len(rawScores) == 0 {
		return 0, 0
	}

	var sum float64
	for _, score := range rawScores {
		sum += score
	}

	avgRawScore := sum / float64(len(rawScores))

	var percentage float64
	if category.MaxScore > 0 {
		percentage = (avgRawScore / category.MaxScore) * 100
	}

	return avgRawScore, percentage
}

type StageComponentResult struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Weight float64 `json:"weight,omitempty"`
}

type StageResult struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Score      float64                `json:"score"`
	Components []StageComponentResult `json:"components"`
}

type FormulaDetailedResult struct {
	FinalScore        float64       `json:"final_score"`
	IsPassing         *bool         `json:"is_passing,omitempty"`
	PassingPercentage *float64      `json:"passing_percentage,omitempty"`
	Stages            []StageResult `json:"stages"`
}

func (gc *GradingCalculator) CalculateFormulaDetailed(
	registrationID uint,
	studentClassID uint,
	courseID uint,
) (*FormulaDetailedResult, error) {
	var studentClass db.StudentClass
	err := db.DB().First(&studentClass, studentClassID).Error
	if err != nil {
		return nil, errors.New("student class not found")
	}

	levelID := studentClass.LevelId

	var formula db.GradingFormula
	query := db.DB().Where("course_id = ? AND is_active = ?", courseID, true)
	if levelID > 0 {
		query = query.Where("level_id = ?", levelID)
	} else {
		query = query.Where("level_id IS NULL")
	}
	err = query.First(&formula).Error
	if err != nil {
		return nil, errors.New("grading formula not found for course and level")
	}

	var categories []db.EvaluationCategory
	err = db.DB().Where("course_id = ? AND is_active = ?", courseID, true).
		Order("display_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, errors.New("no evaluation categories found for course")
	}

	categoryScores := make(db.CategoryScores)
	for _, category := range categories {
		var grades []db.Grade
		err = db.DB().Where("category_id = ? AND student_class_id = ? AND registration_id = ?",
			category.ID, studentClassID, registrationID).Find(&grades).Error
		if err != nil {
			continue
		}

		rawScore, percentage := gc.calculateCategoryScore(grades, category)
		categoryScores[category.Name] = db.CategoryScore{
			CategoryID:   category.ID,
			CategoryName: category.Name,
			Score:        rawScore,
			Percentage:   percentage,
		}
	}

	result := &FormulaDetailedResult{}

	if formula.FormulaType == "CUSTOM" && formula.FormulaConfig != nil {
		stages, finalPercentage, err := gc.calculateMultiStageDetailed(formula.FormulaConfig, categoryScores, registrationID)
		if err != nil {
			return nil, err
		}
		result.FinalScore = math.Round((finalPercentage*0.2)*100) / 100
		result.Stages = stages
	} else {
		var totalScore float64
		var count int
		var components []StageComponentResult
		for name, catScore := range categoryScores {
			totalScore += catScore.Percentage
			count++
			components = append(components, StageComponentResult{
				Name:  name,
				Score: math.Round(catScore.Score*100) / 100,
			})
		}
		var avgPercentage float64
		if count > 0 {
			avgPercentage = totalScore / float64(count)
			result.FinalScore = math.Round((avgPercentage*0.2)*100) / 100
		}
		result.Stages = []StageResult{{
			Name:       "MF",
			Type:       "simple_average",
			Score:      result.FinalScore,
			Components: components,
		}}
	}

	if formula.PassingPercentage != nil {
		result.PassingPercentage = formula.PassingPercentage
		passingScore := *formula.PassingPercentage * 0.2
		isPassing := result.FinalScore >= passingScore
		result.IsPassing = &isPassing
	}

	return result, nil
}

func (gc *GradingCalculator) calculateMultiStageDetailed(
	config db.FormulaConfig,
	categoryScores db.CategoryScores,
	registrationID uint,
) ([]StageResult, float64, error) {
	formulaType, ok := config["type"].(string)
	if !ok {
		return nil, 0, errors.New("invalid formula config: missing or invalid type")
	}

	switch formulaType {
	case "multi_stage":
		return gc.calculateMultiStageStages(config, categoryScores)
	case "multi_stage_level_based":
		var registration db.Registration
		err := db.DB().First(&registration, registrationID).Error
		if err != nil {
			return nil, 0, err
		}
		level := registration.StudentClass.LevelId

		levelsRaw, ok := config["levels"]
		if !ok {
			return nil, 0, errors.New("multi_stage_level_based formula missing 'levels'")
		}
		levels, ok := levelsRaw.(map[string]interface{})
		if !ok {
			return nil, 0, errors.New("invalid levels format")
		}
		levelKey := fmt.Sprintf("%d", level)
		levelConfigRaw, ok := levels[levelKey]
		if !ok {
			return nil, 0, fmt.Errorf("level %d not found in formula config", level)
		}
		levelConfig, ok := levelConfigRaw.(map[string]interface{})
		if !ok {
			return nil, 0, errors.New("invalid level config format")
		}
		stagesRaw, ok := levelConfig["stages"]
		if !ok {
			return nil, 0, errors.New("level config missing 'stages'")
		}
		return gc.calculateMultiStageStages(map[string]interface{}{
			"stages": stagesRaw,
		}, categoryScores)
	case "weighted_average":
		categoriesRaw, ok := config["categories"]
		if !ok {
			return nil, 0, errors.New("weighted_average formula missing 'categories'")
		}
		cats, ok := categoriesRaw.(map[string]interface{})
		if !ok {
			return nil, 0, errors.New("invalid categories format")
		}
		components, percentage := gc.weightedAverageComponents(cats, categoryScores)
		stage := StageResult{
			Name:       "MF",
			Type:       "weighted_average",
			Score:      math.Round((percentage*0.2)*100) / 100,
			Components: components,
		}
		return []StageResult{stage}, percentage, nil
	default:
		return nil, 0, fmt.Errorf("unsupported formula type: %s", formulaType)
	}
}

func (gc *GradingCalculator) calculateMultiStageStages(
	config db.FormulaConfig,
	categoryScores db.CategoryScores,
) ([]StageResult, float64, error) {
	stagesRaw, ok := config["stages"]
	if !ok {
		return nil, 0, errors.New("multi_stage formula missing 'stages'")
	}

	stages, ok := stagesRaw.(map[string]interface{})
	if !ok {
		return nil, 0, errors.New("invalid stages format")
	}

	stageScores := make(map[string]float64)
	var stageResults []StageResult

	for stageName, stageDataRaw := range stages {
		stageData, ok := stageDataRaw.(map[string]interface{})
		if !ok {
			continue
		}

		stageType, _ := stageData["type"].(string)

		switch stageType {
		case "weighted_average":
			categoriesRaw, _ := stageData["categories"].(map[string]interface{})
			components, percentage := gc.weightedAverageComponents(categoriesRaw, categoryScores)
			stageScores[stageName] = percentage
			stageResults = append(stageResults, StageResult{
				Name:       stageName,
				Type:       stageType,
				Score:      math.Round((percentage*0.2)*100) / 100,
				Components: components,
			})

		case "simple_average":
			componentsRaw, _ := stageData["components"].([]interface{})
			var sum float64
			var count int
			var components []StageComponentResult

			for _, compRaw := range componentsRaw {
				compName, _ := compRaw.(string)
				if val, exists := stageScores[compName]; exists {
					sum += val
					count++
					components = append(components, StageComponentResult{
						Name:  compName,
						Score: math.Round((val*0.2)*100) / 100,
					})
				} else if catScore, exists := categoryScores[compName]; exists {
					sum += catScore.Percentage
					count++
					components = append(components, StageComponentResult{
						Name:  compName,
						Score: math.Round(catScore.Score*100) / 100,
					})
				}
			}

			var percentage float64
			if count > 0 {
				percentage = sum / float64(count)
			}
			stageScores[stageName] = percentage
			stageResults = append(stageResults, StageResult{
				Name:       stageName,
				Type:       stageType,
				Score:      math.Round((percentage*0.2)*100) / 100,
				Components: components,
			})
		}
	}

	if mf, exists := stageScores["MF"]; exists {
		return stageResults, mf, nil
	}

	return stageResults, 0, errors.New("MF stage not found in results")
}

func (gc *GradingCalculator) weightedAverageComponents(
	weightConfig map[string]interface{},
	categoryScores db.CategoryScores,
) ([]StageComponentResult, float64) {
	var weightedSum float64
	var totalWeight float64
	var components []StageComponentResult

	for catName, weightRaw := range weightConfig {
		catScore, exists := categoryScores[catName]
		if !exists {
			continue
		}

		var weight float64
		switch w := weightRaw.(type) {
		case float64:
			weight = w
		case map[string]interface{}:
			weight, _ = w["weight"].(float64)
		}

		weightedSum += catScore.Percentage * weight
		totalWeight += weight
		components = append(components, StageComponentResult{
			Name:   catName,
			Score:  math.Round(catScore.Score*100) / 100,
			Weight: weight,
		})
	}

	var score float64
	if totalWeight > 0 {
		score = weightedSum / totalWeight
	}
	return components, score
}

type CategoryGradeStatistics struct {
	CategoryID      uint    `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	Cardinality     int     `json:"cardinality"`
	RecordedGrades  int     `json:"recorded_grades"`
	ExcusedGrades   int     `json:"excused_grades"`
	MissingGrades   int     `json:"missing_grades"`
	AverageScore    float64 `json:"average_score"`
	IncludesMissing bool    `json:"includes_missing"`
}

func (gc *GradingCalculator) GetCategoryStatistics(
	categoryID uint,
	studentClassID uint,
	registrationID uint,
) (*CategoryGradeStatistics, error) {
	var category db.EvaluationCategory
	err := db.DB().First(&category, categoryID).Error
	if err != nil {
		return nil, errors.New("category not found")
	}

	var grades []db.Grade
	err = db.DB().Where("category_id = ? AND student_class_id = ? AND registration_id = ?",
		categoryID, studentClassID, registrationID).Find(&grades).Error
	if err != nil {
		return nil, err
	}

	stats := &CategoryGradeStatistics{
		CategoryID:   category.ID,
		CategoryName: category.Name,
		Cardinality:  category.Cardinality,
	}

	if stats.Cardinality <= 0 {
		stats.Cardinality = 1
	}

	var recordedCount int
	var excusedCount int
	var sum float64

	for _, grade := range grades {
		if grade.IsExcused {
			excusedCount++
		} else if grade.Score != nil && grade.Percentage != nil {
			recordedCount++
			sum += *grade.Score
		}
	}

	stats.RecordedGrades = recordedCount
	stats.ExcusedGrades = excusedCount
	stats.MissingGrades = stats.Cardinality - recordedCount

	if stats.MissingGrades < 0 {
		stats.MissingGrades = 0
	}

	stats.IncludesMissing = stats.MissingGrades > 0

	totalGrades := recordedCount + stats.MissingGrades
	if totalGrades > 0 {
		stats.AverageScore = sum / float64(totalGrades)
	}

	return stats, nil
}

func (gc *GradingCalculator) GetAllCategoryStatistics(
	studentClassID uint,
	registrationID uint,
	courseID uint,
) ([]CategoryGradeStatistics, error) {
	var categories []db.EvaluationCategory
	err := db.DB().Where("course_id = ? AND is_active = ?", courseID, true).
		Order("display_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	var allStats []CategoryGradeStatistics

	for _, category := range categories {
		stats, err := gc.GetCategoryStatistics(category.ID, studentClassID, registrationID)
		if err != nil {
			continue
		}
		allStats = append(allStats, *stats)
	}

	return allStats, nil
}
