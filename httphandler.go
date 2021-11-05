package main

import (
	"context"
	"log"
	"time"

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
	stmt := db.Select("*").From("customers").Limit(1000).Offset(0)
	rawQuery, params, err := stmt.ToSQL()

	if err == nil {
		rows, err := db.QueryContext(ctx, rawQuery, params...)
		if err == nil {
			for rows.Next() {
				cols, _ := rows.Columns()
				for i, v := range cols {
					log.Printf("Key:%v ; Value:%v", i, v)
				}
				app.WaitForConnection(5 * time.Second)
			}
		}
	}
	result := make([]Customer, 0)
	err = stmt.ScanStructsContext(ctx, &result)
	if err != nil {
		return newInternalServerError(c, err)
	}
	nrtx.End()
	ctx.Done()
	return c.JSON(200, result)
}
