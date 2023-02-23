package server

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
	"github.com/synapsecns/sanguine/services/explorer/api/cache"
	"github.com/synapsecns/sanguine/services/explorer/consumer/fetcher"
	"github.com/synapsecns/sanguine/services/explorer/db"
	"github.com/synapsecns/sanguine/services/explorer/graphql/server/graph"
	resolvers "github.com/synapsecns/sanguine/services/explorer/graphql/server/graph/resolver"
	"time"
)

const (
	// GraphqlEndpoint is the base endpoint for graphql and the endpoint for post requests to the graphql service.
	GraphqlEndpoint string = "/graphql"
	// GraphiqlEndpoint is the endpoint for the graphql user interface.
	GraphiqlEndpoint string = "/graphiql"
)

// EnableGraphql enables the scribe graphql service.
func EnableGraphql(engine *gin.Engine, consumerDB db.ConsumerDB, fetcher fetcher.ScribeFetcher, apiCache cache.Service) {
	server := createServer(
		resolvers.NewExecutableSchema(
			resolvers.Config{Resolvers: &graph.Resolver{
				DB:      consumerDB,
				Fetcher: fetcher,
				Cache:   apiCache,
			}},
		),
	)
	engine.GET(GraphqlEndpoint, graphqlHandler(server))
	engine.POST(GraphqlEndpoint, graphqlHandler(server))
	engine.GET(GraphiqlEndpoint, graphiqlHandler())
}

// Create a server without introspection
func createServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	return srv
}
