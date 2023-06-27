package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Post{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Review{}); err != nil {
		return err
	}

	return nil
}
