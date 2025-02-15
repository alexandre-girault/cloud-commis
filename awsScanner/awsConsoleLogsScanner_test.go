package awsScanner

import (
	"os"
	"testing"
)

func TestAwsConsoleLogOutput(t *testing.T) {

	// Test assertions
	assertions := map[string][]string{
		"../testData/console_Amazon_Linux.txt": {"Amazon Linux 2023.6.20241010", "6.1.112-122.189.amzn2023.aarch64"},
		"../testData/console_Ubuntu_2204.txt":  {"Ubuntu 22.04.5 LTS", "vmlinuz-6.8.0-1015-aws"},
		"../testData/console_RHEL_9.4.txt":     {"Red Hat Enterprise Linux 9.4", "5.14.0-427.20.1.el9_4.x86_64"},
	}

	for outputSample, want := range assertions {
		t.Run(outputSample, func(t *testing.T) {
			consoleOutput, err := os.ReadFile(outputSample)
			if err != nil {
				t.Error(err.Error())
			}
			os, kernel := awsConsoleLogParse(consoleOutput)

			if os != want[0] {
				t.Errorf("Expected OS to be : '%s', got'%s'", want[0], os)
			}

			if kernel != want[1] {
				t.Errorf("Expected kernel to be : '%s', got '%s'", want[1], kernel)
			}
		})
	}

}
