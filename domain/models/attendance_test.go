package model

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAttendanceCRUD(t *testing.T) {
	// テスト用のSQLiteデータベースを作成
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// テーブルを作成
	err = db.AutoMigrate(&Attendance{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// テスト用のデータを作成
	startTime1 := time.Now()
	endTime1 := startTime1.Add(8 * time.Hour)
	attendance := Attendance{
		EmployeeID: 1,
		WorkDate:   time.Now(),
		StartTime1: &startTime1,
		EndTime1:   &endTime1,
		StoreID1:   1,
		StatusID:   1,
	}

	// Create: データを保存
	if err := db.Create(&attendance).Error; err != nil {
		t.Fatalf("failed to create attendance: %v", err)
	}

	// Read: 保存したデータを取得
	var result Attendance
	if err := db.First(&result, attendance.ID).Error; err != nil {
		t.Fatalf("failed to read attendance: %v", err)
	}

	// 保存されたデータが期待通りか確認
	if result.EmployeeID != attendance.EmployeeID {
		t.Errorf("expected EmployeeID %v, got %v", attendance.EmployeeID, result.EmployeeID)
	}

	// Update: データを更新
	newStoreID := uint(2)
	if err := db.Model(&result).Update("StoreID1", newStoreID).Error; err != nil {
		t.Fatalf("failed to update attendance: %v", err)
	}

	// 更新が反映されているか確認
	var updatedResult Attendance
	if err := db.First(&updatedResult, attendance.ID).Error; err != nil {
		t.Fatalf("failed to read updated attendance: %v", err)
	}
	if updatedResult.StoreID1 != newStoreID {
		t.Errorf("expected StoreID1 %v, got %v", newStoreID, updatedResult.StoreID1)
	}

	// Delete: データを削除
	if err := db.Delete(&Attendance{}, attendance.ID).Error; err != nil {
		t.Fatalf("failed to delete attendance: %v", err)
	}

	// 削除されたことを確認
	var deletedResult Attendance
	if err := db.First(&deletedResult, attendance.ID).Error; err == nil {
		t.Errorf("expected record to be deleted, but it still exists")
	}
}
