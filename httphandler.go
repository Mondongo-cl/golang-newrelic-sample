package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Handle(c echo.Context) error {

	nrtx := app.StartTransaction("Get Customers")
	ctx := newrelic.NewContext(context.Background(), nrtx)

	tx, db, err := prepareConnection(c, ctx)
	if err != nil {
		return err
	}
	stmt := tx.Select("*").From("customers").Limit(1000).Offset(0)
	rawQuery, params, _ := stmt.ToSQL()
	p, _ := json.Marshal(params)
	log.Printf("starting query:\n%v params %v", rawQuery, string(p))
	result := make([]Customer, 0)
	err = stmt.ScanStructsContext(ctx, &result)
	if err != nil {
		return newInternalServerError(c, err)
	}
	log.Print("End Query Execution")
	db.Close()
	return c.JSON(200, result)
}
