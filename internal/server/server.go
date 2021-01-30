package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/pkg/format"
)

type Server struct {
	http *http.Server
	reg  *registry.Registry
}

func New(addr string, reg *registry.Registry) *Server {
	var server = &Server{
		http: &http.Server{
			Addr: addr,
		},
		reg: reg,
	}
	server.http.Handler = server.setupRouter()

	return server
}

func (s *Server) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/tasks/create", s.createNewTask)
	r.GET("/tasks/get", s.getTasks)
	r.DELETE("/tasks/delete/:uuid", s.deleteTask)
	r.DELETE("/tasks/stop/:uuid", s.stopTask)
	r.GET("/tasks/get/:uuid", s.getTaskByUUID)

	r.POST("/analyzers/create", s.createAnalyzer)
	r.GET("/analyzers/get", s.getAnalyzers)
	r.DELETE("/analyzers/delete/:name", s.deleteAnalyzer)
	r.GET("/analyzers/get/:name", s.getAnalyzerByName)

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

	task := format.Task{}
	err = json.Unmarshal(bodyBytes, &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err.Error()})
		return
	}

	id, err := s.reg.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to create task": err.Error()})
		return
	}

	err = s.reg.SendTask(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to send task": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"uuid": id})
}

func (s *Server) deleteTask(c *gin.Context) {
	id := c.Param("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse task UUID": err.Error()})
		return
	}

	err = s.reg.DeleteTask(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to delete task": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uuid": id})
}

func (s *Server) stopTask(c *gin.Context) {
	id := c.Param("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse task UUID": err.Error()})
		return
	}
	err = s.reg.StopTask(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to delete task": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uuid": id})
}

func (s *Server) getTaskByUUID(c *gin.Context) {
	id := c.Param("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse task UUID": err.Error()})
		return
	}
	task, err := s.reg.GetTaskByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to delete task": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (s *Server) getTasks(c *gin.Context) {
	tasks, err := s.reg.GetTasks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to get tasks": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (s *Server) getAnalyzers(c *gin.Context) {
	analyzers, err := s.reg.GetAnalyzers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to get analyzers": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analyzers)
}

func (s *Server) getAnalyzerByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "received empty analyzer name"})
		return
	}

	analyzer, err := s.reg.GetAnalyzerByName(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to get analyzer": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analyzer)
}

func (s *Server) createAnalyzer(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err.Error()})
		return
	}

	analyzer := format.Analyzer{}
	err = json.Unmarshal(bodyBytes, &analyzer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err.Error()})
		return
	}

	err = s.reg.CreateAnalyzer(analyzer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to create analyzer": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) deleteAnalyzer(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "received empty analyzer name"})
		return
	}

	err := s.reg.DeleteAnalyzer(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to delete analyzer": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
