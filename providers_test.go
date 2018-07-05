package filediscovery

import (
	"os"
	"path"
	"testing"

	"errors"

	"os/user"
	"path/filepath"

	"github.com/stretchr/testify/assert"
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

			assert.Equal(t, testData.ExpectedPath, result)
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

	assert.Equal(t, errorStub, err)
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

			assert.Equal(t, testData.ExpectedPath, result)
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

	assert.Equal(t, errorStub, err)
}

func TestEnvVarFilePathProvider(t *testing.T) {
	testVarName := "TEST-VAR"
	testVarValue := "TEST-VALUE"
	os.Setenv(testVarName, testVarValue)

	provider := EnvVarFilePathProvider(testVarName)

	testFileName := "" // not necessary for this test, since filename comes from env var
	result, err := provider(testFileName)
	if err != nil {
		t.Fatalf("Did not expect provider to return an error, but got: %v", err)
	}

	assert.Equal(t, testVarValue, result)

	os.Unsetenv(testVarName)
}

func TestEnvVarFilePathProvider_error(t *testing.T) {
	testVarName := "TEST-VAR"

	provider := EnvVarFilePathProvider(testVarName)

	testFileName := "" // not necessary for this test, since filename comes from env var
	_, err := provider(testFileName)
	assert.Error(t, err)
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

			assert.Equal(t, testData.ExpectedPath, result)
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

	assert.Exactly(t, errorStub, err)
}
