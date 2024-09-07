package services

import (
	"strconv"
	"strings"
	"time"

	model "github.com/techyoichiro/jobreco-api/domain/models"
	"github.com/techyoichiro/jobreco-api/domain/repositories"
)

func formatDate(date time.Time) string {
	// "2006-01-02"は、フォーマットの基準となる日時
	return date.Format("1/2(日)") // 月/日(曜日) の形式でフォーマット
}

func formatTime(t time.Time) string {
	return t.Format("15:04") // 時:分 の形式でフォーマット
}

type SummaryService struct {
	repo repositories.SummaryRepository
}

type BreakRecordResponse struct {
	ID         uint   `json:"ID"`
	WorkDate   string `json:"WorkDate"`
	BreakStart string `json:"BreakStart"`
	BreakEnd   string `json:"BreakEnd,omitempty"`
}

type SegmentsResponse struct {
	ID        uint   `json:"ID"`
	WorkDate  string `json:"WorkDate"`
	StoreID   uint   `json:"StoreID"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime,omitempty"`
}

type SummaryResponse struct {
	ID            uint                  `json:"ID"`
	WorkDate      string                `json:"WorkDate"`
	StartTime     string                `json:"StartTime"`
	EndTime       string                `json:"EndTime,omitempty"`
	TotalWorkTime float64               `json:"TotalWorkTime"`
	BreakRecords  []BreakRecordResponse `json:"BreakRecords"`
	Overtime      float64               `json:"Overtime"`
	Remarks       string                `json:"Remarks"`
	HourlyPay     int                   `json:"HourlyPay"`
}

type UpdateSummaryRequest struct {
	WorkSegments []SegmentsResponse    `json:"workSegments"`
	BreakRecords []BreakRecordResponse `json:"breakRecords"`
}

func NewSummaryService(repo repositories.SummaryRepository) *SummaryService {
	return &SummaryService{repo: repo}
}

// GetAllEmployee 全従業員の名前を取得するサービス
func (s *SummaryService) GetAllEmployee() ([]model.Employee, error) {
	return s.repo.GetAllEmployee()
}

// 指定した従業員IDの勤怠情報を取得するサービス
func (s *SummaryService) GetSummary(employeeID uint, year int, month int) ([]SummaryResponse, error) {
	summaries, err := s.repo.GetSummary(employeeID, year, month)
	if err != nil {
		return nil, err
	}

	hourlyPay, err := s.repo.GetHourlyPay(employeeID)
	if err != nil {
		return nil, err
	}

	response := []SummaryResponse{}
	for _, summary := range summaries {
		var breakRecords []BreakRecordResponse
		for _, breakRecord := range summary.BreakRecords {
			breakStart := formatTime(breakRecord.BreakStart)
			var breakEnd string
			if breakRecord.BreakEnd != nil {
				breakEnd = formatTime(*breakRecord.BreakEnd)
			}

			breakRecords = append(breakRecords, BreakRecordResponse{
				ID:         breakRecord.ID,
				BreakStart: breakStart,
				BreakEnd:   breakEnd,
			})
		}

		startTime := formatTime(summary.StartTime)
		var endTime string
		if summary.EndTime != nil {
			endTime = formatTime(*summary.EndTime)
		}

		workDate := formatDate(summary.WorkDate)

		response = append(response, SummaryResponse{
			ID:            summary.ID,
			WorkDate:      workDate,
			StartTime:     startTime,
			EndTime:       endTime,
			TotalWorkTime: summary.TotalWorkTime,
			BreakRecords:  breakRecords,
			Overtime:      calculateOvertime(summary),
			Remarks:       generateRemarks(summary),
			HourlyPay:     hourlyPay,
		})
	}

	return response, nil
}

