package util

import (
	"fmt"
	"strings"

	"github.com/lbryio/lbry.go/errors"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

var defaultChannel string
var defaultUsername string
var slackApi *slack.Client

// InitSlack Initializes a slack client with the given token and sets the default channel.
func InitSlack(token string, channel string, username string) {
	slackApi = slack.New(token)
	defaultChannel = channel
	defaultUsername = username
}

// SendToSlackUser Sends message to a specific user.
func SendToSlackUser(user, username, message string) error {
	if !strings.HasPrefix(user, "@") {
		user = "@" + user
	}
	return sendToSlack(user, username, message)
}

// SendToSlackChannel Sends message to a specific channel.
func SendToSlackChannel(channel, username, message string) error {
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}
	return sendToSlack(channel, username, message)
}

// SendToSlack Sends message to the default channel.
func SendToSlack(message string) error {

	if defaultChannel == "" {
		return errors.Err("no default slack channel set")
	}

	return sendToSlack(defaultChannel, defaultUsername, message)
}

// SendErrorToSlack Sends an error message to the default channel and to the process log.
func SendErrorToSlack(format string, a ...interface{}) error {
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	log.Errorln(message)
	return sendToSlack(defaultChannel, defaultUsername, ":sos: "+message)
}

// SendInfoToSlack Sends an info message to the default channel and to the process log.
func SendInfoToSlack(format string, a ...interface{}) error {
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	log.Debugln(message)
	return sendToSlack(defaultChannel, defaultUsername, ":information_source: "+message)
}

func sendToSlack(channel, username, message string) error {
	var err error

	if slackApi == nil {
		err = errors.Err("no slack token provided")
	} else {
		log.Debugln("slack: " + channel + ": " + message)
		_, _, err = slackApi.PostMessage(channel, message, slack.PostMessageParameters{Username: username})
	}

	if err != nil {
		log.Errorln("error sending to slack: " + err.Error())
		return err
	}

	return nil
}
