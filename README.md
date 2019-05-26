# Simple Standup Bot for Telegram #

[![Go Report Card](https://goreportcard.com/badge/github.com/maddevsio/tgsbot)](https://goreportcard.com/report/github.com/maddevsio/tgsbot)

Bot helps to conduct asynchronous daily standup meetings 

## Available commands
```
/help - Display list of available commands
/join - Adds you to standup team of the group
/show - Shows who submit standups
/leave - Removes you from standup team of the group
/edit_deadline - Sets new standup deadline (you can use 10am format or 15:30 format)
/show_deadline - Shows current standup deadline
/remove_deadline - Removes standup deadline at all
/tz - Changes time zone of the group.
```

## Local usage
First you need to set env variables:

```
export TELEGRAM_TOKEN=yourTelegramTokenRecievedFromBotFather
export DEBUG=true
```
Then run. Note, you need `Docker` and `docker-compose` installed on your system
```
make run
```
To run tests: 
```
make test
```