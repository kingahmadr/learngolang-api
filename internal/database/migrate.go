package database

import (
	"fmt"
	"learngolang-api/pkg/models"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		// &models.Product{},
		// &models.Order{},
	)
	if err != nil {
		fmt.Println("Migration failed:", err)
	} else {
		fmt.Println("✅ Database migrated successfully")
	}
}
