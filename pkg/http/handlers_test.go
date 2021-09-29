package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	json "github.com/json-iterator/go"

	"github.com/julb/go/pkg/build"
	"github.com/julb/go/pkg/monitoring"
	"github.com/stretchr/testify/assert"
)

func Test_WhenCallingHealthCheckHandlerAndSystemIsUp_ShouldReturnAppropriateResponse(t *testing.T) {
	actualSystemStatus := monitoring.Up
	monitoring.GetHealthContributor().With("status", actualSystemStatus)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	var healthStatus map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	assert.Nil(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code, "invalid http status code")
	assert.Equal(t, string(actualSystemStatus), healthStatus["status"])
}

func Test_WhenCallingHealthCheckHandlerAndSystemIsPartial_ShouldReturnAppropriateResponse(t *testing.T) {
	actualSystemStatus := monitoring.Partial
	monitoring.GetHealthContributor().With("status", actualSystemStatus)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	var healthStatus map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	assert.Nil(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code, "invalid http status code")
	assert.Equal(t, string(actualSystemStatus), healthStatus["status"])
}

func Test_WhenCallingHealthCheckHandlerAndSystemIsDown_ShouldReturnAppropriateResponse(t *testing.T) {
	actualSystemStatus := monitoring.Down
	monitoring.GetHealthContributor().With("status", actualSystemStatus)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	var healthStatus map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	assert.Nil(t, err)

	// Assertions
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code, "invalid http status code")
	assert.Equal(t, string(actualSystemStatus), healthStatus["status"])
}

func Test_WhenCallingHealthCheckHandlerAndSystemIsOutOfService_ShouldReturnAppropriateResponse(t *testing.T) {
	actualSystemStatus := monitoring.OutOfService
	monitoring.GetHealthContributor().With("status", actualSystemStatus)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	var healthStatus map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	assert.Nil(t, err)

	// Assertions
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code, "invalid http status code")
	assert.Equal(t, string(actualSystemStatus), healthStatus["status"])
}

func Test_WhenCallingHealthCheckHandlerAndSystemIsUnknown_ShouldReturnAppropriateResponse(t *testing.T) {
	actualSystemStatus := monitoring.Unknown
	monitoring.GetHealthContributor().With("status", actualSystemStatus)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	var healthStatus map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	assert.Nil(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code, "invalid http status code")
	assert.Equal(t, string(actualSystemStatus), healthStatus["status"])
}

func Test_WhenCallingOtherHealthCheckHandler_ShouldReturnAppropriateResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	req, err := http.NewRequest(http.MethodPost, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code, "invalid http status code")

	req, err = http.NewRequest(http.MethodPut, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code, "invalid http status code")

	req, err = http.NewRequest(http.MethodDelete, "/health", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code, "invalid http status code")
}

func Test_WhenCallingGetInfoHandler_ShouldReturnAppropriateResponse(t *testing.T) {
	monitoring.GetInfoContributor().With("build", build.Info)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InfoHandler)

	req, err := http.NewRequest(http.MethodGet, "/info", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	var infoStatus map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &infoStatus)
	infoStatusBuild := infoStatus["build"].(map[string]interface{})
	assert.Nil(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code, "invalid http status code")
	assert.Equal(t, build.Info.Arch, infoStatusBuild["arch"])
	assert.Equal(t, build.Info.Artifact, infoStatusBuild["artifact"])
	assert.Equal(t, build.Info.BuildVersion, infoStatusBuild["buildVersion"])
	assert.Equal(t, build.Info.GitRevision, infoStatusBuild["gitRevision"])
	assert.Equal(t, build.Info.GitShortRevision, infoStatusBuild["gitShortRevision"])
	assert.Equal(t, build.Info.Group, infoStatusBuild["group"])
	assert.Equal(t, build.Info.Name, infoStatusBuild["name"])
	assert.Equal(t, build.Info.Time, infoStatusBuild["time"])
	assert.Equal(t, build.Info.Version, infoStatusBuild["version"])
}

func Test_WhenCallingOtherInfoHandler_ShouldReturnAppropriateResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InfoHandler)

	req, err := http.NewRequest(http.MethodPost, "/info", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code, "invalid http status code")

	req, err = http.NewRequest(http.MethodPut, "/info", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code, "invalid http status code")

	req, err = http.NewRequest(http.MethodDelete, "/info", nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code, "invalid http status code")
}
