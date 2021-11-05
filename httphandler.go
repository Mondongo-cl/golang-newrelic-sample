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
	stmt := db.Select("*").From("customers")
	result := make([]Customer, 0)
	err = stmt.ScanStructs(&result)
	if err != nil {
		return newInternalServerError(c, err)
	}
	nrtx.End()
	ctx.Done()
	return c.JSON(200, result)
}
