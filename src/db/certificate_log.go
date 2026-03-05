package db

import "time"

type CertificateLog struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	Template        string    `gorm:"size:50;not null"`
	StudentName     string    `gorm:"size:500;not null"`
	CertDescription string    `gorm:"type:text;not null"`
	CreatedBy       string    `gorm:"size:500;not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}

func (CertificateLog) TableName() string {
	return "CertificateLog"
}

func CreateCertificateLog(template, studentName, certDescription, createdBy string) error {
	log := CertificateLog{
		Template:        template,
		StudentName:     studentName,
		CertDescription: certDescription,
		CreatedBy:       createdBy,
	}
	return db.Create(&log).Error
}
