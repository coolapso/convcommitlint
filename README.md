<p align="center">
  <img src="https://github.com/coolapso/convcommitlint/blob/main/images/logo.png" width="200" >
</p>

# CONVCOMMITLINT

[![Release](https://github.com/coolapso/convcommitlint/actions/workflows/release.yaml/badge.svg?branch=main)](https://github.com/coolapso/convcommitlint/actions/workflows/release.yaml)
![GitHub Tag](https://img.shields.io/github/v/tag/coolapso/convcommitlint?logo=semver&label=semver&labelColor=gray&color=green)
[![Docker image version](https://img.shields.io/docker/v/coolapso/convcommitlint/latest?logo=docker)](https://hub.docker.com/r/coolapso/convcommitlint)
[![Go Report Card](https://goreportcard.com/badge/github.com/coolapso/convcommitlint)](https://goreportcard.com/report/github.com/coolapso/convcommitlint)
![GitHub Sponsors](https://img.shields.io/github/sponsors/coolapso?style=flat&logo=githubsponsors)

A simple, slightly opinionated, but actually *usable* linter for [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/), written in Go.

---

## Motivation

I just wanted something that *works* out of the box. No over-configuration, no headaches—just simple commit linting! This being said, this linter covers the essentials. Contributions are welcome, but I don’t intend to support every possible variation or custom rule!

---

## Features

- **Checks header syntax**
- **Detects common typos** in key keywords: `fix`, `feat`, and `BREAKING CHANGE`
- **GitHub Pull Request Reviews:**
  - Request changes (default)
  - Comment-only mode
- **Lint Modes:**
  - Lint only the current commit
  - Lint all commits
  - Lint only recent commits from the base branch
- **Full environment variable support** for all flags
- **Cross-platform:** Linux, macOS, and Windows
- **GitHub Action** support

---

## How It Works

`convcommitlint` checks your commits against the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard.  
If you use it as a GitHub Action or enable PR review, it will comment or request changes directly on your pull requests, listing any issues it finds.

---

## Usage

```
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
<p align="center">
  <img src="https://raw.githubusercontent.com/coolapso/convcommitlint/refs/heads/main/images/demo.gif">
</p>


### Environment Variables

Every flag can also be set with an environment variable, using the `CONVCOMMITLINT_` prefix, uppercase, and underscores.  
For example, the flag `--lint-all` becomes the variable `CONVCOMMITLINT_LINT_ALL`.

---

## GitHub Action

You can use `convcommitlint` as part of your CI pipeline. Most CLI arguments are supported as action inputs.

```yaml
convcommitlint:
    runs-on: ubuntu-latest
    steps: 
      - uses: actions/checkout@v4
      - uses: coolapso/convcommitlint@v0.3.0
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Base Branch

By default, the base branch is `"main"`. If you use a different default branch, set the `base-branch` input:

```yaml
convcommitlint:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: coolapso/convcommitlint@v0.3.0
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with: 
        base-branch: develop
```

### Git History Depth

By default, the action checks out the reference branch with full history. To limit this, use the `fetch-depth` input (commits beyond this depth will *not* be analyzed):

```yaml
convcommitlint:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: coolapso/convcommitlint@v0.3.0
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with: 
        fetch-depth: 10
```

### Disable Pull Request Reviews

PR reviews are enabled by default. To disable them, set `create-review` to `"false"`:

```yaml
convcommitlint:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: coolapso/convcommitlint@v0.3.0
      with: 
        create-review: "false"
```

### Use a Specific Version

The action uses the latest release by default. To pin a specific version:

```yaml
convcommitlint:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: coolapso/convcommitlint@v0.3.0
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with: 
        version: v0.2.0
```

### Other Features & Flags

The action supports all CLI features and flags.  
Check available inputs in [`action.yaml`](action.yaml).  
If an input is missing, you can always set the corresponding environment variable.

---

## Installation

### Docker

Images are available on both DockerHub and GitHub Container Registry (ghcr.io):

**GitHub Container Registry:**

```sh
docker run -v $(pwd):/data --rm ghcr.io/coolapso/convcommitlint:latest
```

**DockerHub:**

```sh
docker run -v $(pwd):/data --rm coolapso/convcommitlint:latest
```

### Go Install

#### Latest Version

```sh
go install github.com/coolapso/convcommitlint@latest
```

#### Specific Version

```sh
go install github.com/coolapso/convcommitlint@v1.0.0
```

### Arch Linux (AUR)

On Arch Linux, use the AUR package: `convcommitlint-bin`.

### Linux Install Script

Install on any Linux distro using the script:

#### Latest Version

```sh
curl -L http://commitlint.coolapso.sh/install.sh | bash
```

#### Specific Version

```sh
curl -L http://commitlint.coolapso.sh/install.sh | VERSION="v1.1.0" bash
```


### Manual Install

- Download the binary from the [releases page](https://github.com/coolapso/convcommitlint/releases)
- Extract the binary
- Run it

### Build

```sh
go build -o convcommitlint
```

---

## Contributions

Improvements and suggestions are always welcome!  
Check open issues, or open a new Issue or Pull Request.

If you like this project and want to support or contribute in another way, you can [:heart: Sponsor Me](https://github.com/sponsors/coolapso) or:

<a href="https://www.buymeacoffee.com/coolapso" target="_blank">
  <img src="https://cdn.buymeacoffee.com/buttons/default-yellow.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;" />
</a>
