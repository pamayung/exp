package api

import (
	"bytes"
	"encoding/json"
	"exp/core/cache"
	"exp/core/code"
	"exp/core/db"
	"exp/core/response"
	"exp/model"
	"exp/util"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strconv"
)

func InsertTransaction(c echo.Context) error {

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	atm_id := json_map["atm_id"]
	transaction_id := json_map["transaction_id"]
	barcode := json_map["barcode"]
	data := json_map["data"].([]interface{})

	validation := false
	if util.IsNil(atm_id) {
		validation = true
	}
	if util.IsNil(transaction_id) {
		validation = true
	}
	if util.IsNil(barcode) {
		validation = true
	}
	if len(data) <= 0 {
		validation = true
	}

	create_date := util.CurrentTimeMillis()
	expired := util.CurrentTimeMillis() + (3 * 3600000)

	for i := 0; i < len(data); i++ {
		value := data[i].(map[string]interface{})

		slot := value["slot_number"]
		weight := value["gold_weight"]
		qty := value["quantity"]

		if util.IsNil(slot) {
			validation = true
		}
		if util.IsNil(weight) {
			validation = true
		}
		if util.IsNil(qty) {
			validation = true
		}
	}

	if validation {
		return c.JSON(http.StatusBadRequest, response.BadRequest())
	}

	in := bson.M{
		"atm_id":         atm_id,
		"transaction_id": transaction_id,
		"barcode":        barcode,
		"create_date":    create_date,
		"expired":        expired,
		"data":           data,
		"dropped":        nil,
		"status":         "pending",
	}

	result := db.InserRow(in, "transaction")

	if result < 0 {
		return c.JSON(http.StatusInternalServerError, response.InternalError())
	}

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, map[string]interface{}{"status": 200, "message": "Success", "expired": expired})
}

func InsertTask(c echo.Context) error {

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	atm_id := json_map["atm_id"]
	barcode := json_map["barcode"]
	data := json_map["data"].([]interface{})

	validation := false
	if util.IsNil(atm_id) {
		validation = true
	}
	if util.IsNil(barcode) {
		validation = true
	}
	if len(data) <= 0 {
		validation = true
	}

	for i := 0; i < len(data); i++ {
		value := data[i].(map[string]interface{})

		product_name := value["product_name"]
		slot_number := value["slot_number"]
		new_stock := value["new_stock"]
		product_id := value["product_id"]

		if util.IsNil(product_name) {
			validation = true
		}
		if util.IsNil(slot_number) {
			validation = true
		}
		if util.IsNil(new_stock) {
			validation = true
		}
		if util.IsNil(product_id) {
			validation = true
		}
	}

	if validation {
		return c.JSON(http.StatusBadRequest, response.BadRequest())
	}

	task_id := code.TASK_ID + strconv.FormatInt(util.CurrentTimeMillis(), 10)
	create_date := util.CurrentTimeMillis()
	expired := util.CurrentTimeMillis() + (24 * 3600000)

	in := bson.M{
		"atm_id":      atm_id,
		"task_id":     task_id,
		"barcode":     barcode,
		"create_date": create_date,
		"expired":     expired,
		"data":        data,
		"status":      "pending",
	}

	result := db.InserRow(in, "transaction")

	if result < 0 {
		return c.JSON(http.StatusInternalServerError, response.InternalError())
	}

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, map[string]interface{}{"status": 200, "message": "Success", "task_id": task_id, "expired": expired})
}

func UpdateStatusTransaction(c echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	barcode := json_map["barcode"].(string)
	status := json_map["status"].(string)
	dropped := json_map["dropped"].(string)

	validation := false
	if util.IsNil(barcode) {
		validation = true
	}
	if util.IsNil(status) {
		validation = true
	}
	if util.IsNil(dropped) {
		validation = true
	}

	if validation {
		return c.JSON(http.StatusBadRequest, response.BadRequest())
	}

	id := bson.M{"barcode": barcode}
	set := bson.D{{
		"$set", bson.D{{
			"status", status}}},
		{
			"$set", bson.D{{
				"dropped", dropped}}}}

	update := db.UpdateRow(id, set, "transaction")

	if update > 0 {
		return c.JSON(http.StatusOK, response.Success())
	} else {
		return c.JSON(http.StatusBadRequest, response.Failed())
	}

}

