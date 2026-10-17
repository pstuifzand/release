# Release

Release creates a new version, updates the changelog and pushes the code to the repository. Version-only commits include `[CI SKIP]`, so Drone skips the redundant branch build; the subsequent tag build still runs.

## Installation

When installed you have the `release` command.

```shell script
go get -u github.com/pstuifzand/release/cmd/release
```
