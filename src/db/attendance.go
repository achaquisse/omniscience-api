package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Attendance struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	RegistrationID uint      `gorm:"not null;index:idx_registration_id;uniqueIndex:unique_registration_date"`
	Registration   Registration `gorm:"foreignKey:RegistrationID"`
	Date           string    `gorm:"type:date;not null;index:idx_date;uniqueIndex:unique_registration_date"`
	Status         string    `gorm:"size:20;not null;index:idx_status"`
	Remarks        string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (Attendance) TableName() string {
	return "Attendance"
}

func CreateOrUpdateAttendance(registrationID uint, date string, status string, remarks string) error {
	attendance := Attendance{
		RegistrationID: registrationID,
		Date:           date,
		Status:         status,
		Remarks:        remarks,
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "registration_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "remarks", "updated_at"}),
	}).Create(&attendance)

	return result.Error
}

func GetAttendanceByDate(date string, studentClassID uint) []Attendance {
	var attendances []Attendance

	db.Preload("Registration.Student").
		Joins("JOIN Registration ON Registration.id = Attendance.registration_id").
		Where("Attendance.date = ?", date).
		Where("Registration.student_class_id = ?", studentClassID).
		Find(&attendances)

	return attendances
}

type AttendanceReport struct {
	TotalDays    int     `json:"totalDays"`
	PresentCount int     `json:"presentCount"`
	AbsentCount  int     `json:"absentCount"`
	LateCount    int     `json:"lateCount"`
	ExcusedCount int     `json:"excusedCount"`
	Percentage   float64 `json:"percentage"`
}

func GetStudentAttendanceReport(studentID uint, startDate string, endDate string) AttendanceReport {
	var attendances []Attendance

	db.Joins("JOIN Registration ON Registration.id = Attendance.registration_id").
		Where("Registration.student_id = ?", studentID).
		Where("Attendance.date >= ?", startDate).
		Where("Attendance.date <= ?", endDate).
		Find(&attendances)

	report := AttendanceReport{
		TotalDays: len(attendances),
	}

	for _, attendance := range attendances {
		switch attendance.Status {
		case "PRESENT":
			report.PresentCount++
		case "ABSENT":
			report.AbsentCount++
		case "LATE":
			report.LateCount++
		case "EXCUSED":
			report.ExcusedCount++
		}
	}

	if report.TotalDays > 0 {
		report.Percentage = float64(report.PresentCount) / float64(report.TotalDays) * 100
	}

	return report
}

type AttendanceRecord struct {
	Date    string `json:"date"`
	Status  string `json:"status"`
	Remarks string `json:"remarks"`
}

type WeeklyTrend struct {
	Week         string  `json:"week"`
	TotalDays    int     `json:"totalDays"`
	PresentCount int     `json:"presentCount"`
	Percentage   float64 `json:"percentage"`
}

type MonthlyTrend struct {
	Month        string  `json:"month"`
	TotalDays    int     `json:"totalDays"`
	PresentCount int     `json:"presentCount"`
	Percentage   float64 `json:"percentage"`
}

type DetailedAttendanceReport struct {
	Summary       AttendanceReport   `json:"summary"`
	Records       []AttendanceRecord `json:"records"`
	WeeklyTrends  []WeeklyTrend      `json:"weeklyTrends"`
	MonthlyTrends []MonthlyTrend     `json:"monthlyTrends"`
}

type StudentClassAttendanceReport struct {
	StudentClassID   uint             `json:"studentClassId"`
	StudentClassName string           `json:"studentClassName"`
	Summary          AttendanceReport `json:"summary"`
	Records          []AttendanceRecord `json:"records"`
	WeeklyTrends     []WeeklyTrend      `json:"weeklyTrends"`
	MonthlyTrends    []MonthlyTrend     `json:"monthlyTrends"`
}

type AggregatedStudentAttendanceReport struct {
	OverallSummary AttendanceReport               `json:"overallSummary"`
	ByClass        []StudentClassAttendanceReport `json:"byClass"`
}

