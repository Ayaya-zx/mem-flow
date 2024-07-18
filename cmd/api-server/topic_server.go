package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Ayaya-zx/mem-flow/internal/common"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
	repo "github.com/Ayaya-zx/mem-flow/internal/repository"
)

type clientError string

func (e clientError) Error() string {
	return string(e)
}

type topicServer struct {
	topicRepo repo.TopicRepository
}

func newTopicServer(topicRepo repo.TopicRepository) *topicServer {
	return &topicServer{topicRepo: topicRepo}
}

func (s *topicServer) handleError(w http.ResponseWriter, _ *http.Request, err error) {
	fmt.Println(err)
	if _, notExist := err.(common.TopicNotExistsError); notExist {
		w.WriteHeader(404)
		return
	}
	if _, badTitle := err.(common.TopicTitleError); badTitle {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(500)
}

func (s *topicServer) exampleHandler(w http.ResponseWriter, r *http.Request) {
	var topic entity.Topic

	data, err := json.Marshal(&topic)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	w.Write(data)
}

func (s *topicServer) getAllTopicsHandler(w http.ResponseWriter, r *http.Request) {
	topics, err := s.topicRepo.GetAllTopics()
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	data, err := json.Marshal(topics)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
	w.Write(data)
}

func (s *topicServer) createTopicHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	id, err := s.topicRepo.AddTopic(string(data))
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	result, err := json.Marshal(
		&struct {
			Id int `json:"id"`
		}{Id: id})
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	w.Write(result)
}

func (s *topicServer) getTopicHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	topic, err := s.topicRepo.GetTopic(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	data, err := json.Marshal(&topic)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	w.Write(data)
}

func (s *topicServer) repeateTopicHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	topic, err := s.topicRepo.GetTopic(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	topic.Repeat()
}

func (s *topicServer) deleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	err = s.topicRepo.RemoveTopic(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}
