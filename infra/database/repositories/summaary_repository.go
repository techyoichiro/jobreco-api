package repository

import (
	"time"

	model "github.com/techyoichiro/jobreco-api/domain/models"
	"gorm.io/gorm"
)

type SummaryRepositoryImpl struct {
	DB *gorm.DB
}

func NewSummaryRepository(db *gorm.DB) *SummaryRepositoryImpl {
	return &SummaryRepositoryImpl{DB: db}
}

func (r *SummaryRepositoryImpl) GetAllEmployee() ([]model.Employee, error) {
	var employees []model.Employee
	if err := r.DB.Select("id, name").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *SummaryRepositoryImpl) GetAttendance(employeeID uint, year int, month int) ([]model.Attendance, error) {
	var attendances []model.Attendance

	err := r.DB.Where("employee_id = ? AND EXTRACT(YEAR FROM work_date) = ? AND EXTRACT(MONTH FROM work_date) = ?", employeeID, year, month).
		Find(&attendances).Error
	if err != nil {
		return nil, err
	}

	return attendances, nil
}

func (r *SummaryRepositoryImpl) GetHourlyPay(employeeID uint) (int, error) {
	var employee model.Employee
	if err := r.DB.Where("id = ?", employeeID).First(&employee).Error; err != nil {
		return 0, err
	}
	return employee.HourlyPay, nil
}

func (r *SummaryRepositoryImpl) GetAttendanceByID(attedanceID uint) (*model.Attendance, error) {
	var attendance model.Attendance

	err := r.DB.Where("id = ?", attedanceID).First(&attendance).Error
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (r *SummaryRepositoryImpl) UpdateAttendance(attendance *model.Attendance) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// IDでレコードを検索し、更新
	if err := tx.Model(&model.Attendance{}).Where("id = ?", attendance.ID).Updates(attendance).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *SummaryRepositoryImpl) GetWorkDateByID(id uint) (time.Time, error) {
	var attendance model.Attendance
	err := r.DB.Where("id = ?", id).First(&attendance).Error
	if err != nil {
		return time.Time{}, err // errorの場合はゼロ値のtime.Timeを返す
	}

	// attendanceから日付を取得して返す
	return attendance.WorkDate, nil
}