func CheckingBarcode(c echo.Context) error {

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	vending_code := c.Request().Header.Get("vending_code")

	if err != nil {
		return err
	}

	barcode := json_map["barcode"]

	validation := false
	if util.IsNil(barcode) {
		validation = true
	}

	if validation {
		return c.JSON(http.StatusBadRequest, response.BadRequest())
	}

	find := bson.M{"atm_id": "$atm_id",
		"barcode":      "$barcode",
		"expired":      "$expired",
		"data":         "$data",
		"vending_code": "$vending_code"}
	show := bson.M{"atm_id": "$atm_id", "barcode": "$barcode", "expired": "$expired", "data": "$data", "vending_code": "$$vending_code"}
	where := bson.A{
		bson.M{"$eq": bson.A{"$atm_id", "$$atm_id"}},
		bson.M{"$eq": bson.A{"$barcode", barcode}},
		bson.M{"$eq": bson.A{"$status", "pending"}},
		bson.M{"$eq": bson.A{vending_code, "$$vending_code"}}}
	//bson.M{"$gte": bson.A{"$expired", 1605682080162}} }

	result := db.GetRowAgtegate(find, "vending", "transaction", where, show)

	jArray := result["data"].(bson.A)
	if len(jArray) == 0 {
		return c.JSON(http.StatusNotFound, response.NotFound())
	}

	jObject := jArray[0].(bson.M)

	if jObject["expired"].(int64) < util.CurrentTimeMillis() {
		id := bson.M{"barcode": barcode}
		set := bson.D{{
			"$set", bson.D{{
				"status", "failed"}}}}

		update := db.UpdateRow(id, set, "transaction")

		if update > 0 {
			return c.JSON(http.StatusForbidden, response.Custom(http.StatusForbidden, "Barcode Expired"))
		} else {
			return c.JSON(http.StatusBadRequest, response.Failed())
		}

	}

	id := 1
	//resp := make(map[string]interface{})

	var data bson.A
	if jObject["task_id"] == nil {
		data = jObject["data"].(bson.A)

	} else {
		data = jObject["task_list"].(bson.A)
	}

	t := getTokenWA(id)

	if t > 0 {
		return c.JSON(http.StatusOK, data)
	} else {
		return c.JSON(http.StatusInternalServerError, response.InternalError())
	}

}

func getTokenWA(id int) int {
	jsonReq := make(map[string]interface{})
	jsonReq["id"] = id

	requestBody, err := json.Marshal(jsonReq)

	if err != nil {
		log.Panic(err)
	}

	resp, err := http.Post("", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Panic(err)
		return -1
	}

	defer resp.Body.Close()

	//body, err := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))

	json_map := make(map[string]interface{})
	err2 := json.NewDecoder(resp.Body).Decode(&json_map)

	fmt.Println(json_map)

	if err2 != nil {
		log.Panic(err)
		return -1
	}

	if resp.StatusCode == http.StatusOK {
		token := model.OTP{Code: "qqqqqq", Expired: 16057999999999}
		cache.SetOtp("qwertyuiop", &token)
		token2 := model.OTP{Code: "wwwwww", Expired: 1605785604232}
		cache.SetOtp("111111", &token2)

		return 1

	} else {
		return -1
	}
}

func CheckingOTP(c echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	barcode := json_map["barcode"].(string)

	validation := false
	if util.IsNil(barcode) {
		validation = true
	}

	if validation {
		return c.JSON(http.StatusBadRequest, response.BadRequest())
	}

	otp := cache.GetOtp(barcode)

	if otp == nil {
		return c.JSON(http.StatusNotFound, response.NotFound())
	}

	data := otp.(*model.OTP)

	if int64(data.Expired) < util.CurrentTimeMillis() {
		return c.JSON(http.StatusForbidden, response.Custom(http.StatusForbidden, "OTP Expired"))
	} else {
		id := bson.M{"barcode": barcode}
		set := bson.D{{
			"$set", bson.D{{
				"status", "capture"}}}}

		update := db.UpdateRow(id, set, "transaction")

		if update > 0 {
			return c.JSON(http.StatusOK, response.Success())
		} else {
			return c.JSON(http.StatusBadRequest, response.Failed())
		}

	}

}
