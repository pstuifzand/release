package main

import (
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
)

func TestGitCommitSkipsCIForVersionOnlyChange(t *testing.T) {
	cmd := gitCommit(semver.MustParse("1.2.3"))

	assert.Equal(t, []string{
		"git", "commit", "-n", "-m", "Increase version to 1.2.3 [CI SKIP]",
	}, cmd.Args)
}
