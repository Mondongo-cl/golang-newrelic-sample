package main

import (
	"github.com/labstack/echo/v4"
)

func Handle(c echo.Context) error {
	app.StartTransaction("Get Customers")
	db, err := prepareConnection(c)
	if err != nil {
		return err
	}
	stmt := db.Select("*").From("customers").Limit(10).Offset(0)
	result := make([]Customer, 0)
	err = stmt.ScanStructs(&result)
	if err != nil {
		return newInternalServerError(c, err)
	}
	return c.JSON(200, result)
}
