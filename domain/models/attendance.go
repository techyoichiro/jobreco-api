package model

import (
	"time"

	"gorm.io/gorm"
)

type DailyWorkSummary struct {
	gorm.Model
	EmployeeID     uint          `gorm:"not null;index"`
	WorkDate       time.Time     `gorm:"type:date;not null"`
	StartTime      time.Time     `gorm:"type:timestamp;not null"`
	EndTime        *time.Time    `gorm:"type:timestamp"`
	TotalWorkTime  float64       `gorm:"type:decimal(4,2);not null;default:0"` // 総勤務時間（分）
	TotalBreakTime float64       `gorm:"type:decimal(4,2);not null;default:0"` // 総休憩時間（分）
	WorkSegments   []WorkSegment `gorm:"foreignkey:SummaryID"`
	BreakRecords   []BreakRecord `gorm:"foreignKey:SummaryID"`
}

type WorkSegment struct {
	gorm.Model
	SummaryID  uint       `gorm:"not null"`
	EmployeeID uint       `gorm:"not null"`
	StoreID    uint       `gorm:"not null"`
	StartTime  time.Time  `gorm:"type:timestamp;not null"`
	EndTime    *time.Time `gorm:"type:timestamp"`
	StatusID   int        `gorm:"index"`
}

type BreakRecord struct {
	gorm.Model
	SummaryID  uint       `gorm:"not null"`
	BreakStart time.Time  `gorm:"type:timestamp"`
	BreakEnd   *time.Time `gorm:"type:timestamp"`
}
