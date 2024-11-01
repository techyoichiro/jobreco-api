package model

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEmployeeCRUD(t *testing.T) {
	// テスト用のSQLiteデータベースを作成
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// テーブルを作成
	err = db.AutoMigrate(&Employee{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Create: 新しいEmployeeを作成
	employee := Employee{
		Name:      "Test Employee",
		LoginID:   "test_login",
		Password:  "secure_password",
		RoleID:    1,
		HourlyPay: 1000,
	}

	// データを保存
	if err := db.Create(&employee).Error; err != nil {
		t.Fatalf("failed to create employee: %v", err)
	}

	// Read: 保存したデータを取得
	var result Employee
	if err := db.First(&result, employee.ID).Error; err != nil {
		t.Fatalf("failed to read employee: %v", err)
	}

	// 保存されたデータが期待通りか確認
	if result.LoginID != employee.LoginID {
		t.Errorf("expected LoginID %v, got %v", employee.LoginID, result.LoginID)
	}
	if result.Name != employee.Name {
		t.Errorf("expected Name %v, got %v", employee.Name, result.Name)
	}
	if result.HourlyPay != employee.HourlyPay {
		t.Errorf("expected HourlyPay %v, got %v", employee.HourlyPay, result.HourlyPay)
	}

	// Update: データを更新
	newName := "Updated Employee"
	if err := db.Model(&result).Update("Name", newName).Error; err != nil {
		t.Fatalf("failed to update employee: %v", err)
	}

	// 更新されたデータを確認
	var updatedResult Employee
	if err := db.First(&updatedResult, employee.ID).Error; err != nil {
		t.Fatalf("failed to read updated employee: %v", err)
	}
	if updatedResult.Name != newName {
		t.Errorf("expected updated Name %v, got %v", newName, updatedResult.Name)
	}

	// Delete: データを削除
	if err := db.Delete(&Employee{}, employee.ID).Error; err != nil {
		t.Fatalf("failed to delete employee: %v", err)
	}

	// 削除されたことを確認
	var deletedResult Employee
	if err := db.First(&deletedResult, employee.ID).Error; err == nil {
		t.Errorf("expected record to be deleted, but it still exists")
	}
}
