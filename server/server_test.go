package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/inview-team/raptor_team/registry/task"
)

var (
	router = new(Server).setupRouter()
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestTopology(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		reqBody := task.Task{
			CameraIP: "10.11.12.13",
			Jobs: []task.Job{
				{
					Title: "job1",
				},
				{
					Title: "job2",
				},
			},
		}
		body, err := json.Marshal(reqBody)
		require.Nil(t, err)

		w := performRequest(router, "POST", "/newtask", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)
	})
}
