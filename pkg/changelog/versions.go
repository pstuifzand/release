package changelog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/blang/semver"
)

func AddNewVersion(filename string, prevVersion, version semver.Version, date string) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile(".", "release-changelog")
	if err != nil {
		return err
	}

	copyWithNewVersion(body, f, prevVersion, version, date)

	err = os.Rename(f.Name(), filename)
	return err
}

func copyWithNewVersion(body []byte, f io.Writer, prevVersion, version semver.Version, date string) error {
	const HeaderUnreleased = "## [Unreleased]"
	const LinkUnreleased = "[Unreleased]: "

	unreleasedStart := bytes.Index(body, []byte(HeaderUnreleased))
	if unreleasedStart < 0 {
		return errors.New("unreleased level 2 header not found")
	}

	prefix := body[0 : unreleasedStart+len(HeaderUnreleased)]
	suffix := body[unreleasedStart+len(HeaderUnreleased):]

	out := bufio.NewWriter(f)

	linkStart := bytes.Index(suffix, []byte(LinkUnreleased))
	if linkStart < 0 {
		return errors.New("unreleased link start not found")
	}
	linkEnd := bytes.Index(suffix[linkStart:], []byte("\n"))
	if linkEnd < 0 {
		return errors.New("unreleased link end not found")
	}
	linkEnd += linkStart // adjust for linkStart

	var newLink string
	firstVersion := prevVersion.String() == "0.0.0"
	if firstVersion {
		newLink = fmt.Sprintf("[%s]: https://github.com/pstuifzand/release/tag/%s\n", version, version)
	} else {
		newLink = fmt.Sprintf("[%s]: https://github.com/pstuifzand/release/compare/%s...%s\n", version, prevVersion, version)
	}

	parts := []string{
		string(prefix),
		"\n\n## [",
		version.String(),
		"] - ",
		date,
		string(suffix[0:linkStart]),
		LinkUnreleased, "https://github.com/pstuifzand/release/compare/", version.String(), "...HEAD\n",
		// are there more versions?
		newLink,
		string(suffix[linkEnd+1:]),
	}

	for _, s := range parts {
		if _, err := out.WriteString(s); err != nil {
			return err
		}
	}

	if err := out.Flush(); err != nil {
		return err
	}

	return nil
}
