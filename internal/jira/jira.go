package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const contentTypeJson = "application/json"

var ErrAuth = errors.New("authentication error")
var ErrCreatingIssue = errors.New("while creating issue")

func NewJira(login, password, host string) *Jira {
	return &Jira{
		login:    login,
		password: password,
		host:     host,
	}
}

type Jira struct {
	login    string
	password string
	host     string
}

func (j *Jira) Authenticate() ([]*http.Cookie, error) {
	data, err := json.Marshal(map[string]string{
		"username": j.login,
		"password": j.password,
	})
	if err != nil {
		log.Println(err)

		return nil, err
	}

	authResponse, err := http.Post(j.host+"/rest/auth/1/session", contentTypeJson, bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)

		return nil, err
	}

	if authResponse.StatusCode >= 300 {
		log.Println(err)

		return nil, ErrAuth
	}

	return authResponse.Cookies(), nil
}

func (j *Jira) CreateIssue(jiraAuthCookies []*http.Cookie, issue *Issue) (*CreatedIssue, error) {
	marshalled, err := json.Marshal(issue)
	buffered := bytes.NewBuffer(marshalled)

	client := http.Client{}

	request, _ := http.NewRequest("POST", j.host+"/rest/api/2/issue", buffered)
	request.Header.Add("Content-Type", contentTypeJson)

	for _, cookie := range jiraAuthCookies {
		request.AddCookie(&http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, ErrCreatingIssue
	}

	if response.StatusCode >= 300 {
		return nil, ErrCreatingIssue
	}

	jsonDataFromHttp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var createdIssue *CreatedIssue
	err = json.Unmarshal(jsonDataFromHttp, &createdIssue)
	if err != nil {
		return nil, err
	}

	return createdIssue, nil
}
