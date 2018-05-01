package main

import (
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"os"
	"log"
	"net/http"
	_ "github.com/bmizerany/pq"
	"hermesShippingRuleService/models"
	"hermesShippingRuleService/helpers"
	"hermesShippingRuleService/controllers"
)

var db *gorm.DB
var err error

func main() {
	helpers.Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)
	helpers.Info.Println("starting up server...")

	router := mux.NewRouter()

	db, err = gorm.Open(
		"postgres",
		"host="+os.Getenv("PSQL_HOST")+" user="+os.Getenv("PSQL_USER")+
			" dbname="+os.Getenv("PSQL_DBNAME")+" sslmode=disable password="+
			os.Getenv("PSQL_PASSWORD"))

	if err != nil {
		helpers.Error.Println(err)
		helpers.Error.Println("not possible to connect to db, going to die now.... UAAAAAAAAAAAAH!!!!")
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&models.ShippingRule{}, &models.ShippingRules2Countries{})

	controllers.DbConn = db

	router.HandleFunc("/ping", controllers.PingController).Methods("GET")
	router.HandleFunc("/all/{company_id}", controllers.GetShippingRules).Methods("GET")
	router.HandleFunc("/create/{company_id}", controllers.CreateShippingRule).Methods("POST")
	router.HandleFunc("/fetch/{company_id}/{shipping_rule_id}", controllers.GetShippingRule).Methods("GET")
	router.HandleFunc("/delete/{shipping_rule_id}/{company_id}", controllers.DeleteShippingRule).Methods("DELETE")
	router.HandleFunc("/edit/{company_id}/{shipping_rule_id}", controllers.EditShippingRule).Methods("PUT")
	helpers.Info.Println("successfully started server on Port 10000")

	log.Fatal(http.ListenAndServe(":10000", router))
}