package main

import (
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	// BotUserKey the chuck name :)
	BotUserKey = "ChuckNorrisFactsBot"
)

func main() {
	plugin.ClientMain(&Plugin{})
}
