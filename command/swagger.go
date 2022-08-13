package command

import (
	"fmt"
	"strings"

	"github.com/eggmoid/mattermost-chatbot/config"
	"github.com/mattermost/mattermost-server/v5/model"
	fn "github.com/thoas/go-funk"
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
	codes := fn.Map(config.Swagger, func(key string, value interface{}) string {
		return key
	}).([]string)
	envs := []string{"d", "t", "sd", "st"}

	if len(args) > 0 && fn.Contains(envs, strings.ToLower(args[0])) {
		// d, t, sd, st
		// d, t, sd, st
		// d, t, sd, st
		// d, t, sd, st
		return &model.CommandResponse{}
	}

	if len(args) > 0 && fn.Contains(codes, strings.ToLower(args[0])) {
		v := config.Swagger.D(strings.ToLower(args[0]))
		links := v.D("links")

		postCommandResponse(context,
			fmt.Sprintf("#### %s(%s) 서비스의 swagger 목록 조회\n", v.S("name"), strings.ToUpper(args[0]))+
				"| 배포 환경 | URL |\n"+
				"| --- | --- |\n"+
				fmt.Sprintf("| %s | %s |\n", "SI개발", links.S("sd"))+
				fmt.Sprintf("| %s | %s |\n", "SI통시", links.S("st"))+
				fmt.Sprintf("| %s | %s |\n", "SM개발", links.S("d"))+
				fmt.Sprintf("| %s | %s |\n", "SM통시", links.S("t")))
		return &model.CommandResponse{}
	}

	out := invalidCommand + "\n\n"
	out += swaggerHelpText

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         out,
	}
}

func swaggerListCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	postCommandResponse(context,
		"| 구분 | 파트 | 서비스명 | 서비스코드 |\n"+
			"| --- | --- | --- | --- |\n"+
			strings.Join(fn.Map(config.Swagger, func(key string, value interface{}) string {
				v := value.(config.Dict)
				return fmt.Sprintf("| %s | %s | %s | %s |", v.S("category"), v.S("part"), v.S("name"), strings.ToUpper(key))
			}).([]string), "\n"))
	return &model.CommandResponse{}
}

func swaggerHelpCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	postCommandResponse(context, swaggerHelpText)
	return &model.CommandResponse{}
}