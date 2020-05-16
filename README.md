[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/filediscovery)](https://goreportcard.com/report/github.com/Oppodelldog/filediscovery)
[![Coverage Status](https://coveralls.io/repos/github/Oppodelldog/filediscovery/badge.svg)](https://coveralls.io/github/Oppodelldog/filediscovery)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://raw.githubusercontent.com/Oppodelldog/filediscovery/master/LICENSE)
[![Build Status](https://travis-ci.com/Oppodelldog/filediscovery.svg?branch=master)](https://travis-ci.com/Oppodelldog/filediscovery)
[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/Oppodelldog/filediscovery)


# Filediscovery
> this module helps to find a file in various file locations

## Example
See the example in [test/example_test.go](test/example_test.go).

```go
    //noinspection All
    package main

    import "github.com/Oppodelldog/filediscovery"
    func main(){
    
        var envVarName = "MYAPP_CONFIG_FILE"

        fileLocationProviders := []filediscovery.FileLocationProvider{
            filediscovery.WorkingDirProvider(),
            filediscovery.ExecutableDirProvider(),
            filediscovery.EnvVarFilePathProvider(envVarName),
            filediscovery.HomeConfigDirProvider(".config","myapp"),
        }
    
        discovery := filediscovery.New(fileLocationProviders)
    
        filePath, err := discovery.Discover("file_to_discover.yml")
        _ , _= filePath,err
        // filePath - contains the first existing file in sequential order of given file providers
        // err - nil if file was found. if no file was found it displays helpful error information
    }
```

## Advanced
If you'd like to implement a custom file Provider, you just need to
implement the ```FileLoactionProvider``` function type.
Here's a sample demonstration:
```go
    package main

    // type FileLocationProvider func(fileName string) (string, error)
    
    //noinspection All
    func myLocationProvider(filename string)(string, error){
        somePath := "get it from your custom lookup strategy"

        return somePath,nil
    }

```