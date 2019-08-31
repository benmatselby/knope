package version_test

import (
	"testing"

	"github.com/benmatselby/knope/version"
)

func TestGitCommitPresent(t *testing.T) {
	version.GITCOMMIT = "testing"

	if "testing" != version.GITCOMMIT {
		t.Fatalf("expected GITCOMMIT to be testing, got %s", version.GITCOMMIT)
	}
}
