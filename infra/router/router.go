package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	controller "github.com/techyoichiro/jobreco-api/interface/controllers"
)

// SetupRouter sets up the routes for the application.
func SetupRouter(authController *controller.AuthController, attendanceController *controller.AttendanceController, summaryController *controller.SummaryController) *gin.Engine {
	router := gin.Default()

	// CORS設定を手動で追加
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://loaclhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// ルート設定
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/signup", authController.PostSignup)
		authRouter.POST("/login", authController.PostLogin)
	}

	attendanceRouter := router.Group("/attendance")
	{
		attendanceRouter.POST("/clockin", attendanceController.PostClockIn)
		attendanceRouter.POST("/clockout", attendanceController.PostClockOut)
		attendanceRouter.POST("/goout", attendanceController.PostGoOut)
		attendanceRouter.POST("/return", attendanceController.PostReturn)
	}

	summaryRouter := router.Group("/summary")
	{
		summaryRouter.GET("/init", summaryController.GetAllEmployee)
		summaryRouter.GET("/:employeeId/:year/:month", summaryController.GetSummary)
		summaryRouter.GET("/edit/:summaryID", summaryController.GetSummaryBySummaryID)
		// summaryRouter.POST("/edit/:segmentID", summaryController.UpdateSummaryBySegmentID)
	}

	return router
}