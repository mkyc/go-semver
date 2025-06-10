package semver

import (
	"testing"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		name        string
		tag         string
		expected    SemVer
		expectError bool
	}{
		// Valid versions
		{
			name: "Basic version",
			tag:  "1.2.3",
			expected: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expectError: false,
		},
		{
			name: "Version with pre-release",
			tag:  "1.2.3-alpha",
			expected: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha",
			},
			expectError: false,
		},
		{
			name: "Version with build metadata",
			tag:  "1.2.3+build.123",
			expected: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build.123",
			},
			expectError: false,
		},
		{
			name: "Version with pre-release and build metadata",
			tag:  "1.2.3-alpha.1+build.123",
			expected: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1",
				Build:      "build.123",
			},
			expectError: false,
		},
		{
			name: "Version with complex pre-release",
			tag:  "1.2.3-alpha.1.beta.2",
			expected: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1.beta.2",
			},
			expectError: false,
		},
		{
			name: "Version with zero values",
			tag:  "0.0.0",
			expected: SemVer{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			expectError: false,
		},
		{
			name: "Pre-release with numeric identifiers",
			tag:  "1.2.3-0.1.2",
			expected: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "0.1.2",
			},
			expectError: false,
		},
		{
			name: "Pre-release with alphanumeric identifiers",
			tag:  "1.2.3-alpha.1.beta-2",
			expected: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1.beta-2",
			},
			expectError: false,
		},

		// Invalid versions
		{
			name:        "Invalid format - missing patch",
			tag:         "1.2",
			expectError: true,
		},
		{
			name:        "Invalid format - extra version part",
			tag:         "1.2.3.4",
			expectError: true,
		},
		{
			name:        "Invalid major version - not a number",
			tag:         "a.2.3",
			expectError: true,
		},
		{
			name:        "Invalid minor version - not a number",
			tag:         "1.a.3",
			expectError: true,
		},
		{
			name:        "Invalid patch version - not a number",
			tag:         "1.2.a",
			expectError: true,
		},
		{
			name:        "Invalid major version - leading zero",
			tag:         "01.2.3",
			expectError: true,
		},
		{
			name:        "Invalid minor version - leading zero",
			tag:         "1.02.3",
			expectError: true,
		},
		{
			name:        "Invalid patch version - leading zero",
			tag:         "1.2.03",
			expectError: true,
		},
		{
			name:        "Invalid pre-release - empty identifier",
			tag:         "1.2.3-alpha..beta",
			expectError: true,
		},
		{
			name:        "Invalid pre-release - invalid character",
			tag:         "1.2.3-alpha_beta",
			expectError: true,
		},
		{
			name:        "Invalid pre-release - numeric with leading zero",
			tag:         "1.2.3-alpha.01",
			expectError: true,
		},
		{
			name:        "Invalid build metadata - empty identifier",
			tag:         "1.2.3+build..123",
			expectError: true,
		},
		{
			name:        "Invalid build metadata - invalid character",
			tag:         "1.2.3+build_123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			semver, err := Parse(tt.tag)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
				return
			}

			// If we don't expect an error, check the parsed values
			if !tt.expectError {
				if semver.Major != tt.expected.Major {
					t.Errorf("Major version mismatch: got %d, want %d", semver.Major, tt.expected.Major)
				}
				if semver.Minor != tt.expected.Minor {
					t.Errorf("Minor version mismatch: got %d, want %d", semver.Minor, tt.expected.Minor)
				}
				if semver.Patch != tt.expected.Patch {
					t.Errorf("Patch version mismatch: got %d, want %d", semver.Patch, tt.expected.Patch)
				}
				if semver.PreRelease != tt.expected.PreRelease {
					t.Errorf("PreRelease mismatch: got %s, want %s", semver.PreRelease, tt.expected.PreRelease)
				}
				if semver.Build != tt.expected.Build {
					t.Errorf("Build mismatch: got %s, want %s", semver.Build, tt.expected.Build)
				}
			}
		})
	}
}

