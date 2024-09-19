package repository

import (
	model "github.com/techyoichiro/jobreco-api/domain/models"
	"gorm.io/gorm"
)

type AttendanceRepositoryImpl struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepositoryImpl {
	return &AttendanceRepositoryImpl{DB: db}
}

func (r *AttendanceRepositoryImpl) CreateAttendance(attendance *model.Attendance) error {
	return r.DB.Create(attendance).Error
}

func (r *AttendanceRepositoryImpl) FindAttendance(employeeID uint, workDate string) (*model.Attendance, error) {
	var attendance model.Attendance
	err := r.DB.Where("employee_id = ? AND work_date = ?", employeeID, workDate).First(&attendance).Error
	return &attendance, err
}

func (r *AttendanceRepositoryImpl) UpdateAttendance(attendance *model.Attendance) error {
	return r.DB.Save(attendance).Error
}
