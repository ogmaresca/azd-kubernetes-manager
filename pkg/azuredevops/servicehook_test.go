package azuredevops_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
)

func TestDeserialize(t *testing.T) {
	parseServiceHook := func(filename string) (*azuredevops.ServiceHook, error) {
		fileContents, err := ioutil.ReadFile(fmt.Sprintf("../../sample-service-hooks/%s", filename))
		if err != nil {
			return nil, fmt.Errorf("Error reading JSON file \"%s\": %s", filename, err.Error())
		}

		serviceHook := new(azuredevops.ServiceHook)
		if err = json.NewDecoder(bytes.NewReader(fileContents)).Decode(serviceHook); err != nil {
			return nil, fmt.Errorf("Error - could not deserialize JSON file \"%s\": %s", filename, err.Error())
		}

		fmt.Println(fmt.Sprintf("Parsed \"%s\": %#v", filename, *serviceHook))

		return serviceHook, nil
	}

	assertEquals := func(expected interface{}, actual interface{}, message string, t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting equals %s: %v", message, err)
			}
		}()

		if (expected == nil) != (actual == nil) {
			fmt.Println("Nil mismatch")
			t.Errorf("Expected %s to be \"%#v\", but it was actually \"%#v\".", message, expected, actual)
			return
		} else if expected == nil {
			return
		}

		realExpected := expected
		if reflect.ValueOf(expected).Kind() == reflect.Ptr {
			realExpected = reflect.ValueOf(expected).Elem()
		}

		realActual := actual
		if reflect.ValueOf(actual).Kind() == reflect.Ptr {
			realActual = reflect.ValueOf(actual).Elem()
		}

		if actual != expected || realActual != realExpected || !reflect.DeepEqual(realActual, realExpected) {
			t.Errorf("Expected %s to be \"%#v\", but it was actually \"%#v\".", message, realExpected, realActual)
		}
	}

	assertNotNil := func(actual interface{}, message string, t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting not nil %s: %v", message, err)
			}
		}()

		if actual == nil {
			t.Errorf("Expected %s to not be nil.", message)
		}
	}

	assertNil := func(actual interface{}, message string, t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting nil %s: %v", message, err)
			}
		}()

		if actual != nil {
			t.Errorf("Expected %s to be nil.", message)
		}
	}

	assertNotEmpty := func(actual string, message string, t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting not empty %s: %v", message, err)
			}
		}()

		if actual == "" {
			t.Errorf("Expected %s not to be empty.", message)
		}
	}

	assertNotNilOrEmpty := func(actual *string, message string, t *testing.T) {
		if actual == nil {
			assertNotNil(actual, message, t)
		} else {
			assertNotEmpty(*actual, message, t)
		}
	}

	t.Run("test_parse_build.complete.json", func(t *testing.T) {
		serviceHook, err := parseServiceHook("build.complete.json")
		if err != nil {
			t.Error(err)
			return
		}

		assertNotEmpty(serviceHook.ID, "ID", t)
		assertEquals("build.complete", serviceHook.EventType, "Event Type", t)
		assertEquals("tfs", serviceHook.PublisherID, "Publisher ID", t)
		assertEquals("all", serviceHook.Scope, "Scope", t)
		assertNotEmpty(serviceHook.Message.Text, "Message Text", t)
		assertNotEmpty(serviceHook.Message.HTML, "Message TeHTMLxt", t)
		assertNotEmpty(serviceHook.Message.Markdown, "Message Markdown", t)
		assertNotEmpty(serviceHook.DetailedMessage.Text, "Detailed Message Text", t)
		assertNotEmpty(serviceHook.DetailedMessage.HTML, "Detailed Message TeHTMLxt", t)
		assertNotEmpty(serviceHook.DetailedMessage.Markdown, "Detailed Message Markdown", t)
		assertNotNilOrEmpty(serviceHook.Resource.URI, "Resource URI", t)
		assertEquals(1, serviceHook.Resource.ID, "Resource ID", t)
		assertNotNilOrEmpty(serviceHook.Resource.BuildNumber, "Resource Build Number", t)
		assertNotNilOrEmpty(serviceHook.Resource.URL, "Resource URI", t)
		assertNotNil(serviceHook.Resource.StartTime, "Resource Start Time", t)
		assertNotNil(serviceHook.Resource.FinishTime, "Resource Finish Time", t)
		assertNotNilOrEmpty(serviceHook.Resource.DropLocation, "Resource Drop Location", t)
		assertNotNil(serviceHook.Resource.Drop, "Resource Drop", t)
		assertNotEmpty(serviceHook.Resource.Drop.Location, "Resource Drop Location (struct)", t)
		assertNotEmpty(serviceHook.Resource.Drop.Type, "Resource Drop Type", t)
		assertNotEmpty(serviceHook.Resource.Drop.URL, "Resource Drop URL", t)
		assertNotEmpty(serviceHook.Resource.Drop.DownloadURL, "Resource Drop Download URL", t)
		assertNotNil(serviceHook.Resource.Log, "Resource Log", t)
		assertNotEmpty(serviceHook.Resource.Log.Type, "Resource Log Type", t)
		assertNotEmpty(serviceHook.Resource.Log.URL, "Resource Log URL", t)
		assertNotEmpty(serviceHook.Resource.Log.DownloadURL, "Resource Log URI", t)
		assertNotNilOrEmpty(serviceHook.Resource.SourceGetVersion, "Resource Source Get Version", t)
		assertNotNil(serviceHook.Resource.LastChangedBy, "Resource Last Changed By", t)
		assertNotEmpty(serviceHook.Resource.LastChangedBy.ID, "Resource Last Changed By ID", t)
		assertNotEmpty(serviceHook.Resource.LastChangedBy.DisplayName, "Resource Last Changed By Display Name", t)
		assertNotEmpty(serviceHook.Resource.LastChangedBy.UniqueName, "Resource Last Changed By Unique Name", t)
		assertNotNilOrEmpty(serviceHook.Resource.LastChangedBy.URL, "Resource Last Changed By URL", t)
		assertNotNilOrEmpty(serviceHook.Resource.LastChangedBy.ImageURL, "Resource Last Changed By Image URL", t)
	})
}
