package util

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

// TODO: add retry logic
func RunCommand(ctx context.Context, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	klog.Infof("running command: %s %s", command, strings.Join(args, " "))
	return cmd.Run()
}

func RunCommandWithRetry(ctx context.Context, retryCount int, command string, args ...string) (err error) {
	sleep := 1
	for i := 1; i <= retryCount; i++ {
		err = RunCommand(ctx, command, args...)
		if err == nil {
			return
		}
		klog.Errorf("[%d] failed to run command: %s %s, error: %s", i, command, strings.Join(args, " "), err.Error())
		sleep = min(i*2-1, 10)
		time.Sleep(time.Second * time.Duration(sleep))
	}
	return
}

// TODO: add retry logic
func RunCommandWithTimeout(ctx context.Context, timeout time.Duration, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	return cmd.Run()
}

// TODO: add retry logic
func RunCommandWithTimeoutAndReturn(ctx context.Context, timeout time.Duration, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// RunCommandWithOutput 执行命令并返回输出
func RunCommandWithOutput(ctx context.Context, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	klog.Infof("running command: %s %s", command, strings.Join(args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
