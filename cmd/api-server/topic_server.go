package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ayaya-zx/mem-flow/internal/api"
	"github.com/Ayaya-zx/mem-flow/internal/auth"
	"github.com/Ayaya-zx/mem-flow/internal/common"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
	repo "github.com/Ayaya-zx/mem-flow/internal/repository"
)

type clientError string

func (e clientError) Error() string {
	return string(e)
}

type topicServer struct {
	authService   *auth.AuthService
	userTopicRepo repo.UserTopicRepository
}

func newTopicServer(authService *auth.AuthService, userTopicRepo repo.UserTopicRepository) *topicServer {
	return &topicServer{
		authService:   authService,
		userTopicRepo: userTopicRepo,
	}
}

func getToken(bearer string) (string, error) {
	if bearer == "" {
		return "", fmt.Errorf("bearer token not found")
	}
	split := strings.Split(bearer, " ")
	if len(split) != 2 && split[0] != "Bearer" {
		return "", fmt.Errorf("bearer token not found")
	}
	return split[1], nil
}

func (s *topicServer) handleError(w http.ResponseWriter, _ *http.Request, err error) {
	fmt.Println(err)
	if _, notExist := err.(common.TopicNotExistsError); notExist {
		w.WriteHeader(404)
		return
	}

	_, badTitle := err.(common.TopicTitleError)
	_, clientErr := err.(clientError)
	_, invalidAuth := err.(common.InvalidAuthData)
	if badTitle || clientErr || invalidAuth {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(500)
}

func (s *topicServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r.Header.Get("Authorization"))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}
		name, err := s.authService.Validate(token)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}
		ctx := context.WithValue(r.Context(), "username", name)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (s *topicServer) registrationHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	if len(data) == 0 {
		s.handleError(w, r, clientError("empty body"))
		return
	}

	authData := new(auth.AuthData)
	err = json.Unmarshal(data, authData)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	token, err := s.authService.RegUser(authData)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	io.WriteString(w, token)
}

func (s *topicServer) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	authData := new(auth.AuthData)
	err = json.Unmarshal(data, authData)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	token, err := s.authService.AuthUser(authData)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	io.WriteString(w, token)
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
	name := r.Context().Value("username").(string)
	topicRepo, err := s.userTopicRepo.GetUserTopicRepository(name)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	topics, err := topicRepo.GetAllTopics()
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
	name := r.Context().Value("username").(string)
	topicRepo, err := s.userTopicRepo.GetUserTopicRepository(name)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	var req api.CreateTopicRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
	}

	id, err := topicRepo.AddTopic(req.Title)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	resp := api.CreateTopicResponse{Id: id}
	data, _ = json.Marshal(&resp)
	w.Write(data)
}

func (s *topicServer) getTopicHandler(w http.ResponseWriter, r *http.Request) {
	name := r.Context().Value("username").(string)
	topicRepo, err := s.userTopicRepo.GetUserTopicRepository(name)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	raw := r.PathValue("id")
	id, err := strconv.Atoi(raw)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	topic, err := topicRepo.GetTopicById(id)
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
	name := r.Context().Value("username").(string)
	topicRepo, err := s.userTopicRepo.GetUserTopicRepository(name)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	raw := r.PathValue("id")
	id, err := strconv.Atoi(raw)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	topic, err := topicRepo.GetTopicById(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	topic.Repeat()
}

func (s *topicServer) deleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	name := r.Context().Value("username").(string)
	topicRepo, err := s.userTopicRepo.GetUserTopicRepository(name)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	raw := r.PathValue("id")
	id, err := strconv.Atoi(raw)
	if err != nil {
		s.handleError(w, r, clientError(err.Error()))
		return
	}

	err = topicRepo.RemoveTopic(id)
	if err != nil {
		s.handleError(w, r, err)
		return
	}
}
