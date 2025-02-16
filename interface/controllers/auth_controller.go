package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/techyoichiro/jobreco-api/usecase/services"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (ac *AuthController) PostSignup(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		LoginID  string `json:"login_id"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	employee, err := ac.service.Signup(req.Name, req.LoginID, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"employee": employee})
}

func (ac *AuthController) PostLogin(c *gin.Context) {
	var request struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 巡業員に紐づくlogin_idを取得
	LoginID, err := ac.service.GetLoginIDByEmpID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve login_id"})
		return
	}

	// ログイン処理
	emp, err := ac.service.Login(LoginID, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// ログインユーザーに紐づく status_id を取得
	statusID, err := ac.service.GetStatusByEmpID(emp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve status_id"})
		return
	}

	// ユーザー情報と status_id を返す
	c.JSON(http.StatusOK, gin.H{
		"employee": gin.H{
			"ID":          emp.ID,
			"Name":        emp.Name,
			"RoleID":      emp.RoleID,
			"HourlyPay":   emp.HourlyPay,
			"CompetentID": emp.CompetentStoreID,
		},
		"status_id": statusID,
	})
}

// パスワード変更
func (ac *AuthController) PostChangePassword(c *gin.Context) {
	var req struct {
		ID              string `json:"employee_id"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	// JSONをパース
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 従業員IDから暗号化されたlogin_idを取得
	loginID, err := ac.service.GetLoginIDByEmpID(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve login_id"})
		return
	}

	// パスワード更新サービスを呼び出す
	if err := ac.service.UpdatePassword(loginID, req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "パスワードが正常に変更されました"})
}

// アカウント設定更新
func (ac *AuthController) PostUpdateAccount(c *gin.Context) {
	var req struct {
		ID               string `json:"employee_id"`
		Name             string `json:"user_name"`
		HourlyPay        string `json:"hourly_pay"`
		CompetentStoreID string `json:"competent_id"`
	}

	// JSONをパース
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 各フィールドを数値に変換する
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id must be a valid number"})
		return
	}

	hourlyPay, err := strconv.Atoi(req.HourlyPay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hourly_pay must be a valid number"})
		return
	}

	competentStoreID, err := strconv.Atoi(req.CompetentStoreID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "店舗を選択してください"})
		return
	}

	// 数値に変換した値をサービス層に渡して更新処理を実行
	if err := ac.service.UpdateAccount(id, req.Name, hourlyPay, competentStoreID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "アカウント情報が正常に更新されました"})
}
