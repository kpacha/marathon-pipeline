package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kpacha/marathon-pipeline/pipeline"
)

type GinServer struct {
	Port      int
	engine    *gin.Engine
	taskStore pipeline.TaskStore
}

func Default(port int, taskStore pipeline.TaskStore) GinServer {
	return GinServer{port, gin.Default(), taskStore}
}

func (s *GinServer) Run() {
	s.engine.GET("", func(c *gin.Context) {
		tasks, err := s.taskStore.GetAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Error getting the tasks: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, tasks)
	})
	s.engine.POST("", func(c *gin.Context) {
		var task pipeline.Task
		err := c.Bind(&task)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"status": fmt.Sprintf("Error processing the body of the request: %s", err.Error())})
			return
		}
		if task.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error: Task Id must be set!"})
			return
		}
		if _, ok, _ := s.taskStore.Get(task.ID); ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Task already exists: %s", task.ID)})
			return
		}
		if err := s.taskStore.Set(task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Error storing the task: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, task)
	})
	s.engine.GET("/:taskId", func(c *gin.Context) {
		taskId := c.Params.ByName("taskId")
		if t, ok, _ := s.taskStore.Get(taskId); ok {
			c.JSON(http.StatusOK, t)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": fmt.Sprintf("Task not found: %s", taskId)})
		}
	})
	s.engine.DELETE("/:taskId", func(c *gin.Context) {
		taskId := c.Params.ByName("taskId")
		if _, ok, _ := s.taskStore.Get(taskId); !ok {
			c.JSON(http.StatusNotFound, gin.H{"status": fmt.Sprintf("Task not found: %s", taskId)})
			return
		}
		if err := s.taskStore.Delete(taskId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Error deleting the task: %s", err.Error())})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	s.engine.Run(fmt.Sprintf(":%d", s.Port))
}
