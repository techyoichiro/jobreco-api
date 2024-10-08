package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionDB() (*gorm.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// TimeZoneパラメータを追加
	if !strings.Contains(databaseURL, "TimeZone=") {
		databaseURL += "&TimeZone=Asia/Tokyo"
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
