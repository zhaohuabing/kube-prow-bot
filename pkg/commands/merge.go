package commands

import (
	"errors"
	"fmt"

	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
	"github.com/xunzhuo/kube-prow-bot/pkg/utils"
)

func init() {
	registerCommand(mergeCommandName, mergeCommandFunc)
}

var mergeCommandFunc = SafeMerge
var mergeCommandName CommandName = "merge"

func SafeMerge(args ...string) error {
	if config.Get().ISSUE_KIND != "pr" {
		return errors.New("you can only merge PRs")
	}
	if HasLabel("lgtm") && HasLabel("approved") && !HasLabel("do-not-merge") {
		if err := merge(args...); err != nil {
			return err
		}
	} else {
		return errors.New("you can only merge PRs when PR has lgtm, approved label, and without do-not-merge")
	}
	return nil
}

func merge(args ...string) error {
	var action string
	if len(args) == 0 {
		action = "--squash"
	} else {
		action = args[0]
		if action != "rebase" && action != "squash" {
			return errors.New("unsupported merge action, only support: rebase or squash")
		}
		action = fmt.Sprintf("--%s", action)
	}
	return utils.ExecGitHubCmd(
		config.Get().ISSUE_KIND,
		"-R",
		config.Get().GH_REPOSITORY,
		"merge",
		config.Get().ISSUE_NUMBER,
		action)
}
