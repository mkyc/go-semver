package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// SemVer Backus–Naur Form Grammar for Valid Versions
//
// <valid semver> ::= <version core> | <version core> "-" <pre-release> | <version core> "+" <build> | <version core> "-" <pre-release> "+" <build>
//
// <version core> ::= <major> "." <minor> "." <patch>
//
// <major> ::= <numeric identifier>
//
// <minor> ::= <numeric identifier>
//
// <patch> ::= <numeric identifier>
//
// <pre-release> ::= <dot-separated pre-release identifiers>
//
// <dot-separated pre-release identifiers> ::= <pre-release identifier> | <pre-release identifier> "." <dot-separated pre-release identifiers>
//
// <build> ::= <dot-separated build identifiers>
//
// <dot-separated build identifiers> ::= <build identifier> | <build identifier> "." <dot-separated build identifiers>
//
// <pre-release identifier> ::= <alphanumeric identifier> | <numeric identifier>
//
// <build identifier> ::= <alphanumeric identifier> | <digits>
//
// <alphanumeric identifier> ::= <non-digit> | <non-digit> <identifier characters> | <identifier characters> <non-digit> | <identifier characters> <non-digit> <identifier characters>
//
// <numeric identifier> ::= "0" | <positive digit> | <positive digit> <digits>
//
// <identifier characters> ::= <identifier character> | <identifier character> <identifier characters>
//
// <identifier character> ::= <digit> | <non-digit>
//
// <non-digit> ::= <letter> | "-"
//
// <digits> ::= <digit> | <digit> <digits>
//
// <digit> ::= "0" | <positive digit>
//
// <positive digit> ::= "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
//
// <letter> ::= "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | "I" | "J" | "K" | "L" | "M" | "N" | "O" | "P" | "Q" | "R" | "S" | "T" | "U" | "V" | "W" | "X" | "Y" | "Z" | "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i" | "j" | "k" | "l" | "m" | "n" | "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" | "w" | "x" | "y" | "z"
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
