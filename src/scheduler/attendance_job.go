package scheduler

import (
	"fmt"
	"log"
	"os"
	"time"

	"skulla-api/services"

	"github.com/robfig/cron/v3"
)

var AttendanceScheduler *cron.Cron

func InitAttendanceScheduler() {
	scheduleTime := os.Getenv("ATTENDANCE_REPORT_SCHEDULE")
	if scheduleTime == "" {
		scheduleTime = "0 8 * * 6"
	}

	AttendanceScheduler = cron.New()

	_, err := AttendanceScheduler.AddFunc(scheduleTime, SendWeeklyAttendanceReports)
	if err != nil {
		log.Fatalf("Failed to schedule attendance reports: %v", err)
	}

	AttendanceScheduler.Start()
	log.Printf("Attendance report scheduler initialized with schedule: %s", scheduleTime)
}

func SendWeeklyAttendanceReports() {
	log.Println("Starting weekly attendance report job...")

	now := time.Now()
	endDate := now.AddDate(0, 0, -1)
	startDate := endDate.AddDate(0, 0, -6)

	log.Printf("Generating reports for period: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	students, err := services.GetActiveRegistrationsWithAttendance(startDate, endDate)
	if err != nil {
		log.Printf("Error fetching attendance data: %v", err)
		return
	}

	log.Printf("Found %d active registrations to process", len(students))

	if len(students) == 0 {
		log.Println("No active registrations found. Job completed.")
		return
	}

	err = services.SendAttendanceReports(students)
	if err != nil {
		log.Printf("Error sending attendance reports: %v", err)
		return
	}

	log.Println("Weekly attendance report job completed successfully")
}

func StopScheduler() {
	if AttendanceScheduler != nil {
		AttendanceScheduler.Stop()
		log.Println("Attendance report scheduler stopped")
	}
}

func TriggerManualReport() error {
	fmt.Println("Triggering manual attendance report...")
	go SendWeeklyAttendanceReports()
	return nil
}
