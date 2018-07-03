package filediscovery

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"os/user"
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


var homeFolderLookupFunc = user.Current

// HomeConfigDirProvider provides the working directory as a possible file location
func HomeConfigDirProvider(subFolders ...string) FileLocationProvider {

	return func(fileName string) (string, error) {
		usr, err := homeFolderLookupFunc()
		if err != nil {
			return "", err
		}

		subfoldersPath := ""
		for _, subfolder := range subFolders {
			subfoldersPath = path.Join(subfoldersPath, subfolder)
		}

		return path.Join(usr.HomeDir, subfoldersPath, fileName), nil
	}
}