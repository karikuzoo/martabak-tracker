package models

import (
	"fmt"

	libsql "github.com/ekristen/gorm-libsql"
	"gorm.io/gorm"
)

type DBModel struct {
	Order OrderModel
	User  UserModel
	DB    *gorm.DB
}

func InitDB(dataSourceName string) (*DBModel, error) {
	db, err := gorm.Open(libsql.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	err = db.AutoMigrate(&Order{}, &OrderItem{}, &User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database %v", err)
	}

	dbModel := &DBModel{
		DB:    db,
		Order: OrderModel{DB: db},
		User:  UserModel{DB: db},
	}
	return dbModel, nil
}
