package config

import (
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/samber/lo"
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
	return lo.Map(d[k].([]interface{}), func(v interface{}, _ int) Dict {
		return v.(Dict)
	})
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
