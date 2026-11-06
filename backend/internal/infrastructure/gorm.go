package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func dbSetup() {
	var err error

	l := gormLogger.Default.LogMode(gormLogger.Silent)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	config := &gorm.Config{
		Logger: l,
	}

	db, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		panic(err)
	}
}
