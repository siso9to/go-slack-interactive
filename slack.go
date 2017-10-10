package main

import (
	"fmt"
	"log"
	//"strings"

	"github.com/nlopes/slack"
	"strings"
)

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

// LstenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Only response mention to bot. Ignore else.
	/*
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}
	*/

	// Parse message
	/*
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 || m[0] != "error" {
		return fmt.Errorf("invalid message")
	}
	*/

	if ! strings.Contains(ev.Msg.Text,"state is now ALARM:") {
		return fmt.Errorf("invalid message")
	}

	// value is passed to message handler when request is approved.
	attachment := slack.Attachment{
		Text:       "確認しましたか？",
		Color:      "#f9a41b",
		CallbackID: "beer",
		Actions: []slack.AttachmentAction{
			{
				Name: actionSelect,
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "確認して対応する",
						Value: "do",
					},
					{
						Text:  "対応しない",
						Value: "no",
					},
				},
			},

			{
				Name:  actionCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
