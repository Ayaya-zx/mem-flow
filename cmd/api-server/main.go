package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ayaya-zx/mem-flow/internal/auth"
	"github.com/Ayaya-zx/mem-flow/internal/repository/inmem"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var port int

	v := viper.New()
	v.SetDefault("Port", 8765)

	v.AutomaticEnv()
	v.SetEnvPrefix("MEMFLOW")
	v.BindEnv("Port", "port")

	pflag.IntVarP(&port, "port", "p", 8765, "Port")
	pflag.Parse()
	if err := v.BindPFlag("Port", pflag.Lookup("port")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server := newTopicServer(
		auth.NewAuthService(inmem.NewInmemUserRepository()),
		inmem.NewInmemUserTopicRepository(inmem.NewInmemTopicRepositoryFactory()),
	)

	mux := http.NewServeMux()
	mux.Handle("GET /topics", server.authMiddleware(http.HandlerFunc(server.getAllTopicsHandler)))
	mux.Handle("POST /topics", server.authMiddleware(http.HandlerFunc(server.createTopicHandler)))
	mux.Handle("GET /topics/{title}", server.authMiddleware(http.HandlerFunc(server.getTopicHandler)))
	mux.Handle("PATCH /topics/{title}", server.authMiddleware(http.HandlerFunc(server.repeateTopicHandler)))
	mux.Handle("DELETE /topics/{title}", server.authMiddleware(http.HandlerFunc(server.deleteTopicHandler)))
	mux.HandleFunc("GET /example", http.HandlerFunc(server.exampleHandler))
	mux.HandleFunc("POST /registration", server.registrationHandler)
	mux.HandleFunc("POST /auth", server.authenticationHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", v.GetInt("Port")), mux)
	if err != nil {
		log.Fatal(err)
	}
}
