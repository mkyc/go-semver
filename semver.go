package semver

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// SemVer represents a semantic version as defined by the semantic versioning specification.
// It consists of major, minor, and patch version numbers, with optional pre-release and build metadata.
type SemVer struct {
	Major      uint
	Minor      uint
	Patch      uint
	PreRelease string
	Build      string
}

// String returns the string representation of the SemVer struct according to the semantic versioning specification.
func (s SemVer) String() string {
	// Start with the version core (major.minor.patch)
	result := fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)

	// Add pre-release information if present
	if s.PreRelease != "" {
		result += "-" + s.PreRelease
	}

	// Add build metadata if present
	if s.Build != "" {
		result += "+" + s.Build
	}

	return result
}

// IsRelease returns true if the semantic version represents a release version.
// A release version is one that doesn't have a pre-release identifier.
func (s SemVer) IsRelease() bool {
	return s.PreRelease == ""
}

// Parse parses a string tag into a SemVer struct according to the semantic versioning specification.
// It returns an error if the tag does not conform to the semantic versioning format.
func Parse(tag string) (SemVer, error) {
	var semver SemVer

	// Split the tag into version core and optional parts (pre-release and build)
	versionAndMeta := strings.SplitN(tag, "+", 2)
	versionPart := versionAndMeta[0]

	// Check if there's build metadata
	if len(versionAndMeta) > 1 {
		semver.Build = versionAndMeta[1]
	}

	// Split version part into version core and pre-release
	versionAndPreRelease := strings.SplitN(versionPart, "-", 2)
	versionCore := versionAndPreRelease[0]

	// Check if there's pre-release information
	if len(versionAndPreRelease) > 1 {
		semver.PreRelease = versionAndPreRelease[1]
	}

	// Parse version core (major.minor.patch)
	versionParts := strings.Split(versionCore, ".")
	if len(versionParts) != 3 {
		return SemVer{}, fmt.Errorf("invalid version format: %s, expected major.minor.patch", versionCore)
	}

	// Parse major version
	major, err := strconv.ParseUint(versionParts[0], 10, 0)
	if err != nil {
		return SemVer{}, fmt.Errorf("invalid major version: %s", versionParts[0])
	}
	semver.Major = uint(major)

	// Parse minor version
	minor, err := strconv.ParseUint(versionParts[1], 10, 0)
	if err != nil {
		return SemVer{}, fmt.Errorf("invalid minor version: %s", versionParts[1])
	}
	semver.Minor = uint(minor)

	// Parse patch version
	patch, err := strconv.ParseUint(versionParts[2], 10, 0)
	if err != nil {
		return SemVer{}, fmt.Errorf("invalid patch version: %s", versionParts[2])
	}
	semver.Patch = uint(patch)

	// Validate numeric identifiers according to the spec
	if versionParts[0] != "0" && strings.HasPrefix(versionParts[0], "0") {
		return SemVer{}, fmt.Errorf("invalid major version: %s, leading zeros not allowed", versionParts[0])
	}
	if versionParts[1] != "0" && strings.HasPrefix(versionParts[1], "0") {
		return SemVer{}, fmt.Errorf("invalid minor version: %s, leading zeros not allowed", versionParts[1])
	}
	if versionParts[2] != "0" && strings.HasPrefix(versionParts[2], "0") {
		return SemVer{}, fmt.Errorf("invalid patch version: %s, leading zeros not allowed", versionParts[2])
	}

	// Validate pre-release format if present
	if semver.PreRelease != "" {
		preReleaseParts := strings.Split(semver.PreRelease, ".")
		for _, part := range preReleaseParts {
			if part == "" {
				return SemVer{}, fmt.Errorf("invalid pre-release: empty identifier")
			}

			// Check if it's a numeric identifier
			if _, err := strconv.ParseUint(part, 10, 64); err == nil {
				// Numeric identifiers must not have leading zeros unless they are zero
				if part != "0" && strings.HasPrefix(part, "0") {
					return SemVer{}, fmt.Errorf("invalid pre-release: %s, numeric identifiers must not have leading zeros", part)
				}
			} else {
				// Alphanumeric identifiers must only contain alphanumeric characters and hyphens
				for _, c := range part {
					if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-') {
						return SemVer{}, fmt.Errorf("invalid pre-release: %s, contains invalid character", part)
					}
				}
			}
		}
	}

	// Validate build metadata format if present
	if semver.Build != "" {
		buildParts := strings.Split(semver.Build, ".")
		for _, part := range buildParts {
			if part == "" {
				return SemVer{}, fmt.Errorf("invalid build metadata: empty identifier")
			}

			// Build identifiers must only contain alphanumeric characters and hyphens
			for _, c := range part {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-') {
					return SemVer{}, fmt.Errorf("invalid build metadata: %s, contains invalid character", part)
				}
			}
		}
	}

	return semver, nil
}

// Compare compares this version with another version according to semantic versioning precedence rules.
// It returns:
//
//	-1 if this version has lower precedence than the other
//	 0 if this version has equal precedence to the other
//	 1 if this version has higher precedence than the other
func (s SemVer) Compare(other SemVer) int {
	// Compare major version
	if s.Major < other.Major {
		return -1
	}
	if s.Major > other.Major {
		return 1
	}

	// Compare minor version
	if s.Minor < other.Minor {
		return -1
	}
	if s.Minor > other.Minor {
		return 1
	}

	// Compare patch version
	if s.Patch < other.Patch {
		return -1
	}
	if s.Patch > other.Patch {
		return 1
	}

	// At this point, major.minor.patch are equal, so we need to check pre-release identifiers
	// A version without a pre-release has higher precedence
	if s.PreRelease == "" && other.PreRelease != "" {
		return 1
	}
	if s.PreRelease != "" && other.PreRelease == "" {
		return -1
	}
	if s.PreRelease == "" && other.PreRelease == "" {
		return 0
	}

	// Both have pre-release identifiers, compare them
	sPreReleaseParts := strings.Split(s.PreRelease, ".")
	otherPreReleaseParts := strings.Split(other.PreRelease, ".")

	// Compare each pre-release identifier
	minLen := len(sPreReleaseParts)
	if len(otherPreReleaseParts) < minLen {
		minLen = len(otherPreReleaseParts)
	}

	for i := 0; i < minLen; i++ {
		sPart := sPreReleaseParts[i]
		otherPart := otherPreReleaseParts[i]

		// Check if both are numeric
		sNum, sErr := strconv.ParseUint(sPart, 10, 64)
		otherNum, otherErr := strconv.ParseUint(otherPart, 10, 64)

		if sErr == nil && otherErr == nil {
			// Both are numeric, compare numerically
			if sNum < otherNum {
				return -1
			}
			if sNum > otherNum {
				return 1
			}
		} else if sErr != nil && otherErr != nil {
			// Both are non-numeric, compare lexically
			if sPart < otherPart {
				return -1
			}
			if sPart > otherPart {
				return 1
			}
		} else {
			// One is numeric, one is not
			// Numeric identifiers always have lower precedence
			if sErr == nil { // s is numeric
				return -1
			} else { // other is numeric
				return 1
			}
		}
	}

	// If we've compared all identifiers and they're equal up to the length of the shorter one,
	// the version with fewer identifiers has lower precedence
	if len(sPreReleaseParts) < len(otherPreReleaseParts) {
		return -1
	}
	if len(sPreReleaseParts) > len(otherPreReleaseParts) {
		return 1
	}

	// They're completely equal
	return 0
}

// Sort sorts a slice of SemVer objects in ascending order according to semantic versioning precedence rules.
func Sort(versions []SemVer) {
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) < 0
	})
}
