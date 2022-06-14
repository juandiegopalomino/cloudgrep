package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func prepareApiUnitTest(t *testing.T) *mockApi {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	datastoreConfigs := config.Datastore{
		Type:           "sqlite",
		DataSourceName: "file::memory:",
	}
	cfg, err := config.GetDefault()
	require.NoError(t, err)
	cfg.Datastore = datastoreConfigs

	ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
	require.NoError(t, err)

	router := gin.Default()
	resources := testdata.GetResources(t)
	mockApi := mockApi{
		router:    router,
		ds:        ds,
		resources: resources,
	}
	SetupRoutes(router, cfg, logger, ds, mockApi.runEngine)

	//write the resources
	require.NotZero(t, len(resources))
	require.NoError(t, mockApi.runEngine(ctx))
	return &mockApi
}

type mockApi struct {
	router    *gin.Engine
	ds        datastore.Datastore
	resources model.Resources
	//if set calling the engine will return it
	engineErr error
	//incremented every time the engine runs
	engineRuns int
}

//runEngine simulates running the engine by writing resources to the datastore
func (m *mockApi) runEngine(ctx context.Context) error {
	m.engineRuns = m.engineRuns + 1
	if m.engineErr != nil {
		return m.engineErr
	}
	//simulate running the engine
	err := m.ds.WriteEngineStatusStart(ctx, "engine")
	if err != nil {
		return err
	}
	err = m.ds.WriteResources(ctx, m.resources)
	if err != nil {
		return err
	}
	err = m.ds.WriteEngineStatusEnd(ctx, "engine", err)
	return err
}

func TestStatsRoute(t *testing.T) {

	m := prepareApiUnitTest(t)
	path := "/api/stats"

	t.Run("SomeResources", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var body *model.Stats
		err := json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, body.ResourcesCount, 3)
	})
}

func TestResourcesRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/resources"

	t.Run("SomeResources", func(t *testing.T) {
		resources := testdata.GetResources(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var body model.Resources
		err = json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, len(body), 3)
		testingutil.AssertEqualsResources(t, body, resources)
	})
}

func TestResourcesPostRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/resources"

	all_resources := testdata.GetResources(t)
	resourceInst1 := all_resources[0]  //i-123 team:infra, release tag, tag region:us-west-2
	resourceInst2 := all_resources[1]  //i-124 team:dev, no release tag
	resourceBucket := all_resources[2] //s3 bucket without tags

	t.Run("FilterSearch", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := strings.NewReader(`{
  "filter":{
    "$or":[
      {
        "team":"infra"
      },
      {
        "team":"dev"
      }
    ]
  }
}`)
		req, err := http.NewRequest("POST", path, body)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var response model.Resources
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2}, response)
	})

	t.Run("FilterEmpty", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := strings.NewReader(`{
  "filter":{ }
}`)
		req, err := http.NewRequest("POST", path, body)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var response model.Resources
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2, resourceBucket}, response)
	})

	t.Run("NoBody", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := strings.NewReader(``)
		req, err := http.NewRequest("POST", path, body)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		var response model.Resources
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		testingutil.AssertEqualsResources(t, model.Resources{resourceInst1, resourceInst2, resourceBucket}, response)
	})
}

func TestResourceRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/resource"

	t.Run("MissingParam", func(t *testing.T) {
		var body map[string]interface{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Equal(t, body["status"], float64(http.StatusBadRequest))
		require.Equal(t, body["error"], "missing required parameter 'id'")
	})

	t.Run("EmptyParam", func(t *testing.T) {
		var body map[string]interface{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("id", "")
		req.URL.RawQuery = q.Encode()
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Equal(t, body["status"], float64(http.StatusBadRequest))
		require.Equal(t, body["error"], "missing required parameter 'id'")
	})

	t.Run("UnknownParam", func(t *testing.T) {
		var body map[string]interface{}
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("id", "blah")
		req.URL.RawQuery = q.Encode()
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusNotFound, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Equal(t, body["status"], float64(http.StatusNotFound))
		require.Equal(t, body["error"], "can't find resource with id 'blah'")
	})

	t.Run("ValidParam", func(t *testing.T) {
		var actualResource model.Resource
		resources := testdata.GetResources(t)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("id", resources[0].Id)
		req.URL.RawQuery = q.Encode()
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &actualResource))
		testingutil.AssertEqualsResource(t, actualResource, *resources[0])
	})
}

func TestFieldsRoute(t *testing.T) {
	m := prepareApiUnitTest(t)
	path := "/api/fields"

	t.Run("Standard", func(t *testing.T) {
		var response model.FieldGroups
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(w, req)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		require.Equal(t, len(response), 2)
		//check number of groups
		require.Equal(t, 2, len(response))
		//check fields by group
		require.Equal(t, 2, len(response.FindGroup("core").Fields))
		require.Equal(t, 10, len(response.FindGroup("tags").Fields))
	})
}

func TestRefreshPostRoute(t *testing.T) {
	refreshPath := "/api/refresh"
	engineStatusPath := "/api/enginestatus"

	t.Run("Success", func(t *testing.T) {
		m := prepareApiUnitTest(t)
		engineRuns := m.engineRuns
		//trigger a refresh
		record := httptest.NewRecorder()
		req, err := http.NewRequest("POST", refreshPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		require.Equal(t, "", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, record.Code)
		//check that the state is in success
		record = httptest.NewRecorder()
		req, err = http.NewRequest("GET", engineStatusPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		//check engine was run
		require.Equal(t, engineRuns+1, m.engineRuns)
		require.Equal(t, "application/json; charset=utf-8", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusOK, record.Code)
		body := make(map[string]interface{})
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &body))
		require.Equal(t, "success", body["status"])
		require.Equal(t, "", body["errorMessage"])
	})

	t.Run("EngineError", func(t *testing.T) {
		//test an error while refreshing
		m := prepareApiUnitTest(t)
		engineRuns := m.engineRuns
		m.engineErr = fmt.Errorf("There was an engine error")
		record := httptest.NewRecorder()
		req, err := http.NewRequest("POST", refreshPath, nil)
		require.NoError(t, err)
		m.router.ServeHTTP(record, req)
		//check engine was run
		require.Equal(t, engineRuns+1, m.engineRuns)
		require.Equal(t, "application/json; charset=utf-8", record.Header().Get("Content-Type"))
		require.Equal(t, http.StatusBadRequest, record.Code)
		body := make(map[string]interface{})
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &body))
		require.Equal(t, float64(400), body["status"])
		require.Equal(t, "There was an engine error", body["error"])
	})

}
