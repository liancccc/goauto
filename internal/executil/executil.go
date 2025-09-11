package executil

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/projectdiscovery/gologger"
)

func genCommands(cmd string) []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"powershell", "-Command", cmd}
	default:
		return []string{"bash", "-c", cmd}
	}
}

func RunCommandSteamOutput(cmd string, timeoutRaws ...string) (string, error) {
	commands := genCommands(cmd)
	var outputBuilder strings.Builder
	var realCmd *exec.Cmd
	var timeout int
	var ctx context.Context
	var cancel context.CancelFunc

	if len(timeoutRaws) > 0 && timeoutRaws[0] != "" {
		timeout = calcTimeout(timeoutRaws[0])
	}

	gologger.Info().Msgf("Execute: %s", strings.Join(commands, " "))

	if timeout > 0 {
		gologger.Info().Msgf("Timeout: %s", timeoutRaws[0])
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		realCmd = exec.CommandContext(ctx, commands[0], commands[1:]...)
	} else {
		realCmd = exec.Command(commands[0], commands[1:]...)
	}

	// 输出处理
	cmdReader, _ := realCmd.StdoutPipe()
	errReader, _ := realCmd.StderrPipe()
	scanner := bufio.NewScanner(cmdReader)
	errScanner := bufio.NewScanner(errReader)
	go func() {
		for scanner.Scan() {
			out := scanner.Text()
			gologger.Debug().Msg(out)
			outputBuilder.WriteString(out + "\n")
		}
	}()
	go func() {
		for errScanner.Scan() {
			out := errScanner.Text()
			gologger.Debug().Msg(out)
			outputBuilder.WriteString(out + "\n")
		}
	}()

	if err := realCmd.Start(); err != nil {
		return outputBuilder.String(), err
	}

	if err := realCmd.Wait(); err != nil && realCmd.ProcessState != nil {
		if ctx != nil && ctx.Err() != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
			gologger.Info().Msgf("kill this process after %d seconds", timeout)
			_ = realCmd.Process.Kill()
			var killCmd string
			if runtime.GOOS == "windows" {
				killCmd = fmt.Sprintf("taskkill /F /PID %d", realCmd.Process.Pid)
			} else {
				killCmd = fmt.Sprintf("kill -9 %d", realCmd.Process.Pid)
			}
			killCommands := genCommands(killCmd)
			_ = exec.Command(killCommands[0], killCommands[1:]...).Run()
			gologger.Info().Msgf("Execute: %s", strings.Join(killCommands, " "))
			return outputBuilder.String(), fmt.Errorf("command timeout after %d seconds", timeout)
		}
		return outputBuilder.String(), err
	}

	gologger.Info().Msgf("Execute Complete: %s", strings.Split(cmd, " ")[0])
	return outputBuilder.String(), nil
}

// calcTimeout 解析超时参数
func calcTimeout(raw string) int {
	raw = strings.ToLower(strings.TrimSpace(raw))
	seconds := raw
	multiply := 1

	matched, _ := regexp.MatchString(`.*[a-z]`, raw)
	if matched {
		unitTime := fmt.Sprintf("%c", raw[len(raw)-1])
		seconds = raw[:len(raw)-1]
		switch unitTime {
		case "s":
			multiply = 1
		case "m":
			multiply = 60
		case "h":
			multiply = 3600
		}
	}

	timeout, err := strconv.Atoi(seconds)
	if err != nil {
		return 0
	}
	return timeout * multiply
}
