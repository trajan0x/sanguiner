package metrics

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ngrin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/synapsecns/sanguine/core/config"
	"net/http"
	"os"
	"sync"
)

type newRelicHandler struct {
	app       *newrelic.Application
	startMux  sync.Mutex
	buildInfo config.BuildInfo
}

func NewRelicMetricsHandler(buildInfo config.BuildInfo) Handler {
	return &newRelicHandler{
		buildInfo: buildInfo,
	}
}

func (n *newRelicHandler) Gin() gin.HandlerFunc {
	return ngrin.Middleware(n.app)
}

func (n *newRelicHandler) Start(_ context.Context) (err error) {
	n.startMux.Lock()
	defer n.startMux.Unlock()
	if n.app == nil {
		n.app, err = newrelic.NewApplication(
			newrelic.ConfigAppName(n.buildInfo.Name()),
			newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigCodeLevelMetricsEnabled(true),
			func(c *newrelic.Config) {
				c.Labels = map[string]string{
					"version": n.buildInfo.Version(),
					"commit":  n.buildInfo.Commit(),
				}
			},
			// optional overrides
			newrelic.ConfigFromEnvironment(),
		)
		if err != nil {
			return fmt.Errorf("could not create new relic application: %w", err)
		}
	}

	return nil
}

func (n *newRelicHandler) ConfigureHttpClient(client *http.Client) {
	// use the newrelic transport
	nrTransport := newrelic.NewRoundTripper(client.Transport)
	client.Transport = nrRoundTripper{app: n.app, inner: nrTransport}
}

type nrRoundTripper struct {
	inner http.RoundTripper
	app   *newrelic.Application
}

func (n nrRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	txn := newrelic.FromContext(req.Context())
	if txn == nil {
		txn = n.app.StartTransaction(req.URL.String())
		req = newrelic.RequestWithTransactionContext(req, txn)
	}
	// nolint: errrwrap
	return n.inner.RoundTrip(req)
}
