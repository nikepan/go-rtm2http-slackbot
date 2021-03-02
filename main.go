package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

type Config struct {
	SlackToken    string
	HttpPath      string
	BasicUser     string
	BasicPassword string
}

func readJSON(fn string, v interface{}) {
	file, _ := os.Open(fn)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(v)
	if err != nil {
		log.Println("error:", err)
	}
}

var config Config

func main() {

	config = Config{}
	readJSON("config.json", &config)

	api := slack.New(config.SlackToken)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Println("Slackbot started")

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello
			case *slack.ConnectedEvent:
				log.Println("Connection counter:", ev.ConnectionCount)
			case *slack.MessageEvent:
				if ev.Msg.User == "" {
					continue
				}

				// only direct channels reply
				if !strings.HasPrefix(ev.Msg.Channel, "D") {
					continue
				}
				params := slack.OpenConversationParameters{ChannelID: ev.Msg.Channel}
				userChannel, _, _, err := rtm.OpenConversation(&params)
				if err != nil {
					log.Printf("Get channel Error: %v\n", err)
					continue
				}
				if userChannel.ID != ev.Msg.Channel {
					continue
				}

				userInfo, err := rtm.GetUserInfo(ev.Msg.User)
				if err != nil {
					log.Printf("User info Error: %v\n", err)
					continue
				}

				log.Printf("User: %v (%v); Message: %v\n", ev.Msg.User, userInfo.Profile.Email, ev.Msg.Text) // ev.Msg.User, ev.Msg.Channel

				client := &http.Client{}
				parameters := url.Values{}
				parameters.Add("user", ev.Msg.User)
				parameters.Add("message", ev.Msg.Text)
				parameters.Add("email", userInfo.Profile.Email)
				req, err := http.NewRequest("POST", config.HttpPath, strings.NewReader(parameters.Encode()))

				if config.BasicUser != "" {
					req.SetBasicAuth(config.BasicUser, config.BasicPassword)
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Printf("HTTP request Error: %v\n", err)
					continue
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("HTTP read Error %v\n", err)
				} else {
					log.Printf("HTTP response: %+v\n", string(body))
					rtm.SendMessage(rtm.NewOutgoingMessage(string(body), ev.Msg.Channel))
				}

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop

			default:
				// Ignore other events..
			}
		}
	}
}
