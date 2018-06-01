package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

type EnvConfig struct {
	BotToken          string
	VerificationToken string
	BotID             string
	ChannelID         string
}

type Slack struct {
	client    *slack.Client
	botID     string
	channelID string
}

func main() {
	os.Exit(run())
}

func run() int {
	var env EnvConfig
	if err := env.setEnv(); err != nil {
		log.Print(err)
		return 1
	}

	log.Print("Start Slack")
	api := slack.New(env.BotToken)
	s := &Slack{
		client:    api,
		botID:     env.BotID,
		channelID: env.ChannelID,
	}

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			log.Print("[INFO] get message")
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}

	return 0
}

func lookupenv(key string) (string, error) {
	env, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("%s is invalid value", key)
	}
	return env, nil
}

func (e *EnvConfig) setEnv() error {
	var err error
	e.ChannelID, err = lookupenv("CHANNEL_ID")
	if err != nil {
		return err
	}

	e.BotID, err = lookupenv("BOT_ID")
	if err != nil {
		return err
	}

	e.VerificationToken, err = lookupenv("VERIFICATION_TOKEN")
	if err != nil {
		return err
	}

	e.BotToken, err = lookupenv("BOT_TOKEN")
	if err != nil {
		return err
	}

	return nil
}

func (s *Slack) handleMessageEvent(ev *slack.MessageEvent) error {
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	//ここでbotIDを取り除いたテキストがでるここでswitchでもいい
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	switch m[0] {

	default:
		return fmt.Errorf("invalid message")
	}

	attachment := slack.Attachment{}
	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, _, err := s.client.PostMessage(ev.Channel, "hello, help?", params); err != nil {
		return fmt.Errorf("[ERROR] Failed to post message: %s", err)
	}

	return nil
}
