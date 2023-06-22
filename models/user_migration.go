package models

import "gorm.io/gorm"

func MigrateUser(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}

	return nil
}
