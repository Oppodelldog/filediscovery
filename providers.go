package filediscovery

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

var workingDirProviderFunc = os.Getwd

// WorkingDirProvider provides the working directory as a possible file location
func WorkingDirProvider(subFolders ...string) FileLocationProvider {

	return func(fileName string) (string, error) {
		dir, err := workingDirProviderFunc()
		if err != nil {
			return "", err
		}

		subFoldersPath := createPath(subFolders...)

		return path.Join(dir, subFoldersPath, fileName), nil
	}
}

var executableDirProviderFunc = os.Executable

// ExecutableDirProvider provides the executables directory as a possible file location
func ExecutableDirProvider(subFolders ...string) FileLocationProvider {

	return func(fileName string) (string, error) {
		dir, err := executableDirProviderFunc()
		if err != nil {
			return "", err
		}

		executableDir := filepath.Dir(dir)
		subFoldersPath := createPath(subFolders...)

		return path.Join(executableDir, subFoldersPath, fileName), nil
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

var homeFolderLookupFunc = user.Current

// HomeConfigDirProvider provides the working directory as a possible file location
func HomeConfigDirProvider(subFolders ...string) FileLocationProvider {

	return func(fileName string) (string, error) {
		usr, err := homeFolderLookupFunc()
		if err != nil {
			return "", err
		}

		subFoldersPath := createPath(subFolders...)

		return path.Join(usr.HomeDir, subFoldersPath, fileName), nil
	}
}

func createPath(subFolders ...string) string {
	subFoldersPath := ""
	for _, subfolder := range subFolders {
		subFoldersPath = path.Join(subFoldersPath, subfolder)
	}

	return subFoldersPath
}
