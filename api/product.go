package api

import (
	"encoding/json"
	"exp/core/db"
	"exp/core/response"
	"exp/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func InsertSlot(c echo.Context) error {

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	atm_id := json_map["atm_id"]
	product_id := json_map["product_id"]
	slot_number := json_map["slot_number"]
	stock_quantity := json_map["stock_quantity"]

	validation := false
	if util.IsNil(atm_id) {
		validation = true
	}
	if util.IsNil(stock_quantity) {
		validation = true
	}
	if util.IsNil(product_id) {
		validation = true
	}
	if util.IsNil(slot_number) {
		validation = true
	}

	if validation {
		return c.JSON(http.StatusBadRequest, response.BadRequest())
	}

	in := bson.M{
		"atm_id":         atm_id,
		"product_id":     product_id,
		"slot_number":    slot_number,
		"stock_quantity": stock_quantity,
	}

	result := db.InserRow(in, "slot")

	if result < 0 {
		return c.JSON(http.StatusInternalServerError, response.InternalError())
	}

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, response.Success())
}

func GetAtmStock(c echo.Context) error {
	result := db.GetRows(bson.M{"atm_id": c.Param("atm_id")}, "slot", bson.M{"slot_number": 1, "product_id": 1, "stock_quantity": 1, "_id": 0})

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, response.DataArray(http.StatusOK, "OK", result))
}
