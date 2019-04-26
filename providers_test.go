package filediscovery

import (
	"os"
	"path"
	"testing"

	"errors"

	"os/user"
	"path/filepath"
)

func TestWorkingDirProvider(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Did not expect os.Getwd to return an error, but got: %v", err)
	}

	testFileName := "testfile"

	testDataSet := map[string]struct {
		SubFolders   []string
		ExpectedPath string
	}{
		"simple call": {
			SubFolders:   []string{},
			ExpectedPath: path.Join(wd, testFileName),
		},
		"one subdir": {
			SubFolders:   []string{"subdir1"},
			ExpectedPath: path.Join(wd, "subdir1", testFileName),
		},
		"two subdirs": {
			SubFolders:   []string{"subdir1", "subdir2"},
			ExpectedPath: path.Join(wd, "subdir1", "subdir2", testFileName),
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {

			provider := WorkingDirProvider(testData.SubFolders...)
			result, err := provider(testFileName)
			if err != nil {
				t.Fatalf("Did not expect provider to return an error, but got: %v", err)
			}

			if testData.ExpectedPath != result {
				t.Fatalf("expected provider to return '%s', but got: '%s'", testData.ExpectedPath, result)
			}
		})
	}
}

func TestWorkingDirProvider_error(t *testing.T) {
	testFileName := "testfile"

	errorStub := errors.New("error-stub")
	workingDirProviderFunc = func() (string, error) {
		return "", errorStub
	}

	provider := WorkingDirProvider()
	_, err := provider(testFileName)

	if errorStub != err {
		t.Fatalf("did expect provider to return stubbed error %v, but got %v", errorStub, err)
	}
}

func TestExecutableDirProvider(t *testing.T) {
	testFileName := "testfile"

	executableFilePath, err := os.Executable()
	if err != nil {
		t.Fatalf("Did not expect os.Executable to return an error, but got: %v", err)
	}
	executableFilePath = filepath.Dir(executableFilePath)

	testDataSet := map[string]struct {
		SubFolders   []string
		ExpectedPath string
	}{
		"simple call": {
			SubFolders:   []string{},
			ExpectedPath: path.Join(executableFilePath, testFileName),
		},
		"one subdir": {
			SubFolders:   []string{"subdir1"},
			ExpectedPath: path.Join(executableFilePath, "subdir1", testFileName),
		},
		"two subdirs": {
			SubFolders:   []string{"subdir1", "subdir2"},
			ExpectedPath: path.Join(executableFilePath, "subdir1", "subdir2", testFileName),
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			provider := ExecutableDirProvider(testData.SubFolders...)
			result, err := provider(testFileName)
			if err != nil {
				t.Fatalf("Did not expect provider to return an error, but got: %v", err)
			}

			if testData.ExpectedPath != result {
				t.Fatalf("expected provider to return '%s', but got: '%s'", testData.ExpectedPath, result)
			}
		})
	}
}

func TestExecutableDirProvider_error(t *testing.T) {
	testFileName := "testfile"

	errorStub := errors.New("error-stub")
	executableDirProviderFunc = func() (string, error) {
		return "", errorStub
	}

	provider := ExecutableDirProvider()
	_, err := provider(testFileName)

	if errorStub != err {
		t.Fatalf("did expect provider to return stubbed error %v, but got %v", errorStub, err)
	}
}

func TestEnvVarFilePathProvider(t *testing.T) {
	testVarName := "TEST-VAR"
	testVarValue := "TEST-VALUE"
	err := os.Setenv(testVarName, testVarValue)
	if err != nil {
		t.Fatalf("setting env var %s failed with error: %v", envVarName, err)
	}

	provider := EnvVarFilePathProvider(testVarName)

	testFileName := "" // not necessary for this test, since filename comes from env var
	result, err := provider(testFileName)
	if err != nil {
		t.Fatalf("Did not expect provider to return an error, but got: %v", err)
	}

	if testVarValue != result {
		t.Fatalf("expected provider to return value '%s' from env var, but got: '%s'", testVarValue, result)
	}

	err = os.Unsetenv(testVarName)
	if err != nil {
		t.Fatalf("unsetting env var %s failed with error: %v", envVarName, err)
	}
}

func TestEnvVarFilePathProvider_error(t *testing.T) {
	testVarName := "TEST-VAR"

	provider := EnvVarFilePathProvider(testVarName)

	testFileName := "" // not necessary for this test, since filename comes from env var
	_, err := provider(testFileName)
	if err == nil {
		t.Fatalf("expect provider to return and error but got nil")
	}
}

func TestHomeConfigDirProvider(t *testing.T) {
	testFileName := "testfile"

	usr, err := user.Current()
	if err != nil {
		t.Fatalf("Did not expect user.Current to return an error, but got: %v", err)
	}

	testDataSet := map[string]struct {
		SubFolders   []string
		ExpectedPath string
	}{
		"simple call": {
			SubFolders:   []string{},
			ExpectedPath: path.Join(usr.HomeDir, testFileName),
		},
		"one subdir": {
			SubFolders:   []string{"subdir1"},
			ExpectedPath: path.Join(usr.HomeDir, "subdir1", testFileName),
		},
		"two subdirs": {
			SubFolders:   []string{"subdir1", "subdir2"},
			ExpectedPath: path.Join(usr.HomeDir, "subdir1", "subdir2", testFileName),
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			provider := HomeConfigDirProvider(testData.SubFolders...)
			result, err := provider(testFileName)
			if err != nil {
				t.Fatalf("Did not expect provider to return an error, but got: %v", err)
			}

			if testData.ExpectedPath != result {
				t.Fatalf("expected provider to return path '%s', but got: '%s'", testData.ExpectedPath, result)
			}
		})
	}
}

func TestHomeConfigDirProvider_UserLookupReturnsError(t *testing.T) {

	errorStub := errors.New("error-stub")
	homeFolderLookupFunc = func() (*user.User, error) {
		return nil, errorStub
	}

	provider := HomeConfigDirProvider()
	_, err := provider(testFileName)

	if errorStub != err {
		t.Fatalf("expected provider to return %v, but got %v", errorStub, err)
	}
}
