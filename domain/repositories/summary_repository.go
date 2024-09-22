package repositories

import (
	"time"

	model "github.com/techyoichiro/jobreco-api/domain/models"
)

type SummaryRepository interface {
	GetAllEmployee() ([]model.Employee, error)
	GetAttendance(uint, int, int) ([]model.Attendance, error)
	GetHourlyPay(uint) (int, error)
	GetAttendanceByID(uint) (*model.Attendance, error)
	UpdateAttendance(*model.Attendance) error
	GetWorkDateByID(uint) (time.Time, error)
}
