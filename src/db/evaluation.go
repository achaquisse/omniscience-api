package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type EvaluationCategory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CourseID      uint      `gorm:"not null" json:"course_id"`
	Course        Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Description   *string   `gorm:"type:text" json:"description,omitempty"`
	Weight        float64   `gorm:"type:decimal(5,2);not null" json:"weight"`
	MaxScore      float64   `gorm:"type:decimal(10,2);not null;default:20.00" json:"max_score"`
	DropLowest    int       `gorm:"default:0" json:"drop_lowest"`
	IsExtraCredit bool      `gorm:"default:false" json:"is_extra_credit"`
	DisplayOrder  int       `gorm:"default:0" json:"display_order"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     *string   `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy     *string   `gorm:"size:500" json:"updated_by,omitempty"`
}

func (EvaluationCategory) TableName() string {
	return "EvaluationCategory"
}

type GradingScale map[string]GradeRange

type GradeRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func (gs GradingScale) Value() (driver.Value, error) {
	return json.Marshal(gs)
}

func (gs *GradingScale) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, gs)
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
	FormulaType       string        `gorm:"size:50;not null" json:"formula_type"`
	FormulaConfig     FormulaConfig `gorm:"type:json" json:"formula_config,omitempty"`
	PassingPercentage *float64      `gorm:"type:decimal(5,2)" json:"passing_percentage,omitempty"`
	GradingScale      GradingScale  `gorm:"type:json" json:"grading_scale,omitempty"`
	IsActive          bool          `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	CreatedBy         *string       `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy         *string       `gorm:"size:500" json:"updated_by,omitempty"`
}

func (GradingFormula) TableName() string {
	return "GradingFormula"
}

type EvaluationItem struct {
	ID             uint               `gorm:"primaryKey" json:"id"`
	CategoryID     uint               `gorm:"not null" json:"category_id"`
	Category       EvaluationCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	StudentClassID uint               `gorm:"not null" json:"student_class_id"`
	StudentClass   StudentClass       `gorm:"foreignKey:StudentClassID" json:"student_class,omitempty"`
	Name           string             `gorm:"size:200;not null" json:"name"`
	Description    *string            `gorm:"type:text" json:"description,omitempty"`
	Date           *time.Time         `gorm:"type:date" json:"date,omitempty"`
	DueDate        *time.Time         `gorm:"type:datetime" json:"due_date,omitempty"`
	MaxScore       float64            `gorm:"type:decimal(10,2);not null" json:"max_score"`
	WeightOverride *float64           `gorm:"type:decimal(5,2)" json:"weight_override,omitempty"`
	Status         string             `gorm:"size:20;default:'DRAFT'" json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	CreatedBy      *string            `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy      *string            `gorm:"size:500" json:"updated_by,omitempty"`
}

func (EvaluationItem) TableName() string {
	return "EvaluationItem"
}

type Grade struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	EvaluationItemID uint           `gorm:"not null" json:"evaluation_item_id"`
	EvaluationItem   EvaluationItem `gorm:"foreignKey:EvaluationItemID" json:"evaluation_item,omitempty"`
	RegistrationID   uint           `gorm:"not null" json:"registration_id"`
	Score            *float64       `gorm:"type:decimal(10,2)" json:"score,omitempty"`
	Percentage       *float64       `gorm:"type:decimal(5,2)" json:"percentage,omitempty"`
	LetterGrade      *string        `gorm:"size:5" json:"letter_grade,omitempty"`
	IsExcused        bool           `gorm:"default:false" json:"is_excused"`
	IsLate           bool           `gorm:"default:false" json:"is_late"`
	LatePenalty      float64        `gorm:"type:decimal(5,2);default:0.00" json:"late_penalty"`
	Comments         *string        `gorm:"type:text" json:"comments,omitempty"`
	GradedBy         *string        `gorm:"size:500" json:"graded_by,omitempty"`
	GradedAt         *time.Time     `json:"graded_at,omitempty"`
	Status           string         `gorm:"size:20;default:'DRAFT'" json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	CreatedBy        *string        `gorm:"size:500" json:"created_by,omitempty"`
	UpdatedBy        *string        `gorm:"size:500" json:"updated_by,omitempty"`
}

func (Grade) TableName() string {
	return "Grade"
}

type GradeHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	GradeID      uint      `gorm:"not null" json:"grade_id"`
	Grade        Grade     `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	OldScore     *float64  `gorm:"type:decimal(10,2)" json:"old_score,omitempty"`
	NewScore     *float64  `gorm:"type:decimal(10,2)" json:"new_score,omitempty"`
	OldStatus    *string   `gorm:"size:20" json:"old_status,omitempty"`
	NewStatus    *string   `gorm:"size:20" json:"new_status,omitempty"`
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
	Weight       float64 `json:"weight"`
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
	LetterGrade          *string        `gorm:"size:5" json:"letter_grade,omitempty"`
	IsPassing            *bool          `json:"is_passing,omitempty"`
	CategoryScores       CategoryScores `gorm:"type:json" json:"category_scores,omitempty"`
	CalculationDate      *time.Time     `json:"calculation_date,omitempty"`
	OverrideScore        *float64       `gorm:"type:decimal(10,2)" json:"override_score,omitempty"`
	OverrideReason       *string        `gorm:"type:text" json:"override_reason,omitempty"`
	Status               string         `gorm:"size:20;default:'DRAFT'" json:"status"`
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
	EvaluationItemID     *uint               `json:"evaluation_item_id,omitempty"`
	EvaluationItem       *EvaluationItem     `gorm:"foreignKey:EvaluationItemID" json:"evaluation_item,omitempty"`
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
