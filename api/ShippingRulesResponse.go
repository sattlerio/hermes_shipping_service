package api

import "hermesShippingRuleService/models"

type ShippingRuleResponse struct {
	Response *Response			   `json:"response"`
	Data 	 []models.ShippingRule `json:"body"`
}
