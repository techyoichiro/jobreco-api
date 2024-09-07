package repositories

import (
	model "github.com/techyoichiro/jobreco-api/domain/models"
)

type SummaryRepository interface {
	GetAllEmployee() ([]model.Employee, error)
	GetSummary(uint, int, int) ([]model.DailyWorkSummary, error)
	GetWorkSegments(uint) []model.WorkSegment
	GetHourlyPay(uint) (int, error)
	GetSummaryBySummaryID(uint) (*model.DailyWorkSummary, error)
	UpdateSummary(*model.DailyWorkSummary) error
	FindWorkSegmentsBySummaryID(uint) ([]model.WorkSegment, error)
	FindBreakRecords(uint) (*model.BreakRecord, error)
	FindSummaryBySegmentID(segmentID uint) (*model.DailyWorkSummary, error)
	UpdateWorkSegment(segment *model.WorkSegment) error
	UpdateBreakRecord(breakRecord *model.BreakRecord) error
	FindWorkSegmentByID(ID uint) ([]model.WorkSegment, error)
}
