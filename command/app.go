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

func postEnvAppList(context *model.CommandArgs, app []config.Dict, env string, args ...string) {
	v, ok := fn.Find(app, func(d config.Dict) bool {
		return d.S("code") == strings.ToLower(args[0])
	}).(config.Dict)
	if !ok || v == nil || v["links"] == nil {
		return
	}
	links := v.A("links")

	message := fmt.Sprintf("### %s(%s) 서비스 관련 사이트 목록 조회 (%s)\n\n", v.S("name"), strings.ToUpper(args[0]), env)

	for _, link := range links {
		message += fmt.Sprintf("#### %s\n", link.S("name"))
		message += "| 서비스 정보 | 서비스명 | 링크 | 접속 방법 |\n"
		message += "| --- | --- | --- | --- |\n"
		for _, app := range link.A("item") {
			message += fmt.Sprintf("| %s | %s | %s | %s |\n", app.S("type"), app.S("name"), app.S("link"), app.S("info"))
		}
	}
	postCommandResponse(context, message)
}

func executeAppDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
	svc := append(config.Service.D("app").A("public"), config.Service.D("app").A("private")...)
	codes := fn.Map(svc, func(d config.Dict) string {
		return d.S("code")
	}).([]string)
	envs := []string{"d", "t", "sd", "st", "e"}

	if len(args) > 0 && fn.Contains(envs, strings.ToLower(args[0])) {
		// TODO: 개발환경별 app 목록 조회
		return &model.CommandResponse{}
	}

	if len(args) > 0 && fn.Contains(codes, strings.ToLower(args[0])) {
		if len(args) > 1 && fn.Contains(envs, strings.ToLower(args[1])) {
			switch strings.ToLower(args[1]) {
			case "p":
				postEnvAppList(context, config.AppP.A("app"), "운영", args...)
			case "sd":
				postEnvAppList(context, config.AppSD.A("app"), "SI개발", args...)
			case "st":
				postEnvAppList(context, config.AppST.A("app"), "SI통시", args...)
			case "d":
				postEnvAppList(context, config.AppD.A("app"), "SM개발", args...)
			case "t":
				postEnvAppList(context, config.AppT.A("app"), "SM통시", args...)
			case "e":
				postEnvAppList(context, config.AppE.A("app"), "교육", args...)
			}
		} else {
			postEnvAppList(context, config.AppP.A("app"), "운영", args...)
			postEnvAppList(context, config.AppSD.A("app"), "SI개발", args...)
			postEnvAppList(context, config.AppST.A("app"), "SI통시", args...)
			postEnvAppList(context, config.AppD.A("app"), "SM개발", args...)
			postEnvAppList(context, config.AppT.A("app"), "SM통시", args...)
			postEnvAppList(context, config.AppE.A("app"), "교육", args...)
		}
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
		"### Public 클라우드(AWS)\n\n"+
			"| 서비스명 | 서비스코드 |\n"+
			"| --- | --- |\n"+
			strings.Join(fn.Map(config.Service.D("app").A("public"), func(d config.Dict) string {
				return fmt.Sprintf("| %s | %s |", d.S("name"), strings.ToUpper(d.S("code")))
			}).([]string), "\n"))
	postCommandResponse(context,
		"### Private 클라우드(Openshift, VMWare)\n\n"+
			"| 서비스명 | 서비스코드 |\n"+
			"| --- | --- |\n"+
			strings.Join(fn.Map(config.Service.D("app").A("private"), func(d config.Dict) string {
				return fmt.Sprintf("| %s | %s |", d.S("name"), strings.ToUpper(d.S("code")))
			}).([]string), "\n")+
			"\n| 컨피그서버 | NUPF-CNF |")
	return &model.CommandResponse{}
}

func appHelpCommand(context *model.CommandArgs, args ...string) *model.CommandResponse {
	postCommandResponse(context, appHelpText)
	return &model.CommandResponse{}
}
