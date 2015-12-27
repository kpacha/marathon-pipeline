package server

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kpacha/marathon-pipeline/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	taskStore := pipeline.MemoryTaskStore{
		Store: &map[string]pipeline.Task{"supu": pipeline.Task{Name: "Name to display", ID: "supu"}},
	}
	testGinServer := GinServer{
		Port:      9876,
		engine:    gin.Default(),
		taskStore: taskStore,
	}
	go func() {
		testGinServer.Run()
	}()
	time.Sleep(5 * time.Millisecond)

	testServerGetMethods(t, testGinServer.engine)
	testServerPostMethods(t, testGinServer.engine)
	testServerDeleteMethods(t, testGinServer.engine)
}

func testServerGetMethods(t *testing.T, gs *gin.Engine) {
	getCases := []struct {
		url    string
		code   int
		result string
	}{
		{"/", 200, "{\"supu\":{\"ID\":\"supu\",\"Name\":\"Name to display\",\"Filter\":null,\"Params\":null}}\n"},
		{"/supu", 200, "{\"ID\":\"supu\",\"Name\":\"Name to display\",\"Filter\":null,\"Params\":null}\n"},
		{"/unknown", 404, "{\"status\":\"Task not found: unknown\"}\n"},
	}

	for _, c := range getCases {
		w := performGetRequest(gs, c.url)
		assert.Equal(t, c.code, w.Code)
		assert.Equal(t, c.result, w.Body.String())
	}
}

func testServerPostMethods(t *testing.T, gs *gin.Engine) {
	postCases := []struct {
		payload string
		code    int
		result  string
	}{
		{
			"{\"id\":\"new\",\"name\":\"new\",\"filter\":null,\"params\":null}",
			200,
			"{\"ID\":\"new\",\"Name\":\"new\",\"Filter\":null,\"Params\":null}\n",
		},
		{
			"{\"id\":\"supu\",\"name\":\"\",\"filter\":null,\"params\":null}",
			400,
			"{\"status\":\"Task already exists: supu\"}\n",
		},
		{
			"{",
			400,
			"{\"status\":\"Error processing the body of the request: unexpected EOF\"}\n",
		},
		{
			"{\"Name\":\"Name to display\"}",
			400,
			"{\"status\":\"Error: Task Id must be set!\"}\n",
		},
	}

	for _, c := range postCases {
		w := performRequest(gs, "POST", "/", c.payload)
		assert.Equal(t, c.code, w.Code)
		assert.Equal(t, c.result, w.Body.String())
	}
}

func testServerDeleteMethods(t *testing.T, gs *gin.Engine) {
	w := performDeleteRequest(gs, "/new")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"status\":\"ok\"}\n", w.Body.String())

	w = performDeleteRequest(gs, "/new")
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "{\"status\":\"Task not found: new\"}\n", w.Body.String())
}

func performRequest(r http.Handler, method, path, payload string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return recordServerRequest(r, req)
}

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	return performEmptyRequest(r, "GET", path)
}

func performDeleteRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	return performEmptyRequest(r, "DELETE", path)
}

func performEmptyRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		log.Fatal(err)
	}
	return recordServerRequest(r, req)
}

func recordServerRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
