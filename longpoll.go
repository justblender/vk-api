package vk_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const flagNewMessage = 4

// What should I say here?
type flag int

const (
	UNREAD flag = 1 << iota
	OUTBOX
	REPLIED
	IMPORTANT
	CHAT
	FRIENDS
	SPAM
	DELETED
	FIXED
	MEDIA
)

// LongPoll is a structure of data for long poll to properly work
type LongPoll struct {
	client *Client

	Server string `json:"server"`
	Key    string `json:"key"`
	TS     int64  `json:"ts"`
}

// Message is a structure of data that holds contents from received message
type Message struct {
	ID, PeerID, Timestamp, Flags int64
	Subject, Text                string
	Attachments                  map[string]interface{}
}

// NewLongPoll creates a new LongPoll object,
// use LongPoll#Poll method to start listening for incoming messages
func (client *Client) NewLongPoll() (*LongPoll, error) {
	longPoll := &LongPoll{client: client}
	if err := longPoll.update(); err != nil {
		return nil, err
	}
	return longPoll, nil
}

func (longPoll *LongPoll) update() error {
	response, err := longPoll.client.Request("messages.getLongPollServer", RequestParameters{
		"need_pts":   0, // not sure what this does, but I'm probably going to implement the use of this
		"lp_version": 2,
	})

	if err != nil {
		return err
	}
	if err = json.Unmarshal(response, &longPoll); err != nil {
		return err
	}
	return nil
}

// Poll waits for an incoming message and returns it when received
// Returns an error if something goes wrong while waiting
func (longPoll *LongPoll) Poll() ([]Message, error) {
	query := url.Values{
		"act": {"a_check"},
		"key": {longPoll.Key},
		"ts":  {fmt.Sprint(longPoll.TS)},

		"version": {"2"},
		"wait":    {"25"},
		"mode":    {"2"},
	}

	response, err := http.Get(fmt.Sprintf("https://%s?%s", longPoll.Server, query.Encode()))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var longPollResponse struct {
		Failed  int64           `json:"failed"`
		TS      int64           `json:"ts"`
		Updates [][]interface{} `json:"updates"`
	}

	if err = json.Unmarshal(body, &longPollResponse); err != nil {
		return nil, err
	}

	var messages []Message
	switch longPollResponse.Failed {
	case 0:
		for _, update := range longPollResponse.Updates {
			if update[0].(float64) != flagNewMessage {
				continue
			}

			var message Message
			message.ID = int64(update[1].(float64))
			message.Flags = int64(update[2].(float64))
			message.PeerID = int64(update[3].(float64))
			message.Timestamp = int64(update[4].(float64))
			message.Text = update[5].(string)
			message.Attachments = update[6].(map[string]interface{})
			messages = append(messages, message)
		}
		fallthrough
	case 1:
		longPoll.TS = longPollResponse.TS
	case 2, 3:
		if err = longPoll.update(); err != nil {
			return nil, err
		}
	}
	return messages, nil
}

// HasFlag returns true if message contains specified flag and false if not
func (message Message) HasFlag(flag flag) bool {
	return int(message.Flags)&int(flag) != 1
}
