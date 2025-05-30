<p align="center">
  <img src="https://github.com/coolapso/convcommitlint/blob/main/images/logo.png" width="200" >
</p>

# CONVCOMMITLINT
[![Release](https://github.com/coolapso/convcommitlint/actions/workflows/release.yaml/badge.svg?branch=main)](https://github.com/coolapso/convcommitlint/actions/workflows/release.yaml)
![GitHub Tag](https://img.shields.io/github/v/tag/coolapso/convcommitlint?logo=semver&label=semver&labelColor=gray&color=green)
[![Docker image version](https://img.shields.io/docker/v/coolapso/convcommitlint/latest?logo=docker)](https://hub.docker.com/r/coolapso/convcommitlint)
[![Go Report Card](https://goreportcard.com/badge/github.com/coolapso/convcommitlint)](https://goreportcard.com/report/github.com/coolapso/convcommitlint)
![GitHub Sponsors](https://img.shields.io/github/sponsors/coolapso?style=flat&logo=githubsponsors)

A simple, slightly opinionated, yet usable linter for [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/), written in Go.
This linter covers the convention essentials, contributions are wellcome but, no! I don't intend to support all kinds of variations and custom rules! 

## Motivation

I just wanted something simple that works, not configure all kinds of things just to get a linter working!

## Features

* Checks for header syntax
* Checks for common typos on the most important keywords: fix, feat and BREAKING CHANGE
* Supports creating pull github pull request reviews
    * Supports request for changes (default)
    * Supports commenting only
* Supports linting only the current commit
* Supports linting all commits
* Supports linting only most recent commits from base branch

## How it works

convcommitlint will check your commits against the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) standard, and let you know if there are any issues, when running on github actions and if you enable creation of pull request reviews it will create a pull request review for you. 

## Installation 

### Github action

Github action is available to use with your CI pipelines and most arguments from the CLI application are available.

```
  convcommitlint:
    runs-on: ubuntu-latest
    steps:
      - uses: coolapso/convcommitlint@v0
```

You can find the available inputs in the [action.yaml](action.yaml) file, if a input is not available you still should be able to use it by setting the associated ENV variable.

### Docker

You can run convcommitlint with docker, containers are provided both on dockerhub and ghcr.io 

**ghcr.io**

```
docker run -v $(pwd):/data --rm ghcr.io/coolapso/convcommitlint:latest
```

**dockerhub:**

```
docker run -v $(pwd):/data --rm coolapso/convcommitlint:latest
```

### Go Install

#### Latest version 

`go install github.com/coolapso/convcommitlint@latest`

#### Specific version

`go install github.com/coolapso/convcommitlint@v1.0.0`

### AUR

On Arch linux you can use the AUR `convcommitlint-bin`

### Linux Script

It is also possible to install on any linux distro with the installation script

#### Latest version

```
curl -L http://commitlint.coolapso.sh/install.sh | bash
```

#### Specific version

```
curl -L http://commitlint.coolapso.sh/install.sh | VERSION="v1.1.0" bash
```

### Manual install

* Grab the binary from the [releases page](https://github.com/coolapso/convcommitlint/releases).
* Extract the binary
* Execute it


## Usage

```
Usage:
  convcommitlint [flags]
  convcommitlint [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print convcommitlint version

Flags:
  -b, --base-branch string   The base branch to check commits from (default "main")
      --comment-only         Pull request reviews will only comment instead of requesting changes
  -r, --create-review        Creates review on github pull request
  -c, --current              Lint only the current commit
  -h, --help                 help for convcommitlint
  -a, --lint-all             Lint all repository commits
  -p, --path string          Git repository path (default "./")
      --pr-number int        The number of pull request to create the review
      --repository string    The github repository in owner/name format ex: coolapso/convcommitlint

Use "convcommitlint [command] --help" for more information about a command.
```

## Build 

`go build -o convcommitlint

# Contributions

Improvements and suggestions are always welcome, feel free to check for any open issues, open a new Issue or Pull Request

If you like this project and want to support / contribute in a different way you can always [:heart: Sponsor Me](https://github.com/sponsors/coolapso) or

<a href="https://www.buymeacoffee.com/coolapso" target="_blank">
  <img src="https://cdn.buymeacoffee.com/buttons/default-yellow.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;" />
</a>



Also consider supporting [tapio/live-server](https://github.com/tapio/live-server) which inspired this project

