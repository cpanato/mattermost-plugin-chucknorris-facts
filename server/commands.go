package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Chuck struct {
	IconURL string `json:"icon_url,omitempty"`
	ID      string `json:"id,omitempty"`
	URL     string `json:"url,omitempty"`
	Value   string `json:"value,omitempty"`
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:      "chuck-norris-facts",
		DisplayName:  "Chuck Norris Facts",
		Description:  "Chuck Norris Facts",
		AutoComplete: true,
	}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	resp, err := makeRequest("GET", "https://api.chucknorris.io/jokes/random", nil)
	if err != nil {
		return &model.CommandResponse{}, model.NewAppError("", "", nil, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var chuckFact Chuck
	err = json.NewDecoder(resp.Body).Decode(&chuckFact)
	if err != nil && err != io.EOF {
		return &model.CommandResponse{}, model.NewAppError("", "", nil, err.Error(), http.StatusInternalServerError)
	}

	if chuckFact.Value == "" {
		chuckPost := &model.Post{
			UserId:    p.botUserID,
			ChannelId: args.ChannelId,
			Message:   "Looks like even Chuck Norris fail, API server is down or anyother issue",
		}

		p.CreateBotPost(chuckPost)
		return &model.CommandResponse{}, nil
	}

	chuckPost := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   ">" + chuckFact.Value,
	}

	p.CreateBotPost(chuckPost)

	return &model.CommandResponse{}, nil
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Type:         model.POST_DEFAULT,
	}
}

func (p *Plugin) CreateBotPost(post *model.Post) (*model.Post, *model.AppError) {
	created, err := p.API.CreatePost(post)
	if err != nil {
		p.API.LogError("Couldn't send bot message", "err", err)
		return nil, err
	}

	return created, nil
}

func makeRequest(method, url string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
