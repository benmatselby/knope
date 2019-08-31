package cmd_test

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/benmatselby/knope/client"
	"github.com/benmatselby/knope/cmd"
	"github.com/golang/mock/gomock"
)

func TestNewListBuildsForProjectCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := client.NewMockAPI(ctrl)

	cmd := cmd.NewListBuildsForProjectCommand(client)

	use := "builds"
	short := "List all the builds for a given project"

	if cmd.Use != use {
		t.Fatalf("expected use: %s; got %s", use, cmd.Use)
	}

	if cmd.Short != short {
		t.Fatalf("expected use: %s; got %s", short, cmd.Short)
	}
}

type testBuild struct {
	Status string
	Source string
	Start  time.Time
	Finish time.Time
}

func TestDisplayBuildsForProject(t *testing.T) {
	tt := []struct {
		name         string
		project      string
		builds       []testBuild
		expected     string
		listBuildErr error
		getBuildErr  error
	}{
		{name: "can return a succeeded build", project: "project-one", builds: []testBuild{testBuild{
			Status: "SUCCEEDED",
			Source: "my-branch",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name      Branch           Finished
✅       my-branch 19-07-2019 23:00 19-07-2019 23:10
`, listBuildErr: nil, getBuildErr: nil},
		{name: "can return a failed build", project: "project-one", builds: []testBuild{testBuild{
			Status: "FAILED",
			Source: "my-branch",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name      Branch           Finished
❌       my-branch 19-07-2019 23:00 19-07-2019 23:10
`, listBuildErr: nil, getBuildErr: nil},
		{name: "can return a fault build", project: "project-one", builds: []testBuild{testBuild{
			Status: "FAULT",
			Source: "my-branch",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name      Branch           Finished
❌       my-branch 19-07-2019 23:00 19-07-2019 23:10
`, listBuildErr: nil, getBuildErr: nil},
		{name: "can return an in progress build", project: "project-one", builds: []testBuild{testBuild{
			Status: "IN_PROGRESS",
			Source: "my-branch",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name      Branch           Finished
🏗       my-branch 19-07-2019 23:00 19-07-2019 23:10
`, listBuildErr: nil, getBuildErr: nil},
		{name: "can return a stopped build build", project: "project-one", builds: []testBuild{testBuild{
			Status: "STOPPED",
			Source: "my-branch",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name      Branch           Finished
🕳       my-branch 19-07-2019 23:00 19-07-2019 23:10
`, listBuildErr: nil, getBuildErr: nil},
		{name: "can return a timed out build build", project: "project-one", builds: []testBuild{testBuild{
			Status: "STOPPED",
			Source: "my-branch",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name      Branch           Finished
🕳       my-branch 19-07-2019 23:00 19-07-2019 23:10
`, listBuildErr: nil, getBuildErr: nil},
		{
			name:         "unable to list builds for project",
			project:      "project-one",
			builds:       nil,
			expected:     "",
			listBuildErr: errors.New("there was an error"),
			getBuildErr:  nil,
		},
		{
			name:         "unable to get builds for project",
			project:      "project-one",
			builds:       nil,
			expected:     "",
			listBuildErr: nil,
			getBuildErr:  errors.New("there was an error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := client.NewMockAPI(ctrl)

			buildProjectOutput := codebuild.ListBuildsForProjectOutput{
				Ids: []*string{&tc.project},
			}

			var builds []*codebuild.Build
			for index, _ := range tc.builds {
				builds = append(builds, &codebuild.Build{
					StartTime:             &tc.builds[index].Start,
					EndTime:               &tc.builds[index].Finish,
					BuildStatus:           &tc.builds[index].Status,
					ResolvedSourceVersion: &tc.builds[index].Source,
				})
			}
			buildOutput := codebuild.BatchGetBuildsOutput{
				Builds: builds,
			}

			client.
				EXPECT().
				ListBuildsForProject(gomock.Any()).
				Return(&buildProjectOutput, tc.listBuildErr).
				AnyTimes()

			client.
				EXPECT().
				BatchGetBuilds(gomock.Any()).
				Return(&buildOutput, tc.getBuildErr).
				AnyTimes()

			var b bytes.Buffer
			writer := bufio.NewWriter(&b)

			opt := cmd.ListBuildForProjectOptions{
				Project: tc.project,
			}

			err := cmd.DisplayBuildsForProject(client, opt, writer)
			writer.Flush()

			if b.String() != tc.expected {
				t.Fatalf("expected '%s'; got '%s'", tc.expected, b.String())
			}

			if tc.listBuildErr != nil {
				if err != tc.listBuildErr {
					t.Fatalf("expected err to be %v; got %v", tc.listBuildErr, err)
				}
			}

			if tc.getBuildErr != nil {
				if err != tc.getBuildErr {
					t.Fatalf("expected err to be %v; got %v", tc.getBuildErr, err)
				}
			}
		})
	}
}
