# vk-api
vk-api is a simple Golang package that provides tools to interact with [VK API](https://vk.com/dev).
Supports authentication using login and password or by passing an access token right into the app, as well as long polling.
More features yet to come, probably gonna turn this into a full SDK.

## Installation:
`go get -t github.com/justblender/vk-api`

## Example:
```go
package main

import (
    "os"
    "github.com/justblender/vk-api"
)

func main() {
    client, err := vk_api.Authenticate(os.Getenv("login"), os.Getenv("password"), vk_api.ANDROID)
    if err != nil {
        panic("Couldn't authenticate, sad!")
    }

    parameters := vk_api.RequestParameters{
        "user_id": 1,
        "message": "Pasha, bring back the wall!",
    }

    if _, err = client.Request("messages.send", parameters); err != nil {
        panic("Pasha couldn't bring back the wall")
    }
}
```

## Using long polling:
```go
// don't forget about the previously created Client!..
longPoll, err := client.NewLongPoll()
if err != nil {
    panic("Couldn't create a new LongPoll")
}

for {
    messages, err := longPoll.Poll()
    if err != nil {
        panic("Some bad error occurred")
    }

    for _, message := range messages {
        if message.HasFlag(vk_api.CHAT) {
            fmt.Printf("New message from chat %d: %s\n", message.PeerID, message.Text)
        }
    }
}
```