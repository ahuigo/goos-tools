// BEGIN: abpxx6d04wxr
package shell

import (
	"testing"
	"time"
)

func TestExecCommand(t *testing.T) {
	// stdout, stderr, errno := ExecCommand("sh", "-c", `echo abc;ls | wc;no-such-command`)
	stdout, stderr, errno := ExecCommand("lsx", "-c", `echo abc;ls | wc;no-such-command`)
	t.Log("stdout:", stdout)
	t.Log("stderr:", stderr)
	t.Log("errno:", errno)
	if stderr==""{
		t.Fatal("stderr is empty")
	}
}

func TestExecCommandWithEnvsFirstLine(t *testing.T) {
	envs := []string{"ENV1=value1", "ENV2=value2"}
	command := "sh"
	args := []string{"-c","echo apple >&2;echo apple ; sleep 0.5; echo stderr1 >&2; echo stdout1; exit 1"}
	// args := []string{"-c","touch ../a.log;sleep 20;"}

	output,job_id, err := ExecCommandWithEnvsFindLineErr1(envs, command, args...)
	if job_id == "" {
		t.Fatal("job_id is empty")
	}
	t.Log("output:", output)
	t.Log("err:", err)
	time.Sleep(1000*time.Millisecond)
}