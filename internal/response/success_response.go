package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool        `json:"success"`
	Response  interface{} `json:"response,omitempty"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Response:  data,
		Message:   "Request successful",
		Timestamp: time.Now().Unix(),
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success:   true,
		Response:  data,
		Message:   "Resource created successfully",
		Timestamp: time.Now().Unix(),
	})
}

func Deleted(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Message:   "Resource deleted successfully",
		Timestamp: time.Now().Unix(),
	})
}

func Updated(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Message:   "Resource updated successfully",
		Timestamp: time.Now().Unix(),
	})
}
