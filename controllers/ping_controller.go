package controllers

import (
	"net/http"
	"encoding/json"
	"os"
	"hermesShippingRuleService/api"
	"hermesShippingRuleService/helpers"
)

func PingController(w http.ResponseWriter, r *http.Request) {
	helpers.Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)

	helpers.Info.Println("new request for ping pong route")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-id")

	if len(transactionId) <= 0 || len(userId) <= 0 {
		helpers.Warning.Println("got no user Id and no transaction id in the header")
		w.WriteHeader(400)
		response := api.Response{Status:"ERROR", StatusCode: 400, Message:"you have to be logged in to use this service", TransactionId:transactionId}

		json.NewEncoder(w).Encode(response)
		return
	}

	helpers.Info.Println(transactionId + ": got new request for ping route")
	json.NewEncoder(w).Encode("pong")
	return
}
