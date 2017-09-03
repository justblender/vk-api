package vk_api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type clientType struct{ id, secret string }

var (
	ANDROID = &clientType{"2274003", "hHbZxrka2uZ6jB1inYsH"}
	WINDOWS = &clientType{"3697615", "AlVXZFMUqyrnABp8ncuU"}
	IOS     = &clientType{"3140623", "VeWdmVclDCtn6ihuP1nt"}
)

// Authenticate returns a VK clientType that can be accessed with login and password,
// error is returned if direct authentication to get access token has failed
func Authenticate(login, password string, clientType *clientType) (*Client, error) {
	if clientType == nil {
		return nil, errors.New("clientType can't be nil, use one of the prepared client types")
	}

	query := url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientType.id},
		"client_secret": {clientType.secret},
		"username":      {login},
		"password":      {password},
	}

	response, err := http.Get("https://oauth.vk.com/token?" + query.Encode())
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var authResponse struct {
		AccessToken      string `json:"access_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err := json.Unmarshal(body, &authResponse); err != nil {
		return nil, err
	}

	if authResponse.Error != "" {
		return nil, errors.New(authResponse.ErrorDescription)
	}
	return NewClient(authResponse.AccessToken)
}
