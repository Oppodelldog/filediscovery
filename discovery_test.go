package filediscovery

import (
	"bytes"
	"errors"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"fmt"
	"io/ioutil"
)

func TestNew(t *testing.T) {
	interfaceType := reflect.TypeOf(new(FileDiscoverer)).Elem()
	object := New([]FileLocationProvider{})
	if !reflect.TypeOf(object).Implements(interfaceType) {
		t.Fatalf("%T must implement %v", object, interfaceType)
	}
}

func TestFileDiscovery_Discover_callsFileLocationProviders(t *testing.T) {

	mock1, provider1 := newFileLocationProviderMock()
	mock2, provider2 := newFileLocationProviderMock()

	providers := []FileLocationProvider{
		provider1,
		provider2,
	}

	discovery := New(providers)
	_, _ = discovery.Discover("")

	if !mock1.WasCalled() {
		t.Fatalf("expected mock1 to be called, but it was not")
	}
	if !mock2.WasCalled() {
		t.Fatalf("expected mock2 to be called, but it was not")
	}
}

func TestFileDiscovery_Discover_callsFileLocationProvidersWithFilename(t *testing.T) {

	mock, provider := newFileLocationProviderMock()
	providers := []FileLocationProvider{provider}

	discovery := New(providers)
	testFilename := "test-file"
	_, _ = discovery.Discover(testFilename)

	if testFilename != mock.GetCalledFilenameParameter() {
		t.Fatalf("expected %s would have been passed to LocationProviderMock, but it got: %s", testFilename, mock.GetCalledFilenameParameter())
	}
}

func TestFileDiscovery_Discover_providerErrorsAreAppendedToError(t *testing.T) {

	errorMessage := "stub-error"
	errStub := errors.New(errorMessage)
	mock := &fileLocationProviderMock{}
	provider := mock.GetFunc("", errStub)

	providers := []FileLocationProvider{provider}

	discovery := New(providers)
	testFilename := "test-file"
	_, err := discovery.Discover(testFilename)

	if err == nil {
		t.Fatalf("expected error to be not nil, but it was")
	}
	if !strings.Contains(err.Error(), errorMessage) {
		t.Fatalf("expected error %s to contain %s", err.Error(), errorMessage)
	}
}

func TestFileDiscovery_Discover_ifFileNotFoundReturnsError(t *testing.T) {

	mock := &fileLocationProviderMock{}
	provider := mock.GetFunc("", nil)

	providers := []FileLocationProvider{provider}

	discovery := New(providers)
	testFilename := "test-file"
	_, err := discovery.Discover(testFilename)

	expectedError := bytes.NewBufferString("could not find config file at ''")
	expectedError.WriteString("\n")

	if err == nil {
		t.Fatalf("expected error to be not nil, but it was")
	}
	if !strings.Contains(expectedError.String(), err.Error()) {
		t.Fatalf("expected error %s to contain %s", expectedError.String(), err.Error())
	}
}

func TestFileDiscovery_DiscoverMultipleProviders_ifFileNotFoundReturnsMultipleErrorLines(t *testing.T) {

	mock := &fileLocationProviderMock{}
	provider := mock.GetFunc("", nil)

	providers := []FileLocationProvider{provider, provider}

	discovery := New(providers)
	testFilename := "test-file"
	_, err := discovery.Discover(testFilename)

	expectedError := bytes.NewBufferString("could not find config file at ''\ncould not find config file at ''")
	expectedError.WriteString("\n")

	if err == nil {
		t.Fatalf("expected error to be not nil, but it was")
	}
	if !strings.Contains(expectedError.String(), err.Error()) {
		t.Fatalf("expected error %s to contain %s", expectedError.String(), err.Error())
	}
}

func TestFileDiscovery_Discover_ifFileWasFoundReturnsFilePath(t *testing.T) {

	testFilename := "test-file"
	testFilePath := path.Join(os.TempDir(), testFilename)
	f, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("did not expect os.Create to return an error, but got: %v", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("did not expect f.Close to return an error, but got: %v", err)
	}

	tempDirProvider := func(fileName string) (string, error) {
		return path.Join(os.TempDir(), testFilename), nil
	}

	providers := []FileLocationProvider{tempDirProvider}

	discovery := New(providers)

	result, err := discovery.Discover(testFilename)
	if err != nil {
		t.Fatalf("did not expect discovery.Discover to return an error, but got: %v", err)
	}

	if testFilePath != result {
		t.Fatalf("expected '%s' to match '%s'", testFilePath, result)
	}

	err = os.Remove(testFilePath)
	if err != nil {
		t.Fatalf("did not expect os.Remove to return an error, but got: %v", err)
	}
}

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
	discovery := New(
		[]FileLocationProvider{
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
