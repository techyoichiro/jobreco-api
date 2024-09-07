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

func (r *AttendanceRepositoryImpl) CreateDailyWorkSummary(summary *model.DailyWorkSummary) error {
	return r.DB.Create(summary).Error
}

func (r *AttendanceRepositoryImpl) FindDailyWorkSummary(employeeID uint, workDate string) (*model.DailyWorkSummary, error) {
	var summary model.DailyWorkSummary
	err := r.DB.Where("employee_id = ? AND work_date = ?", employeeID, workDate).First(&summary).Error
	return &summary, err
}

func (r *AttendanceRepositoryImpl) CreateWorkSegment(segment *model.WorkSegment) error {
	return r.DB.Create(segment).Error
}

func (r *AttendanceRepositoryImpl) UpdateWorkSegment(segment *model.WorkSegment) error {
	return r.DB.Save(segment).Error
}

func (r *AttendanceRepositoryImpl) FindLatestWorkSegment(employeeID uint) (*model.WorkSegment, error) {
	var segment model.WorkSegment
	err := r.DB.Where("employee_id = ? AND end_time IS NULL", employeeID).Order("start_time DESC").First(&segment).Error
	return &segment, err
}

func (r *AttendanceRepositoryImpl) FindBreakRecords(summaryID uint) ([]model.BreakRecord, error) {
	var breakRecords []model.BreakRecord
	err := r.DB.Where("summary_id = ? AND break_end IS NULL", summaryID).Find(&breakRecords).Error
	return breakRecords, err
}

func (r *AttendanceRepositoryImpl) FindWorkSegmentToReturn(employeeID uint) (*model.WorkSegment, error) {
	var segment model.WorkSegment
	err := r.DB.Where("employee_id = ? AND status_id = 2 AND end_time IS NULL", employeeID).Order("start_time DESC").First(&segment).Error
	return &segment, err
}

func (r *AttendanceRepositoryImpl) UpdateDailyWorkSummary(summary *model.DailyWorkSummary) error {
	return r.DB.Save(summary).Error
}

func (r *AttendanceRepositoryImpl) CreateBreakRecord(record *model.BreakRecord) error {
	return r.DB.Create(record).Error
}

func (r *AttendanceRepositoryImpl) FindWorkSegmentsByDate(employeeID uint, workDate string) ([]model.WorkSegment, error) {
	var segments []model.WorkSegment
	err := r.DB.Where("employee_id = ? AND DATE(start_time) = ?", employeeID, workDate).Find(&segments).Error
	return segments, err
}

func (r *AttendanceRepositoryImpl) UpdateBreakRecord(record *model.BreakRecord) error {
	return r.DB.Save(record).Error
}