func GetDetailedStudentAttendanceReport(studentID uint, startDate string, endDate string, studentClassID *uint) DetailedAttendanceReport {
	var attendances []Attendance

	query := db.Joins("JOIN Registration ON Registration.id = Attendance.registration_id").
		Where("Registration.student_id = ?", studentID).
		Where("Attendance.date >= ?", startDate).
		Where("Attendance.date <= ?", endDate)

	if studentClassID != nil {
		query = query.Where("Registration.student_class_id = ?", *studentClassID)
	}

	query.Order("Attendance.date ASC").Find(&attendances)

	summary := AttendanceReport{
		TotalDays: len(attendances),
	}

	var records []AttendanceRecord
	weeklyMap := make(map[string]*WeeklyTrend)
	monthlyMap := make(map[string]*MonthlyTrend)

	for _, attendance := range attendances {
		switch attendance.Status {
		case "PRESENT":
			summary.PresentCount++
		case "ABSENT":
			summary.AbsentCount++
		case "LATE":
			summary.LateCount++
		case "EXCUSED":
			summary.ExcusedCount++
		}

		records = append(records, AttendanceRecord{
			Date:    attendance.Date,
			Status:  attendance.Status,
			Remarks: attendance.Remarks,
		})

		parsedDate, err := time.Parse("2006-01-02", attendance.Date)
		if err == nil {
			year, week := parsedDate.ISOWeek()
			weekKey := fmt.Sprintf("%d-W%02d", year, week)
			if weeklyMap[weekKey] == nil {
				weeklyMap[weekKey] = &WeeklyTrend{Week: weekKey}
			}
			weeklyMap[weekKey].TotalDays++
			if attendance.Status == "PRESENT" {
				weeklyMap[weekKey].PresentCount++
			}

			monthKey := parsedDate.Format("2006-01")
			if monthlyMap[monthKey] == nil {
				monthlyMap[monthKey] = &MonthlyTrend{Month: monthKey}
			}
			monthlyMap[monthKey].TotalDays++
			if attendance.Status == "PRESENT" {
				monthlyMap[monthKey].PresentCount++
			}
		}
	}

	if summary.TotalDays > 0 {
		summary.Percentage = float64(summary.PresentCount) / float64(summary.TotalDays) * 100
	}

	var weeklyTrends []WeeklyTrend
	for _, trend := range weeklyMap {
		if trend.TotalDays > 0 {
			trend.Percentage = float64(trend.PresentCount) / float64(trend.TotalDays) * 100
		}
		weeklyTrends = append(weeklyTrends, *trend)
	}

	var monthlyTrends []MonthlyTrend
	for _, trend := range monthlyMap {
		if trend.TotalDays > 0 {
			trend.Percentage = float64(trend.PresentCount) / float64(trend.TotalDays) * 100
		}
		monthlyTrends = append(monthlyTrends, *trend)
	}

	return DetailedAttendanceReport{
		Summary:       summary,
		Records:       records,
		WeeklyTrends:  weeklyTrends,
		MonthlyTrends: monthlyTrends,
	}
}

func GetAggregatedStudentAttendanceReport(studentID uint, startDate string, endDate string) AggregatedStudentAttendanceReport {
	var attendances []Attendance

	db.Joins("JOIN Registration ON Registration.id = Attendance.registration_id").
		Preload("Registration").
		Where("Registration.student_id = ?", studentID).
		Where("Attendance.date >= ?", startDate).
		Where("Attendance.date <= ?", endDate).
		Order("Attendance.date ASC").
		Find(&attendances)

	overallSummary := AttendanceReport{
		TotalDays: len(attendances),
	}

	classMap := make(map[uint]*StudentClassAttendanceReport)

	for _, attendance := range attendances {
		classID := attendance.Registration.StudentClassID

		if classMap[classID] == nil {
			var studentClass StudentClass
			db.First(&studentClass, classID)

			classMap[classID] = &StudentClassAttendanceReport{
				StudentClassID:   classID,
				StudentClassName: studentClass.Name,
				Summary:          AttendanceReport{},
				Records:          []AttendanceRecord{},
				WeeklyTrends:     []WeeklyTrend{},
				MonthlyTrends:    []MonthlyTrend{},
			}
		}

		classReport := classMap[classID]
		classReport.Summary.TotalDays++

		switch attendance.Status {
		case "PRESENT":
			overallSummary.PresentCount++
			classReport.Summary.PresentCount++
		case "ABSENT":
			overallSummary.AbsentCount++
			classReport.Summary.AbsentCount++
		case "LATE":
			overallSummary.LateCount++
			classReport.Summary.LateCount++
		case "EXCUSED":
			overallSummary.ExcusedCount++
			classReport.Summary.ExcusedCount++
		}

		classReport.Records = append(classReport.Records, AttendanceRecord{
			Date:    attendance.Date,
			Status:  attendance.Status,
			Remarks: attendance.Remarks,
		})
	}

	if overallSummary.TotalDays > 0 {
		overallSummary.Percentage = float64(overallSummary.PresentCount) / float64(overallSummary.TotalDays) * 100
	}

	var byClass []StudentClassAttendanceReport
	for _, classReport := range classMap {
		weeklyMap := make(map[string]*WeeklyTrend)
		monthlyMap := make(map[string]*MonthlyTrend)

		for _, record := range classReport.Records {
			parsedDate, err := time.Parse("2006-01-02", record.Date)
			if err == nil {
				year, week := parsedDate.ISOWeek()
				weekKey := fmt.Sprintf("%d-W%02d", year, week)
				if weeklyMap[weekKey] == nil {
					weeklyMap[weekKey] = &WeeklyTrend{Week: weekKey}
				}
				weeklyMap[weekKey].TotalDays++
				if record.Status == "PRESENT" {
					weeklyMap[weekKey].PresentCount++
				}

				monthKey := parsedDate.Format("2006-01")
				if monthlyMap[monthKey] == nil {
					monthlyMap[monthKey] = &MonthlyTrend{Month: monthKey}
				}
				monthlyMap[monthKey].TotalDays++
				if record.Status == "PRESENT" {
					monthlyMap[monthKey].PresentCount++
				}
			}
		}

		if classReport.Summary.TotalDays > 0 {
			classReport.Summary.Percentage = float64(classReport.Summary.PresentCount) / float64(classReport.Summary.TotalDays) * 100
		}

		for _, trend := range weeklyMap {
			if trend.TotalDays > 0 {
				trend.Percentage = float64(trend.PresentCount) / float64(trend.TotalDays) * 100
			}
			classReport.WeeklyTrends = append(classReport.WeeklyTrends, *trend)
		}

		for _, trend := range monthlyMap {
			if trend.TotalDays > 0 {
				trend.Percentage = float64(trend.PresentCount) / float64(trend.TotalDays) * 100
			}
			classReport.MonthlyTrends = append(classReport.MonthlyTrends, *trend)
		}

		byClass = append(byClass, *classReport)
	}

	return AggregatedStudentAttendanceReport{
		OverallSummary: overallSummary,
		ByClass:        byClass,
	}
}

