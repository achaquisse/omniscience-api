package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"skulla-api/db"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testTeacherEmail = "teacher@test.com"
const testTeacherEmail2 = "teacher2@test.com"

func setupTestDB() (*gorm.DB, error) {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = testDB.AutoMigrate(
		&db.Course{},
		&db.Period{},
		&db.StudentClass{},
		&db.Student{},
		&db.Registration{},
		&db.Attendance{},
	)
	if err != nil {
		return nil, err
	}

	return testDB, nil
}

func seedTestData(testDB *gorm.DB) error {
	now := time.Now()

	courses := []db.Course{
		{ID: 1, Name: "Mathematics", TeacherEmail: testTeacherEmail},
		{ID: 2, Name: "Physics", TeacherEmail: testTeacherEmail},
		{ID: 3, Name: "Chemistry", TeacherEmail: testTeacherEmail2},
	}

	for _, course := range courses {
		if err := testDB.Create(&course).Error; err != nil {
			return err
		}
	}

	periods := []db.Period{
		{ID: 1, Start: now.AddDate(0, -2, 0), End: now.AddDate(0, 2, 0)},
		{ID: 2, Start: now.AddDate(0, -6, 0), End: now.AddDate(0, -4, 0)},
		{ID: 3, Start: now.AddDate(0, 3, 0), End: now.AddDate(0, 5, 0)},
	}
	for _, period := range periods {
		if err := testDB.Create(&period).Error; err != nil {
			return err
		}
	}

	studentClasses := []db.StudentClass{
		{ID: 1, Name: "Math 101", CourseID: 1, PeriodId: 1},
		{ID: 2, Name: "Math 102", CourseID: 1, PeriodId: 2},
		{ID: 3, Name: "Physics 101", CourseID: 2, PeriodId: 1},
		{ID: 4, Name: "Chemistry 101", CourseID: 3, PeriodId: 1},
	}
	for _, class := range studentClasses {
		if err := testDB.Create(&class).Error; err != nil {
			return err
		}
	}

	students := []db.Student{
		{ID: 1, FirstName: "John", LastName: "Doe"},
		{ID: 2, FirstName: "Jane", LastName: "Smith"},
		{ID: 3, FirstName: "Bob", LastName: "Johnson"},
	}
	for _, student := range students {
		if err := testDB.Create(&student).Error; err != nil {
			return err
		}
	}

	registrations := []db.Registration{
		{ID: 1, StudentID: 1, StudentClassID: 1, Status: "ACTIVE"},
		{ID: 2, StudentID: 2, StudentClassID: 1, Status: "ACTIVE"},
		{ID: 3, StudentID: 3, StudentClassID: 1, Status: "ACTIVE"},
		{ID: 4, StudentID: 1, StudentClassID: 3, Status: "ACTIVE"},
		{ID: 5, StudentID: 2, StudentClassID: 3, Status: "ACTIVE"},
	}
	for _, reg := range registrations {
		if err := testDB.Create(&reg).Error; err != nil {
			return err
		}
	}

	attendances := []db.Attendance{
		{ID: 1, RegistrationID: 1, Date: "2024-01-15", Status: "PRESENT", Remarks: "On time"},
		{ID: 2, RegistrationID: 1, Date: "2024-01-16", Status: "ABSENT", Remarks: "Sick"},
		{ID: 3, RegistrationID: 1, Date: "2024-01-17", Status: "PRESENT", Remarks: ""},
		{ID: 4, RegistrationID: 2, Date: "2024-01-15", Status: "PRESENT", Remarks: ""},
		{ID: 5, RegistrationID: 2, Date: "2024-01-16", Status: "LATE", Remarks: "Traffic"},
		{ID: 6, RegistrationID: 4, Date: "2024-01-15", Status: "PRESENT", Remarks: ""},
		{ID: 7, RegistrationID: 4, Date: "2024-01-16", Status: "EXCUSED", Remarks: "Medical appointment"},
	}
	for _, att := range attendances {
		if err := testDB.Create(&att).Error; err != nil {
			return err
		}
	}

	return nil
}

func setupTestApp(t *testing.T) *fiber.App {
	err := os.Setenv("TEST_MODE", "true")
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	testDB, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	if err := seedTestData(testDB); err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	db.SetDB(testDB)

	app := fiber.New()
	Init(app)

	return app
}

func createTestJWT(email string) string {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))
	return tokenString
}

func makeRequest(app *fiber.App, method, path, authEmail string, body interface{}) (*httptest.ResponseRecorder, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	if authEmail != "" {
		req.Header.Set("Authorization", "Bearer "+createTestJWT(authEmail))
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		return nil, err
	}

	rec := httptest.NewRecorder()
	bodyBytes, _ := io.ReadAll(resp.Body)
	rec.Body.Write(bodyBytes)
	rec.Code = resp.StatusCode

	return rec, nil
}
