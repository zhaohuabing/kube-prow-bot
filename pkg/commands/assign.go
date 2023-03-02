package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
)

func init() {
	registerCommand(assignCommandName, assignCommandFunc)
	registerCommand(unassignCommandName, unassignCommand)
}

var assignCommandFunc = assign
var assignCommandName CommandName = "assign"

type assignees struct {
	Assignees []string `json:"assignees"`
}

func assign(args ...string) error {
	var assignee []string
	if len(args) == 0 {
		assignee = []string{config.Get().LOGIN}
	} else {
		assignee = formatUserIDs(args)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/assignees", config.Get().GH_REPOSITORY, config.Get().ISSUE_NUMBER)
	data, err := json.Marshal(assignees{Assignees: assignee})
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

var unassignCommand = unassign
var unassignCommandName CommandName = "unassign"

func unassign(args ...string) error {
	var assignee []string
	if len(args) == 0 {
		assignee = []string{config.Get().LOGIN}
	} else {
		assignee = formatUserIDs(args)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/assignees", config.Get().GH_REPOSITORY, config.Get().ISSUE_NUMBER)
	data, err := json.Marshal(assignees{Assignees: assignee})
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

func formatUserIDs(names []string) []string {
	formatedIDs := []string{}
	for _, name := range names {
		formatedIDs = append(formatedIDs, strings.TrimPrefix(name, "@"))
	}
	return formatedIDs
}
