package model

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	EmployeeID uint       `gorm:"not null;index"`                                   // 外部キー：employees テーブル
	WorkDate   time.Time  `gorm:"type:date;not null;uniqueIndex:unique_attendance"` // 勤務日、一意制約に含める
	StartTime1 *time.Time `gorm:"type:timestamp"`                                   // 勤務開始時間1
	EndTime1   *time.Time `gorm:"type:timestamp"`                                   // 勤務終了時間1
	StartTime2 *time.Time `gorm:"type:timestamp"`                                   // 勤務開始時間2（オプション）
	EndTime2   *time.Time `gorm:"type:timestamp"`                                   // 勤務終了時間2（オプション）
	BreakStart *time.Time `gorm:"type:timestamp"`                                   // 休憩開始時間
	BreakEnd   *time.Time `gorm:"type:timestamp"`                                   // 休憩終了時間
	StoreID1   uint       `gorm:"not null"`                                         // 外部キー：stores テーブル
	StoreID2   *uint      `gorm:""`                                                 // 外部キー（オプション）：stores テーブル
	StatusID   int        `gorm:"not null"`                                         // 勤務ステータスID（オプション）

	// 一意制約
	UniqueIndex string `gorm:"uniqueIndex:unique_attendance,employee_id,work_date"`
}

type AttendanceResponse struct {
	ID            uint    `json:"ID"`
	WorkDate      string  `json:"WorkDate"`
	StartTime1    string  `json:"StartTime1"`
	EndTime1      string  `json:"EndTime1,omitempty"`
	StartTime2    string  `json:"StartTime2"`
	EndTime2      string  `json:"EndTime2"`
	BreakStart    string  `json:"BreakStart"`
	BreakEnd      string  `json:"BreakEnd"`
	TotalWorkTime float64 `json:"TotalWorkTime"`
	Overtime      float64 `json:"Overtime"`
	StoreID1      uint    `json:"StoreID1"`
	StoreID2      *uint   `json:"StoreID2"`
	Remarks       string  `json:"Remarks"`
	HourlyPay     int     `json:"HourlyPay"`
}
