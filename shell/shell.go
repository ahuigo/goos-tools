package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// command = exec.Command("sh", "-c", cmd) // waring!!!
func ExecCommand(command string, args ...string) (output string, errmsg string, errno int) {
	return ExecCommandWithEnvs(nil, command, args...)
}

func ExecCommandWithEnvs(envs []string, command string, args ...string) (output string, errmsg string, errno int) {
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stderr = &stderr
	// err := cmd.Run() // err: exit status 1
	outputBytes, err := cmd.Output()
	output = string(outputBytes)
	errmsg = stderr.String()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(interface{ ExitStatus() int }); ok {
				errmsg = err.Error() + " " + errmsg
				errno = status.ExitStatus()
			}
		} else if execErr, ok := err.(*exec.Error); ok {
			errmsg = execErr.Error() + " " + errmsg
			errno = 127
		}
	}
	return
}

func ExecCommandPipe(cmd string, stdinBytes []byte, args ...string) (output []byte, errmsg string, errno int) {
	var stdout, stderr bytes.Buffer
	var err error
	command := exec.Command(cmd, args...)
	if pipe, err := command.StdinPipe(); err != nil {
		return nil, err.Error(), 1
	} else {
		if _, err = pipe.Write(stdinBytes); err != nil {
			return nil, err.Error(), 1
		}
	}
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run() // err: exit status 1
	output = bytes.TrimSpace(stdout.Bytes())
	errmsg = strings.TrimSpace(stderr.String())
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(interface{ ExitStatus() int }); ok {
				errno = status.ExitStatus()
			}
		}
	}
	return
}

func parseShellCmdErrno(err error) (errno int) {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(interface{ ExitStatus() int }); ok {
			errno = status.ExitStatus()
		}
	} else if _, ok := err.(*exec.Error); ok {
		errno = 127
	}
	return
}

func ExecCommandWithEnvsFindLineErr1(envs []string, command string, args ...string) (outmsg, job_id string, err error) {
	var stdoutbuf bytes.Buffer
	var stderrbuf bytes.Buffer // 从stderr读取
	var stderr io.ReadCloser   // 按行读取
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), envs...)
	// if true {
	cmd.Stdout = &stdoutbuf
	stderr, _ = cmd.StderrPipe()
	// } else {
	// 	cmd.Stderr = &stdoutbuf
	// 	stderr, err = cmd.StdoutPipe()
	// }
	if err = cmd.Start(); err != nil {
		return outmsg, "", err
	}
	reader := bufio.NewReader(stderr)
	var line string
	for i := 0; i < 100; i++ {
		line, err = reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return outmsg, job_id, err
		}
		stderrbuf.WriteString(line)
		if line == "apple\n" {
			job_id = line
			outmsg = ""
			go func() {
				// ReadAll 必须在 Wait 之前，否则Wait 结束后，stderr 会被关闭(无法再读取stderr的内容)
				stderr, err2 := io.ReadAll(reader)
				err := cmd.Wait()
				errno := parseShellCmdErrno(err)
				stdout := stdoutbuf.String()
				fmt.Printf("async errno:%d, err:%v, stdout:%s, stderr:%s, err2:%v\n", errno, err, stdout, stderr, err2)
			}()
			return outmsg, job_id, err
		}
	}

	// go func() {
	err = cmd.Wait()
	outmsg = stdoutbuf.String() + " " + stderrbuf.String()
	errno := parseShellCmdErrno(err)
	if errno != 0 {
		err = errors.Wrapf(err, "errno:%d", errno)
	}
	return
}
