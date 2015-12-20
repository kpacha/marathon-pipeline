package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kpacha/marathon-pipeline/pipeline"
)

type GinServer struct {
	Port   int
	engine *gin.Engine
	jobs   map[string]pipeline.Job
}

func Default(port int) GinServer {
	return GinServer{port, gin.Default(), map[string]pipeline.Job{}}
}

func (s *GinServer) Run() {
	s.engine.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, s.jobs)
	})
	s.engine.POST("", func(c *gin.Context) {
		var job pipeline.Job
		err := c.Bind(&job)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"status": fmt.Sprintf("Error processing the body of the request: %s", err.Error())})
			return
		}
		if job.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error: Job Id must be set!"})
			return
		}
		if _, ok := s.jobs[job.ID]; ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Job already exists: %s", job.ID)})
			return
		}
		s.jobs[job.ID] = job
		c.JSON(http.StatusOK, s.jobs[job.ID])
	})
	s.engine.GET("/:jobId", func(c *gin.Context) {
		jobId := c.Params.ByName("jobId")
		if j, ok := s.jobs[jobId]; ok {
			c.JSON(http.StatusOK, j)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": fmt.Sprintf("Job not found: %s", jobId)})
		}
	})
	s.engine.DELETE("/:jobId", func(c *gin.Context) {
		jobId := c.Params.ByName("jobId")
		if _, ok := s.jobs[jobId]; ok {
			delete(s.jobs, jobId)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": fmt.Sprintf("Job not found: %s", jobId)})
		}
	})
	s.engine.Run(fmt.Sprintf(":%d", s.Port))
}
