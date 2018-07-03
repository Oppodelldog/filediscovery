package filediscovery

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"os/user"
)

func TestWorkingDirProvider(t *testing.T) {
	testFileName := "testfile"
	provider := WorkingDirProvider()
	result, err := provider(testFileName)
	if err != nil {
		t.Fatalf("Did not expect provider to return an error, but got: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Did not expect os.Getwd to return an error, but got: %v", err)
	}

	assert.Equal(t, path.Join(wd, testFileName), result)
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

	provider := ExecutableDirProvider()
	result, err := provider(testFileName)
	if err != nil {
		t.Fatalf("Did not expect provider to return an error, but got: %v", err)
	}

	executableFilePath, err := os.Executable()
	if err != nil {
		t.Fatalf("Did not expect os.Executable to return an error, but got: %v", err)
	}

	expectedFilePath := path.Join(filepath.Dir(executableFilePath), testFileName)

	assert.Equal(t, expectedFilePath, result)
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

	subfolder1 := ".config"
	subfolder2 := "some-project"
	provider := HomeConfigDirProvider(subfolder1, subfolder2)
	result, err := provider(testFileName)
	if err != nil {
		t.Fatalf("Did not expect provider to return an error, but got: %v", err)
	}

	usr, err := user.Current()
	if err != nil {
		t.Fatalf("Did not expect user.Current to return an error, but got: %v", err)
	}

	expectedFilepath := path.Join(usr.HomeDir, subfolder1, subfolder2, testFileName)
	assert.Equal(t, expectedFilepath, result)
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
