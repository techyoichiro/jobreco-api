package repositories

import (
	model "github.com/techyoichiro/jobreco-api/domain/models"
)

// AttendanceRepository
type AttendanceRepository interface {
	CreateDailyWorkSummary(summary *model.DailyWorkSummary) error
	FindDailyWorkSummary(employeeID uint, workDate string) (*model.DailyWorkSummary, error)
	CreateWorkSegment(segment *model.WorkSegment) error
	UpdateWorkSegment(segment *model.WorkSegment) error
	FindLatestWorkSegment(employeeID uint) (*model.WorkSegment, error)
	FindBreakRecords(summaryID uint) ([]model.BreakRecord, error)
	FindWorkSegmentToReturn(employeeID uint) (*model.WorkSegment, error)
	UpdateDailyWorkSummary(summary *model.DailyWorkSummary) error
	CreateBreakRecord(record *model.BreakRecord) error
	FindWorkSegmentsByDate(employeeID uint, workDate string) ([]model.WorkSegment, error)
	UpdateBreakRecord(record *model.BreakRecord) error
}
