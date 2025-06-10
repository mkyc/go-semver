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
			semver, err := ParseSemVer(tt.tag)

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
