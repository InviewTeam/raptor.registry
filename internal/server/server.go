package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/task"
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
	r.POST("/create", s.createNewTask)
	r.GET("/tasks", s.getTasks)
	r.DELETE("/delete/:uuid", s.deleteTask)
	r.DELETE("/stop/:uuid", s.stopTask)
	r.GET("/tasks/:uuid", s.getTaskByUUID)

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

	id, err := s.reg.CreateTask(&task)
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
