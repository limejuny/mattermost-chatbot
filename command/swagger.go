package command

import (
	"fmt"
	"strings"

	"github.com/limejuny/mattermost-chatbot/config"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/samber/lo"
)

const (
	swaggerHelpText = commonHelpText +
		"* `/swagger list` - 서비스 또는 배포환경 목록 조회.\n" +
		"* `/swagger help` - 도움말.\n" +
		"예시) /swagger nuxx - NUXX의 swagger 목록 조회\n" +
		"예시) /swagger d - SM개발 swagger 목록 조회\n"
)

var SwaggerHandler = Handler{
	handlers: map[string]HandlerFunc{
		"/swagger/list": swaggerListCommand,
		"/swagger/help": swaggerHelpCommand,
	},
	defaultHandler: executeSwaggerDefault,
}

func executeSwaggerDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
	swagger := config.Swagger.A("swagger")
	codes := lo.Map(swagger, func(d config.Dict, _ int) string {
		return d.S("code")
	})
	envs := []string{"d", "t", "sd", "st", "e"}

	if len(args) > 0 && lo.Contains(envs, strings.ToLower(args[0])) {
		list := lo.Filter(swagger, func(d config.Dict, _ int) bool {
			_, ok := d.D("links")[args[0]]
			return ok
		})
		env := map[string]string{
			"sd": "SI개발",
			"st": "SI통시",
			"d":  "SM개발",
			"t":  "SM통시",
			"e":  "교육",
		}

		postCommandResponse(context,
			fmt.Sprintf("#### %s기 swagger 목록 조회\n", env[args[0]])+
				"| 구분 | 파트 | 서비스명 | 서비스코드 | 링크 |\n"+
				"| --- | --- | --- | --- | --- |\n"+
				strings.Join(lo.Map(list, func(d config.Dict, _ int) string {
					return fmt.Sprintf("| %s | %s | %s | %s | %s |", d.S("category"), d.S("part"), d.S("name"), strings.ToUpper(d.S("code")), d.D("links").S(args[0]))
				}), "\n"))
		return &model.CommandResponse{}
	}

	if len(args) > 0 && lo.Contains(codes, strings.ToLower(args[0])) {
		v, _ := lo.Find(swagger, func(d config.Dict) bool {
			return d.S("code") == strings.ToLower(args[0])
		})
		links := v.D("links")

		message := fmt.Sprintf("#### %s(%s) 서비스의 swagger 목록 조회\n", v.S("name"), strings.ToUpper(args[0])) +
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
	out += swaggerHelpText

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         out,
	}
}

func swaggerListCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	postCommandResponse(context,
		"| 구분 | 파트 | 서비스명 | 서비스코드 |\n"+
			"| --- | --- | --- | --- |\n"+
			strings.Join(lo.Map(config.Swagger.A("swagger"), func(d config.Dict, _ int) string {
				return fmt.Sprintf("| %s | %s | %s | %s |", d.S("category"), d.S("part"), d.S("name"), strings.ToUpper(d.S("code")))
			}), "\n"))
	return &model.CommandResponse{}
}

func swaggerHelpCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	postCommandResponse(context, swaggerHelpText)
	return &model.CommandResponse{}
}
