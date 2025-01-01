package config

import (
	"os"
	"testing"

	"github.com/knadh/koanf/v2"
)

func TestRead(t *testing.T) {

	var mockConfig = koanf.New(".")

	// Create a temporary YAML config file for testing
	mockConfigFile := "test_config.yaml"
	yamlContent := `
loglevel: "debug"
scanIntervalMin: 30
disable_ui: true
storage: "s3"
s3BucketName: "test-bucket"
s3BucketPath: "test/path"
httpPort: 9090
`
	err := os.WriteFile(mockConfigFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(mockConfigFile)

	// Set environment variables for testing
	// we want to override the loglevel to error
	os.Setenv("CC_LOGLEVEL", "error")
	defer os.Unsetenv("CC_LOGLEVEL")

	os.Args = []string{"cmd", "-config", mockConfigFile}
	Read(mockConfig)

	// Test assertions
	assertions := map[string]string{
		"loglevel":        "error",
		"scanIntervalMin": "30",
		"disable_ui":      "true",
		"storage":         "s3",
		"s3BucketName":    "test-bucket",
		"s3BucketPath":    "test/path",
		"httpPort":        "9090",
	}
	for setting, want := range assertions {
		t.Run(setting, func(t *testing.T) {
			if mockConfig.String(setting) != want {
				t.Errorf("Expected "+setting+" to be '"+want+
					"', got '%s'", mockConfig.String(setting))
			}
		})
	}
}
