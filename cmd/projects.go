package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/spf13/cobra"
)

// NewListProjectsCommand creates a new `projects` command
func NewListProjectsCommand(svc *codebuild.CodeBuild) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List all the projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return DisplayProjects(svc, os.Stdout)
		},
	}
	return cmd
}

// DisplayProjects will render the projects you have access to
func DisplayProjects(svc *codebuild.CodeBuild, w io.Writer) error {
	projects, err := svc.ListProjects(&codebuild.ListProjectsInput{SortOrder: aws.String("ASCENDING")})
	if err != nil {
		return err
	}

	var sorted []string
	for _, name := range projects.Projects {
		sorted = append(sorted, *name)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	for _, name := range sorted {
		fmt.Fprintf(w, "%s\n", name)
	}

	return nil
}
