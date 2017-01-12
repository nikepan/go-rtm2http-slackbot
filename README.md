# go-rtm2http-slackbot
This bot can receive direst RTM messages from slack, send it to http server and send reply to to slack user

Bot sends GET-queries to server url with params: **user**, **message** and **email**.

Set settings in config.json:

`"SlackToken": "xoxb-xxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxx"`

`"HttpPath": "http://127.0.0.1/slackbot/"`

Bot can use basic http auth of http-server

`"BasicUser": "slackbot" // "" if no auth`

`"BasicPassword": "slackbot"`

Uses slack access package "github.com/nlopes/slack"