func calculateOvertime(summary model.DailyWorkSummary) float64 {
	var overtime float64

	// 勤務開始時間と終了時間を取得
	startTime := summary.StartTime
	endTime := summary.EndTime

	// 勤務時間がある場合、時間外労働を計算
	if !startTime.IsZero() && endTime != nil {
		// 勤務時間を計算
		workDuration := endTime.Sub(startTime).Hours()

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

func generateRemarks(summary model.DailyWorkSummary) string {
	var remarks []string

	// 勤務セグメントがある場合、備考欄を生成
	for _, segment := range summary.WorkSegments {
		// 時間をフォーマットし、StoreID を string 型に変換
		startTime := segment.StartTime.Format("15:04")
		endTime := "-"
		if segment.EndTime != nil {
			endTime = segment.EndTime.Format("15:04")
		}
		storeID := strconv.FormatUint(uint64(segment.StoreID), 10)

		// フォーマットした文字列を作成
		segmentRemark := startTime + "-" + endTime + " " + storeID
		remarks = append(remarks, segmentRemark)
	}

	// 備考欄をカンマで連結
	return strings.Join(remarks, ", ")
}

// サマリ１件を取得
func (s *SummaryService) GetSummaryBySummaryID(summaryID uint) (*SummaryResponse, error) {
	summary, err := s.repo.GetSummaryBySummaryID(summaryID)
	if err != nil {
		return nil, err
	}

	var breakRecords []BreakRecordResponse
	for _, breakRecord := range summary.BreakRecords {
		breakStart := formatTime(breakRecord.BreakStart)
		var breakEnd string
		if breakRecord.BreakEnd != nil {
			breakEnd = formatTime(*breakRecord.BreakEnd)
		}

		breakRecords = append(breakRecords, BreakRecordResponse{
			ID:         breakRecord.ID,
			BreakStart: breakStart,
			BreakEnd:   breakEnd,
		})
	}

	remarks := generateRemarks(*summary)

	response := SummaryResponse{
		ID:            summary.ID,
		WorkDate:      formatDate(summary.WorkDate),
		StartTime:     formatTime(summary.StartTime),
		EndTime:       formatTimeIfNotNil(summary.EndTime),
		TotalWorkTime: summary.TotalWorkTime,
		BreakRecords:  breakRecords,
		Overtime:      calculateOvertime(*summary),
		Remarks:       remarks,
	}

	return &response, nil
}

// nil チェックとデリファレンスを行う関数
func formatTimeIfNotNil(t *time.Time) string {
	if t == nil {
		return "" // または適切なデフォルト値
	}
	return formatTime(*t)
}

// GetSegmentsBySummaryID returns formatted work segments for a given summary ID
func (s *SummaryService) GetSegmentsBySummaryID(summaryID uint) ([]SegmentsResponse, error) {
	segments, err := s.repo.FindWorkSegmentsBySummaryID(summaryID)
	if err != nil {
		return nil, err
	}

	var response []SegmentsResponse
	for _, segment := range segments {
		startTime := formatTime(segment.StartTime)
		endTime := ""
		if segment.EndTime != nil {
			endTime = formatTime(*segment.EndTime)
		}

		workDate := formatDate(segment.StartTime) // Using StartTime to determine WorkDate

		response = append(response, SegmentsResponse{
			ID:        segment.ID,
			WorkDate:  workDate,
			StoreID:   segment.StoreID,
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	return response, nil
}

// GetBreakRecordsBySummaryID returns formatted break records for a given summary ID
func (s *SummaryService) GetBreakRecordBySummaryID(summaryID uint) (*BreakRecordResponse, error) {
	breakRecord, err := s.repo.FindBreakRecords(summaryID)
	if err != nil {
		return nil, err
	}

	if breakRecord == nil {
		return nil, nil
	}

	breakStart := formatTime(breakRecord.BreakStart)
	breakEnd := ""
	if breakRecord.BreakEnd != nil {
		breakEnd = formatTime(*breakRecord.BreakEnd)
	}

	workDate := formatDate(breakRecord.BreakStart)

	response := BreakRecordResponse{
		ID:         breakRecord.ID,
		WorkDate:   workDate,
		BreakStart: breakStart,
		BreakEnd:   breakEnd,
	}

	return &response, nil
}

// セグメントIDで指定された勤怠情報を更新するサービス
// func (s *SummaryService) UpdateSummary(employeeID uint, storeID uint, segmentUpdates []SegmentUpdate, breakUpdate BreakUpdate) error {
// 	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
// 	workDate := now.Format("2006-01-02")

// 	// まず、セグメントの更新を行う
// 	var earliestSegment, latestSegment *model.WorkSegment
// 	for _, update := range segmentUpdates {
// 		segment, err := s.repo.FindWorkSegmentByID(update.ID)
// 		if err != nil {
// 			return err
// 		}

// 		segment.StartTime = update.StartTime
// 		if update.EndTime != nil {
// 			segment.EndTime = update.EndTime
// 		}
// 		segment.StatusID = update.StatusID

// 		if err := s.repo.UpdateWorkSegment(segment); err != nil {
// 			return err
// 		}

// 		// 最も早いセグメントと最も遅いセグメントを特定
// 		if earliestSegment == nil || segment.StartTime.Before(earliestSegment.StartTime) {
// 			earliestSegment = segment
// 		}
// 		if latestSegment == nil || (segment.EndTime != nil && segment.EndTime.After(latestSegment.EndTime)) {
// 			latestSegment = segment
// 		}
// 	}

// 	// 次に、休憩記録を更新する
// 	breakRecord, err := s.repo.FindBreakRecordByID(breakUpdate.ID)
// 	if err != nil {
// 		return err
// 	}

// 	breakRecord.BreakStart = breakUpdate.BreakStart
// 	if breakUpdate.BreakEnd != nil {
// 		breakRecord.BreakEnd = breakUpdate.BreakEnd
// 	}

// 	if err := s.repo.UpdateBreakRecord(breakRecord); err != nil {
// 		return err
// 	}

// 	// セグメントと休憩記録が更新されたので、サマリーを更新する
// 	summary, err := s.repo.FindDailyWorkSummary(employeeID, workDate)
// 	if err != nil {
// 		return err
// 	}

// 	// 総休憩時間を計算
// 	breakRecords, err := s.repo.FindBreakRecords(summary.ID)
// 	if err != nil {
// 		return err
// 	}

// 	var totalBreakTime time.Duration
// 	for _, record := range breakRecords {
// 		if record.BreakEnd != nil {
// 			totalBreakTime += record.BreakEnd.Sub(record.BreakStart)
// 		}
// 	}

// 	// 総勤務時間を計算
// 	if latestSegment.EndTime == nil {
// 		latestSegment.EndTime = &now
// 	}
// 	workDuration := latestSegment.EndTime.Sub(earliestSegment.StartTime)
// 	totalWorkTime := workDuration - totalBreakTime

// 	// 5分単位で切り下げる
// 	const roundTo = 5 * time.Minute
// 	totalWorkTimeTruncated := totalWorkTime.Truncate(roundTo)

// 	summary.StartTime = earliestSegment.StartTime
// 	summary.EndTime = latestSegment.EndTime
// 	summary.TotalBreakTime = totalBreakTime.Seconds() / 3600        // hours
// 	summary.TotalWorkTime = totalWorkTimeTruncated.Seconds() / 3600 // hours

// 	return s.repo.UpdateSummary(summary)
// }
