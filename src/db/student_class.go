package db

import (
	"time"
)

type StudentClass struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	CourseID uint   `gorm:"foreignKey:CourseID"`
	Course   Course `gorm:"foreignKey:CourseID"`
	PeriodId uint   `gorm:"foreignKey:PeriodID"`
	Period   Period `gorm:"foreignKey:PeriodID"`
	Disabled bool   `gorm:"column:disabled;type:TINYINT(1);default:0"`
}

func (StudentClass) TableName() string {
	return "StudentClass"
}

type Period struct {
	ID    uint      `gorm:"primaryKey"`
	Start time.Time `gorm:"type:datetime"`
	End   time.Time `gorm:"type:datetime"`
}

func (Period) TableName() string {
	return "Period"
}

func ListStudentClasses(courseIds []uint, startDate *time.Time, endDate *time.Time) []StudentClass {
	var studentClass []StudentClass

	query := db.Preload("Course").Preload("Period")
	query = query.Where("disabled = ?", false).Where("course_id IN ?", courseIds)

	if startDate != nil || endDate != nil {
		query = query.Joins("Period")

		if startDate != nil {
			query = query.Where("Period.end >= ?", startDate)
		}
		if endDate != nil {
			query = query.Where("Period.start <= ?", endDate)
		}
	}

	query.Find(&studentClass)
	return studentClass
}

func GetStudentClassCourseID(studentClassID uint) (uint, error) {
	var studentClass StudentClass
	err := db.First(&studentClass, studentClassID).Error
	if err != nil {
		return 0, err
	}
	return studentClass.CourseID, nil
}
