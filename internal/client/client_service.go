package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Ayaya-zx/mem-flow/internal/api"
	"github.com/Ayaya-zx/mem-flow/internal/auth"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
)

type ClientService struct {
	serverURL string
	token     string
}

func NewClientService(serverURL string) *ClientService {
	return &ClientService{serverURL: serverURL}
}

func (cs *ClientService) Register(authData auth.AuthData) error {
	return cs.getAuthInfo("/registration", authData)
}

func (cs *ClientService) Auth(authData auth.AuthData) error {
	return cs.getAuthInfo("/auth", authData)
}

func (cs *ClientService) GetAllTopics() ([]entity.Topic, error) {
	data, err := cs.sendGet("/topics")
	if err != nil {
		return nil, err
	}

	var result []entity.Topic
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *ClientService) GetTopicById(id int) (*entity.Topic, error) {
	data, err := cs.sendGet("/topics" + fmt.Sprintf("/%d", id))
	if err != nil {
		return nil, err
	}

	var result entity.Topic
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (cs *ClientService) AddTopic(title string) error {
	_, err := cs.sendPost("/topics",
		api.CreateTopicRequest{Title: title})
	return err
}

func (cs *ClientService) RepeatTopic(id int) error {
	_, err := cs.sendPatch("/topics" + fmt.Sprintf("/%d", id))
	return err
}

func (cs *ClientService) RemoveTopic(id int) error {
	_, err := cs.sendDelete("/topics" + fmt.Sprintf("/%d", id))
	return err
}

func (cs *ClientService) addAuthData(r *http.Request) {
	if cs.token == "" {
		panic("not authorized")
	}
	r.Header.Add(
		"Authorization",
		"Bearer "+cs.token,
	)
}

func (cs *ClientService) sendGet(path string) ([]byte, error) {
	URL := cs.serverURL + path

	req, err := http.NewRequest(
		"GET",
		URL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return cs.sendRequest(req)
}

func (cs *ClientService) sendPost(path string, data any) ([]byte, error) {
	var err error
	var body []byte

	URL := cs.serverURL + path

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(
		"POST",
		URL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	return cs.sendRequest(req)
}

func (cs *ClientService) sendPatch(path string) ([]byte, error) {
	URL := cs.serverURL + path

	req, err := http.NewRequest(
		"PATCH",
		URL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return cs.sendRequest(req)
}

func (cs *ClientService) sendDelete(path string) ([]byte, error) {
	URL := cs.serverURL + path

	req, err := http.NewRequest(
		"DELETE",
		URL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return cs.sendRequest(req)
}

func (cs *ClientService) sendRequest(req *http.Request) ([]byte, error) {
	cs.addAuthData(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, apiError(resp.StatusCode)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *ClientService) getAuthInfo(path string, authData auth.AuthData) error {
	data, err := json.Marshal(&authData)
	if err != nil {
		return err
	}

	URL := cs.serverURL + path

	resp, err := http.Post(
		URL,
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return apiError(resp.StatusCode)
	}

	tokenData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body.Close()

	cs.token = string(tokenData)
	return nil
}

func apiError(code int) error {
	return fmt.Errorf("api status code %d", code)
}
