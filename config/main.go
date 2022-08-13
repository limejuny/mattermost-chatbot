package config

import (
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Dict map[string]interface{}

func (d Dict) D(k string) Dict {
	return d[k].(Dict)
}

func (d Dict) S(k string) string {
	return d[k].(string)
}

func (d Dict) I(k string) int {
	return int(d[k].(float64))
}

func (d Dict) A(k string) []Dict {
	return d[k].([]Dict)
}

var (
	Mattermost plugin.API
	BotUserID  string
	Swagger    Dict
)
