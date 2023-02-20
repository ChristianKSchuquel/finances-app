package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	TotalIncome   []Income
	TotalExpenses []Expense
}

type Income struct {
	gorm.Model
	UserID uint
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type Expense struct {
	gorm.Model
	UserID uint
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}
