package models

import "github.com/jinzhu/gorm"

type Transaction struct {
	gorm.Model
	Type   uint
	Amount int
}
