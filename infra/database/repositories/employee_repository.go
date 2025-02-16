package repository

import (
	"log"
	"time"

	model "github.com/techyoichiro/jobreco-api/domain/models"
	"gorm.io/gorm"
)

type EmployeeRepositoryImpl struct {
	DB *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepositoryImpl {
	return &EmployeeRepositoryImpl{DB: db}
}

// ログイン
func (r *EmployeeRepositoryImpl) FindEmpByLoginID(loginID string) (*model.Employee, error) {
	var employee model.Employee
	// ログインIDに紐づくユーザーを取得
	if err := r.DB.Where("login_id = ?", loginID).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	log.Print(employee)
	return &employee, nil
}

// ユーザー作成
func (r *EmployeeRepositoryImpl) CreateEmp(employee *model.Employee) error {
	return r.DB.Create(employee).Error
}

// ステータス取得
func (r *EmployeeRepositoryImpl) GetStatusByEmpID(employeeID uint) (int, error) {
	var attendance model.Attendance

	// 現在の日本時間を取得
	jst, _ := time.LoadLocation("Asia/Tokyo")
	today := time.Now().In(jst).Format("2006-01-02")

	// 従業員IDに紐づいた今日のステータスを取得
	if err := r.DB.Where("employee_id = ? AND DATE(created_at) = ?", employeeID, today).
		Order("created_at DESC").
		First(&attendance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // レコードが見つからない場合は 0 を返す
		}
		return 0, err
	}
	return attendance.StatusID, nil
}

// ログインID取得
func (r *EmployeeRepositoryImpl) GetLoginIDByEmpID(employeeID string) (string, error) {
	var employee model.Employee

	// 従業員IDに紐づいたログインIDを取得
	if err := r.DB.Where("id = ?", employeeID).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return employee.LoginID, nil
}

// パスワード更新
func (r *EmployeeRepositoryImpl) UpdateEmpPassword(employee *model.Employee) error {
	return r.DB.Model(employee).Update("password", employee.Password).Error
}

// 従業員取得
func (r *EmployeeRepositoryImpl) FindEmpByEmpID(employeeID int) (*model.Employee, error) {
	var employee model.Employee
	if err := r.DB.Where("id = ?", employeeID).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &employee, nil
}

// 従業員更新
func (r *EmployeeRepositoryImpl) UpdateEmployee(employee *model.Employee) error {
	return r.DB.Model(employee).Updates(employee).Error
}
