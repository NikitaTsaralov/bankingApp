package models

import "github.com/jinzhu/gorm"

type TransactionType struct {
	gorm.Model
	Name string
}
