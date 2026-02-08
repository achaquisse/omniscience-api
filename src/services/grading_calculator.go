package services

import (
	"errors"
	"math"
	"skulla-api/db"
	"sort"
	"time"
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

func (gc *GradingCalculator) GetLetterGrade(percentage float64, gradingScale db.GradingScale) string {
	for grade, gradeRange := range gradingScale {
		if percentage >= gradeRange.Min && percentage <= gradeRange.Max {
			return grade
		}
	}
	return "N/A"
}

func (gc *GradingCalculator) CalculateFinalGrade(
	registrationID uint,
	studentClassID uint,
	courseID uint,
) (*db.FinalGrade, error) {
	var formula db.GradingFormula
	err := db.DB().Where("course_id = ? AND is_active = ?", courseID, true).First(&formula).Error
	if err != nil {
		return nil, errors.New("grading formula not found for course")
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
	var totalWeightedScore float64
	var totalWeight float64

	for _, category := range categories {
		var evaluationItems []db.EvaluationItem
		err = db.DB().Where("category_id = ? AND student_class_id = ? AND status = ?",
			category.ID, studentClassID, "PUBLISHED").Find(&evaluationItems).Error
		if err != nil {
			continue
		}

		if len(evaluationItems) == 0 {
			continue
		}

		var grades []db.Grade
		for _, item := range evaluationItems {
			var grade db.Grade
			err = db.DB().Where("evaluation_item_id = ? AND registration_id = ? AND status = ?",
				item.ID, registrationID, "PUBLISHED").First(&grade).Error
			if err == nil && !grade.IsExcused {
				grade.EvaluationItem = item
				grades = append(grades, grade)
			}
		}

		if len(grades) == 0 {
			continue
		}

		categoryScore := gc.calculateCategoryScore(grades, category)
		categoryScores[category.Name] = db.CategoryScore{
			CategoryID:   category.ID,
			CategoryName: category.Name,
			Score:        categoryScore,
			Percentage:   categoryScore,
			Weight:       category.Weight,
		}

		if !category.IsExtraCredit {
			totalWeightedScore += categoryScore * (category.Weight / 100)
			totalWeight += category.Weight
		} else {
			totalWeightedScore += categoryScore * (category.Weight / 100)
		}
	}

	var finalPercentage float64
	if totalWeight > 0 {
		finalPercentage = totalWeightedScore
	}

	finalPercentage = math.Round(finalPercentage*100) / 100

	letterGrade := gc.GetLetterGrade(finalPercentage, formula.GradingScale)

	isPassing := false
	if formula.PassingPercentage != nil {
		isPassing = finalPercentage >= *formula.PassingPercentage
	}

	now := time.Now()
	finalGrade := &db.FinalGrade{
		RegistrationID:       registrationID,
		StudentClassID:       studentClassID,
		CalculatedPercentage: &finalPercentage,
		LetterGrade:          &letterGrade,
		IsPassing:            &isPassing,
		CategoryScores:       categoryScores,
		CalculationDate:      &now,
		Status:               "PUBLISHED",
	}

	return finalGrade, nil
}

func (gc *GradingCalculator) calculateCategoryScore(grades []db.Grade, category db.EvaluationCategory) float64 {
	if len(grades) == 0 {
		return 0
	}

	scores := make([]float64, 0, len(grades))
	for _, grade := range grades {
		if grade.Score != nil && grade.Percentage != nil {
			scores = append(scores, *grade.Percentage)
		}
	}

	if len(scores) == 0 {
		return 0
	}

	if category.DropLowest > 0 && len(scores) > category.DropLowest {
		sort.Float64s(scores)
		scores = scores[category.DropLowest:]
	}

	var sum float64
	for _, score := range scores {
		sum += score
	}

	return sum / float64(len(scores))
}

func (gc *GradingCalculator) SaveOrUpdateFinalGrade(finalGrade *db.FinalGrade) error {
	var existing db.FinalGrade
	err := db.DB().Where("registration_id = ? AND student_class_id = ?",
		finalGrade.RegistrationID, finalGrade.StudentClassID).First(&existing).Error

	if err != nil {
		return db.DB().Create(finalGrade).Error
	}

	finalGrade.ID = existing.ID
	return db.DB().Save(finalGrade).Error
}

func (gc *GradingCalculator) RecalculateAllFinalGrades(studentClassID uint) error {
	var registrations []db.Registration
	err := db.DB().Where("student_class_id = ? AND status = ?", studentClassID, "ACTIVE").
		Find(&registrations).Error
	if err != nil {
		return err
	}

	courseID, err := db.GetStudentClassCourseID(studentClassID)
	if err != nil {
		return err
	}

	for _, registration := range registrations {
		finalGrade, err := gc.CalculateFinalGrade(registration.ID, studentClassID, courseID)
		if err != nil {
			continue
		}

		err = gc.SaveOrUpdateFinalGrade(finalGrade)
		if err != nil {
			continue
		}
	}

	return nil
}
