package repository

import (
	"errors"

	model "github.com/techyoichiro/jobreco-api/domain/models"
	"gorm.io/gorm"
)

type SummaryRepositoryImpl struct {
	DB *gorm.DB
}

func NewSummaryRepository(db *gorm.DB) *SummaryRepositoryImpl {
	return &SummaryRepositoryImpl{DB: db}
}

// 全従業員の名前を取得するリポジトリメソッド
func (r *SummaryRepositoryImpl) GetAllEmployee() ([]model.Employee, error) {
	var employees []model.Employee
	if err := r.DB.Select("id, name").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

// GetSummary
func (r *SummaryRepositoryImpl) GetSummary(employeeID uint, year int, month int) ([]model.DailyWorkSummary, error) {
	var summaries []model.DailyWorkSummary

	// 休憩記録と勤怠セグメントをロード
	err := r.DB.Preload("BreakRecords").
		Preload("WorkSegments").
		Where("employee_id = ? AND EXTRACT(YEAR FROM work_date) = ? AND EXTRACT(MONTH FROM work_date) = ?", employeeID, year, month).
		Find(&summaries).Error
	if err != nil {
		return nil, err
	}

	// 店舗IDとそのセグメントをマップで取得
	workSegmentMap := make(map[uint][]model.WorkSegment)
	var workSegments []model.WorkSegment
	err = r.DB.Where("employee_id = ? AND EXTRACT(YEAR FROM start_time) = ? AND EXTRACT(MONTH FROM start_time) = ?", employeeID, year, month).
		Find(&workSegments).Error
	if err != nil {
		return nil, err
	}

	// 勤怠セグメントをマップに追加
	for _, segment := range workSegments {
		if _, exists := workSegmentMap[segment.SummaryID]; !exists {
			workSegmentMap[segment.SummaryID] = []model.WorkSegment{}
		}
		workSegmentMap[segment.SummaryID] = append(workSegmentMap[segment.SummaryID], segment)
	}

	// DailyWorkSummaryにWorkSegmentsを割り当て
	for i := range summaries {
		if segments, exists := workSegmentMap[summaries[i].ID]; exists {
			summaries[i].WorkSegments = segments
		}
	}

	return summaries, nil
}

// 指定した勤怠記録のセグメントを取得するリポジトリメソッド
func (r *SummaryRepositoryImpl) GetWorkSegments(summaryID uint) []model.WorkSegment {
	var segments []model.WorkSegment
	r.DB.Where("summary_id = ?", summaryID).Find(&segments)
	return segments
}

// 従業員IDから時給を取得する
func (r *SummaryRepositoryImpl) GetHourlyPay(employeeID uint) (int, error) {
	var employee model.Employee
	if err := r.DB.Where("id = ?", employeeID).First(&employee).Error; err != nil {
		return 0, err
	}
	return employee.HourlyPay, nil
}

// サマリIDから勤怠情報を取得する
func (r *SummaryRepositoryImpl) GetSummaryBySummaryID(summaryID uint) (*model.DailyWorkSummary, error) {
	var summary model.DailyWorkSummary

	// 休憩記録と勤怠セグメントをロード
	err := r.DB.Preload("BreakRecords").
		Preload("WorkSegments").
		Where("id = ?", summaryID).
		First(&summary).Error
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

// UpdateSummary 勤怠記録を更新するリポジトリメソッド
func (r *SummaryRepositoryImpl) UpdateSummary(summary *model.DailyWorkSummary) error {
	// Begin a new transaction
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update the DailyWorkSummary
	if err := tx.Save(summary).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// サマリIDから勤怠セグメントを取得する
func (r *SummaryRepositoryImpl) FindWorkSegmentsBySummaryID(summaryID uint) ([]model.WorkSegment, error) {
	var segments []model.WorkSegment

	//
	err := r.DB.Where("summary_id = ?", summaryID).Find(&segments).Error
	if err != nil {
		return nil, err
	}

	return segments, nil
}

func (r *SummaryRepositoryImpl) FindBreakRecords(summaryID uint) (*model.BreakRecord, error) {
	var breakRecord model.BreakRecord
	// レコードが見つからない場合に nil を返すために First を使用
	err := r.DB.Where("summary_id = ?", summaryID).First(&breakRecord).Error
	if err != nil {
		// レコードが存在しない場合は nil を返す
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &breakRecord, nil
}

func (repo *SummaryRepositoryImpl) FindSummaryBySegmentID(segmentID uint) (*model.DailyWorkSummary, error) {
	var summary model.DailyWorkSummary
	err := repo.DB.Joins("JOIN work_segments ON work_segments.summary_id = daily_work_summaries.id").
		Where("work_segments.id = ?", segmentID).First(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (repo *SummaryRepositoryImpl) UpdateWorkSegment(segment *model.WorkSegment) error {
	return repo.DB.Save(segment).Error
}

func (repo *SummaryRepositoryImpl) UpdateBreakRecord(breakRecord *model.BreakRecord) error {
	return repo.DB.Save(breakRecord).Error
}

func (r *SummaryRepositoryImpl) FindWorkSegmentByID(ID uint) ([]model.WorkSegment, error) {
	var segments []model.WorkSegment
	if err := r.DB.Where("id = ?", ID).Find(&segments).Error; err != nil {
		return nil, err
	}
	return segments, nil
}
