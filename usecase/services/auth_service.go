package services

import (
	"errors"
	"log"

	"github.com/techyoichiro/jobreco-api/crypto"
	model "github.com/techyoichiro/jobreco-api/domain/models"
	"github.com/techyoichiro/jobreco-api/domain/repositories"
)

type AuthService struct {
	repo repositories.EmployeeRepository
}

func NewAuthService(repo repositories.EmployeeRepository) *AuthService {
	return &AuthService{repo: repo}
}

// サインアップ
func (s *AuthService) Signup(name, loginID, password string) (*model.Employee, error) {
	// メールアドレスの存在チェック
	existingEmployee, err := s.repo.FindEmpByLoginID(loginID)
	if err != nil {
		log.Printf("Error finding employee by loginID: %v", err)
		return nil, err
	}
	if existingEmployee != nil {
		log.Printf("Employee already exists: %v", existingEmployee)
		return nil, errors.New("同一の従業員IDが既に登録されています")
	}

	// パスワードの暗号化
	encryptedPw, err := crypto.PasswordEncrypt(password)
	if err != nil {
		log.Printf("Error encrypting password: %v", err)
		return nil, err
	}

	// メールアドレスの暗号化
	encryptedEmail, err := crypto.EncryptEmail(loginID)
	if err != nil {
		log.Printf("Error encrypting email: %v", err)
		return nil, err
	}

	employee := &model.Employee{
		Name:      name,
		LoginID:   encryptedEmail,
		Password:  encryptedPw,
		RoleID:    1,    // 初期値には従業員権限を付与
		HourlyPay: 1112, // 初期値には時給1112円を付与
	}

	err = s.repo.CreateEmp(employee)
	if err != nil {
		log.Printf("Error creating employee: %v", err)
		return nil, err
	}

	return employee, nil
}

// ログイン
func (s *AuthService) Login(loginID, password string) (*model.Employee, error) {
	emp, err := s.repo.FindEmpByLoginID(loginID)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, errors.New("ログインIDが一致するユーザーが存在しません。")
	}

	err = crypto.CompareHashAndPassword(emp.Password, password)
	if err != nil {
		return nil, errors.New("パスワードが一致しませんでした。")
	}

	return emp, nil
}

// employee_id に紐づく status_id を取得
func (s *AuthService) GetStatusByEmpID(employeeID uint) (int, error) {
	statusID, err := s.repo.GetStatusByEmpID(employeeID)
	if err != nil {
		return 0, err
	}
	return statusID, nil
}

// employee_id に紐づく login_id を取得
func (s *AuthService) GetLoginIDByEmpID(employeeID string) (string, error) {
	loginID, err := s.repo.GetLoginIDByEmpID(employeeID)
	if err != nil {
		return "", err
	}
	return loginID, nil
}
