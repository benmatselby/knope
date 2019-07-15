package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/benmatselby/knope/client"
	"github.com/benmatselby/knope/ui"
	"github.com/spf13/cobra"
)

// ListBuildForProjectOptions defines what arguments/options the user can provide
type ListBuildForProjectOptions struct {
	Args    []string
	Project string
}

// NewListBuildsForProjectCommand creates a new `builds` command
func NewListBuildsForProjectCommand(client client.API) *cobra.Command {
	var opts ListBuildForProjectOptions

	cmd := &cobra.Command{
		Use:   "builds",
		Short: "List all the builds for a given project",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return DisplayBuildsForProject(client, opts, os.Stdout)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Project, "project", "", "Name of the project to list builds for")
	return cmd
}

// DisplayBuildsForProject will render the projects you have access to
func DisplayBuildsForProject(client client.API, opts ListBuildForProjectOptions, w io.Writer) error {
	if opts.Project == "" {
		return fmt.Errorf("please specify a project name")
	}

	projectBuilds, err := client.ListBuildsForProject(&codebuild.ListBuildsForProjectInput{
		ProjectName: &opts.Project,
	})
	if err != nil {
		return err
	}

	builds, err := client.BatchGetBuilds(&codebuild.BatchGetBuildsInput{Ids: projectBuilds.Ids})
	if err != nil {
		return err
	}

	tr := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(tr, "%s\t%s\t%s\t%s\n", "", "Name", "Branch", "Finished")
	for _, build := range builds.Builds {
		start := build.StartTime.Format(ui.AppDateTimeFormat)

		finish := ""

		if build.EndTime != nil {
			finish = build.EndTime.Format(ui.AppDateTimeFormat)
		}

		result := ""
		if aws.StringValue(build.BuildStatus) == "FAILED" || aws.StringValue(build.BuildStatus) == "FAULT" {
			result = ui.AppFailure
		} else if aws.StringValue(build.BuildStatus) == "IN_PROGRESS" {
			result = ui.AppProgress
		} else if aws.StringValue(build.BuildStatus) == "STOPPED" || aws.StringValue(build.BuildStatus) == "TIMED_OUT" {
			result = ui.AppStale
		} else {
			result = ui.AppSuccess
		}

		fmt.Fprintf(tr, "%s \t%s\t%s\t%s\n", result, aws.StringValue(build.ResolvedSourceVersion), start, finish)
	}
	tr.Flush()

	return nil
}
