package cmd_test

import (
	"testing"

	"github.com/benmatselby/knope/client"
	"github.com/benmatselby/knope/cmd"
	"github.com/golang/mock/gomock"
)

func TestNewRootCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := client.NewMockAPI(ctrl)

	cmd := cmd.NewRootCommand(client)

	use := "knope"
	short := "CLI tool for retrieving data from AWS CodeBuild"

	if cmd.Use != use {
		t.Fatalf("expected use: %s; got %s", use, cmd.Use)
	}

	if cmd.Short != short {
		t.Fatalf("expected use: %s; got %s", short, cmd.Short)
	}
}
