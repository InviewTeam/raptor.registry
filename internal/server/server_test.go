package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/internal/config"
	"gitlab.com/inview-team/raptor_team/registry/pkg/format"
	"gitlab.com/inview-team/raptor_team/registry/tests"
)

type Response struct {
	UUID string `json:"uuid"`
}

var (
	srv = Server{
		reg: registry.New(&config.Settings{}, tests.NewDB(), &tests.Publisher{}),
	}
	router = srv.setupRouter()

	task1 = format.Task{
		CameraIP: "10.11.12.13",
		Job:      "job1",
	}

	task2 = format.Task{
		CameraIP: "92.138.141.54",
		Job:      "some_job_name",
	}
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestCreate(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		body, err := json.Marshal(task1)
		require.Nil(t, err)

		w := performRequest(router, "POST", "/api/tasks", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)

		resp := Response{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.Nil(t, err)
		id, err := uuid.Parse(resp.UUID)
		require.Nil(t, err)
		task1.UUID = id

		body, err = json.Marshal(task2)
		require.Nil(t, err)

		w = performRequest(router, "POST", "/api/tasks", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.Nil(t, err)
		id, err = uuid.Parse(resp.UUID)
		require.Nil(t, err)
		task2.UUID = id
	})
}

func TestGet(t *testing.T) {
	t.Run("get all", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/tasks", nil)
		require.Equal(t, http.StatusOK, w.Code)

		received := []format.Task{}
		err := json.Unmarshal(w.Body.Bytes(), &received)
		require.Nil(t, err)

		require.Equal(t, 2, len(received))

		task1.Status = "in work"
		task2.Status = "in work"

		if diff := deep.Equal(task1, received[0]); diff != nil {
			t.Error(diff)
		}

		if diff := deep.Equal(task2, received[1]); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("get by UUID", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/tasks/"+task1.UUID.String(), nil)
		require.Equal(t, http.StatusOK, w.Code)
		received := format.Task{}
		err := json.Unmarshal(w.Body.Bytes(), &received)
		require.Nil(t, err)
		if diff := deep.Equal(task1, received); diff != nil {
			t.Error(diff)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		w := performRequest(router, "DELETE", "/api/tasks/"+task1.UUID.String(), nil)
		require.Equal(t, http.StatusOK, w.Code)

		w = performRequest(router, "GET", "/api/tasks", nil)
		require.Equal(t, http.StatusOK, w.Code)

		received := []format.Task{}
		err := json.Unmarshal(w.Body.Bytes(), &received)
		require.Nil(t, err)

		require.Equal(t, 1, len(received))

		if diff := deep.Equal(task2, received[0]); diff != nil {
			t.Error(diff)
		}
	})
}
