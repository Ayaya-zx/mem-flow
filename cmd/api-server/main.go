package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	server := newTopicServer(inmem.NewInmemTopicRepository())
	mux := http.NewServeMux()

	mux.HandleFunc("GET /topics", server.getAllTopicsHandler)
	mux.HandleFunc("POST /topics", server.createTopicHandler)
	mux.HandleFunc("GET /topics/{id}", server.getTopicHandler)
	mux.HandleFunc("PATCH /topics/{id}", server.repeateTopicHandler)
	mux.HandleFunc("DELETE /topics/{id}", server.deleteTopicHandler)
	mux.HandleFunc("GET /example", server.exampleHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", v.GetInt("Port")), mux)
	if err != nil {
		log.Fatal(err)
	}
}
