package Entities

import "gorm.io/gorm"

type Quotation struct {
	Value float64
	gorm.Model
}
