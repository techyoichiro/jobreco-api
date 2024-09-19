package repositories

import (
	model "github.com/techyoichiro/jobreco-api/domain/models"
)

// AttendanceRepository
type AttendanceRepository interface {
	CreateAttendance(summary *model.Attendance) error
	FindAttendance(employeeID uint, workDate string) (*model.Attendance, error)
	UpdateAttendance(summary *model.Attendance) error
}
