package filediscovery

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const testFileName = "test_config.yml"

const envVarName = "TEST_FILE_PATH"

var providersUnderTest = []FileLocationProvider{
	WorkingDirProvider(),
	ExecutableDirProvider(),
	EnvVarFilePathProvider(envVarName),
}

type TestCase struct {
	no             int
	description    string
	prepareFunc    func() error
	cleanupFunc    func() error
	expectedString string
	expectError    bool
}

func TestIntegrationOfDiscoveryWithProviders(t *testing.T) {

	workingDir, _ := os.Getwd()
	workingDirTestFilePath := path.Join(workingDir, testFileName)

	executableFilePath, _ := os.Executable()
	executableDir := filepath.Dir(executableFilePath)
	executableDirTestFilePath := path.Join(executableDir, testFileName)

	envVarTestFolder := "/tmp/filediscovery-test/envVarTest/"
	envVarTestFilePath := path.Join(envVarTestFolder, testFileName)

	testCases := []TestCase{
		{
			no:             1,
			description:    "Without preparation, there will be no filepath returned, but an error",
			prepareFunc:    func() error { return nil },
			cleanupFunc:    func() error { return nil },
			expectedString: "",
			expectError:    true,
		},
		{
			no:             2,
			description:    "File created under working dir, expect filepath and no error",
			prepareFunc:    func() error { return ioutil.WriteFile(workingDirTestFilePath, []byte("workingDir"), 0777) },
			cleanupFunc:    func() error { return os.Remove(workingDirTestFilePath) },
			expectedString: workingDirTestFilePath,
			expectError:    false,
		},
		{
			no:             3,
			description:    "File created under executable dir, expect filepath and no error",
			prepareFunc:    func() error { return ioutil.WriteFile(executableDirTestFilePath, []byte("executableDir"), 0777) },
			cleanupFunc:    func() error { return os.Remove(executableDirTestFilePath) },
			expectedString: executableDirTestFilePath,
			expectError:    false,
		},
		{
			no:          4,
			description: "File created under envVar filepath, expect filepath and no error",
			prepareFunc: func() error {
				err := os.MkdirAll(envVarTestFolder, 0777)
				if err != nil {
					t.Fatalf("mkdir failed with an error: %v", err)
				}
				err = os.Setenv(envVarName, envVarTestFilePath)
				if err != nil {
					t.Fatalf("setting env var %s failed with an error: %v", envVarName, err)
				}
				return ioutil.WriteFile(envVarTestFilePath, []byte("envVarFilePath"), 0777)
			},
			cleanupFunc: func() error {
				err := os.Unsetenv(envVarName)
				if err != nil {
					t.Fatalf("unsetting env var %s failed with an error: %v", envVarName, err)
				}
				return os.RemoveAll(envVarTestFolder)
			},
			expectedString: envVarTestFilePath,
			expectError:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("test %v", testCase.no), func(t *testing.T) {
			err := testCase.prepareFunc()
			if err != nil {
				t.Fatalf("prepare func returned error: %v", err)
			}

			discovery := New(providersUnderTest)
			result, err := discovery.Discover(testFileName)

			if testCase.expectedString != result {
				t.Fatalf("expecxted provider to return '%s', but got '%s'", testCase.expectedString, result)
			}
			if testCase.expectError {
				if err == nil {
					t.Fatalf("expecxted provider to return an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expecxted provider to return no error, but got %v", err)
				}
			}

			err = testCase.cleanupFunc()
			if err != nil {
				t.Fatalf("cleanup func returned error: %v", err)
			}
		})
	}
}
