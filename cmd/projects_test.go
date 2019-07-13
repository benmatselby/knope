package cmd_test

import (
	"testing"

	"github.com/benmatselby/knope/client"
	"github.com/benmatselby/knope/cmd"
	"github.com/golang/mock/gomock"
)

func TestNewListProjectsCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := client.NewMockAPI(ctrl)

	cmd := cmd.NewListProjectsCommand(client)

	use := "projects"
	short := "List all the projects"

	if cmd.Use != use {
		t.Fatalf("expected use: %s; got %s", use, cmd.Use)
	}

	if cmd.Short != short {
		t.Fatalf("expected use: %s; got %s", short, cmd.Short)
	}
}
