# Knope

![GitHub Badge](https://github.com/benmatselby/knope/workflows/Go/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=knope&metric=alert_status)](https://sonarcloud.io/dashboard?id=knope)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/knope)](https://goreportcard.com/report/github.com/benmatselby/knope)

_I am super chill all the time!_

```text
CLI tool for retrieving data from AWS CodeBuild

Usage:
  knope [command]

Available Commands:
  builds      List all the builds for a given project
  help        Help about any command
  overview    Will provide an overview of the last build per project
  projects    List all the projects

Flags:
      --config string   config file (default is $HOME/.benmatselby/knope.yaml)
  -h, --help            help for knope

Use "knope [command] --help" for more information about a command.
```

## Requirements

If you are wanting to build and develop this, you will need the following items installed.

- Go version 1.12+

## Configuration

You will need the following environment variables defining:

```shell
export AWS_DEFAULT_REGION=""
export AWS_PROFILE=""
```

## Installation via Git

```shell
git clone git@github.com:benmatselby/knope.git
cd knope
make all
./knope
```

You can also install into your `$GOPATH/bin` by running `make build && go install`.

## Testing

To generate the code used to mock away the CodeBuild interaction, run the following command.

```shell
mockgen -source client/client.go
```

This will generate you some source code you can copy into `client/mock_client.go`. You will need to change the package to `client`.
