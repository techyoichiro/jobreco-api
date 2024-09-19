package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	model "github.com/techyoichiro/jobreco-api/domain/models"
	"github.com/techyoichiro/jobreco-api/usecase/services"
)

type SummaryController struct {
	service *services.SummaryService
}

func NewSummaryController(service *services.SummaryService) *SummaryController {
	return &SummaryController{service: service}
}

// 返却用
type EmployeeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateSummaryRequest struct {
	EmployeeID uint `json:"employeeID"`
	StoreID    uint `json:"storeID"`
}

// 全従業員のIDと名前を取得するハンドラー
func (sc *SummaryController) GetAllEmployee(c *gin.Context) {
	employees, err := sc.service.GetAllEmployee()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 必要なフィールドだけを抽出
	var employeeResponses []EmployeeResponse
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, EmployeeResponse{
			ID:   employee.ID,
			Name: employee.Name,
		})
	}

	// JSONとして返却
	c.JSON(http.StatusOK, employeeResponses)
}

// GetSummaryByEmpID 指定した従業員IDの勤怠情報を取得するハンドラー
func (sc *SummaryController) GetAttendance(c *gin.Context) {
	employeeIDStr := c.Param("employeeId")
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format"})
		return
	}

	response, err := sc.service.GetAttendance(uint(employeeID), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// 指定したattendanceIDの勤怠情報を取得するハンドラー
func (sc *SummaryController) GetAttendanceByID(c *gin.Context) {
	attendanceIDStr := c.Param("attendanceID")

	attendanceID, err := strconv.ParseUint(attendanceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance ID"})
		return
	}

	// uint64からuintに変換
	attendanceIDUint := uint(attendanceID)

	// サービスメソッド呼び出し
	attendance, err := sc.service.GetAttendanceByID(attendanceIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	response := gin.H{
		"attendance": attendance,
	}

	c.JSON(http.StatusOK, response)
}

// セグメントIDで指定された勤怠情報を更新するハンドラー
func (sc *SummaryController) UpdateAttendance(c *gin.Context) {
	var request *model.AttendanceResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sc.service.UpdateAttendance(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work segment updated successfully"})
}
