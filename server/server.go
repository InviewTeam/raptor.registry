package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/task"
)

type Server struct {
	http *http.Server
}

func New(addr string) *Server {
	var server = &Server{
		http: &http.Server{
			Addr: addr,
		},
	}
	server.http.Handler = server.setupRouter()

	return server
}

func (s *Server) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/newtask", s.createNewTask)

	return r
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.http.Close()
}

func (s *Server) createNewTask(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err.Error()})
		return
	}

	task := task.Task{}
	err = json.Unmarshal(bodyBytes, &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": uuid.New()})
}
