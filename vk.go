package vk_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	version   = 5.68
	methodURL = "https://api.vk.com/method/"
)

// Client is a structure that holds access token that allows VK API to allow us to request methods and do some cool stuff. Neat!
type Client struct {
	AccessToken string
}

// RequestParameters is an alias for map[string]interface{}
type RequestParameters map[string]interface{}

// NewClient returns an API structure and takes access token, error is returned if something goes wrong
func NewClient(authType authentication) (*Client, error) {
	accessToken, err := authType.retrieveAccessToken()
	if err != nil {
		return nil, err
	}
	return &Client{accessToken}, nil
}

// Request function makes an API request and returns a slice of bytes that represent JSON requestResponse
func (client *Client) Request(method string, parameters RequestParameters) ([]byte, error) {
	query := url.Values{
		"access_token": {client.AccessToken},
		"v":            {fmt.Sprint(version)},
	}

	// add parameters to query
	for key, value := range parameters {
		query.Set(key, fmt.Sprint(value))
	}

	response, err := http.Get(fmt.Sprintf("%s?%s", methodURL+method, query.Encode()))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var requestResponse map[string]*json.RawMessage
	if err = json.Unmarshal(body, &requestResponse); err != nil {
		return nil, err
	}

	if err, ok := requestResponse["error"]; ok {
		var errorData struct {
			Code    int    `json:"error_code"`
			Message string `json:"error_msg"`
		}

		json.Unmarshal(*err, &errorData)
		return nil, fmt.Errorf("error #%d: %s", errorData.Code, errorData.Message)
	}

	if response, ok := requestResponse["response"]; ok {
		return response.MarshalJSON()
	}
	return nil, errors.New("no response returned")

}
