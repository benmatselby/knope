package cmd_test

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/benmatselby/knope/client"
	"github.com/benmatselby/knope/cmd"
	"github.com/golang/mock/gomock"
)

func TestNewOverviewCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := client.NewMockAPI(ctrl)

	cmd := cmd.NewOverviewCommand(client)

	use := "overview"
	short := "Will provide an overview of the last build per project"

	if cmd.Use != use {
		t.Fatalf("expected use: %s; got %s", use, cmd.Use)
	}

	if cmd.Short != short {
		t.Fatalf("expected use: %s; got %s", short, cmd.Short)
	}
}

type testOverviewBuild struct {
	Status string
	Start  time.Time
	Finish time.Time
}

func TestDisplayOverview(t *testing.T) {
	tt := []struct {
		name     string
		projects []string
		builds   []testOverviewBuild
		expected string
		err      error
	}{
		{name: "can return a build per project", projects: []string{"a"}, builds: []testOverviewBuild{testOverviewBuild{
			Status: "SUCCEEDED",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name Branch           Finished
✅       a    19-07-2019 23:00 19-07-2019 23:10
`, err: nil},
		{name: "can order the projects", projects: []string{"a", "d", "c"}, builds: []testOverviewBuild{testOverviewBuild{
			Status: "SUCCEEDED",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name Branch           Finished
✅       a    19-07-2019 23:00 19-07-2019 23:10
✅       c    19-07-2019 23:00 19-07-2019 23:10
✅       d    19-07-2019 23:00 19-07-2019 23:10
`, err: nil},
		{name: "can return a failed build per project", projects: []string{"a"}, builds: []testOverviewBuild{testOverviewBuild{
			Status: "FAILED",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name Branch           Finished
❌       a    19-07-2019 23:00 19-07-2019 23:10
`, err: nil},
		{name: "can return a faulted build per project", projects: []string{"a"}, builds: []testOverviewBuild{testOverviewBuild{
			Status: "FAILED",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name Branch           Finished
❌       a    19-07-2019 23:00 19-07-2019 23:10
`, err: nil},
		{name: "can return an in progress build per project", projects: []string{"a"}, builds: []testOverviewBuild{testOverviewBuild{
			Status: "IN_PROGRESS",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name Branch           Finished
🏗       a    19-07-2019 23:00 19-07-2019 23:10
`, err: nil},
		{name: "can return a stopped build per project", projects: []string{"a"}, builds: []testOverviewBuild{testOverviewBuild{
			Status: "STOPPED",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name Branch           Finished
🕳       a    19-07-2019 23:00 19-07-2019 23:10
`, err: nil},
		{name: "can return a timed out build per project", projects: []string{"a"}, builds: []testOverviewBuild{testOverviewBuild{
			Status: "TIMED_OUT",
			Start:  time.Date(2019, time.July, 19, 23, 0, 0, 0, time.UTC),
			Finish: time.Date(2019, time.July, 19, 23, 10, 0, 0, time.UTC),
		}},
			expected: `Status  Name Branch           Finished
🕳       a    19-07-2019 23:00 19-07-2019 23:10
`, err: nil},
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
			projectOutput := codebuild.ListProjectsOutput{
				Projects: projects,
			}

			buildProjectOutput := codebuild.ListBuildsForProjectOutput{
				Ids: projects,
			}

			var builds []*codebuild.Build
			for index, _ := range tc.builds {
				builds = append(builds, &codebuild.Build{
					StartTime:   &tc.builds[index].Start,
					EndTime:     &tc.builds[index].Finish,
					BuildStatus: &tc.builds[index].Status,
				})
			}
			buildOutput := codebuild.BatchGetBuildsOutput{
				Builds: builds,
			}

			client.
				EXPECT().
				ListProjects(gomock.Any()).
				Return(&projectOutput, tc.err).
				AnyTimes()

			client.
				EXPECT().
				ListBuildsForProject(gomock.Any()).
				Return(&buildProjectOutput, tc.err).
				AnyTimes()

			client.
				EXPECT().
				BatchGetBuilds(gomock.Any()).
				Return(&buildOutput, tc.err).
				AnyTimes()

			var b bytes.Buffer
			writer := bufio.NewWriter(&b)

			cmd.DisplayOverview(client, writer)
			writer.Flush()

			if b.String() != tc.expected {
				t.Fatalf("expected '%s'; got '%s'", tc.expected, b.String())
			}
		})
	}
}
