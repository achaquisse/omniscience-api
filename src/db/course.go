package db

import "fmt"

type Course struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:255;not null"`
	TeacherEmail string `gorm:"column:teacher_email;size:255"`
}

func (Course) TableName() string {
	return "Course"
}

func ListCoursesByTeacherEmail(email string) []Course {
	var courses []Course
	db.
		Where("teacher_email LIKE ?", fmt.Sprintf("%%%s%%", email)).
		Find(&courses)
	return courses
}

func IsTeacherEmailBelongToCourse(email string, courseId int) bool {
	var courses []Course
	db.
		Where("teacher_email LIKE ?", fmt.Sprintf("%%%s%%", email)).
		Where("id = ?", courseId).
		Find(&courses)
	return len(courses) > 0
}
