package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/eggmoid/mattermost-chatbot/command"
	"github.com/eggmoid/mattermost-chatbot/config"
	"github.com/eggmoid/mattermost-chatbot/util"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	botUserName    = "chatbot"
	botDisplayName = "Chatbot"
	botDescription = "Chatbot"
)

type Plugin struct {
	plugin.MattermostPlugin
}

func (p *Plugin) OnActivate() error {
	config.Mattermost = p.API

	if err := p.setUpBotUser(); err != nil {
		config.Mattermost.LogError("Failed to create a bot user", "Error", err.Error())
		return err
	}

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	cmds := command.GetCommands(p.API)
	for _, cmd := range cmds {
		p.API.RegisterCommand(cmd)
	}

	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	if config.Mattermost == nil {
		return nil
	}

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return err
	}

	swaggerFile, err := ioutil.ReadFile(filepath.Join(bundlePath, "site", "swagger.yaml"))
	if err != nil {
		return err
	}
	yaml.Unmarshal(swaggerFile, &config.Swagger)

	appFile, err := ioutil.ReadFile(filepath.Join(bundlePath, "site", "app.yaml"))
	if err != nil {
		return err
	}
	yaml.Unmarshal(appFile, &config.App)

	return nil
}

func (p *Plugin) setUpBotUser() error {
	botUserID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
		Description: botDescription,
	})
	if err != nil {
		config.Mattermost.LogError("Error in setting up bot user", "Error", err.Error())
		return err
	}

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return err
	}

	profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "icon.png"))
	if err != nil {
		return err
	}

	if appErr := p.API.SetProfileImage(botUserID, profileImage); appErr != nil {
		return errors.Wrap(appErr, "couldn't set profile image")
	}

	config.BotUserID = botUserID
	return nil
}

func (p *Plugin) ExecuteCommand(context *plugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	args, argErr := util.SplitArgs(commandArgs.Command)
	if argErr != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         argErr.Error(),
		}, nil
	}
	return command.CommandHandler.Handle(commandArgs, args...), nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}
