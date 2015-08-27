package minishop

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Validate(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

}

func AddProduct(c *gin.Context) {
	var p Product

	if c.BindJSON(&p) == nil {

	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
