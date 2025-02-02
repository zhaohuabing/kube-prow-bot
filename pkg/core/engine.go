package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/tetratelabs/multierror"
	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
	"github.com/xunzhuo/kube-prow-bot/pkg/commands"
	"github.com/xunzhuo/kube-prow-bot/pkg/utils"
	"gopkg.in/yaml.v3"
	"k8s.io/klog"
)

var CommandRegex = regexp.MustCompile(`\/.+`)

var (
	MEMBERS_PLUGINS     = []string{}
	REVIEWERS_PLUGINS   = []string{}
	MAINTAINERS_PLUGINS = []string{}
	APPROVERS_PLUGINS   = []string{}
	AUTHOR_PLUGINS      = []string{}
	COMMON_PLUGINS      = []string{}
)

var (
	ROLES = Roles{
		Maintainers: []string{},
		Approvers:   []string{},
		Reviewers:   []string{},
	}
)

func init() {
	constructPlugins()
	constructRoles()
}

type Roles struct {
	Maintainers []string `yaml:"maintainers"`
	Approvers   []string `yaml:"approvers"`
	Reviewers   []string `yaml:"reviewers"`
}

func constructPlugins() {
	plugins := os.Getenv("COMMON_PLUGINS")
	COMMON_PLUGINS = strings.Split(plugins, "\n")

	plugins = os.Getenv("AUTHOR_PLUGINS")
	AUTHOR_PLUGINS = strings.Split(plugins, "\n")

	plugins = os.Getenv("MEMBERS_PLUGINS")
	MEMBERS_PLUGINS = strings.Split(plugins, "\n")

	plugins = os.Getenv("REVIEWERS_PLUGINS")
	REVIEWERS_PLUGINS = strings.Split(plugins, "\n")

	plugins = os.Getenv("APPROVERS_PLUGINS")
	APPROVERS_PLUGINS = strings.Split(plugins, "\n")

	plugins = os.Getenv("MAINTAINERS_PLUGINS")
	MAINTAINERS_PLUGINS = strings.Split(plugins, "\n")
}

func constructRoles() {
	valid := false
	if rs, err := constructOWNERRoles(); err == nil && rs != nil {
		ROLES = *rs
		data, _ := json.Marshal(ROLES)
		if len(ROLES.Reviewers) != 0 ||
			len(ROLES.Approvers) != 0 ||
			len(ROLES.Maintainers) != 0 {
			valid = true
			klog.Info("PROJECT OWNER ROLES: \n", string(data))
		}
	} else {
		klog.Error(err)
	}

	if valid {
		return
	}

	ROLES = constructEnvRoles()
	data, _ := json.Marshal(ROLES)
	klog.Info("PROJECT ENV ROLES: \n", string(data))
}

func constructEnvRoles() Roles {
	roleList := Roles{
		Maintainers: []string{},
		Approvers:   []string{},
		Reviewers:   []string{},
	}

	roles := os.Getenv("REVIEWERS")
	REVIEWERS := strings.Split(roles, "\n")
	roleList.Reviewers = append(roleList.Reviewers, REVIEWERS...)

	roles = os.Getenv("APPROVERS")
	APPROVERS := strings.Split(roles, "\n")
	roleList.Approvers = append(roleList.Approvers, APPROVERS...)

	roles = os.Getenv("MAINTAINERS")
	MAINTAINERS := strings.Split(roles, "\n")
	roleList.Maintainers = append(roleList.Maintainers, MAINTAINERS...)

	return roleList
}

func constructOWNERRoles() (*Roles, error) {
	roles := &Roles{
		Maintainers: []string{},
		Approvers:   []string{},
		Reviewers:   []string{},
	}

	branch, err := utils.ExecGitHubCmdWithOutput("api", fmt.Sprintf("/repos/%s", os.Getenv("GH_REPOSITORY")), "-q", ".default_branch")
	if err != nil {
		return nil, err
	}
	ir := strings.Split(branch, "\n")
	for _, i := range ir {
		if strings.TrimSpace(i) != "" && strings.TrimSpace(i) != "\n" {
			branch = strings.TrimSpace(i)
			break
		}
	}
	r, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/OWNERS", os.Getenv("GH_REPOSITORY"), strings.TrimSuffix(branch, "\n")))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal([]byte(body), roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func belongTo(name string, groups []string) bool {
	for _, group := range groups {
		if strings.TrimSpace(group) == name {
			return true
		}
	}

	return false
}

