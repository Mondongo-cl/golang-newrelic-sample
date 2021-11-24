package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/labstack/echo/v4"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Handle(c echo.Context) error {
	ctx := context.Background()
	cnn, err := prepareConnection(c, ctx)
	if err != nil {
		return err
	}
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("MySQL HeathCheck"),
		newrelic.ConfigLicense(os.Getenv("NEWRELIC_LICENCE")),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		err = nrpkgerrors.Wrap(err)
		log.Fatalf("INIT:: New Relics initialization fails, message %s\n", err.Error())
	}
	app.WaitForConnection(5 * time.Second)
	txn := app.StartTransaction("get-goqu-customers")
	db := goqu.New(DatabaseDialect, cnn)

	ctx = newrelic.NewContext(ctx, txn)

	if err != nil {
		return newInternalServerError(c, err)
	}

	stmt := db.Select("*").From("customers").Limit(1000).Offset(0)
	s, p, _ := stmt.ToSQL()
	log.Print(s)
	log.Print(p)
	result := make([]Customer, 0)

	err = stmt.ScanStructsContext(ctx, &result)
	if err != nil {
		return newInternalServerError(c, err)
	}

	txn.End()

	app.Shutdown(5 * time.Second)
	log.Print("End Query Execution")
	return c.JSON(200, result)
}
