package main

import (
	"exp/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/transaction", api.InsertTransaction)
	e.POST("/api/v1/task", api.InsertTask)
	e.POST("/api/v1/check/barcode", api.CheckingBarcode)
	e.POST("/api/v1/check/otp", api.CheckingOTP)
	e.POST("/api/v1/update/transaction", api.UpdateStatusTransaction)
	e.POST("/api/v1/atm", api.InsertAtm)
	e.POST("/api/v1/slot", api.InsertSlot)

	e.GET("/api/v1/atm/city", api.GetCity)
	e.GET("/api/v1/atm", api.GetAtm)

	e.GET("/api/v1/atm/:city/area", api.GetVendingByCity)
	e.GET("/api/v1/atm/:atm_id", api.GetAtmStock)

	e.Logger.Fatal(e.Start(":8080"))
}
