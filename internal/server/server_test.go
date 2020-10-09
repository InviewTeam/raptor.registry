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
	"gitlab.com/inview-team/raptor_team/registry/task"
	"gitlab.com/inview-team/raptor_team/registry/tests"
)

type Response struct {
	UUID string `json:"uuid"`
}

var (
	srv = Server{
		reg: registry.New(tests.New()),
	}
	router = srv.setupRouter()

	task1 = task.Task{
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

	task2 = task.Task{
		CameraIP: "92.138.141.54",
		Jobs: []task.Job{
			{
				Title: "some_job_name",
			},
		},
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

		w := performRequest(router, "POST", "/create", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)

		resp := Response{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.Nil(t, err)
		id, err := uuid.Parse(resp.UUID)
		require.Nil(t, err)
		task1.UUID = id

		body, err = json.Marshal(task2)
		require.Nil(t, err)

		w = performRequest(router, "POST", "/create", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.Nil(t, err)
		id, err = uuid.Parse(resp.UUID)
		require.Nil(t, err)
		task2.UUID = id
	})
}

func TestGet(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		w := performRequest(router, "GET", "/tasks", nil)
		require.Equal(t, http.StatusOK, w.Code)

		received := []task.Task{}
		err := json.Unmarshal(w.Body.Bytes(), &received)
		require.Nil(t, err)

		require.Equal(t, 2, len(received))

		if diff := deep.Equal(task1, received[0]); diff != nil {
			t.Error(diff)
		}

		if diff := deep.Equal(task2, received[1]); diff != nil {
			t.Error(diff)
		}
	})
}
