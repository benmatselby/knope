package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/benmatselby/knope/client"
	"github.com/benmatselby/knope/ui"

	"github.com/spf13/cobra"
)

// OverviewOptions defines what arguments/options the user can provide
type OverviewOptions struct {
	Args   []string
	Filter string
}

// BuildRecord gives us a struct to store records
type BuildRecord struct {
	Project string
	Status  string
	Start   string
	Finish  string
}

// NewOverviewCommand creates a new `overview` command
func NewOverviewCommand(client client.API) *cobra.Command {
	var opts OverviewOptions

	cmd := &cobra.Command{
		Use:   "overview",
		Short: "Will provide an overview of the last build per project",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return DisplayOverview(client, opts, os.Stdout)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Filter, "filter", ".*", "Regex to filter the projects displayed")

	return cmd
}

// DisplayOverview will render each project asked for and the last build value
func DisplayOverview(client client.API, opts OverviewOptions, w io.Writer) error {
	projects, err := client.ListProjects(&codebuild.ListProjectsInput{SortOrder: aws.String("ASCENDING")})
	if err != nil {
		return err
	}

	records := make(chan BuildRecord)
	var wg sync.WaitGroup
	wg.Add(len(projects.Projects))

	go func() {
		wg.Wait()
		close(records)
	}()

	for _, project := range projects.Projects {
		go func(project *string) {
			defer wg.Done()

			matched, err := regexp.MatchString(opts.Filter, *project)
			if !matched {
				return
			}

			projectBuilds, err := client.ListBuildsForProject(&codebuild.ListBuildsForProjectInput{
				ProjectName: project,
			})
			if err != nil {
				records <- BuildRecord{Project: *project, Status: ui.AppUnknown, Start: "", Finish: ""}
				return
			}

			if len(projectBuilds.Ids) == 0 {
				records <- BuildRecord{Project: *project, Status: ui.AppEmpty, Start: "", Finish: ""}
				return
			}

			builds, err := client.BatchGetBuilds(&codebuild.BatchGetBuildsInput{Ids: projectBuilds.Ids})
			if err != nil {
				records <- BuildRecord{Project: *project, Status: ui.AppUnknown, Start: "", Finish: ""}
				return
			}

			build := builds.Builds[0]
			start := build.StartTime.Format(ui.AppDateTimeFormat)

			finish := ""

			if build.EndTime != nil {
				finish = build.EndTime.Format(ui.AppDateTimeFormat)
			}

			records <- BuildRecord{
				Project: *project,
				Status:  getBuildIcon(build.BuildStatus),
				Start:   start,
				Finish:  finish,
			}
		}(project)
	}

	var builds []BuildRecord
	for r := range records {
		builds = append(builds, r)
	}

	sort.Slice(builds, func(i, j int) bool { return builds[i].Project < builds[j].Project })

	tr := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(tr, "%s \t%s\t%s\t%s\n", "Status", "Name", "Branch", "Finished")

	for _, build := range builds {
		fmt.Fprintf(tr, "%s \t%s\t%s\t%s\n", build.Status, build.Project, build.Start, build.Finish)
	}

	tr.Flush()

	return nil
}

func getBuildIcon(status *string) string {
	result := ""
	if aws.StringValue(status) == "FAILED" || aws.StringValue(status) == "FAULT" {
		result = ui.AppFailure
	} else if aws.StringValue(status) == "IN_PROGRESS" {
		result = ui.AppProgress
	} else if aws.StringValue(status) == "STOPPED" || aws.StringValue(status) == "TIMED_OUT" {
		result = ui.AppStale
	} else {
		result = ui.AppSuccess
	}

	return result
}
