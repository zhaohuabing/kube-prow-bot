package commands

import (
	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
	"github.com/xunzhuo/kube-prow-bot/pkg/utils"
)

func init() {
	registerCommand(closeCommandName, closeCommandFunc)
	registerCommand(reopenCommandName, reopenCommandFunc)
}

var closeCommandFunc = close
var closeCommandName CommandName = "close"

func close(args ...string) error {
	return utils.ExecGitHubCmd(
		config.Get().ISSUE_KIND,
		"-R",
		config.Get().GH_REPOSITORY,
		"close",
		config.Get().ISSUE_NUMBER)
}

var reopenCommandFunc = reopen
var reopenCommandName CommandName = "reopen"

func reopen(args ...string) error {
	return utils.ExecGitHubCmd(
		config.Get().ISSUE_KIND,
		"-R",
		config.Get().GH_REPOSITORY,
		"reopen",
		config.Get().ISSUE_NUMBER)
}
