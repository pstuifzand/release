package changelog

import (
	"bytes"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
)

func TestCopyWithNewVersionNewRelease(t *testing.T) {
	var buf bytes.Buffer

	err := copyWithNewVersion([]byte(`## [Unreleased]

### Added

- hello world

[Unreleased]: https://github.com/pstuifzand/release/tree/master
`), &buf, semver.MustParse("0.0.0"), semver.MustParse("1.0.0"), "2021-01-01")

	if assert.NoError(t, err) {
		assert.Equal(t, `## [Unreleased]

## [1.0.0] - 2021-01-01

### Added

- hello world

[Unreleased]: https://github.com/pstuifzand/release/compare/1.0.0...HEAD
[1.0.0]: https://github.com/pstuifzand/release/tag/1.0.0
`, buf.String())
	}
}

func TestCopyWithNewVersionNewReleaseAndOldRelease(t *testing.T) {
	var buf bytes.Buffer

	err := copyWithNewVersion([]byte(`## [Unreleased]

### Added

- test hello

## [1.0.0] - 2021-01-01

### Added

- hello world

[Unreleased]: https://github.com/pstuifzand/release/compare/1.0.0...HEAD
[1.0.0]: https://github.com/pstuifzand/release/tag/1.0.0
`), &buf, semver.MustParse("1.0.0"), semver.MustParse("1.1.0"), "2021-02-01")

	if assert.NoError(t, err) {
		assert.Equal(t, `## [Unreleased]

## [1.1.0] - 2021-02-01

### Added

- test hello

## [1.0.0] - 2021-01-01

### Added

- hello world

[Unreleased]: https://github.com/pstuifzand/release/compare/1.1.0...HEAD
[1.1.0]: https://github.com/pstuifzand/release/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/pstuifzand/release/tag/1.0.0
`, buf.String())
	}
}
