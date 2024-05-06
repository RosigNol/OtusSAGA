package db

import (
	"database/sql"
	"fmt"
	"item-service/config"
	"item-service/internal/models"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormModule = fx.Provide(NewDatabase)

type Database struct {
	*gorm.DB
	logger *zap.Logger
}

func NewDatabase(config config.Config, logger *zap.Logger) *gorm.DB {
	var err error
	var sqlDB *sql.DB

	logger.Info("Connecting to database...")
	gormDB, err := getDatabaseInstance(config)
	logger.Info("Database connected")
	db := &Database{gormDB, logger}

	db.RegisterTables()

	if err != nil {
		logger.Fatal("Database connection error", zap.Error(err))
	}
	sqlDB, err = db.DB.DB()
	if err != nil {
		logger.Fatal("sqlDB connection error", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db.DB
}

func getDatabaseInstance(config config.Config) (db *gorm.DB, err error) {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			config.Database.Host, config.Database.Username, config.Database.Password, config.Database.Name, config.Database.Port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

		if err != nil {
			return nil, fmt.Errorf("failed to connect database: %w", err)
		}
	return db, nil
}

func (d *Database) RegisterTables() {
	fmt.Println(d.DB)
	err := d.DB.AutoMigrate(
		models.Item{},
	)

	d.Exec("ALTER TABLE items ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('english', name)) STORED;")

	if err != nil {
		d.logger.Fatal("Database migration error", zap.Error(err))
		os.Exit(0)
	}
	d.logger.Info("Database migration success")
}
