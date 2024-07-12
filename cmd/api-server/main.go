package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Ayaya-zx/mem-flow/internal/entity"
	"github.com/Ayaya-zx/mem-flow/internal/store"
	"github.com/Ayaya-zx/mem-flow/internal/store/inmem"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type topicServer struct {
	topicStore store.TopicStore
}

func newTopicServer(topicStore store.TopicStore) *topicServer {
	return &topicServer{topicStore: topicStore}
}

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

	server := newTopicServer(inmem.NewInmemTopicStore())
	mux := http.NewServeMux()

	mux.HandleFunc("GET /topics", server.getAllTopicsHandler)
	mux.HandleFunc("POST /topics", server.createTopicHandler)
	mux.HandleFunc("GET /topics/{id}", server.getTopicHandler)
	mux.HandleFunc("PATCH /topics/{id}", server.repeateTopicHandler)
	mux.HandleFunc("DELETE /topics/{id}", server.deleteTopicHandler)
	mux.HandleFunc("GET /example", exampleHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", v.GetInt("Port")), mux)
	if err != nil {
		log.Fatal(err)
	}
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	var topic entity.Topic

	data, err := json.Marshal(&topic)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
}

func (s *topicServer) getAllTopicsHandler(w http.ResponseWriter, r *http.Request) {
	topics, err := s.topicStore.GetAllTopics()
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(500)
		return
	}
	data, err := json.Marshal(topics)
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}

func (s *topicServer) createTopicHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(400)
		return
	}

	id, err := s.topicStore.AddTopic(string(data))
	if _, ok := err.(store.TopicTitleError); ok {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(500)
		return
	}
	io.WriteString(w, strconv.Itoa(id))
}

func (s *topicServer) getTopicHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	topic, err := s.topicStore.GetTopic(id)
	if _, ok := err.(store.TopicNotExistsError); ok {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(&topic)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
}

func (s *topicServer) repeateTopicHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	err = s.topicStore.TopicRepeated(id)
	if _, ok := err.(store.TopicNotExistsError); ok {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

func (s *topicServer) deleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	err = s.topicStore.RemoveTopic(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
