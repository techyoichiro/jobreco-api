package repositories

import (
	model "github.com/techyoichiro/jobreco-api/domain/models"
)

// EmployeeRepository
type EmployeeRepository interface {
	FindEmpByLoginID(loginID string) (*model.Employee, error)
	CreateEmp(employee *model.Employee) error
	GetStatusByEmpID(employeeID uint) (int, error)
	GetLoginIDByEmpID(employeeID string) (string, error)
}
