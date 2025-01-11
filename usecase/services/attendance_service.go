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

	// 打刻日の勤怠記録があるか確認
	attendance, err := s.repo.FindAttendance(employeeID, workDate)
	// 新規作成
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			attendance = &model.Attendance{
				EmployeeID: employeeID,
				WorkDate:   now,
				StartTime1: &now,
				StoreID1:   storeID,
				StatusID:   1, // 出勤
			}
			if err := s.repo.CreateAttendance(attendance); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// 退勤
func (s *AttendanceService) ClockOut(employeeID uint, storeID uint) error {
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	workDate := now.Format("2006-01-02")

	attendance, err := s.repo.FindAttendance(employeeID, workDate)
	if err != nil {
		// 見つからなかった場合は、前日の日付を求めて再検索
		// (あるいは「最終の未完了レコード」を探す方式に切り替える、など)
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
		attendance, err = s.repo.FindAttendance(employeeID, yesterday)
		if err != nil {
			return err
		}
		if attendance == nil {
			return fmt.Errorf("退勤対象の勤怠が見つかりません")
		}
	}

	// リクエストのstoreIDと最新の勤怠記録のStoreIDが異なり、かつStatusIDが3以外の場合にエラーを返す
	if attendance.StartTime2 == nil {
		if attendance.StoreID1 != storeID && attendance.StatusID != 3 {
			return fmt.Errorf("打刻する店舗が違います。")
		} else {
			attendance.EndTime1 = &now
			attendance.StatusID = 3 // 退勤
		}
	} else {
		if *attendance.StoreID2 != storeID && attendance.StatusID != 3 {
			return fmt.Errorf("打刻する店舗が違います。")
		} else {
			attendance.EndTime2 = &now
			attendance.StatusID = 3 // 退勤
		}
	}

	return s.repo.UpdateAttendance(attendance)
}

// 外出
func (s *AttendanceService) GoOut(employeeID uint, storeID uint) error {
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	workDate := now.Format("2006-01-02")

	attendance, err := s.repo.FindAttendance(employeeID, workDate)
	if err != nil {
		return err
	}

	// StoreID が異なる場合にエラーを返す
	if attendance.StoreID1 != storeID {
		return fmt.Errorf("打刻する店舗が違います。")
	}

	attendance.BreakStart = &now
	attendance.StatusID = 2 // 外出
	return s.repo.UpdateAttendance(attendance)
}

// 戻り
func (s *AttendanceService) Return(employeeID uint, storeID uint) error {
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	workDate := now.Format("2006-01-02")

	attendance, err := s.repo.FindAttendance(employeeID, workDate)
	if err != nil {
		return err
	}

	attendance.StatusID = 4 // 休憩戻り
	attendance.BreakEnd = &now

	// 外出時と同じ店舗に戻る場合
	if attendance.StoreID1 == storeID {
		return s.repo.UpdateAttendance(attendance)
		// 外出時と別の店舗に戻る場合
	} else {
		attendance.EndTime1 = attendance.BreakStart
		attendance.StartTime2 = &now
		attendance.StoreID2 = &storeID
		return s.repo.UpdateAttendance(attendance)
	}
}