func TestIsRelease(t *testing.T) {
	tests := []struct {
		name     string
		semver   SemVer
		expected bool
	}{
		{
			name: "Release version",
			semver: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: true,
		},
		{
			name: "Release version with build metadata",
			semver: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build.123",
			},
			expected: true,
		},
		{
			name: "Non-release version with pre-release",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha",
			},
			expected: false,
		},
		{
			name: "Non-release version with pre-release and build metadata",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1",
				Build:      "build.123",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.semver.IsRelease()
			if result != tt.expected {
				t.Errorf("IsRelease() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		semver   SemVer
		expected string
	}{
		{
			name: "Basic version",
			semver: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: "1.2.3",
		},
		{
			name: "Version with pre-release",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha",
			},
			expected: "1.2.3-alpha",
		},
		{
			name: "Version with build metadata",
			semver: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build.123",
			},
			expected: "1.2.3+build.123",
		},
		{
			name: "Version with pre-release and build metadata",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1",
				Build:      "build.123",
			},
			expected: "1.2.3-alpha.1+build.123",
		},
		{
			name: "Version with complex pre-release",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1.beta.2",
			},
			expected: "1.2.3-alpha.1.beta.2",
		},
		{
			name: "Version with zero values",
			semver: SemVer{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			expected: "0.0.0",
		},
		{
			name: "Pre-release with numeric identifiers",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "0.1.2",
			},
			expected: "1.2.3-0.1.2",
		},
		{
			name: "Pre-release with alphanumeric identifiers",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1.beta-2",
			},
			expected: "1.2.3-alpha.1.beta-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.semver.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		version1 SemVer
		version2 SemVer
		expected int
	}{
		// Different major versions
		{
			name: "Major version: 1.0.0 < 2.0.0",
			version1: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			version2: SemVer{
				Major: 2,
				Minor: 0,
				Patch: 0,
			},
			expected: -1,
		},
		{
			name: "Major version: 2.0.0 > 1.0.0",
			version1: SemVer{
				Major: 2,
				Minor: 0,
				Patch: 0,
			},
			version2: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: 1,
		},

		// Different minor versions
		{
			name: "Minor version: 1.0.0 < 1.1.0",
			version1: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			version2: SemVer{
				Major: 1,
				Minor: 1,
				Patch: 0,
			},
			expected: -1,
		},
		{
			name: "Minor version: 1.2.0 > 1.1.0",
			version1: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 0,
			},
			version2: SemVer{
				Major: 1,
				Minor: 1,
				Patch: 0,
			},
			expected: 1,
		},

		// Different patch versions
		{
			name: "Patch version: 1.0.0 < 1.0.1",
			version1: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			version2: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 1,
			},
			expected: -1,
		},
		{
			name: "Patch version: 1.0.2 > 1.0.1",
			version1: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 2,
			},
			version2: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 1,
			},
			expected: 1,
		},

		// Pre-release vs. no pre-release
		{
			name: "Pre-release vs. no pre-release: 1.0.0-alpha < 1.0.0",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			version2: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: -1,
		},
		{
			name: "No pre-release vs. pre-release: 1.0.0 > 1.0.0-alpha",
			version1: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			expected: 1,
		},
		{
			name: "Higher version with pre-release vs. lower version: 2.0.0-alpha < 1.0.0",
			version1: SemVer{
				Major:      2,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			version2: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: -1,
		},

		// Different pre-release identifiers
		{
			name: "Pre-release identifiers: 1.0.0-alpha < 1.0.0-beta",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "beta",
			},
			expected: -1,
		},
		{
			name: "Pre-release identifiers: 1.0.0-beta > 1.0.0-alpha",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "beta",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			expected: 1,
		},

		// Numeric vs. non-numeric pre-release identifiers
		{
			name: "Numeric vs. non-numeric pre-release: 1.0.0-1 < 1.0.0-alpha",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "1",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			expected: -1,
		},
		{
			name: "Non-numeric vs. numeric pre-release: 1.0.0-alpha > 1.0.0-1",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "1",
			},
			expected: 1,
		},

		// Numeric pre-release identifiers compared numerically
		{
			name: "Numeric pre-release comparison: 1.0.0-1 < 1.0.0-2",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "1",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "2",
			},
			expected: -1,
		},
		{
			name: "Numeric pre-release comparison: 1.0.0-10 > 1.0.0-2",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "10",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "2",
			},
			expected: 1,
		},

		// Multiple pre-release identifiers
		{
			name: "Multiple pre-release identifiers: 1.0.0-alpha.1 < 1.0.0-alpha.2",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha.1",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha.2",
			},
			expected: -1,
		},
		{
			name: "Multiple pre-release identifiers: 1.0.0-alpha.beta < 1.0.0-beta",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha.beta",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "beta",
			},
			expected: -1,
		},
		{
			name: "Multiple pre-release identifiers: 1.0.0-alpha.1 < 1.0.0-alpha.1.1",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha.1",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha.1.1",
			},
			expected: -1,
		},

		// Equal versions
		{
			name: "Equal versions: 1.0.0 = 1.0.0",
			version1: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			version2: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: 0,
		},
		{
			name: "Equal versions with pre-release: 1.0.0-alpha = 1.0.0-alpha",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			expected: 0,
		},

		// Build metadata is ignored in precedence
		{
			name: "Build metadata ignored: 1.0.0+build = 1.0.0",
			version1: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
				Build: "build",
			},
			version2: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: 0,
		},
		{
			name: "Build metadata ignored: 1.0.0-alpha+build = 1.0.0-alpha",
			version1: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
				Build:      "build",
			},
			version2: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version1.Compare(tt.version2)
			if result != tt.expected {
				t.Errorf("Compare() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		versions []SemVer
		expected []SemVer
	}{
		{
			name: "Sort mixed versions",
			versions: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
				{Major: 1, Minor: 0, Patch: 0},
				{Major: 2, Minor: 0, Patch: 0},
				{Major: 0, Minor: 9, Patch: 9},
				{Major: 1, Minor: 0, Patch: 1},
				{Major: 1, Minor: 1, Patch: 0},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"},
			},
			expected: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"},
				{Major: 0, Minor: 9, Patch: 9},
				{Major: 1, Minor: 0, Patch: 0},
				{Major: 1, Minor: 0, Patch: 1},
				{Major: 1, Minor: 1, Patch: 0},
				{Major: 2, Minor: 0, Patch: 0},
			},
		},
		{
			name: "Sort pre-release versions",
			versions: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "rc.1"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta.11"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta.2"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.1"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.beta"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.1.beta"},
				{Major: 1, Minor: 0, Patch: 0},
			},
			expected: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.1"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.1.beta"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.beta"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta.2"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta.11"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "rc.1"},
				{Major: 1, Minor: 0, Patch: 0},
			},
		},
		{
			name: "Sort with build metadata (ignored for sorting)",
			versions: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, Build: "build.2"},
				{Major: 1, Minor: 0, Patch: 0, Build: "build.1"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha", Build: "build.3"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha", Build: "build.1"},
			},
			expected: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha", Build: "build.3"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha", Build: "build.1"},
				{Major: 1, Minor: 0, Patch: 0, Build: "build.2"},
				{Major: 1, Minor: 0, Patch: 0, Build: "build.1"},
			},
		},
		{
			name: "Sort already sorted versions",
			versions: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"},
				{Major: 1, Minor: 0, Patch: 0},
				{Major: 1, Minor: 0, Patch: 1},
				{Major: 1, Minor: 1, Patch: 0},
				{Major: 2, Minor: 0, Patch: 0},
			},
			expected: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"},
				{Major: 1, Minor: 0, Patch: 0},
				{Major: 1, Minor: 0, Patch: 1},
				{Major: 1, Minor: 1, Patch: 0},
				{Major: 2, Minor: 0, Patch: 0},
			},
		},
		{
			name: "Sort reverse sorted versions",
			versions: []SemVer{
				{Major: 2, Minor: 0, Patch: 0},
				{Major: 1, Minor: 1, Patch: 0},
				{Major: 1, Minor: 0, Patch: 1},
				{Major: 1, Minor: 0, Patch: 0},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
			},
			expected: []SemVer{
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"},
				{Major: 1, Minor: 0, Patch: 0},
				{Major: 1, Minor: 0, Patch: 1},
				{Major: 1, Minor: 1, Patch: 0},
				{Major: 2, Minor: 0, Patch: 0},
			},
		},
		{
			name: "Higher pre release identifier",
			versions: []SemVer{
				{Major: 0, Minor: 9, Patch: 8},
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
				{Major: 0, Minor: 9, Patch: 9},
			},
			expected: []SemVer{
				// ie.: 1.0.0 > 0.11.0 > 1.0.0-beta > 1.0.0-alpha > 0.11.0-rc.1
				{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
				{Major: 0, Minor: 9, Patch: 8},
				{Major: 0, Minor: 9, Patch: 9},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the input slice to avoid modifying the original
			versions := make([]SemVer, len(tt.versions))
			copy(versions, tt.versions)

			// Sort the versions
			Sort(versions)

			// Check if the sorted versions match the expected order
			if len(versions) != len(tt.expected) {
				t.Errorf("Sort() resulted in slice of length %d, want %d", len(versions), len(tt.expected))
				return
			}

			for i := range versions {
				if versions[i].Compare(tt.expected[i]) != 0 {
					t.Errorf("Sort() at index %d = %v, want %v", i, versions[i].String(), tt.expected[i].String())
				}
			}
		})
	}
}
