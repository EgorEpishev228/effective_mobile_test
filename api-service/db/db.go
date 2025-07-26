package db

import (
	"emtest/api-service/config"
	"emtest/api-service/subscription"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func FormatDNS(db *config.Database) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		db.DNS.Host, db.DNS.User, db.DNS.Password, db.DNS.DbName, db.DNS.Port, db.DNS.SslMode,
	)
}

var DB *gorm.DB

func InitDB(database config.Database) error {
	var err error
	DB, err = gorm.Open(
		postgres.Open(FormatDNS(&database)),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}

	// err = DB.Migrator().DropTable(&subscription.Subscription{})
	// if err != nil {
	// 	logrus.Printf("Failed to drop table: %v", err)
	// }

	err = DB.AutoMigrate(&subscription.Subscription{})
	if err != nil {
		return err
	}

	logrus.Info("Database initiated")
	return nil
}
