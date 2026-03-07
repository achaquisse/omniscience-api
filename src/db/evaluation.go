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
	LevelID      *uint     `gorm:"column:level_id" json:"level_id,omitempty"`
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

func FindActiveCategories(courseID uint, levelId int) []EvaluationCategory {
	var categories []EvaluationCategory
	query := db.Where("course_id = ? AND is_active = ?", courseID, true)
	if levelId > 0 {
		query = query.Where("level_id = ?", levelId)
	} else {
		query = query.Where("level_id IS NULL")
	}
	query.Order("display_order ASC").Find(&categories)
	return categories
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

type CategoryScores map[string]CategoryScore

type CategoryScore struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Score        float64 `json:"score"`
	Percentage   float64 `json:"percentage"`
}
