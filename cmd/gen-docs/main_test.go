package main

import (
	"os"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	dir := t.TempDir()
	args := []string{"--man-page", "--website", "--doc-path", dir}
	err := run(args)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	manPage, err := os.ReadFile(dir + "/goctl-issue-create.1")
	if err != nil {
		t.Fatalf("error reading `goctl-issue-create.1`: %v", err)
	}
	if !strings.Contains(string(manPage), `\fBgoctl issue create`) {
		t.Fatal("man page corrupted")
	}

	markdownPage, err := os.ReadFile(dir + "/gh_issue_create.md")
	if err != nil {
		t.Fatalf("error reading `gh_issue_create.md`: %v", err)
	}
	if !strings.Contains(string(markdownPage), `## goctl issue create`) {
		t.Fatal("markdown page corrupted")
	}
}
