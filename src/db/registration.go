package db

type Registration struct {
	ID             uint    `gorm:"primaryKey"`
	Status         string  `gorm:"size:255;not null"`
	StudentID      uint    `gorm:"foreignKey:StudentID"`
	Student        Student `gorm:"foreignKey:StudentID"`
	StudentClassID uint    `gorm:"column:student_class_id"`
}

func (Registration) TableName() string {
	return "Registration"
}

type Student struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"column:firstName;size:255"`
	LastName  string `gorm:"column:lastName;size:255"`
}

func (Student) TableName() string {
	return "Student"
}

func ListRegistrations(studentClassId int) []Registration {
	var registration []Registration

	db.Preload("Student").
		Joins("JOIN Student ON Student.ID = Registration.student_id").
		Where("student_class_id = ?", studentClassId).
		Order("Student.firstName ASC, Student.lastName ASC").
		Find(&registration)

	return registration
}

func RegistrationExists(registrationID uint) bool {
	var count int64
	db.Model(&Registration{}).Where("id = ?", registrationID).Count(&count)
	return count > 0
}