func appendPlugins(plugins map[string]struct{}, target []string) map[string]struct{} {
	for _, t := range target {
		plugins[t] = struct{}{}
	}
	return plugins
}

func constructOwnPlugins() map[string]struct{} {
	var plugins = map[string]struct{}{}
	plugins = appendPlugins(plugins, COMMON_PLUGINS)

	own := config.Get().LOGIN
	if own == os.Getenv("AUTHOR") {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
	}
	if os.Getenv("AUTHOR_ASSOCIATION") != "NONE" && os.Getenv("AUTHOR_ASSOCIATION") != "" {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
		plugins = appendPlugins(plugins, MEMBERS_PLUGINS)
	}
	if belongTo(own, ROLES.Reviewers) {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
		plugins = appendPlugins(plugins, MEMBERS_PLUGINS)
		plugins = appendPlugins(plugins, REVIEWERS_PLUGINS)
	}
	if belongTo(own, ROLES.Approvers) {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
		plugins = appendPlugins(plugins, MEMBERS_PLUGINS)
		plugins = appendPlugins(plugins, REVIEWERS_PLUGINS)
		plugins = appendPlugins(plugins, APPROVERS_PLUGINS)
	}
	if belongTo(own, ROLES.Maintainers) {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
		plugins = appendPlugins(plugins, MEMBERS_PLUGINS)
		plugins = appendPlugins(plugins, REVIEWERS_PLUGINS)
		plugins = appendPlugins(plugins, APPROVERS_PLUGINS)
		plugins = appendPlugins(plugins, MAINTAINERS_PLUGINS)
	}

	return plugins
}

func listPlugins(cmds map[string]struct{}) []string {
	cmdList := []string{}
	for cmd := range cmds {
		cmdList = append(cmdList, cmd)
	}
	return cmdList
}

func RunCommands() error {
	var errs error

	messages := os.Getenv("MESSAGE")
	prState := os.Getenv("PR_STATE")

	if messages == "" && prState != "approved" {
		return nil
	}

	ownerPlugins := constructOwnPlugins()

	klog.Info("Available commands for @", config.Get().LOGIN, ":\n", listPlugins(ownerPlugins))

	hasRunApprove := false

	if prState == "approved" {
		if _, ok := ownerPlugins["approve"]; ok {
			cfunc, found := commands.GetCommand(commands.CommandName("approve"))
			if found {
				klog.Info("Running command: ", "approve")
				hasRunApprove = true
				if err := cfunc(); err != nil {
					errs = multierror.Append(errs, err)
				}
			}
		}
	}

	for _, message := range strings.Split(messages, "\n") {
		cmd := CommandRegex.Find([]byte(message))
		if cmd != nil {
			c := strings.TrimSpace(string(cmd))
			c = strings.TrimPrefix(string(c), "/")
			c = strings.TrimSpace(c)
			cm := strings.Split(c, " ")
			if len(cm) == 1 {
				commandName := cm[0]
				if _, ok := ownerPlugins[commandName]; !ok {
					klog.Info("User: ", config.Get().LOGIN, " does not have this plugin: ", commandName, " privilege.")
					continue
				}
				if commandName == "approve" && hasRunApprove {
					continue
				}
				cfunc, found := commands.GetCommand(commands.CommandName(commandName))
				if found {
					klog.Info("Running command: ", commandName)
					if err := cfunc(); err != nil {
						errs = multierror.Append(errs, err)
					}
				}
			} else if len(cm) > 1 {
				commandName := cm[0]
				commandInput := cm[1:]
				if _, ok := ownerPlugins[commandName]; !ok {
					klog.Info("User: ", config.Get().LOGIN, " does not have this plugin: ", commandName, " privilege.")
					continue
				}
				cfunc, found := commands.GetCommand(commands.CommandName(commandName))
				if found {
					klog.Info("Running command: ", commandName)
					if err := cfunc(commandInput...); err != nil {
						errs = multierror.Append(errs, err)
					}
				}
			}
		}
	}

	commands.SafeMerge()

	return errs
}
