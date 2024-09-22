package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techyoichiro/jobreco-api/usecase/services"
)

type AttendanceController struct {
	service *services.AttendanceService
}

func NewAttendanceController(service *services.AttendanceService) *AttendanceController {
	return &AttendanceController{
		service: service,
	}
}

// 出勤
func (ac *AttendanceController) PostClockIn(c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id"`
		StoreID    uint `json:"store_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ac.service.ClockIn(req.EmployeeID, req.StoreID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusID": 1})
}

// 退勤
func (ac *AttendanceController) PostClockOut(c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id"`
		StoreID    uint `json:"store_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ac.service.ClockOut(req.EmployeeID, req.StoreID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusID": 3})
}

// 外出
func (ac *AttendanceController) PostGoOut(c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id"`
		StoreID    uint `json:"store_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ac.service.GoOut(req.EmployeeID, req.StoreID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusID": 2})
}

// 戻り
func (ac *AttendanceController) PostReturn(c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id"`
		StoreID    uint `json:"store_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ac.service.Return(req.EmployeeID, req.StoreID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusID": 4})
}
