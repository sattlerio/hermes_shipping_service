package models

import "github.com/jinzhu/gorm"

type ShippingRules2Countries struct {
	gorm.Model

	CountryId		string	`gorm:"size:3;not null" json:"country_id"`
	ShippingRuleId 	int64	`gorm:"not null" json:"-"`
}

