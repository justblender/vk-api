package vk_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type device struct {
	id, secret string
}

var (
	ANDROID = &device{"2274003", "hHbZxrka2uZ6jB1inYsH"}
	WINDOWS = &device{"3697615", "AlVXZFMUqyrnABp8ncuU"}
	IOS     = &device{"3140623", "VeWdmVclDCtn6ihuP1nt"}
)

type authentication interface {
	retrieveAccessToken() (string, error)
}

// NoAuthentication is a type of authentication that doesn't do any authentication at all (lol), it just creates a Client from an AccessToken field
type NoAuthentication struct {
	authentication

	AccessToken string
}

// DirectAuthentication is a type of authentication that allows to access VK API using user login and password
type DirectAuthentication struct {
	authentication

	Username, Password string
	Device             *device
}

// ClientCredentialsFlow is a type of authentication that allows to access VK API (plus "secure" methods)
// using client ID and client secret that can be found in apps settings page
type ClientCredentialsFlow struct {
	authentication

	ClientID, ClientSecret string
}

func (no NoAuthentication) retrieveAccessToken() (accessToken string, err error) {
	if no.AccessToken != "" {
		return no.AccessToken, nil
	}
	err = errors.New("invalid access token")
	return
}

func (direct DirectAuthentication) retrieveAccessToken() (accessToken string, err error) {
	switch {
	case direct.Device == nil:
		err = errors.New("device can't be nil, use one of the prepared client types")
		return
	case direct.Username == "", direct.Password == "":
		err = errors.New("invalid login or password credentials")
		return
	}

	return retrieveAccessToken("token", url.Values{
		"grant_type":    {"password"},
		"client_id":     {direct.Device.id},
		"client_secret": {direct.Device.secret},
		"username":      {direct.Username},
		"password":      {direct.Password},
	})
}

func (client ClientCredentialsFlow) retrieveAccessToken() (accessToken string, err error) {
	if client.ClientID == "" || client.ClientSecret == "" {
		err = errors.New("invalid client credentials")
		return
	}

	return retrieveAccessToken("access_token", url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
	})
}

func retrieveAccessToken(typ string, values url.Values) (accessToken string, err error) {
	response, err := http.Get(fmt.Sprintf("https://oauth.vk.com/%s?%s", typ, values.Encode()))
	if err != nil {
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var authResponse struct {
		AccessToken      string `json:"access_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err = json.Unmarshal(body, &authResponse); err != nil {
		return
	}

	switch {
	case authResponse.Error != "":
		err = errors.New(authResponse.ErrorDescription)
		return
	case authResponse.AccessToken == "":
		err = errors.New("invalid access token")
		return

	default:
		return authResponse.AccessToken, nil
	}
}
