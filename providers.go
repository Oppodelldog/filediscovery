package filediscovery

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var workingDirProviderFunc = os.Getwd

// WorkingDirProvider provides the working directory as a possible file location
func WorkingDirProvider() FileLocationProvider {

	return func(fileName string) (string, error) {
		dir, err := workingDirProviderFunc()
		if err != nil {
			return "", err
		}

		return path.Join(dir, fileName), nil
	}
}

var executableDirProviderFunc = os.Executable

// ExecutableDirProvider provides the executables directory as a possible file location
func ExecutableDirProvider() FileLocationProvider {

	return func(fileName string) (string, error) {
		dir, err := executableDirProviderFunc()
		if err != nil {
			return "", err
		}

		return path.Join(filepath.Dir(dir), fileName), nil
	}
}

var envVarFilePathLookupFunc = os.LookupEnv

// EnvVarFilePathProvider provides a filePath in the given environment variable.
// In contrast to other FileLocationProviders, this file location provider expects a complete filePath in the given
// environment variable.
func EnvVarFilePathProvider(envVar string) FileLocationProvider {
	return func(fileName string) (string, error) {
		_ = fileName
		if envConfigFile, ok := envVarFilePathLookupFunc(envVar); ok {
			return envConfigFile, nil
		}

		return "", fmt.Errorf("env var '%s' not defined", envVar)
	}
}
