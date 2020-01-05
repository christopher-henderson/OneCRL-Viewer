package git

import (
	"fmt"
	"strings"
	"testing"
)

const repo = `H:\TestRepo`

var TestRepo = NewRepo(repo)

func TestGit(t *testing.T) {
	t.Log(TestRepo.Add("."))
	commit := &strings.Builder{}
	fmt.Fprintln(commit, "CERT ADD")
	fmt.Fprintln(commit, "added 1234")
	fmt.Fprintln(commit, "added 45678")
	t.Log(TestRepo.Commit(commit.String()))
	t.Log(TestRepo.HEAD())
}

func TestRemote(t *testing.T) {
	t.Log(TestRepo.Remote())
}
