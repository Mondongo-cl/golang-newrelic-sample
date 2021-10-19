package configuration

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

func init() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("AWS Lambda - Health Check"),
		newrelic.ConfigLicense("105222bb0b2377911c642609243ac8518c7dNRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	app.RecordCustomEvent("Init Config module", map[string]interface{}{"Test": "Data"})
	if nil != err {
		panic(err)
	}
}
