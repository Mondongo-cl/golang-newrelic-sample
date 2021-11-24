package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
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

func prepareConnection(c echo.Context, ctx context.Context) (*sql.DB, error) {
	connectionstring, err := cfg.Parse()

	if err != nil {
		return nil, newInternalServerError(c, err)
	}

	cnn, err := sql.Open(DatabaseDriver, connectionstring)
	if err != nil {
		return nil, newInternalServerError(c, err)
	}
	cnn.SetConnMaxLifetime(time.Minute * 3)
	cnn.SetMaxOpenConns(10)
	cnn.SetMaxIdleConns(10)
	err = cnn.Ping()
	if err != nil {
		return nil, newInternalServerError(c, err)
	}

	return cnn, nil
}

func newInternalServerError(c echo.Context, err error) error {
	err = nrpkgerrors.Wrap(err)
	c.JSON(500, err.Error())
	return err
}
