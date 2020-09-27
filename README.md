[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/filediscovery)](https://goreportcard.com/report/github.com/Oppodelldog/filediscovery)
[![Coverage Status](https://coveralls.io/repos/github/Oppodelldog/filediscovery/badge.svg)](https://coveralls.io/github/Oppodelldog/filediscovery)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://raw.githubusercontent.com/Oppodelldog/filediscovery/master/LICENSE)
[![Build Status](https://travis-ci.com/Oppodelldog/filediscovery.svg?branch=master)](https://travis-ci.com/Oppodelldog/filediscovery)
[![pkg.go.dev](https://img.shields.io/badge/pkg.go.dev-reference-%23007d9c.svg)](https://pkg.go.dev/github.com/Oppodelldog/filediscovery)


# Filediscovery
> this module helps to find a file in various file locations

## Example
```go
    import "github.com/Oppodelldog/filediscovery"

	fileLocationProviders := []filediscovery.FileLocationProvider{
		WorkingDirProvider(),
		ExecutableDirProvider(),
		EnvVarFilePathProvider(envVarName),
		HomeConfigDirProvider(".config","myapp"),
	}

	discovery := filediscovery.New(fileLocationProviders)

	filePath, err := discovery.Discover("file_to_discover.yml")

	// filePath - contains the first existing file in sequential order of given file providers
	// err - nil if file was found. if no file was found it displays helpful error information
```

## Advanced
If you'd like to implement a custom file Provider, you just need to
implement the ```FileLoactionProvider``` function type.
Here's a sample demonstration:
```go
    // type FileLocationProvider func(fileName string) (string, error)

    func myLocationProvider(filename string)(string, error){
        somePath := myCustomLookupStrategy()
        return somePath,nil
    }

```