package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Handle(c echo.Context) error {
	nrtx := app.StartTransaction("Get Customers")
	ctx := newrelic.NewContext(context.Background(), nrtx)
	db, err := prepareConnection(c, ctx)
	if err != nil {
		return err
	}
	stmt := db.Select("*").From("customers").Limit(10).Offset(0)
	result := make([]Customer, 10)
	err = stmt.ScanStructs(&result)
	if err != nil {
		return newInternalServerError(c, err)
	}
	return c.JSON(200, result)
}
