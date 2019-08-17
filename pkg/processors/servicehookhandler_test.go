package processors_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/args"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/config"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetes"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetesmock"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/processors"
)

func TestHTTPMethod(t *testing.T) {
	args := args.ServiceHookArgs{
		Username: "testusername",
		Password: "VeryStrongP@$$W0RD",
	}

	handler := processors.NewServiceHookHandler(args, []config.ServiceHook{}, processors.NewRuleHandler(kubernetes.MakeFromClient(kubernetesmock.MockK8sClient{})))

	httpMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

	for _, httpMethod := range httpMethods {
		t.Run(fmt.Sprintf("httpmethod_test_%s", httpMethod), func(t *testing.T) {
			req, err := http.NewRequest(httpMethod, "/serviceHooks", bytes.NewBufferString("{ \"eventType\": \"mock\" }"))
			req.SetBasicAuth(args.Username, args.Password)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			expectedStatusCode := http.StatusMethodNotAllowed
			if httpMethod == "POST" {
				expectedStatusCode = http.StatusOK
			}

			if recorder.Code != expectedStatusCode {
				t.Errorf("Expected HTTP status %d but received %d", expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestBasicAuthentication(t *testing.T) {
	args := args.ServiceHookArgs{
		Username: "testusername",
		Password: "VeryStrongP@$$W0RD",
	}

	handler := processors.NewServiceHookHandler(args, []config.ServiceHook{}, processors.NewRuleHandler(kubernetes.MakeFromClient(kubernetesmock.MockK8sClient{})))

	t.Run("basicauthentication_test_good", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/serviceHooks", bytes.NewBufferString("{ \"eventType\": \"mock\" }"))
		req.SetBasicAuth(args.Username, args.Password)
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected HTTP status %d but received %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("basicauthentication_test_bad_noauth", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/serviceHooks", bytes.NewBufferString("{ \"eventType\": \"mock\" }"))
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected HTTP status %d but received %d", http.StatusUnauthorized, recorder.Code)
		}
	})

	t.Run("basicauthentication_test_bad_invalidusername", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/serviceHooks", bytes.NewBufferString("{ \"eventType\": \"mock\" }"))
		req.SetBasicAuth("invalidusername", args.Password)
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected HTTP status %d but received %d", http.StatusUnauthorized, recorder.Code)
		}
	})

	t.Run("basicauthentication_test_bad_invalidpassword", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/serviceHooks", bytes.NewBufferString("{ \"eventType\": \"mock\" }"))
		req.SetBasicAuth(args.Username, "invalidpassword")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected HTTP status %d but received %d", http.StatusUnauthorized, recorder.Code)
		}
	})

	t.Run("basicauthentication_test_bad_invalidusernameandpassword", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/serviceHooks", bytes.NewBufferString("{ \"eventType\": \"mock\" }"))
		req.SetBasicAuth("invalidusername", "invalidpassword")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected HTTP status %d but received %d", http.StatusUnauthorized, recorder.Code)
		}
	})
}
