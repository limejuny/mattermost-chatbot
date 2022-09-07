package config

import (
	"github.com/mattermost/mattermost-server/v5/plugin"

	fn "github.com/thoas/go-funk"
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
	return fn.Map(d[k].([]interface{}), func(v interface{}) Dict {
		return v.(Dict)
	}).([]Dict)
}

var (
	Mattermost plugin.API
	BotUserID  string
	Service    Dict
	Swagger    Dict
	AppP       Dict
	AppSD      Dict
	AppST      Dict
	AppD       Dict
	AppT       Dict
	AppE       Dict
)
