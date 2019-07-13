# Knope

[![Build Status](https://travis-ci.com/benmatselby/knope.svg?branch=master)](https://travis-ci.com/benmatselby/knope)
[![codecov](https://codecov.io/gh/benmatselby/knope/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/knope)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/knope)](https://goreportcard.com/report/github.com/benmatselby/knope)

_I am super chill all the time!_

```text
CLI tool for retrieving data from AWS CodeBuild

Usage:
  knope [command]

Available Commands:
  builds      List all the builds for a given project
  help        Help about any command
  projects    List all the projects

Flags:
      --config string   config file (default is $HOME/.benmatselby/knope.yaml)
  -h, --help            help for knope

Use "knope [command] --help" for more information about a command.
```

## Requirements

If you are wanting to build and develop this, you will need the following items installed. If, however, you just want to run the application I recommend using the docker container (See below).

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
