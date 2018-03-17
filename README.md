[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/filediscovery)](https://goreportcard.com/report/github.com/Oppodelldog/filediscovery)
[![Coverage Status](https://coveralls.io/repos/github/Oppodelldog/filediscovery/badge.svg)](https://coveralls.io/github/Oppodelldog/filediscovery)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://raw.githubusercontent.com/Oppodelldog/filediscovery/master/LICENSE)
[![Linux build](http://nulldog.de:12080/api/badges/Oppodelldog/filediscovery/status.svg)](http://nulldog.de:12080/Oppodelldog/filediscovery)
[![Windows build](https://ci.appveyor.com/api/projects/status/qpe2889fbk1bw7lf/branch/master?svg=true)](https://ci.appveyor.com/project/Oppodelldog/filediscovery/branch/master)

# Filediscovery
> this module helps to find a file in various file locations

## Example
```go
	fileLocationProviders := []FileLocationProvider{
		WorkingDirProvider(),
		ExecutableDirProvider(),
		EnvVarFilePathProvider(envVarName),
	}

	discovery := New(fileLocationProviders)

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