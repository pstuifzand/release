package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver"

	"github.com/pstuifzand/release/pkg/changelog"
)

var dryRun = flag.Bool("dry-run", false, "dry run")
var shouldPatch = flag.Bool("patch", true, "increment patch")
var shouldMinor = flag.Bool("minor", false, "increment minor")
var shouldMajor = flag.Bool("major", false, "increment major")

func main() {
	flag.Parse()

	versions, err := currentVersions()
	if err != nil {
		log.Fatalf("currentVersions: %s", err)
	}

	latest, err := latestVersion(versions)
	if err != nil {
		log.Fatalf("latestVersion: %s", err)
	}

	nextVersion := latest
	if *shouldMajor {
		nextVersion.Major++
		nextVersion.Minor = 0
		nextVersion.Patch = 0
	} else if *shouldMinor {
		nextVersion.Minor++
		nextVersion.Patch = 0
	} else if *shouldPatch {
		nextVersion.Patch++
	}

	fmt.Printf("Current version: %q\n", latest)
	fmt.Printf("Next version: %q\n", nextVersion)

	if *dryRun {
		fmt.Printf("Update changelog\n")
	} else {
		err = changelog.AddNewVersion("CHANGELOG.md", latest, nextVersion, time.Now().In(time.Local).Format("2006-01-02"))
		if err != nil {
			log.Fatal(err)
		}
	}

	var cmds []*exec.Cmd

	cmds = append(cmds,
		gitAdd("CHANGELOG.md"),
		gitCommit(nextVersion),
		gitTag(nextVersion),
		gitPush(),
		gitPushTags(),
	)

	runCommands(cmds)
}

func runCommands(cmds []*exec.Cmd) {
	for _, cmd := range cmds {
		if *dryRun {
			err := outputCommand(cmd)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := outputCommand(cmd)
			if err != nil {
				log.Fatal(err)
			}
			err = execCommand(cmd)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func execCommand(cmd *exec.Cmd) error {
	return cmd.Run()
}

func outputCommand(cmd *exec.Cmd) error {
	_, err := fmt.Printf("Running command: %s\n", strings.Join(cmd.Args, " "))
	return err
}

func gitPushTags() *exec.Cmd {
	return exec.Command("git", "push", "--tags")
}

func gitAdd(filename string) *exec.Cmd {
	return exec.Command("git", "add", filename)
}

func gitCommit(version semver.Version) *exec.Cmd {
	return exec.Command("git", "commit", "-n", "-m", fmt.Sprintf("Increase version to %s", version))
}

func gitTag(version semver.Version) *exec.Cmd {
	return exec.Command("git", "tag", version.String())
}

func gitPush() *exec.Cmd {
	return exec.Command("git", "push")
}

func currentVersions() ([]semver.Version, error) {
	var buf bytes.Buffer

	cmd := exec.Command("git", "tag", "-l")
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(&buf)
	scanner.Split(bufio.ScanLines)

	var versions []semver.Version

	for scanner.Scan() {
		current, err := semver.Make(scanner.Text())
		if err != nil {
			log.Printf("Parse failed of %s", scanner.Text())
			continue
		}

		versions = append(versions, current)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[j].LT(versions[i])
	})

	return versions, nil
}

// latestVersion returns the maximum version from versions
func latestVersion(versions []semver.Version) (semver.Version, error) {
	latest, _ := semver.Make("0.0.0")

	for _, current := range versions {
		if current.GT(latest) {
			latest = current
		}
	}

	return latest, nil
}
