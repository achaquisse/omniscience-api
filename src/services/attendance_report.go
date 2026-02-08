package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"skulla-api/db"
)

type StudentAttendanceInfo struct {
	RegistrationID uint
	StudentID      uint
	StudentName    string
	PhoneNumber    string
	CourseName     string
	Percentage     float64
}

func GetActiveRegistrationsWithAttendance(startDate, endDate time.Time) ([]StudentAttendanceInfo, error) {
	var registrations []db.Registration
	now := time.Now()

	err := db.DB().
		Preload("Student").
		Joins("JOIN Student ON Student.id = Registration.student_id").
		Joins("JOIN StudentClass ON StudentClass.id = Registration.student_class_id").
		Joins("JOIN Period ON Period.id = StudentClass.period_id").
		Joins("JOIN Course ON Course.id = StudentClass.course_id").
		Where("Registration.status = ?", "ACTIVE").
		Where("Period.start <= ?", now).
		Where("Period.end >= ?", now).
		Find(&registrations).Error

	if err != nil {
		return nil, err
	}

	var results []StudentAttendanceInfo

	for _, reg := range registrations {
		var studentClass db.StudentClass
		err := db.DB().Preload("Course").First(&studentClass, reg.StudentClassID).Error
		if err != nil {
			continue
		}

		var attendances []db.Attendance
		err = db.DB().
			Where("registration_id = ?", reg.ID).
			Where("date >= ?", startDate.Format("2006-01-02")).
			Where("date <= ?", endDate.Format("2006-01-02")).
			Find(&attendances).Error

		if err != nil {
			continue
		}

		totalDays := len(attendances)
		presentCount := 0

		for _, att := range attendances {
			if att.Status == "PRESENT" {
				presentCount++
			}
		}

		var percentage float64
		if totalDays > 0 {
			percentage = math.Round((float64(presentCount) / float64(totalDays)) * 100)
		}

		if reg.Student.FirstName == "" && reg.Student.LastName == "" {
			continue
		}

		if reg.Student.Phone == "" {
			continue
		}

		studentName := strings.TrimSpace(reg.Student.FirstName + " " + reg.Student.LastName)

		results = append(results, StudentAttendanceInfo{
			RegistrationID: reg.ID,
			StudentID:      reg.StudentID,
			StudentName:    studentName,
			PhoneNumber:    reg.Student.Phone,
			CourseName:     studentClass.Course.Name,
			Percentage:     percentage,
		})
	}

	return results, nil
}

func GetMessageTemplate(percentage float64, courseName string) string {
	templates, err := LoadMessageTemplates()
	if err != nil {
		return ""
	}

	lang := "pt"

	var template string
	var category string

	if percentage == 100 {
		category = "excellent"
		template = templates.Excellent.PT
	} else if percentage >= 80 {
		category = "good"
		template = templates.Good.PT
	} else if percentage >= 50 {
		category = "warning"
		template = templates.Warning.PT
	} else {
		category = "critical"
		template = templates.Critical.PT
	}

	fmt.Printf("Selected template category: %s, language: %s for percentage: %.0f%%, course: %s\n",
		category, lang, percentage, courseName)

	return template
}

func FormatMessage(template string, studentName string, courseName string, percentage float64) string {
	message := template
	message = strings.ReplaceAll(message, "[Nome]", studentName)
	message = strings.ReplaceAll(message, "[Name]", studentName)
	message = strings.ReplaceAll(message, "[Curso]", courseName)
	message = strings.ReplaceAll(message, "[Course]", courseName)
	message = strings.ReplaceAll(message, "[X]", fmt.Sprintf("%.0f", percentage))
	return message
}
