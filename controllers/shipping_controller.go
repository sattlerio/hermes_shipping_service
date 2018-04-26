package controllers

import (
	"net/http"
	"hermesShippingRuleService/helpers"
	"github.com/jinzhu/gorm"
	"hermesShippingRuleService/models"
	"hermesShippingRuleService/api"
	"encoding/json"
	"io/ioutil"
	"github.com/rs/xid"
)

var DbConn *gorm.DB
var err error

func CreateShippingRule(w http.ResponseWriter, r *http.Request) {
	helpers.Info.Println("new request to create shipping rule")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-id")

	if len(transactionId) <= 0 || len(userId) <= 0 {
		helpers.Warning.Println("got no user Id and no transaction id in the header")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode: 400, Message:"you have to be logged in to use this service", TransactionId:transactionId}

		json.NewEncoder(w).Encode(response)
		return
	}

	var shippingRule models.ShippingRule
	body, _ := ioutil.ReadAll(r.Body)


	err = json.Unmarshal(body, &shippingRule)

	if err != nil {
		helpers.Info.Println(transactionId + ": no valid shipping rule in Post body")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode:400,
			Message: "please submit a valid shipping rule object", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	if shippingRule.Name == nil || shippingRule.SelfShipping == nil {
		helpers.Info.Println(transactionId + ": no self shipping and / or no shipping rule name provided")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode:400,
			Message:"please provide a valid shipping rule object", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	guid := xid.New().String()
	shippingRule.ShippingRuleId = guid

	helpers.Info.Println(transactionId + ": successfully generate " + guid + " as id for the new shipping rule")

	err = DbConn.Create(&shippingRule).Error
	if err != nil {
		helpers.Warning.Println(err)
		helpers.Info.Println(transactionId + ": error with database communication abort transaction")
		w.WriteHeader(500)
		response := api.Response{Status:"ERROR", StatusCode: 500,
			Message: "internal server error", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	helpers.Info.Println(transactionId + userId)

	helpers.Info.Println("eerrr")

}
