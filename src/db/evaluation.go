package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type EvaluationCategory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CourseID     uint      `gorm:"not null" json:"course_id"`
	Course       Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	MaxScore     float64   `gorm:"type:decimal(10,2);not null;default:20.00" json:"max_score"`
	Cardinality  int       `gorm:"default:1" json:"cardinality"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    *string   `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy    *string   `gorm:"size:500" json:"updated_by,omitempty"`
}

func (EvaluationCategory) TableName() string {
	return "EvaluationCategory"
}

type FormulaConfig map[string]interface{}

func (fc FormulaConfig) Value() (driver.Value, error) {
	return json.Marshal(fc)
}

func (fc *FormulaConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, fc)
}

type GradingFormula struct {
	ID                uint          `gorm:"primaryKey" json:"id"`
	CourseID          uint          `gorm:"not null" json:"course_id"`
	Course            Course        `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	LevelID           *uint         `gorm:"column:level_id" json:"level_id,omitempty"`
	FormulaType       string        `gorm:"size:50;not null" json:"formula_type"`
	FormulaConfig     FormulaConfig `gorm:"type:json" json:"formula_config,omitempty"`
	PassingPercentage *float64      `gorm:"type:decimal(5,2)" json:"passing_percentage,omitempty"`
	IsActive          bool          `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	CreatedBy         *string       `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy         *string       `gorm:"size:500" json:"updated_by,omitempty"`
}

func (GradingFormula) TableName() string {
	return "EvaluationGradingFormula"
}

type Grade struct {
	ID             uint               `gorm:"primaryKey" json:"id"`
	CategoryID     uint               `gorm:"not null" json:"category_id"`
	Category       EvaluationCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	StudentClassID uint               `gorm:"not null" json:"student_class_id"`
	StudentClass   StudentClass       `gorm:"foreignKey:StudentClassID" json:"student_class,omitempty"`
	RegistrationID uint               `gorm:"not null" json:"registration_id"`
	Name           string             `gorm:"size:200;not null" json:"name"`
	Description    *string            `gorm:"type:text" json:"description,omitempty"`
	Date           *time.Time         `gorm:"type:date" json:"date,omitempty"`
	DueDate        *time.Time         `gorm:"type:datetime" json:"due_date,omitempty"`
	MaxScore       float64            `gorm:"type:decimal(10,2);not null" json:"max_score"`
	Score          *float64           `gorm:"type:decimal(10,2)" json:"score,omitempty"`
	Percentage     *float64           `gorm:"type:decimal(5,2)" json:"percentage,omitempty"`
	IsExcused      bool               `gorm:"default:false" json:"is_excused"`
	IsLate         bool               `gorm:"default:false" json:"is_late"`
	LatePenalty    float64            `gorm:"type:decimal(5,2);default:0.00" json:"late_penalty"`
	Comments       *string            `gorm:"type:text" json:"comments,omitempty"`
	GradedBy       *string            `gorm:"size:500" json:"graded_by,omitempty"`
	GradedAt       *time.Time         `json:"graded_at,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	CreatedBy      *string            `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy      *string            `gorm:"size:500" json:"updated_by,omitempty"`
}

func (Grade) TableName() string {
	return "EvaluationGrade"
}

type GradeHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	GradeID      uint      `gorm:"not null" json:"grade_id"`
	Grade        Grade     `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	OldScore     *float64  `gorm:"type:decimal(10,2)" json:"old_score,omitempty"`
	NewScore     *float64  `gorm:"type:decimal(10,2)" json:"new_score,omitempty"`
	ChangeReason *string   `gorm:"type:text" json:"change_reason,omitempty"`
	ChangedBy    string    `gorm:"size:500;not null" json:"changed_by"`
	ChangedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"changed_at"`
}

func (GradeHistory) TableName() string {
	return "GradeHistory"
}

type CategoryScores map[string]CategoryScore

type CategoryScore struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Score        float64 `json:"score"`
	Percentage   float64 `json:"percentage"`
}

func (cs CategoryScores) Value() (driver.Value, error) {
	return json.Marshal(cs)
}

func (cs *CategoryScores) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, cs)
}

type FinalGrade struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	RegistrationID       uint           `gorm:"not null" json:"registration_id"`
	StudentClassID       uint           `gorm:"not null" json:"student_class_id"`
	StudentClass         StudentClass   `gorm:"foreignKey:StudentClassID" json:"student_class,omitempty"`
	CalculatedScore      *float64       `gorm:"type:decimal(10,2)" json:"calculated_score,omitempty"`
	CalculatedPercentage *float64       `gorm:"type:decimal(5,2)" json:"calculated_percentage,omitempty"`
	IsPassing            *bool          `json:"is_passing,omitempty"`
	CategoryScores       CategoryScores `gorm:"type:json" json:"category_scores,omitempty"`
	CalculationDate      *time.Time     `json:"calculation_date,omitempty"`
	OverrideScore        *float64       `gorm:"type:decimal(10,2)" json:"override_score,omitempty"`
	OverrideReason       *string        `gorm:"type:text" json:"override_reason,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	CreatedBy            *string        `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy            *string        `gorm:"size:500" json:"updated_by,omitempty"`
}

func (FinalGrade) TableName() string {
	return "FinalGrade"
}

type GradeRubric struct {
	ID                   uint                `gorm:"primaryKey" json:"id"`
	EvaluationCategoryID *uint               `json:"evaluation_category_id,omitempty"`
	EvaluationCategory   *EvaluationCategory `gorm:"foreignKey:EvaluationCategoryID" json:"evaluation_category,omitempty"`
	GradeID              *uint               `json:"grade_id,omitempty"`
	Grade                *Grade              `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	CriteriaName         string              `gorm:"size:200;not null" json:"criteria_name"`
	Description          *string             `gorm:"type:text" json:"description,omitempty"`
	MaxPoints            float64             `gorm:"type:decimal(10,2);not null" json:"max_points"`
	DisplayOrder         int                 `gorm:"default:0" json:"display_order"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

func (GradeRubric) TableName() string {
	return "GradeRubric"
}

type GradeRubricScore struct {
	ID       uint        `gorm:"primaryKey" json:"id"`
	GradeID  uint        `gorm:"not null" json:"grade_id"`
	Grade    Grade       `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	RubricID uint        `gorm:"not null" json:"rubric_id"`
	Rubric   GradeRubric `gorm:"foreignKey:RubricID" json:"rubric,omitempty"`
	Score    float64     `gorm:"type:decimal(10,2);not null" json:"score"`
	Comments *string     `gorm:"type:text" json:"comments,omitempty"`
}

func (GradeRubricScore) TableName() string {
	return "GradeRubricScore"
}
