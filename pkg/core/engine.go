package core

import (
	"os"
	"regexp"
	"strings"

	"github.com/tetratelabs/multierror"
	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
	"github.com/xunzhuo/kube-prow-bot/pkg/commands"
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
	REVIEWERS   = []string{}
	APPROVERS   = []string{}
	MAINTAINERS = []string{}
)

func init() {
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

	roles := os.Getenv("REVIEWERS")
	REVIEWERS = strings.Split(roles, "\n")

	roles = os.Getenv("APPROVERS")
	APPROVERS = strings.Split(roles, "\n")

	roles = os.Getenv("MAINTAINERS")
	MAINTAINERS = strings.Split(roles, "\n")
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
	if belongTo(own, REVIEWERS) {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
		plugins = appendPlugins(plugins, MEMBERS_PLUGINS)
		plugins = appendPlugins(plugins, REVIEWERS_PLUGINS)
	}
	if belongTo(own, APPROVERS) {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
		plugins = appendPlugins(plugins, MEMBERS_PLUGINS)
		plugins = appendPlugins(plugins, REVIEWERS_PLUGINS)
		plugins = appendPlugins(plugins, APPROVERS_PLUGINS)
	}
	if belongTo(own, MAINTAINERS) {
		plugins = appendPlugins(plugins, AUTHOR_PLUGINS)
		plugins = appendPlugins(plugins, MEMBERS_PLUGINS)
		plugins = appendPlugins(plugins, REVIEWERS_PLUGINS)
		plugins = appendPlugins(plugins, APPROVERS_PLUGINS)
		plugins = appendPlugins(plugins, MAINTAINERS_PLUGINS)
	}

	return plugins
}

func RunCommands() error {
	messages := os.Getenv("MESSAGE")
	if messages == "" {
		return nil
	}

	var errs error
	for _, message := range strings.Split(messages, "\n") {
		cmd := CommandRegex.Find([]byte(message))
		if cmd != nil {
			c := strings.TrimSpace(string(cmd))
			c = strings.TrimPrefix(string(c), "/")
			c = strings.TrimSpace(c)
			cm := strings.Split(c, " ")
			if len(cm) == 1 {
				commandName := cm[0]
				if _, ok := constructOwnPlugins()[commandName]; !ok {
					klog.Info("User: ", config.Get().LOGIN, " does have this plugin: ", commandName, " privilege.")
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
				if _, ok := constructOwnPlugins()[commandName]; !ok {
					klog.Info("User: ", config.Get().LOGIN, " does have this plugin: ", commandName, " privilege.")
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
