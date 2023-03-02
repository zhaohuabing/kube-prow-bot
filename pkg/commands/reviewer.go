package commands

import (
	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
	"github.com/xunzhuo/kube-prow-bot/pkg/utils"
)

func init() {
	registerCommand(ccCommandName, ccCommandFunc)
	registerCommand(unCCCommandName, unCCCommand)
}

var ccCommandFunc = cc
var ccCommandName CommandName = "cc"

func cc(args ...string) error {
	var revs []string
	if len(args) == 0 {
		revs = []string{config.Get().LOGIN}
	} else {
		revs = formatUserIDs(args)
	}

	return addReviewers(revs)
}

var unCCCommand = unCC
var unCCCommandName CommandName = "uncc"

func unCC(args ...string) error {
	var revs []string
	if len(args) == 0 {
		revs = []string{config.Get().LOGIN}
	} else {
		revs = formatUserIDs(args)
	}

	return removeReviewers(revs)
}

func addReviewers(IDs []string) error {
	ids := []string{
		config.Get().ISSUE_KIND,
		"-R",
		config.Get().GH_REPOSITORY,
		"edit",
		config.Get().ISSUE_NUMBER,
	}
	for _, id := range IDs {
		ids = append(ids, "--add-reviewer", id)
	}
	return utils.ExecGitHubCmd(
		ids...)
}

func removeReviewers(IDs []string) error {
	ids := []string{
		config.Get().ISSUE_KIND,
		"-R",
		config.Get().GH_REPOSITORY,
		"edit",
		config.Get().ISSUE_NUMBER,
	}
	for _, id := range IDs {
		ids = append(ids, "--remove-reviewer", id)
	}
	return utils.ExecGitHubCmd(
		ids...)
}
