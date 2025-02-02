package utils

import (
	"os/exec"
	"strings"

	"github.com/xunzhuo/kube-prow-bot/cmd/kube-prow-bot/config"
	"k8s.io/klog"
)

const GitHubCMD = "gh"

func ExecGitHubCommonCmd(args ...string) error {
	options := append([]string{config.Get().ISSUE_KIND, "-R", config.Get().GH_REPOSITORY}, args...)
	cmd := exec.Command(GitHubCMD, options...)
	cmdOutput, err := cmd.CombinedOutput()
	klog.Info("command: ", "gh ", strings.Join(options, " "), "\n", string(cmdOutput), "\n")
	if err != nil {
		klog.Error(err, "\n")
		return err
	}
	klog.Info(string(cmdOutput))
	// klog.Info("\nEnvs:", cmd.Environ())
	return nil
}

func ExecGitHubCmd(args ...string) error {
	cmd := exec.Command(GitHubCMD, args...)
	cmdOutput, err := cmd.CombinedOutput()
	klog.Info("command: ", "gh ", strings.Join(args, " "), "\n", string(cmdOutput), "\n")
	if err != nil {
		klog.Error(err, "\n")
		return err
	}

	klog.Info(string(cmdOutput))
	// klog.Info("\nEnvs:", cmd.Environ())
	return nil
}

func ExecGitHubCmdWithOutput(args ...string) (string, error) {
	cmd := exec.Command(GitHubCMD, args...)
	cmdOutput, err := cmd.CombinedOutput()
	klog.Info("command: ", "gh ", strings.Join(args, " "), "\n", string(cmdOutput), "\n")
	if err != nil {
		klog.Error(err, "\n")
		return "", err
	}

	// klog.Info("\nEnvs:", cmd.Environ())
	return string(cmdOutput), nil
}
