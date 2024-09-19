package services

import (
	"fmt"
	"strconv"
	"time"

	model "github.com/techyoichiro/jobreco-api/domain/models"
	"github.com/techyoichiro/jobreco-api/domain/repositories"
)

type SummaryService struct {
	repo repositories.SummaryRepository
}

func NewSummaryService(repo repositories.SummaryRepository) *SummaryService {
	return &SummaryService{repo: repo}
}

// GetAllEmployee 全従業員の名前を取得するサービス
func (s *SummaryService) GetAllEmployee() ([]model.Employee, error) {
	return s.repo.GetAllEmployee()
}

// 指定した従業員IDの勤怠情報を取得するサービス
func (s *SummaryService) GetAttendance(employeeID uint, year int, month int) ([]model.AttendanceResponse, error) {
	attendances, err := s.repo.GetAttendance(employeeID, year, month)
	if err != nil {
		return nil, err
	}

	hourlyPay, err := s.repo.GetHourlyPay(employeeID)
	if err != nil {
		return nil, err
	}

	response := []model.AttendanceResponse{}
	for _, attendance := range attendances {

		// 勤務取得
		workDate := formatDate(&attendance.WorkDate)

		// 勤務時刻取得
		startTime1 := formatTime(attendance.StartTime1)
		endTime1 := formatTime(attendance.EndTime1)
		startTime2 := formatTime(attendance.StartTime2)
		endTime2 := formatTime(attendance.EndTime2)

		// 休憩記録取得
		var breakStart string
		var breakEnd string
		if attendance.BreakStart != nil {
			breakStart = formatTime(attendance.BreakStart)
		}
		if attendance.BreakEnd != nil {
			breakEnd = formatTime(attendance.BreakEnd)
		}

		response = append(response, model.AttendanceResponse{
			ID:            attendance.ID,
			WorkDate:      workDate,
			StartTime1:    startTime1,
			EndTime1:      endTime1,
			StartTime2:    startTime2,
			EndTime2:      endTime2,
			TotalWorkTime: calculateWorkTime(attendance),
			BreakStart:    breakStart,
			BreakEnd:      breakEnd,
			Overtime:      calculateOvertime(attendance),
			Remarks:       generateRemarks(attendance),
			HourlyPay:     hourlyPay,
		})
	}

	return response, nil
}

// 勤務時間を計算
func calculateWorkTime(attendance model.Attendance) float64 {
	// 勤務開始時間と終了時間を取得
	startTime := attendance.StartTime1
	var endTime *time.Time
	if attendance.EndTime2 != nil {
		endTime = attendance.EndTime1
	} else {
		endTime = attendance.EndTime2
	}

	// 勤務開始時間または終了時間がnilの場合は0時間を返却
	if startTime == nil || endTime == nil {
		return 0.0
	}

	// 時間を5分単位で切り下げる
	const roundTo = 5 * time.Minute
	startTimeRounded := startTime.Truncate(roundTo)
	endTimeRounded := endTime.Truncate(roundTo)

	// 勤務時間を計算する
	workDuration := endTimeRounded.Sub(startTimeRounded)

	// 勤務時間を時間単位で返却する
	return workDuration.Seconds() / 3600
}

// 22時以降の勤務時間を計算
func calculateOvertime(attendance model.Attendance) float64 {
	var overtime float64

	// 勤務開始時間と終了時間を取得
	startTime := attendance.StartTime1
	var endTime *time.Time
	if attendance.EndTime2 != nil {
		endTime = attendance.EndTime1
	} else {
		endTime = attendance.EndTime2
	}

	// 勤務時間がある場合、時間外労働を計算
	if !startTime.IsZero() && endTime != nil {
		startTimeValue := *startTime
		// 勤務時間を計算
		workDuration := endTime.Sub(startTimeValue).Hours()

		// 22:00を超える部分を時間外労働として計算
		if endTime.Hour() > 22 {
			threshold := 22.0
			// 勤務終了時間が22:00を超える場合の計算
			overTimeStart := float64(startTime.Hour()) + float64(startTime.Minute())/60.0
			overtime = workDuration - (threshold - overTimeStart)

			// 15分刻みに丸める（切り下げ）
			overtime = float64(int(overtime*4)) / 4.0
		}

	}

	return overtime
}

// サマリ１件を取得
func (s *SummaryService) GetAttendanceByID(attendanceID uint) (*model.AttendanceResponse, error) {
	attendance, err := s.repo.GetAttendanceByID(attendanceID)
	if err != nil {
		return nil, err
	}

	remarks := generateRemarks(*attendance)

	response := model.AttendanceResponse{
		ID:         attendance.ID,
		WorkDate:   formatDate(&attendance.WorkDate),
		StartTime1: formatTime(attendance.StartTime1),
		EndTime1:   formatTimeIfNotNil(attendance.EndTime1),
		StartTime2: formatTime(attendance.StartTime2),
		EndTime2:   formatTime(attendance.EndTime2),
		BreakStart: formatTime(attendance.EndTime2),
		BreakEnd:   formatTime(attendance.EndTime2),
		StoreID1:   attendance.StoreID1,
		StoreID2:   attendance.StoreID2,
		Remarks:    remarks,
	}

	return &response, nil
}

// IDで指定された勤怠情報を更新する
func (s *SummaryService) UpdateAttendance(attendanceResponse *model.AttendanceResponse) error {
	return s.repo.UpdateAttendance(attendanceResponse)
}

// 備考欄生成
func generateRemarks(attendance model.Attendance) string {

	// 時間をフォーマットし、StoreID を string 型に変換
	startTime1 := attendance.StartTime1.Format("15:04")
	endTime1 := "-"
	if attendance.EndTime1 != nil {
		endTime1 = attendance.EndTime1.Format("15:04")
	}
	storeID1 := strconv.FormatUint(uint64(attendance.StoreID1), 10)

	remark1 := startTime1 + "-" + endTime1 + " " + storeID1
	// 2店舗で勤務していた場合
	if attendance.StartTime2 != nil {
		startTime2 := attendance.StartTime2.Format("15:04")
		endTime2 := "-"
		if attendance.EndTime2 != nil {
			endTime2 = attendance.EndTime2.Format("15:04")
		}
		storeID2 := strconv.FormatUint(uint64(*attendance.StoreID2), 10)
		remark2 := startTime2 + "-" + endTime2 + " " + storeID2
		return formatRemarks(remark1, remark2)
	}

	// 備考欄をカンマで連結
	return remark1
}

// 日付フォーマット
func formatDate(date *time.Time) string {
	// "2006-01-02"は、フォーマットの基準となる日時
	return date.Format("1/2(日)") // 月/日(曜日) の形式でフォーマット
}

// 時刻フォーマット
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04") // 時:分 の形式でフォーマット
}

// 時刻の文字列を time.Time 型に変換するヘルパー関数
func parseTime(timeStr string) *time.Time {
	t, _ := time.Parse("15:04", timeStr)
	return &t
}

// nil チェックとデリファレンスを行う関数
func formatTimeIfNotNil(t *time.Time) string {
	if t == nil {
		return "" // または適切なデフォルト値
	}
	return formatTime(t)
}

// カンマ区切り
func formatRemarks(segmentRemark1, segmentRemark2 string) string {
	return fmt.Sprintf("%s, %s", segmentRemark1, segmentRemark2)
}
