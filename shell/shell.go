package shell

import (
	"bytes"
	"os/exec"
	"strings"
)

// command = exec.Command("sh", "-c", cmd) // waring!!!
func ExecCommand(cmd string, args ...string) (output string, errmsg string, errno int) {
	var stdout, stderr bytes.Buffer
	command := exec.Command(cmd, args...)
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run() // err: exit status 1
	output = strings.TrimSpace(stdout.String())
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

func ExecCommandPipe(cmd string,stdin []byte, args ...string) (output []byte, errmsg string, errno int) {
	var stdout, stderr bytes.Buffer
	var err error
	command := exec.Command(cmd, args...)
	if pipe, err := command.StdinPipe() ; err!=nil{
		return nil, err.Error(), 1
	}else{
		if _, err = pipe.Write(stdin); err!=nil{
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