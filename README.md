# Knope

[![Build Status](https://travis-ci.org/benmatselby/knope.png?branch=master)](https://travis-ci.org/benmatselby/knope)
[![codecov](https://codecov.io/gh/benmatselby/knope/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/knope)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/knope)](https://goreportcard.com/report/github.com/benmatselby/knope)

_I am super chill all the time!_

CLI application for getting build information out of AWS CodeBuild.

```text

```

## Requirements

If you are wanting to build and develop this, you will need the following items installed. If, however, you just want to run the application I recommend using the docker container (See below).

- Go version 1.12+

## Configuration

### Environment variables

You will need the following environment variables defining:

```shell
export AWS_DEFAULT_REGION=""
```

### Application configuration file

Long term this may not be required, but right now we need a configuration file (by default, `~/.benmatselby/knope.yaml`).

An example:

```yml
```

## Installation via Git

```shell
git clone git@github.com:benmatselby/knope.git
cd knope
make all
./knope
```

You can also install into your `$GOPATH/bin` by running `make build && go install`.