type BulkAttendanceRecord struct {
	RegistrationID uint
	Date           string
	Status         string
	Remarks        string
}

func CreateOrUpdateBulkAttendance(records []BulkAttendanceRecord) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			attendance := Attendance{
				RegistrationID: record.RegistrationID,
				Date:           record.Date,
				Status:         record.Status,
				Remarks:        record.Remarks,
			}

			result := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "registration_id"}, {Name: "date"}},
				DoUpdates: clause.AssignmentColumns([]string{"status", "remarks", "updated_at"}),
			}).Create(&attendance)

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

type StudentAttendanceSummary struct {
	StudentID    uint    `json:"studentId"`
	StudentName  string  `json:"studentName"`
	TotalDays    int     `json:"totalDays"`
	PresentCount int     `json:"presentCount"`
	AbsentCount  int     `json:"absentCount"`
	LateCount    int     `json:"lateCount"`
	ExcusedCount int     `json:"excusedCount"`
	Percentage   float64 `json:"percentage"`
}

type DailyAttendance struct {
	Date         string `json:"date"`
	TotalStudents int    `json:"totalStudents"`
	PresentCount int     `json:"presentCount"`
	AbsentCount  int     `json:"absentCount"`
	LateCount    int     `json:"lateCount"`
	ExcusedCount int     `json:"excusedCount"`
	Percentage   float64 `json:"percentage"`
}

type ClassAttendanceReport struct {
	Period           string                     `json:"period"`
	StartDate        string                     `json:"startDate"`
	EndDate          string                     `json:"endDate"`
	TotalStudents    int                        `json:"totalStudents"`
	OverallSummary   AttendanceReport           `json:"overallSummary"`
	StudentSummaries []StudentAttendanceSummary `json:"studentSummaries"`
	DailyData        []DailyAttendance          `json:"dailyData,omitempty"`
	WeeklyData       []WeeklyTrend              `json:"weeklyData,omitempty"`
	MonthlyData      []MonthlyTrend             `json:"monthlyData,omitempty"`
}

