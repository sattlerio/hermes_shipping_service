package client

import (
	"hermesShippingRuleService/helpers"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

var err error

type GuardianClient struct {
	Host 		string
	CompanyId 	string
	UserId 		string
}

type GuardianResponse struct {
	Status		string 				 `json:"status"`
	Data		GuardianResponseData `json:"data"`

}

type GuardianResponseData struct {
	UserPermission int `json:"user_permission"`
}

func CheckCompanyAndPermissionFromGuardian(client GuardianClient, permission int) (bool, error) {
	helpers.Info.Println("starting to communicate with guardian")

	url := fmt.Sprintf("%s/%s/%s", client.Host, client.UserId, client.CompanyId)
	response, err := http.Get(url)
	if err != nil {
		helpers.Info.Println("not possible to communicate with guardian")
		return false, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		helpers.Info.Println("error with guardian response")
		return false, err
	}
	if response.StatusCode != 200 {
		helpers.Info.Println("abort transaction because guardian answered with status " + strconv.Itoa(response.StatusCode))
		return false, nil
	}

	guardianData := GuardianResponse{}
	jsonErr := json.Unmarshal(body, &guardianData)

	if jsonErr != nil {
		helpers.Info.Println("not possible to parse JSON response from guardian")
		return false, jsonErr
	}

	if guardianData.Status != "OK" {
		helpers.Info.Println("error with guardian communication, guardian has answered with status: " + guardianData.Status)
		return false, nil
	}

	helpers.Info.Println("guardian answered with status " + guardianData.Status + " and with permission " + strconv.Itoa(guardianData.Data.UserPermission))

	if guardianData.Data.UserPermission < 0 || guardianData.Data.UserPermission > permission {
		helpers.Info.Println("user has permission " + strconv.Itoa(guardianData.Data.UserPermission) + " but required permission is " + strconv.Itoa(permission))
		return false, nil
	}

	helpers.Info.Println("user has the required permission user level {" + strconv.Itoa(guardianData.Data.UserPermission) + "} required permission is {" + strconv.Itoa(permission) + "}")
	return true, err
}