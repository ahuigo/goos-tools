// BEGIN: abpxx6d04wxr
package shell

import (
	"testing"
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
