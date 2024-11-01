package router

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	controller "github.com/techyoichiro/jobreco-api/interface/controllers"
)

func TestSetupRouter(t *testing.T) {
	type args struct {
		authController       *controller.AuthController
		attendanceController *controller.AttendanceController
		summaryController    *controller.SummaryController
	}
	tests := []struct {
		name string
		args args
		want *gin.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetupRouter(tt.args.authController, tt.args.attendanceController, tt.args.summaryController); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetupRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
