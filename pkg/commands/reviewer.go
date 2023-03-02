package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
)

func init() {
	registerCommand(ccCommandName, ccCommandFunc)
	registerCommand(unCCCommandName, unCCCommand)
}

var ccCommandFunc = cc
var ccCommandName CommandName = "cc"

type reviewers struct {
	Reviewers     []string `json:"reviewers"`
	TeamReviewers []string `json:"team_reviewers"`
}

func cc(args ...string) error {
	var revs []string
	if len(args) == 0 {
		revs = []string{config.Get().LOGIN}
	} else {
		revs = formatUserIDs(args)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s/requested_reviewers", config.Get().GH_REPOSITORY, config.Get().ISSUE_NUMBER)
	data, err := json.Marshal(reviewers{Reviewers: revs})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.Get().GH_TOKEN))
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
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

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s/requested_reviewers", config.Get().GH_REPOSITORY, config.Get().ISSUE_NUMBER)
	data, err := json.Marshal(reviewers{Reviewers: revs})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.Get().GH_TOKEN))
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
