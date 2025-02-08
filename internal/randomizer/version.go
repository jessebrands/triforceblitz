package randomizer

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	versionRegexp = regexp.MustCompile(
		`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)-([a-z][a-z0-9]*)-(0|[1-9]\d*)\.(0|[1-9]\d*)$`,
	)
)

type Version struct {
	Major       int
	Minor       int
	Patch       int
	Branch      string
	BranchMajor int
	BranchMinor int
}

func ValidVersion(s string) bool {
	return versionRegexp.MatchString(s)
}

func VersionFromString(s string) (Version, error) {
	if !ValidVersion(s) {
		return Version{}, ErrInvalidVersion
	}
	matches := versionRegexp.FindStringSubmatch(s)
	var major, minor, patch int
	var branchMajor, branchMinor int
	major, _ = strconv.Atoi(matches[1])
	minor, _ = strconv.Atoi(matches[2])
	patch, _ = strconv.Atoi(matches[3])
	branchMajor, _ = strconv.Atoi(matches[5])
	branchMinor, _ = strconv.Atoi(matches[6])
	return Version{
		Major:       major,
		Minor:       minor,
		Patch:       patch,
		Branch:      matches[4],
		BranchMajor: branchMajor,
		BranchMinor: branchMinor,
	}, nil
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s-%d.%d",
		v.Major,
		v.Minor,
		v.Patch,
		v.Branch,
		v.BranchMajor,
		v.BranchMinor,
	)
}

func (v *Version) Set(s string) error {
	version, err := VersionFromString(s)
	if err != nil {
		return err
	}
	*v = version
	return nil
}

func (v *Version) Equal(o Version) bool {
	return (v.Major == o.Major &&
		v.Minor == o.Minor &&
		v.Patch == o.Patch &&
		v.Branch == o.Branch &&
		v.BranchMajor == o.BranchMajor &&
		v.BranchMinor == o.BranchMinor)
}
