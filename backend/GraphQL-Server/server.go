// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/playground"
// 	"github.com/andru100/Graphql-Social-Network/graph"
//     "github.com/andru100/Graphql-Social-Network/graph/social"
// )

// const defaultPort = "8080"

// func main() {

//     go social.UploadEndpoint() // start upload listener

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultPort
// 	}

// 	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", srv)

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }

package main

import (
    "net/http"
    "fmt"
    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/handler/transport"
    "github.com/99designs/gqlgen/graphql/playground"
    "github.com/andru100/Social-Network-Microservice/backend/GraphQL-Server/graph"
    //"github.com/andru100/Graphql-Social-Network/graph/social"
    "github.com/go-chi/chi"
    "github.com/gorilla/websocket"
    "github.com/rs/cors"
)

func main() {

    fmt.Println("GQL Server is running!")
    router := chi.NewRouter()

    // Add CORS middleware around every request
    // See https://github.com/rs/cors for full option listing
    router.Use(cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowCredentials: true,
        Debug:            true,
    }).Handler)

    // srv := handler.New(starwars.NewExecutableSchema(starwars.NewResolver())) // MODIFIED THIS.
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
    // router.Handle("/", handler.Playground("Starwars", "/query")) // MODIFIED THIS.
    router.Handle("/query", srv)

    //go social.UploadEndpoint()

    err := http.ListenAndServe(":8080", router)
    if err != nil {
        panic(err)
    }
}