func GetClassAttendanceReport(studentClassID uint, startDate string, endDate string, period string) ClassAttendanceReport {
	var registrations []Registration
	db.Where("student_class_id = ?", studentClassID).Preload("Student").Find(&registrations)

	var attendances []Attendance
	db.Joins("JOIN Registration ON Registration.id = Attendance.registration_id").
		Where("Registration.student_class_id = ?", studentClassID).
		Where("Attendance.date >= ?", startDate).
		Where("Attendance.date <= ?", endDate).
		Preload("Registration.Student").
		Order("Attendance.date ASC").
		Find(&attendances)

	overallSummary := AttendanceReport{
		TotalDays: len(attendances),
	}

	studentMap := make(map[uint]*StudentAttendanceSummary)
	for _, reg := range registrations {
		studentMap[reg.StudentID] = &StudentAttendanceSummary{
			StudentID:   reg.StudentID,
			StudentName: fmt.Sprintf("%s %s", reg.Student.FirstName, reg.Student.LastName),
		}
	}

	dailyMap := make(map[string]*DailyAttendance)
	weeklyMap := make(map[string]*WeeklyTrend)
	monthlyMap := make(map[string]*MonthlyTrend)

	for _, attendance := range attendances {
		switch attendance.Status {
		case "PRESENT":
			overallSummary.PresentCount++
		case "ABSENT":
			overallSummary.AbsentCount++
		case "LATE":
			overallSummary.LateCount++
		case "EXCUSED":
			overallSummary.ExcusedCount++
		}

		studentID := attendance.Registration.StudentID
		if summary, exists := studentMap[studentID]; exists {
			summary.TotalDays++
			switch attendance.Status {
			case "PRESENT":
				summary.PresentCount++
			case "ABSENT":
				summary.AbsentCount++
			case "LATE":
				summary.LateCount++
			case "EXCUSED":
				summary.ExcusedCount++
			}
		}

		if period == "day" || period == "all" {
			if dailyMap[attendance.Date] == nil {
				dailyMap[attendance.Date] = &DailyAttendance{
					Date:          attendance.Date,
					TotalStudents: len(registrations),
				}
			}
			dailyMap[attendance.Date].TotalStudents++
			switch attendance.Status {
			case "PRESENT":
				dailyMap[attendance.Date].PresentCount++
			case "ABSENT":
				dailyMap[attendance.Date].AbsentCount++
			case "LATE":
				dailyMap[attendance.Date].LateCount++
			case "EXCUSED":
				dailyMap[attendance.Date].ExcusedCount++
			}
		}

		parsedDate, err := time.Parse("2006-01-02", attendance.Date)
		if err == nil {
			if period == "week" || period == "all" {
				year, week := parsedDate.ISOWeek()
				weekKey := fmt.Sprintf("%d-W%02d", year, week)
				if weeklyMap[weekKey] == nil {
					weeklyMap[weekKey] = &WeeklyTrend{Week: weekKey}
				}
				weeklyMap[weekKey].TotalDays++
				if attendance.Status == "PRESENT" {
					weeklyMap[weekKey].PresentCount++
				}
			}

			if period == "month" || period == "all" {
				monthKey := parsedDate.Format("2006-01")
				if monthlyMap[monthKey] == nil {
					monthlyMap[monthKey] = &MonthlyTrend{Month: monthKey}
				}
				monthlyMap[monthKey].TotalDays++
				if attendance.Status == "PRESENT" {
					monthlyMap[monthKey].PresentCount++
				}
			}
		}
	}

	if overallSummary.TotalDays > 0 {
		overallSummary.Percentage = float64(overallSummary.PresentCount) / float64(overallSummary.TotalDays) * 100
	}

	var studentSummaries []StudentAttendanceSummary
	for _, summary := range studentMap {
		if summary.TotalDays > 0 {
			summary.Percentage = float64(summary.PresentCount) / float64(summary.TotalDays) * 100
		}
		studentSummaries = append(studentSummaries, *summary)
	}

	var dailyData []DailyAttendance
	if period == "day" {
		for _, daily := range dailyMap {
			if daily.TotalStudents > 0 {
				daily.Percentage = float64(daily.PresentCount) / float64(daily.TotalStudents) * 100
			}
			dailyData = append(dailyData, *daily)
		}
	}

	var weeklyData []WeeklyTrend
	if period == "week" {
		for _, trend := range weeklyMap {
			if trend.TotalDays > 0 {
				trend.Percentage = float64(trend.PresentCount) / float64(trend.TotalDays) * 100
			}
			weeklyData = append(weeklyData, *trend)
		}
	}

	var monthlyData []MonthlyTrend
	if period == "month" {
		for _, trend := range monthlyMap {
			if trend.TotalDays > 0 {
				trend.Percentage = float64(trend.PresentCount) / float64(trend.TotalDays) * 100
			}
			monthlyData = append(monthlyData, *trend)
		}
	}

	return ClassAttendanceReport{
		Period:           period,
		StartDate:        startDate,
		EndDate:          endDate,
		TotalStudents:    len(registrations),
		OverallSummary:   overallSummary,
		StudentSummaries: studentSummaries,
		DailyData:        dailyData,
		WeeklyData:       weeklyData,
		MonthlyData:      monthlyData,
	}
}
