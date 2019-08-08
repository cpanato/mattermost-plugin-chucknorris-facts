package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

// OnActivate initialize the plugin
func (p *Plugin) OnActivate() error {
	p.API.LogDebug("Activating Chuck Norris Facts plugin")

	if err := p.ensureBotExists(); err != nil {
		return errors.Wrap(err, "failed to ensure bot user exists")
	}

	p.API.RegisterCommand(getCommand())

	p.API.LogDebug("Chuck Norris Facts plugin activated")

	return nil
}

func (p *Plugin) ensureBotExists() error {
	// Attempt to find an existing bot
	botUserIDBytes, err := p.API.KVGet(BotUserKey)
	if err != nil {
		return err
	}

	if botUserIDBytes == nil {
		// Create a bot since one doesn't exist
		p.API.LogDebug("Creating bot for chuck norris facts plugin")

		bot, err := p.API.CreateBot(&model.Bot{
			Username:    "chucknorrisbot",
			DisplayName: "Chuck Norris Fact Bot",
			Description: "Created by the Chuck Norris Facts plugin.",
		})
		if err != nil {
			return err
		}

		// Give it a profile picture
		err = p.API.SetProfileImage(bot.UserId, profileImage)
		if err != nil {
			p.API.LogError("Failed to set profile image for bot", "err", err)
		}

		p.API.LogDebug("Bot created for Chuck Norris Facts plugin")

		// Save the bot ID
		err = p.API.KVSet(BotUserKey, []byte(bot.UserId))
		if err != nil {
			return err
		}

		p.botUserID = bot.UserId
	} else {
		p.botUserID = string(botUserIDBytes)
	}

	return nil
}
