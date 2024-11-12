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
func calculateWorkTime(attendance model.Attendance) string {
	// 勤務開始時間と終了時間を取得
	startTime := attendance.StartTime1
	var endTime *time.Time
	if attendance.EndTime2 != nil {
		endTime = attendance.EndTime2
	} else {
		endTime = attendance.EndTime1
	}

	// 勤務開始時間または終了時間がnilの場合は0時間を返却
	if startTime == nil || endTime == nil {
		return "0.0"
	}

	// 時間を5分単位で切り下げるための定数
	const roundTo = 5 * time.Minute

	// 勤務時間を5分単位で丸める
	startTimeRounded := startTime.Truncate(roundTo)
	endTimeRounded := endTime.Truncate(roundTo)

	// 勤務時間を計算
	workDuration := endTimeRounded.Sub(startTimeRounded)

	// 休憩時間を計算
	var breakDuration time.Duration
	if attendance.BreakStart != nil && attendance.BreakEnd != nil {
		breakStartRounded := attendance.BreakStart.Truncate(roundTo)
		breakEndRounded := attendance.BreakEnd.Truncate(roundTo)

		// 休憩時間を5分単位で丸めた後に計算
		breakDuration = breakEndRounded.Sub(breakStartRounded)
	}

	// 実勤務時間（勤務時間 - 休憩時間）を計算
	actualWorkDuration := workDuration - breakDuration

	// 実勤務時間を時間単位で返却する
	hours := actualWorkDuration.Seconds() / 3600
	formattedHours := fmt.Sprintf("%.2f", hours)
	return formattedHours
}

// 22時以降の勤務時間を計算
func calculateOvertime(attendance model.Attendance) float64 {
	const overtimeThresholdHour = 22
	const roundTo = 5 * time.Minute // 5分刻み

	// 勤務終了時間を取得 (EndTime2が存在する場合はそれを使用)
	var endTime *time.Time
	if attendance.EndTime2 != nil {
		endTime = attendance.EndTime2
	} else {
		endTime = attendance.EndTime1
	}

	// 終了時間がnilの場合、0を返却
	if endTime == nil {
		return 0.0
	}

	// 勤務終了時間が22時を超えている場合の処理
	overtime := 0.0
	if endTime.Hour() >= overtimeThresholdHour {
		// 22時以降の時間を計算
		overtimeStart := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), overtimeThresholdHour, 0, 0, 0, endTime.Location())
		if endTime.After(overtimeStart) {
			overtimeDuration := endTime.Sub(overtimeStart)

			// 5分刻みに切り下げ
			roundedOvertime := overtimeDuration.Truncate(roundTo)

			// 時間単位で返却
			overtime = roundedOvertime.Seconds() / 3600
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
		BreakStart: formatTime(attendance.BreakStart),
		BreakEnd:   formatTime(attendance.BreakEnd),
		StoreID1:   attendance.StoreID1,
		StoreID2:   attendance.StoreID2,
		Remarks:    remarks,
	}

	return &response, nil
}

// IDで指定された勤怠情報を更新する
func (s *SummaryService) UpdateAttendance(attendanceResponse *model.AttendanceResponse) error {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return fmt.Errorf("failed to load location: %w", err)
	}

	// 受け取ったデータのwork_dateを取得
	workDate, err := s.repo.GetWorkDateByID(attendanceResponse.ID)
	if err != nil {
		return fmt.Errorf("failed to load workDate: %w", err)
	}

	// workDate を time.Time 型から受け取っている前提で、文字列に変換
	workDateStr := workDate.Format("2006-01-02")

	// 受け取ったデータをパースして、DBに合わせた形式に変換
	startTime1Str := workDateStr + " " + attendanceResponse.StartTime1
	startTime1, err := time.ParseInLocation("2006-01-02 15:04", startTime1Str, loc)
	if err != nil {
		return fmt.Errorf("invalid start time 1: %w", err)
	}

	var endTime1 *time.Time
	if attendanceResponse.EndTime1 != "" {
		endTime1Str := workDateStr + " " + attendanceResponse.EndTime1
		t, err := time.ParseInLocation("2006-01-02 15:04", endTime1Str, loc)
		if err != nil {
			return fmt.Errorf("invalid end time 1: %w", err)
		}
		endTime1 = &t
	}

	// 他の時間も同様に処理
	var startTime2 *time.Time
	if attendanceResponse.StartTime2 != "" {
		startTime2Str := workDateStr + " " + attendanceResponse.StartTime2
		t, err := time.ParseInLocation("2006-01-02 15:04", startTime2Str, loc)
		if err != nil {
			return fmt.Errorf("invalid start time 2: %w", err)
		}
		startTime2 = &t
	}

	var endTime2 *time.Time
	if attendanceResponse.EndTime2 != "" {
		endTime2Str := workDateStr + " " + attendanceResponse.EndTime2
		t, err := time.ParseInLocation("2006-01-02 15:04", endTime2Str, loc)
		if err != nil {
			return fmt.Errorf("invalid end time 2: %w", err)
		}
		endTime2 = &t
	}

	var breakStart *time.Time
	if attendanceResponse.BreakStart != "" {
		breakStartStr := workDateStr + " " + attendanceResponse.BreakStart
		t, err := time.ParseInLocation("2006-01-02 15:04", breakStartStr, loc)
		if err != nil {
			return fmt.Errorf("invalid break start: %w", err)
		}
		breakStart = &t
	}

	var breakEnd *time.Time
	if attendanceResponse.BreakEnd != "" {
		breakEndStr := workDateStr + " " + attendanceResponse.BreakEnd
		t, err := time.ParseInLocation("2006-01-02 15:04", breakEndStr, loc)
		if err != nil {
			return fmt.Errorf("invalid break end: %w", err)
		}
		breakEnd = &t
	}

	// Attendanceモデルのインスタンスを作成
	attendance := &model.Attendance{
		ID:         attendanceResponse.ID,
		StartTime1: &startTime1,
		EndTime1:   endTime1,
		StartTime2: startTime2,
		EndTime2:   endTime2,
		BreakStart: breakStart,
		BreakEnd:   breakEnd,
		StoreID1:   attendanceResponse.StoreID1,
		StoreID2:   attendanceResponse.StoreID2,
	}

	return s.repo.UpdateAttendance(attendance)
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
	// 曜日を日本語にマッピング
	daysOfWeek := []string{"日", "月", "火", "水", "木", "金", "土"}
	weekday := daysOfWeek[date.Weekday()]
	return date.Format("1/2") + "(" + weekday + ")" // 月/日(曜日) の形式でフォーマット
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
