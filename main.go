//go:generate prisma generate
//go:generate go run gqlgen/cmd.go -c gqlgen/gqlgen.yml

package main

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/steebchen/graphql/api"
	"github.com/steebchen/graphql/gqlgen"
	"github.com/steebchen/graphql/lib/auth"
	"github.com/steebchen/graphql/lib/config"
	"github.com/steebchen/graphql/lib/handler_adapter"
	"github.com/steebchen/graphql/prisma"
	"log"
	"net/http"
)

const defaultPort = "4000"

func main() {
	port := config.GetOrDefault(config.Port, defaultPort)

	client := prisma.New(&prisma.Options{
		Endpoint: "http://localhost:4466/graphql/dev",
		Secret:   "",
	})

	resolver := api.New(client)

	http.Handle("/", handler.Playground("GraphQL Playground", "/query"))
	schema := gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: resolver})
	http.Handle("/query", &auth.Handler{
		Prisma: client,
		Next: &handler_adapter.HandlerFuncAdapter{
			NextFunc: handler.GraphQL(schema),
		},
	})

	log.Printf("Server is running on http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
