package cmd_test

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/codebuild"
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

func TestDisplayProjects(t *testing.T) {
	tt := []struct {
		name     string
		projects []string
		expected string
		err      error
	}{
		{name: "can return the projects as a list", projects: []string{"a", "b", "c", "d"}, expected: "a\nb\nc\nd\n", err: nil},
		{name: "can return the projects as an ordered list", projects: []string{"zebra", "brilliance", "woodland", "blue-sea"}, expected: "blue-sea\nbrilliance\nwoodland\nzebra\n", err: nil},
		{name: "does not return anything if there is an error", projects: []string{}, expected: "", err: errors.New("there was an error")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := client.NewMockAPI(ctrl)

			var projects []*string

			for index, _ := range tc.projects {
				projects = append(projects, &tc.projects[index])
			}

			output := codebuild.ListProjectsOutput{
				Projects: projects,
			}

			client.
				EXPECT().
				ListProjects(gomock.Any()).
				Return(&output, tc.err).
				AnyTimes()

			var b bytes.Buffer
			writer := bufio.NewWriter(&b)

			cmd.DisplayProjects(client, writer)
			writer.Flush()

			if b.String() != tc.expected {
				t.Fatalf("expected '%s'; got '%s'", tc.expected, b.String())
			}
		})
	}
}
