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

	getActual := func(value interface{}) interface{} {
		reflectValue := reflect.ValueOf(value)
		if reflectValue.Kind() == reflect.Ptr {
			if reflectValue.IsNil() {
				return nil
			} else {
				return reflectValue.Elem().Interface()
			}
		}

		return value
	}

	wrap := func(funcs ...func()) {
		for pos, function := range funcs {
			func() {
				defer func() {
					if err := recover(); err != nil {
						t.Errorf("Panic when executing function %d: %v", pos, err)
					}
				}()

				function()
			}()
		}
	}

	assertEquals := func(expected interface{}, actual interface{}, message string, t *testing.T) bool {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting equals %s: %v", message, err)
			}
		}()

		realActual := getActual(actual)
		realExpected := getActual(expected)
		if realActual != realExpected && !reflect.DeepEqual(realActual, realExpected) {
			t.Errorf(
				"Expected %s to be \"%#v\" (%s / %s), but it was actually \"%#v\" (%s / %s).",
				message,
				realExpected, reflect.ValueOf(expected).Kind(), reflect.ValueOf(realExpected).Kind(),
				realActual, reflect.ValueOf(actual).Kind(), reflect.ValueOf(realActual).Kind(),
			)
			return false
		}
		return true
	}

	assertNotNil := func(actual interface{}, message string, t *testing.T) bool {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting not nil %s: %v", message, err)
			}
		}()

		realActual := getActual(actual)
		if realActual == nil {
			t.Errorf("Expected %s to not be nil.", message)
			return false
		}
		return true
	}

	assertNil := func(actual interface{}, message string, t *testing.T) bool {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting nil %s: %v", message, err)
			}
		}()

		realActual := getActual(actual)
		if realActual != nil {
			t.Errorf("Expected %s (%v, type %s) to be nil.", message, realActual, reflect.ValueOf(actual).Kind())
			return false
		}

		return true
	}

	assertNotEmpty := func(actual string, message string, t *testing.T) bool {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic when asserting not empty %s: %v", message, err)
			}
		}()

		if actual == "" {
			t.Errorf("Expected %s not to be empty.", message)
			return false
		}
		return true
	}

	assertNotNilOrEmpty := func(actual *string, message string, t *testing.T) bool {
		return assertNotNil(actual, message, t) && assertNotEmpty(*actual, message, t)
	}

	t.Run("test_parse_build.complete.json", func(t *testing.T) {
		serviceHook, err := parseServiceHook("build.complete.json")
		if err != nil {
			t.Error(err)
			return
		}

		// ServiceHook
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

		// ServiceHookResource
		assertNotNilOrEmpty(serviceHook.Resource.URI, "Resource URI", t)
		assertEquals(1, serviceHook.Resource.ID, "Resource ID", t)

		// ServiceHookResourceBuildComplete
		assertNotNilOrEmpty(serviceHook.Resource.BuildNumber, "Resource Build Number", t)
		assertNotNilOrEmpty(serviceHook.Resource.URL, "Resource URI", t)
		assertNotNil(serviceHook.Resource.StartTime, "Resource Start Time", t)
		assertNotNil(serviceHook.Resource.FinishTime, "Resource Finish Time", t)
		assertEquals(string(azuredevops.ReasonManual), serviceHook.Resource.Reason, "Resource Reason", t)
		assertEquals(string(azuredevops.StatusSucceeded), serviceHook.Resource.Status, "Resource Status", t)
		assertNotNilOrEmpty(serviceHook.Resource.DropLocation, "Resource Drop Location", t)
		wrap(
			func() { assertNotNil(serviceHook.Resource.Drop, "Resource Drop", t) },
			func() { assertNotEmpty(serviceHook.Resource.Drop.Location, "Resource Drop Location (struct)", t) },
			func() { assertNotEmpty(serviceHook.Resource.Drop.Type, "Resource Drop Type", t) },
			func() { assertNotEmpty(serviceHook.Resource.Drop.URL, "Resource Drop URL", t) },
			func() { assertNotEmpty(serviceHook.Resource.Drop.DownloadURL, "Resource Drop Download URL", t) },
		)
		wrap(
			func() { assertNotNil(serviceHook.Resource.Log, "Resource Log", t) },
			func() { assertNotEmpty(serviceHook.Resource.Log.Type, "Resource Log Type", t) },
			func() { assertNotEmpty(serviceHook.Resource.Log.URL, "Resource Log URL", t) },
			func() { assertNotEmpty(serviceHook.Resource.Log.DownloadURL, "Resource Log URI", t) },
		)
		assertNotNilOrEmpty(serviceHook.Resource.SourceGetVersion, "Resource Source Get Version", t)
		assertNotNil(serviceHook.Resource.LastChangedBy, "Resource Last Changed By", t)
		assertNotEmpty(serviceHook.Resource.LastChangedBy.ID, "Resource Last Changed By ID", t)
		assertNotEmpty(serviceHook.Resource.LastChangedBy.DisplayName, "Resource Last Changed By Display Name", t)
		assertNotEmpty(serviceHook.Resource.LastChangedBy.UniqueName, "Resource Last Changed By Unique Name", t)
		assertNotNilOrEmpty(serviceHook.Resource.LastChangedBy.URL, "Resource Last Changed By URL", t)
		assertNotNilOrEmpty(serviceHook.Resource.LastChangedBy.ImageURL, "Resource Last Changed By Image URL", t)
		assertEquals(false, serviceHook.Resource.RetainIndefinitely, "Resource Retain Indefinitely", t)
		assertEquals(true, serviceHook.Resource.HasDiagnostics, "Resource Has Diagnostics", t)

		// ServiceHookResource
		assertNil(serviceHook.Resource.Approval, "Resource Approval", t)
		assertNil(serviceHook.Resource.Release, "Resource Release", t)
		assertNil(serviceHook.Resource.Environment, "Resource Environment", t)
		assertNil(serviceHook.Resource.Project, "Resource Project", t)
		assertNil(serviceHook.Resource.Repository, "Resource Repository", t)
		assertEquals(0, len(serviceHook.Resource.Commits), "Resource Commits", t)

		// ServiceHookResourceCodeCheckedIn
		assertNil(serviceHook.Resource.ChangesetID, "Resource Changeset ID", t)
		assertNil(serviceHook.Resource.Author, "Resource Author", t)
		assertNil(serviceHook.Resource.CheckedInBy, "Resource Checked In By", t)
		assertNil(serviceHook.Resource.CreatedDate, "Resource Created Date", t)
		assertNil(serviceHook.Resource.Comment, "Resource Comment", t)

		// ServiceHookResourceCodePushed
		assertEquals(0, len(serviceHook.Resource.RefUpdates), "Resource Ref Updates", t)
		assertNil(serviceHook.Resource.PushedBy, "Resource Pushed By", t)
		assertNil(serviceHook.Resource.PushID, "Resource Push ID", t)
		assertNil(serviceHook.Resource.Date, "Resource Date", t)

		// ServiceHookResourcePullRequest
		assertNil(serviceHook.Resource.PullRequestID, "Resource Pull Request ID", t)
		assertNil(serviceHook.Resource.CreatedBy, "Resource CreatedBy", t)
		assertNil(serviceHook.Resource.CreationDate, "Resource Creation Date", t)
		assertNil(serviceHook.Resource.ClosedDate, "Resource Closed Date", t)
		assertNil(serviceHook.Resource.Title, "Resource Title", t)
		assertNil(serviceHook.Resource.Description, "Resource Description", t)
		assertNil(serviceHook.Resource.SourceRefName, "Resource Source Ref Name", t)
		assertNil(serviceHook.Resource.TargetRefName, "Resource Target Ref Name", t)
		assertNil(serviceHook.Resource.MergeStatus, "Resource Merge Status", t)
		assertNil(serviceHook.Resource.MergeID, "Resource Merge ID", t)
		assertNil(serviceHook.Resource.LastMergeSourceCommit, "Resource Last Merge Source Commit", t)
		assertNil(serviceHook.Resource.LastMergeTargetCommit, "Resource Last Merge Target Commit", t)
		assertNil(serviceHook.Resource.LastMergeCommit, "Resource Last Merge Commit", t)
		assertEquals(0, len(serviceHook.Resource.Reviewers), "Resource Reviewers", t)

		// ServiceHookResourceWorkItems
		assertNil(serviceHook.Resource.Rev, "Resource Rev", t)
		assertEquals(0, len(serviceHook.Resource.Fields), "Resource Fields", t)

		// ServiceHookResourceWorkItemsUpdatedRevision
		assertNil(serviceHook.Resource.WorkItemID, "Resource Work Item ID", t)
		assertNil(serviceHook.Resource.RevisedBy, "Resource Revised By", t)
		assertNil(serviceHook.Resource.RevisedDate, "Resource Revised Date", t)
		assertNil(serviceHook.Resource.Revision, "Resource Revision", t)

	})
}
