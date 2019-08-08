package main

import (
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	BOT_USER_KEY = "ChuckNorrisFactsBot"
)

func main() {
	plugin.ClientMain(&Plugin{})
}
