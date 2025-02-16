package repositories

import (
	model "github.com/techyoichiro/jobreco-api/domain/models"
)

// EmployeeRepository
type EmployeeRepository interface {
	FindEmpByLoginID(loginID string) (*model.Employee, error)
	FindEmpByEmpID(employeeID int) (*model.Employee, error)
	CreateEmp(employee *model.Employee) error
	GetStatusByEmpID(employeeID uint) (int, error)
	GetLoginIDByEmpID(employeeID string) (string, error)
	UpdateEmpPassword(employee *model.Employee) error
	UpdateEmployee(employee *model.Employee) error
}
