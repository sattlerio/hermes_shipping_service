package models

import "github.com/jinzhu/gorm"

type ShippingRule struct {
	gorm.Model

	Name 			*string `gorm:"size:255;not null" json:"name"`
	ShippingRuleId  string  `gorm:"size:255;not null;unique" json:"shipping_rule_id"`
	
	CompanyId 		string	`gorm:"size:255;not null;" json:"company_id"`
	SelfShipping	*bool   `gorm:"not null" json:"self_shipping"`
	
}
