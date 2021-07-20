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
	r.GET("/api/tasks", s.getTasks)
	r.POST("/api/tasks", s.createNewTask)
	r.GET("/api/tasks/:uuid", s.getTaskByUUID)
	r.GET("/api/tasks/:uuid/report", s.getReportByUUID)
	r.POST("/api/tasks/:uuid/report", s.addReport)
	r.PATCH("/api/tasks/:uuid", s.stopTask)
	r.DELETE("/api/tasks/:uuid", s.deleteTask)

	r.GET("/api/analyzers", s.getAnalyzers)
	r.POST("/api/analyzers", s.createAnalyzer)
	r.GET("/api/analyzers/:uuid", s.getTaskByUUID)
	r.DELETE("/api/analyzers/:uuid", s.deleteAnalyzer)
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
		c.JSON(http.StatusBadRequest, gin.H{"failed to stop task": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"failed to get task": err.Error()})
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

func (s *Server) addReport(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err.Error()})
		return
	}

	report := format.Report{}
	err = json.Unmarshal(bodyBytes, &report)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err.Error()})
		return
	}

	err = s.reg.AddReport(report)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to add report": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) getReportByUUID(c *gin.Context) {
	id := c.Param("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse task UUID": err.Error()})
		return
	}

	rep, err := s.reg.GetReport(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to get report": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rep)
}
