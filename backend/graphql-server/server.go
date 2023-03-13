package main

import (
    "net/http"
    "fmt"

    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/handler/transport"
    "github.com/99designs/gqlgen/graphql/playground"
    "github.com/andru100/Social-Network-Microservice/backend/graphql-server/graph"
    "github.com/go-chi/chi"
    "github.com/gorilla/websocket"
    "github.com/rs/cors"
)


func main() {

    fmt.Println("GQL Server is running!")

    router := chi.NewRouter()

    // Add CORS middleware around every request
    router.Use(cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowCredentials: true,
        Debug:            true,
    }).Handler)

   
    srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

    srv.AddTransport(&transport.Websocket{
        Upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                // Check against your desired domains here
                return r.Host == "example.org"
            },
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        },
    })

    router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
    router.Handle("/query", srv)

    err := http.ListenAndServe(":8080", router)
    if err != nil {
        panic(err)
    }
}