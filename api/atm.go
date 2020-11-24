package api

import (
	"encoding/json"
	"exp/core/code"
	"exp/core/db"
	"exp/core/response"
	"exp/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

func InsertAtm(c echo.Context) error {

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		return err
	}

	//atm_id := json_map["atm_id"]
	vending_code := json_map["vending_code"]
	city := json_map["city"]
	area := json_map["area"]
	description := json_map["description"]

	validation := false
	//if util.IsNil(atm_id) {
	//	validation = true
	//}
	if util.IsNil(vending_code) {
		validation = true
	}
	if util.IsNil(city) {
		validation = true
	}
	if util.IsNil(area) {
		validation = true
	}
	if util.IsNil(description) {
		description = ""
	}

	if validation {
		return c.JSON(http.StatusBadRequest, response.BadRequest())
	}

	in := bson.M{
		"atm_id":       code.ATM_ID + strconv.FormatInt(util.CurrentTimeMillis(), 10),
		"vending_code": vending_code,
		"city":         city,
		"area":         area,
		"description":  description,
	}

	result := db.InserRow(in, "vending")

	if result < 0 {
		return c.JSON(http.StatusInternalServerError, response.InternalError())
	}

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, response.Success())
}

func GetAtm(c echo.Context) error {
	result := db.GetRows(bson.M{}, "vending", bson.M{"atm_id": 1, "area": 1, "city": 1, "description": 1, "_id": 0})

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, response.DataArray(http.StatusOK, "OK", result))
}

func GetCity(c echo.Context) error {
	result := db.GetRows(bson.M{}, "vending", bson.M{"city": 1, "_id": 0})

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, response.DataArray(http.StatusOK, "OK", result))
}

func GetVendingByCity(c echo.Context) error {
	result := db.GetRows(bson.M{"city": c.Param("city")}, "vending", bson.M{"atm_id": 1, "area": 1, "description": 1, "_id": 0})

	if len(result) == 0 {
		return c.JSON(http.StatusNotFound, response.NotFound())
	}

	c.Response().Header().Set("content-type", "application/json")
	return c.JSON(http.StatusOK, response.DataArray(http.StatusOK, "OK", result))
}
