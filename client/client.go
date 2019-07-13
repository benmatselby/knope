package client

import (
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/spf13/viper"
)

// API defines the client interface
type API interface {
	BatchGetBuilds(input *codebuild.BatchGetBuildsInput) (*codebuild.BatchGetBuildsOutput, error)
	ListBuildsForProject(input *codebuild.ListBuildsForProjectInput) (*codebuild.ListBuildsForProjectOutput, error)
	ListProjects(input *codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error)
}

// Client is the content implementation of the API we are using in the app
type Client struct {
	codebuild *codebuild.CodeBuild
}

// NewClient will return a internal codebuild client.
func NewClient() Client {
	sess, _ := session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 viper.GetString("AWS_PROFILE"),
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	})

	svc := codebuild.New(sess)

	client := Client{
		codebuild: svc,
	}

	return client
}

// BatchGetBuilds will call the same function on the codebuild client
func (c *Client) BatchGetBuilds(input *codebuild.BatchGetBuildsInput) (*codebuild.BatchGetBuildsOutput, error) {
	return c.codebuild.BatchGetBuilds(input)
}

// ListBuildsForProject will call the same function on the codebuild client
func (c *Client) ListBuildsForProject(input *codebuild.ListBuildsForProjectInput) (*codebuild.ListBuildsForProjectOutput, error) {
	return c.codebuild.ListBuildsForProject(input)
}

// ListProjects will call the same function on the codebuild client
func (c *Client) ListProjects(input *codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error) {
	return c.codebuild.ListProjects(input)
}
