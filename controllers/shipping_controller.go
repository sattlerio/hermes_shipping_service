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
	"github.com/gorilla/mux"
	client2 "hermesShippingRuleService/client"
	"os"
)

var DbConn *gorm.DB
var err error

var RequiredPermission int = 3

func DeleteShippingRule(w http.ResponseWriter, r *http.Request) {
	helpers.Info.Println("new request to delete shipping rule")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-uuid")

	if len(transactionId) <= 0 || len(userId) <= 0 {
		helpers.Warning.Println("got no user Id and no transaction id in the header")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode: 400, Message:"you have to be logged in to use this service", TransactionId:transactionId}

		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)
	company_id := params["company_id"]
	shippingRuleId := params["shipping_rule_id"]

	if len(company_id) <= 0 || len(shippingRuleId) <= 0 {
		helpers.Info.Println(transactionId + ": no company id or shipping rule id in request")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode:400, Message:"you have to submit a valid company id and shippingRuleId as url param", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	guardianClient := client2.GuardianClient{Host:os.Getenv("GUARDIAN_URL"), CompanyId:company_id, UserId:userId}

	guardianResponse, err := client2.CheckCompanyAndPermissionFromGuardian(guardianClient, RequiredPermission)
	if err != nil {
		helpers.Info.Println(transactionId + ": guardian host responded with errror, abort transaction")
		helpers.Info.Println(err)
		w.WriteHeader(500)
		response := api.Response{Status:"ERROR", StatusCode:500, Message:"internal server error", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	if !guardianResponse {
		helpers.Info.Println(transactionId + ": user is not allowed to access company / settings")
		w.WriteHeader(401)
		response := api.Response{Status:"ERROR", StatusCode:401, Message:"not allowed to access", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}
	helpers.Info.Println(transactionId + " user is allowed to access, continue with request")

	helpers.Info.Println(transactionId + ": trying to delete shipping rule with uuid and company id {" + shippingRuleId  + " / " + company_id + "}")
	dbResult := DbConn.Where("shipping_rule_id = ? AND company_id = ?", shippingRuleId, company_id).Delete(&models.ShippingRule{})

	if dbResult.Error != nil || dbResult.RowsAffected != 1 {
		helpers.Info.Println(transactionId + ": not possible to delete rows no result affected or db raised error")
		helpers.Info.Println(dbResult.Error)
		helpers.Info.Println(dbResult.RowsAffected)
		w.WriteHeader(500)
		response := api.Response{Status:"ERROR", StatusCode:500, Message:"error not possible to delete entry", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	helpers.Info.Println(transactionId + ": successfully deleted entry")
	w.WriteHeader(200)
	response := api.Response{Status:"OK", StatusCode:200, Message:"successfully deleted entry", TransactionId:transactionId}
	json.NewEncoder(w).Encode(response)
	return
}

func GetShippingRules(w http.ResponseWriter, r *http.Request) {
	helpers.Info.Println("new request to fetch shipping rules")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-uuid")

	if len(transactionId) <= 0 || len(userId) <= 0 {
		helpers.Warning.Println("got no user Id and no transaction id in the header")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode: 400, Message:"you have to be logged in to use this service", TransactionId:transactionId}

		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)
	company_id := params["company_id"]

	if len(company_id) <= 0{
		helpers.Info.Println(transactionId + ": no company id in request")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode:400, Message:"you have to submit a valid company id as url param", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	guardianClient := client2.GuardianClient{Host:os.Getenv("GUARDIAN_URL"), CompanyId:company_id, UserId:userId}

	guardianResponse, err := client2.CheckCompanyAndPermissionFromGuardian(guardianClient, RequiredPermission)
	if err != nil {
		helpers.Info.Println(transactionId + ": guardian host responded with errror, abort transaction")
		helpers.Info.Println(err)
		w.WriteHeader(500)
		response := api.Response{Status:"ERROR", StatusCode:500, Message:"internal server error", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	if !guardianResponse {
		helpers.Info.Println(transactionId + ": user is not allowed to access company / settings")
		w.WriteHeader(401)
		response := api.Response{Status:"ERROR", StatusCode:401, Message:"not allowed to access", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}
	helpers.Info.Println(transactionId + " user is allowed to access, continue with request")

	var shippingRules []models.ShippingRule

	DbConn.Where("company_id = ?", company_id).Find(&shippingRules)

	response := api.Response{Status:"OK", StatusCode:200, Message:"successfully fetched shipping rules", TransactionId:transactionId}
	shipping_response := api.ShippingRuleResponse{Response:&response, Data:shippingRules}

	json.NewEncoder(w).Encode(&shipping_response)
	return

}

func CreateShippingRule(w http.ResponseWriter, r *http.Request) {
	helpers.Info.Println("new request to create shipping rule")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-uuid")

	if len(transactionId) <= 0 || len(userId) <= 0 {
		helpers.Warning.Println("got no user Id and no transaction id in the header")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode: 400, Message:"you have to be logged in to use this service", TransactionId:transactionId}

		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)
	company_id := params["company_id"]

	if len(company_id) <= 0{
		helpers.Info.Println(transactionId + ": no company id in request")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode:400, Message:"you have to submit a valid company id as url param", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	guardianClient := client2.GuardianClient{Host:os.Getenv("GUARDIAN_URL"), CompanyId:company_id, UserId:userId}

	guardianResponse, err := client2.CheckCompanyAndPermissionFromGuardian(guardianClient, RequiredPermission)
	if err != nil {
		helpers.Info.Println(transactionId + ": guardian host responded with errror, abort transaction")
		helpers.Info.Println(err)
		w.WriteHeader(500)
		response := api.Response{Status:"ERROR", StatusCode:500, Message:"internal server error", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	if !guardianResponse {
		helpers.Info.Println(transactionId + ": user is not allowed to access company / settings")
		w.WriteHeader(401)
		response := api.Response{Status:"ERROR", StatusCode:401, Message:"not allowed to access", TransactionId:transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}
	helpers.Info.Println(transactionId + " user is allowed to access, continue with request")

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

	shippingRule.CompanyId = company_id

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

	helpers.Info.Println(transactionId + ": successfully created shipping rule")

	w.WriteHeader(200)
	response := api.Response{Status:"OK", StatusCode: 200,
		Message:"successfully created shipping rule", TransactionId:transactionId}
	json.NewEncoder(w).Encode(response)
	return
}
