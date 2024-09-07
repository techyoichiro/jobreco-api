package services

import (
	"errors"
	"fmt"
	"time"

	model "github.com/techyoichiro/jobreco-api/domain/models"
	"github.com/techyoichiro/jobreco-api/domain/repositories"
	"gorm.io/gorm"
)

type AttendanceService struct {
	repo repositories.AttendanceRepository
}

func NewAttendanceService(repo repositories.AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

// 出勤
func (s *AttendanceService) ClockIn(employeeID uint, storeID uint) error {
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	workDate := now.Format("2006-01-02")

	summary, err := s.repo.FindDailyWorkSummary(employeeID, workDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			summary = &model.DailyWorkSummary{
				EmployeeID:     employeeID,
				WorkDate:       now,
				StartTime:      now,
				TotalWorkTime:  0,
				TotalBreakTime: 0,
			}
			if err := s.repo.CreateDailyWorkSummary(summary); err != nil {
				return err
			}

			summary, err = s.repo.FindDailyWorkSummary(employeeID, workDate)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	workSegment := model.WorkSegment{
		SummaryID:  summary.ID,
		EmployeeID: employeeID,
		StoreID:    storeID,
		StartTime:  now,
		StatusID:   1, // 出勤
	}

	if err := s.repo.CreateWorkSegment(&workSegment); err != nil {
		return err
	}

	return nil
}

// 退勤
func (s *AttendanceService) ClockOut(employeeID uint, storeID uint) error {
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	workDate := now.Format("2006-01-02")

	workSegments, err := s.repo.FindWorkSegmentsByDate(employeeID, workDate)
	if err != nil {
		return err
	}
	if len(workSegments) == 0 {
		return gorm.ErrRecordNotFound
	}

	var earliestSegment, latestSegment *model.WorkSegment
	for i := range workSegments {
		if earliestSegment == nil || workSegments[i].StartTime.Before(earliestSegment.StartTime) {
			earliestSegment = &workSegments[i]
		}
		if latestSegment == nil || workSegments[i].StartTime.After(latestSegment.StartTime) {
			latestSegment = &workSegments[i]
		}
	}

	// リクエストの storeID と最新セグメントの StoreID が異なり、かつ StatusID が 3 以外の場合にエラーを返す
	if latestSegment.StoreID != storeID && latestSegment.StatusID != 3 {
		return fmt.Errorf("打刻する店舗が違います。")
	}

	latestSegment.EndTime = &now
	latestSegment.StatusID = 3 // 退勤
	if err := s.repo.UpdateWorkSegment(latestSegment); err != nil {
		return err
	}

	summary, err := s.repo.FindDailyWorkSummary(employeeID, workDate)
	if err != nil {
		return err
	}

	breakRecords, err := s.repo.FindBreakRecords(summary.ID)
	if err != nil {
		return err
	}

	var totalBreakTime time.Duration
	for _, record := range breakRecords {
		if record.BreakEnd != nil {
			totalBreakTime += record.BreakEnd.Sub(record.BreakStart)
		}
	}

	workDuration := latestSegment.EndTime.Sub(earliestSegment.StartTime)
	totalWorkTime := workDuration - totalBreakTime

	// 5分単位で切り下げる処理
	const roundTo = 5 * time.Minute
	totalWorkTimeTruncated := totalWorkTime.Truncate(roundTo)

	summary.StartTime = earliestSegment.StartTime
	summary.EndTime = latestSegment.EndTime
	summary.TotalBreakTime += totalBreakTime.Seconds() / 3600       // hours
	summary.TotalWorkTime = totalWorkTimeTruncated.Seconds() / 3600 // hours
	return s.repo.UpdateDailyWorkSummary(summary)
}

// 外出
func (s *AttendanceService) GoOut(employeeID uint, storeID uint) error {
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	workDate := now.Format("2006-01-02")

	summary, err := s.repo.FindDailyWorkSummary(employeeID, workDate)
	if err != nil {
		return err
	}

	workSegment, err := s.repo.FindLatestWorkSegment(employeeID)
	if err != nil {
		return err
	}

	// StoreID が異なる場合にエラーを返す
	if workSegment.StoreID != storeID {
		return fmt.Errorf("打刻する店舗が違います。")
	}

	breakRecord := model.BreakRecord{
		SummaryID:  summary.ID,
		BreakStart: now,
	}
	if err := s.repo.CreateBreakRecord(&breakRecord); err != nil {
		return err
	}

	workSegment.StatusID = 2 // 外出
	return s.repo.UpdateWorkSegment(workSegment)
}

// 戻り
func (s *AttendanceService) Return(employeeID uint, storeID uint) error {
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))

	workSegment, err := s.repo.FindWorkSegmentToReturn(employeeID)
	if err != nil {
		return err
	}

	if workSegment.StoreID == storeID {
		workSegment.StatusID = 1 // 出勤

		// サマリIDに紐づく休憩記録を更新
		breakRecords, err := s.repo.FindBreakRecords(workSegment.SummaryID)
		if err != nil {
			return err
		}
		for _, record := range breakRecords {
			if record.BreakEnd == nil {
				record.BreakEnd = &now
				if err := s.repo.UpdateBreakRecord(&record); err != nil {
					return err
				}
			}
		}

		return s.repo.UpdateWorkSegment(workSegment)
	}

	// 既存のセグメントのステータスと終了時間を更新
	existingSegments, err := s.repo.FindWorkSegmentsByDate(employeeID, now.Format("2006-01-02"))
	if err != nil {
		return err
	}

	for _, segment := range existingSegments {
		segment.StatusID = 3 // 他店舗での勤務終了としてステータスを更新
		segment.EndTime = &now
		if err := s.repo.UpdateWorkSegment(&segment); err != nil {
			return err
		}
	}

	// 新しいセグメントを作成
	newWorkSegment := model.WorkSegment{
		SummaryID:  workSegment.SummaryID,
		EmployeeID: employeeID,
		StoreID:    storeID,
		StartTime:  now,
		StatusID:   1, // 出勤
	}

	// サマリIDに紐づく休憩記録を更新
	breakRecords, err := s.repo.FindBreakRecords(workSegment.SummaryID)
	if err != nil {
		return err
	}
	for _, record := range breakRecords {
		if record.BreakEnd == nil {
			record.BreakEnd = &now
			if err := s.repo.UpdateBreakRecord(&record); err != nil {
				return err
			}
		}
	}

	return s.repo.CreateWorkSegment(&newWorkSegment)
}
