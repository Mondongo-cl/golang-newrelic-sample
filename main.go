package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	nrecho "github.com/newrelic/go-agent/v3/integrations/nrecho-v4"

	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	DatabaseDialect = "mysql"
	DatabaseDriver  = "mysql"
)

var (
	cfg                    Configuration
	app                    *newrelic.Application
	errUserNameNotFound    error = errors.New("username is required")
	errInvalidPortNumber   error = errors.New("port number must be greather than 0")
	errInvalidHostName     error = errors.New("hostname can not be empty")
	errInvalidDatabaseName error = errors.New("database name can not be null")
)

func main() {

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("MySQL HeathCheck"),
		newrelic.ConfigLicense(os.Getenv("NEWRELIC_LICENCE")),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		log.Fatalf("INIT:: New Relics initialization fails, message %s\n", err.Error())
	}

	app.WaitForConnection(5 * time.Second)

	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		port = 3306
	}
	cfg = Configuration{
		username: os.Getenv("MYSQL_USERNAME"),
		password: os.Getenv("MYSQL_PASSWORD"),
		port:     int(port),
		host:     os.Getenv("MYSQL_HOST"),
		database: os.Getenv("MYSQL_DATABASE"),
	}

	e := echo.New()
	e.Use(nrecho.Middleware(app))
	e.GET("/customers", Handle)
	err = e.Start(":3000")
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
	app.Shutdown(5 * time.Second)
}
