package main

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"

	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (c Configuration) Parse() (string, error) {
	if c.username == "" {
		return "", errUserNameNotFound
	}
	if c.port <= 0 {
		return "", errInvalidPortNumber
	}
	if c.host == "" {

		return "", errInvalidHostName
	}
	if c.database == "" {

		return "", errInvalidDatabaseName
	}
	cnnstr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", c.username, c.password, c.host, c.port, c.database)
	return cnnstr, nil
}

func prepareConnection(c echo.Context) (*goqu.TxDatabase, error) {
	connectionstring, err := cfg.Parse()

	if err != nil {
		return nil, newInternalServerError(c, err)
	}

	cnn, err := sqlx.Open(DatabaseDriver, connectionstring)
	if err != nil {
		return nil, newInternalServerError(c, err)
	}
	err = cnn.Ping()
	if err != nil {
		return nil, newInternalServerError(c, err)
	}
	db := goqu.New(DatabaseDialect, cnn)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, newInternalServerError(c, err)
	}
	return tx, nil
}

func newInternalServerError(c echo.Context, err error) error {
	return c.JSON(500, err)
}
