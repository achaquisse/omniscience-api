package rest

import (
	"fmt"
	"skulla-api/db"
	"skulla-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type RecordAttendanceRequest struct {
	RegistrationID uint   `json:"registration_id"`
	Date           string `json:"date"`
	Status         string `json:"status"`
	Remarks        string `json:"remarks"`
}

var validStatuses = map[string]bool{
	"PRESENT": true,
	"ABSENT":  true,
	"LATE":    true,
	"EXCUSED": true,
}

func validateAttendanceRecord(req RecordAttendanceRequest, index int) error {
	if req.RegistrationID == 0 {
		return fmt.Errorf("record %d: registration_id is required", index)
	}

	if req.Date == "" {
		return fmt.Errorf("record %d: date is required", index)
	}

	if !validStatuses[req.Status] {
		return fmt.Errorf("record %d: status must be one of: PRESENT, ABSENT, LATE, EXCUSED", index)
	}

	if !db.RegistrationExists(req.RegistrationID) {
		return fmt.Errorf("record %d: registration not found", index)
	}

	return nil
}

func RecordAttendance(c *fiber.Ctx) error {
	var req RecordAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return ReturnBadRequest(c, "Invalid request body")
	}

	if err := validateAttendanceRecord(req, 0); err != nil {
		if err.Error() == "record 0: registration not found" {
			return ReturnNotFound(c, "registration not found")
		}
		return ReturnBadRequest(c, err.Error())
	}

	userEmail, err := GetUserEmailFromToken(c)
	if err != nil {
		return ReturnUnauthorized(c, "Unable to extract user email from token")
	}

	err = db.CreateOrUpdateAttendance(req.RegistrationID, req.Date, req.Status, req.Remarks, userEmail)
	if err != nil {
		log.Error(err)
		return ReturnInternalError(c, "Failed to record attendance")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":         "Attendance recorded successfully",
		"registration_id": req.RegistrationID,
		"date":            req.Date,
		"status":          req.Status,
		"remarks":         req.Remarks,
	})
}

func RecordBulkAttendance(c *fiber.Ctx) error {
	var requests []RecordAttendanceRequest
	if err := c.BodyParser(&requests); err != nil {
		return ReturnBadRequest(c, "Invalid request body. Expected JSON array of attendance records")
	}

	if len(requests) == 0 {
		return ReturnBadRequest(c, "At least one attendance record is required")
	}

	for i, req := range requests {
		if err := validateAttendanceRecord(req, i); err != nil {
			return ReturnBadRequest(c, err.Error())
		}
	}

	userEmail, err := GetUserEmailFromToken(c)
	if err != nil {
		return ReturnUnauthorized(c, "Unable to extract user email from token")
	}

	var bulkRecords []db.BulkAttendanceRecord
	for _, req := range requests {
		bulkRecords = append(bulkRecords, db.BulkAttendanceRecord{
			RegistrationID: req.RegistrationID,
			Date:           req.Date,
			Status:         req.Status,
			Remarks:        req.Remarks,
			UserEmail:      userEmail,
		})
	}

	err = db.CreateOrUpdateBulkAttendance(bulkRecords)
	if err != nil {
		log.Error(err)
		return ReturnInternalError(c, "Failed to record bulk attendance. All records have been rolled back.")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":           "Bulk attendance recorded successfully",
		"records_processed": len(requests),
	})
}

func GetStudentAttendanceReport(c *fiber.Ctx) error {
	studentID, err := ParseUintQueryParam(c, "student_id", true)
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	studentClassID, err := ParseOptionalUintQueryParam(c, "student_class_id")
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	startDate, endDate = GetDateRangeWithDefaults(startDate, endDate)

	if err := ValidateDateString(startDate, "start_date"); err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	if err := ValidateDateString(endDate, "end_date"); err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	if studentClassID != nil {
		report := db.GetDetailedStudentAttendanceReport(studentID, startDate, endDate, studentClassID)
		return c.JSON(report)
	}

	aggregatedReport := db.GetAggregatedStudentAttendanceReport(studentID, startDate, endDate)
	return c.JSON(aggregatedReport)
}

func GetClassAttendanceReport(c *fiber.Ctx) error {
	studentClassID, err := ParseUintQueryParam(c, "student_class_id", true)
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	startDate, endDate = GetDateRangeWithDefaults(startDate, endDate)

	if err := ValidateDateString(startDate, "start_date"); err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	if err := ValidateDateString(endDate, "end_date"); err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	period := c.Query("period")
	if period == "" {
		period = "all"
	}

	validPeriods := map[string]bool{
		"day":   true,
		"week":  true,
		"month": true,
		"all":   true,
	}

	if !validPeriods[period] {
		return ReturnBadRequest(c, "period must be one of: day, week, month, all")
	}

	report := db.GetClassAttendanceReport(studentClassID, startDate, endDate, period)

	return c.JSON(report)
}

func TriggerAttendanceReports(c *fiber.Ctx) error {
	dryRun := c.QueryBool("dry_run", false)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" || endDateStr == "" {
		now := time.Now()
		weekday := int(now.Weekday())

		daysToSubtract := weekday
		if weekday == 0 {
			daysToSubtract = 7
		}

		lastSunday := now.AddDate(0, 0, -daysToSubtract)
		lastSaturday := lastSunday.AddDate(0, 0, 6)

		startDate = time.Date(lastSunday.Year(), lastSunday.Month(), lastSunday.Day(), 0, 0, 0, 0, lastSunday.Location())
		endDate = time.Date(lastSaturday.Year(), lastSaturday.Month(), lastSaturday.Day(), 23, 59, 59, 0, lastSaturday.Location())
	} else {
		if err := ValidateDateString(startDateStr, "start_date"); err != nil {
			return ReturnBadRequest(c, err.Error())
		}
		if err := ValidateDateString(endDateStr, "end_date"); err != nil {
			return ReturnBadRequest(c, err.Error())
		}

		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return ReturnBadRequest(c, "Invalid start_date format")
		}

		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return ReturnBadRequest(c, "Invalid end_date format")
		}

		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
	}

	students, err := services.GetActiveRegistrationsWithAttendance(startDate, endDate)
	if err != nil {
		log.Error(err)
		return ReturnInternalError(c, "Failed to fetch attendance data")
	}

	result, err := services.SendAttendanceReportsWithResult(students, dryRun)
	if err != nil {
		log.Error(err)
		return ReturnInternalError(c, "Failed to process attendance reports")
	}

	response := fiber.Map{
		"dry_run":        dryRun,
		"start_date":     startDate.Format("2006-01-02"),
		"end_date":       endDate.Format("2006-01-02"),
		"total_messages": result.TotalMessages,
	}

	if dryRun {
		response["messages_planned"] = result.MessagesPlanned
	} else {
		response["success_count"] = result.SuccessCount
		response["failure_count"] = result.FailureCount
		response["messages_sent"] = result.MessagesSent
		response["messages_failed"] = result.MessagesFailed
	}

	return c.JSON(response)
}
