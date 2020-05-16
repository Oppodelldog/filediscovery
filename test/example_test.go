package test

import (
	"fmt"
	"github.com/Oppodelldog/filediscovery"
	"io/ioutil"
	"os"
	"path"
)

func ExampleFileDiscovery_Discover() {

	// for this demonstration we create a test file in /tmp
	testFilePath := "/tmp/test-file.yml"
	err := ioutil.WriteFile(testFilePath, []byte("test"), 0666)
	if err != nil {
		panic("error writing test file")
	}

	// Discovery needs at least one FileLocationProvider which provides a file location to search for.
	// There are already some providers available, but let's create a new one, for the sake of completion.
	// In this case the FileLocation provided will be the /tmp folder.
	tempDirLocationProvider := func(fileName string) (string, error) {
		someFileLocation := "/tmp"
		suggestedFilePath := path.Join(someFileLocation, fileName)
		return suggestedFilePath, nil
	}

	// create the discovery
	discovery := filediscovery.New(
		[]filediscovery.FileLocationProvider{
			tempDirLocationProvider,
		},
	)

	// search the file
	fileFoundAtPath, _ := discovery.Discover("test-file.yml")

	// we receive the full path of the file found in /tmp
	fmt.Println(fileFoundAtPath)
	// Output: /tmp/test-file.yml

	err = os.Remove(testFilePath)
	if err != nil {
		panic("error removing test file")
	}
}
