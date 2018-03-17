package filediscovery

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

type (
	// FileDiscoverer defines logic to discover a file.
	FileDiscoverer interface {

		// Discover tries to find the given fileName in all FileLocationProviders. The providers are checked in given sequence.
		// the first matching result will be returned. If the file could not be found and error is returned as if any other
		// error occurs.
		Discover(fileName string) (string, error)
	}

	FileDiscovery struct {
		fileLocationProviders []FileLocationProvider
	}

	// FileLocationProvider provides a possible file location to FileDiscoverer
	FileLocationProvider func(fileName string) (string, error)
)

// New creates a new FileDiscoverer and takes a list of FileLocationProviders which specify possible location a given file
// will be searched in.
func New(fileLocationProviders []FileLocationProvider) FileDiscoverer {
	return &FileDiscovery{
		fileLocationProviders: fileLocationProviders,
	}
}

// Discover tries to find the given fileName in all FileLocationProviders. The providers are checked in given sequence.
// the first matching result will be returned. If the file could not be found and error is returned as if any other
// error occurs.
func (fd *FileDiscovery) Discover(fileName string) (string, error) {

	var possibleConfigFilePaths []string

	errorString := bytes.NewBufferString("")

	for _, configFileProvider := range fd.fileLocationProviders {
		possibleConfigFile, err := configFileProvider(fileName)
		if err != nil {
			errorString.WriteString(err.Error())
			errorString.WriteString("\n")
		}

		possibleConfigFilePaths = append(possibleConfigFilePaths, possibleConfigFile)
	}

	for _, configFilePath := range possibleConfigFilePaths {
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			errorString.WriteString(fmt.Sprintf("could not find config file at '%s'\n", configFilePath))
			continue
		}

		if info, err := os.Stat(configFilePath); err == nil && !info.IsDir() {
			return configFilePath, nil
		}
	}

	return "", errors.New(errorString.String())
}
