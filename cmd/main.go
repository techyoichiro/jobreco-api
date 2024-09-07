package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/techyoichiro/jobreco-api/infra/database"
	repository "github.com/techyoichiro/jobreco-api/infra/database/repositories"
	"github.com/techyoichiro/jobreco-api/infra/router"
	controller "github.com/techyoichiro/jobreco-api/interface/controllers"
	"github.com/techyoichiro/jobreco-api/usecase/services"
)

func initialize() (*gin.Engine, *controller.AuthController, *controller.AttendanceController, *controller.SummaryController) {
	// データベース接続の設定
	db, err := database.ConnectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// リポジトリの初期化
	empRepo := repository.NewEmployeeRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	summaryRepo := repository.NewSummaryRepository(db)

	// サービス層の初期化
	authService := services.NewAuthService(empRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	summaryService := services.NewSummaryService(summaryRepo)

	// コントローラの初期化
	authController := controller.NewAuthController(authService)
	attendanceController := controller.NewAttendanceController(attendanceService)
	summaryController := controller.NewSummaryController(summaryService)

	// ルータの設定
	engine := router.SetupRouter(authController, attendanceController, summaryController)
	return engine, authController, attendanceController, summaryController
}

func main() {
	engine, _, _, _ := initialize()

	// サーバを8080ポートで起動
	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
