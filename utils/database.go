package utils

import (
	"fmt"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/structs"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	databaseType := GetStringConfig(consts.DatabaseType)
	if databaseType == "postgres" {
		DB = connectToPostgres()
	} else {
		DB = connectToSQLite()
	}
}

func connectToSQLite() *gorm.DB {
	sqliteDatabase := GetStringConfig(consts.SQLiteDatabase)
	db, err := gorm.Open(sqlite.Open(sqliteDatabase), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func connectToPostgres() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		GetStringConfig(consts.PostgresHost),
		GetStringConfig(consts.PostgresUser),
		GetStringConfig(consts.PostgresPassword),
		GetStringConfig(consts.PostgresDB),
		GetStringConfig(consts.PostgresPort),
		GetStringConfig(consts.PostgresTimezone),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func AutoMigrate() {
	DB.AutoMigrate(&structs.Candlestick{})
	DB.AutoMigrate(&structs.TradingStatus{})
	DB.AutoMigrate(&structs.Operation{})

	DB.AutoMigrate(&structs.User{})
	DB.AutoMigrate(&structs.Parameters{})
	DB.AutoMigrate(&structs.BotParameters{})
}

func InitDatabase() {
	ConnectDatabase()
	AutoMigrate()
}
