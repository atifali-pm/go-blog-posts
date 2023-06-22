package models

import "gorm.io/gorm"

func MigratePost(db *gorm.DB) error {
	if err := db.AutoMigrate(&Post{}); err != nil {
		return err
	}

	return nil
}
