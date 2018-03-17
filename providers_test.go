package filediscovery

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
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
