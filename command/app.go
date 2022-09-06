package command

import (
	"fmt"
	"strings"

	"github.com/eggmoid/mattermost-chatbot/config"
	"github.com/mattermost/mattermost-server/v5/model"
	fn "github.com/thoas/go-funk"
)

const (
	appHelpText = commonHelpText +
		"* `/app list` - 서비스 또는 배포환경 목록 조회.\n" +
		"* `/app help` - 도움말.\n" +
		"예시) /app nuxx - NUXX의 app 목록 조회\n" +
		"예시) /app d - SM개발 app 목록 조회\n"
)

var AppHandler = Handler{
	handlers: map[string]HandlerFunc{
		"/app/list": appListCommand,
		"/app/help": appHelpCommand,
	},
	defaultHandler: executeAppDefault,
}

func executeAppDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
	app := config.App.A("app")
	codes := fn.Map(app, func(d config.Dict) string {
		return d.S("code")
	}).([]string)
	envs := []string{"d", "t", "sd", "st", "e"}

	if len(args) > 0 && fn.Contains(envs, strings.ToLower(args[0])) {
		list := fn.Filter(app, func(d config.Dict) bool {
			_, ok := d.D("links")[args[0]]
			return ok
		}).([]config.Dict)
		env := map[string]string{
			"sd": "SI개발",
			"st": "SI통시",
			"d":  "SM개발",
			"t":  "SM통시",
			"e":  "교육",
		}

		postCommandResponse(context,
			fmt.Sprintf("#### %s기 app 목록 조회\n", env[args[0]])+
				"| 구분 | 파트 | 서비스명 | 서비스코드 | 링크 |\n"+
				"| --- | --- | --- | --- | --- |\n"+
				strings.Join(fn.Map(list, func(d config.Dict) string {
					return fmt.Sprintf("| %s | %s | %s | %s | %s |", d.S("category"), d.S("part"), d.S("name"), strings.ToUpper(d.S("code")), d.D("links").S(args[0]))
				}).([]string), "\n"))
		return &model.CommandResponse{}
	}

	if len(args) > 0 && fn.Contains(codes, strings.ToLower(args[0])) {
		v := fn.Find(app, func(d config.Dict) bool {
			return d.S("code") == strings.ToLower(args[0])
		}).(config.Dict)
		links := v.D("links")

		message := fmt.Sprintf("#### %s(%s) 서비스의 app 목록 조회\n", v.S("name"), strings.ToUpper(args[0])) +
			"| 배포 환경 | URL |\n" +
			"| --- | --- |\n"
		if val, ok := links["sd"]; ok {
			message += fmt.Sprintf("| %s | %s |\n", "SI개발", val)
		}
		if val, ok := links["st"]; ok {
			message += fmt.Sprintf("| %s | %s |\n", "SI통시", val)
		}
		if val, ok := links["d"]; ok {
			message += fmt.Sprintf("| %s | %s |\n", "SM개발", val)
		}
		if val, ok := links["t"]; ok {
			message += fmt.Sprintf("| %s | %s |\n", "SM통시", val)
		}
		if val, ok := links["e"]; ok {
			message += fmt.Sprintf("| %s | %s |\n", "교육", val)
		}
		postCommandResponse(context, message)
		return &model.CommandResponse{}
	}

	out := invalidCommand + "\n\n"
	out += appHelpText

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         out,
	}
}

func appListCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	postCommandResponse(context,
		"| 구분 | 파트 | 서비스명 | 서비스코드 |\n"+
			"| --- | --- | --- | --- |\n"+
			strings.Join(fn.Map(config.App.A("app"), func(d config.Dict) string {
				return fmt.Sprintf("| %s | %s | %s | %s |", d.S("category"), d.S("part"), d.S("name"), strings.ToUpper(d.S("code")))
			}).([]string), "\n"))
	return &model.CommandResponse{}
}

func appHelpCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	postCommandResponse(context, appHelpText)
	return &model.CommandResponse{}
}
