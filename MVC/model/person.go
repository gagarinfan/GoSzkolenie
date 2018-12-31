package model

import "github.com/jinzhu/gorm"

type Person struct {
	gorm.Model
	// ID       string `json:"id,omitempty" gorm:"INDEX"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}
