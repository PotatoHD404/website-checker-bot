package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (env *Env) Test(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": os.Getenv("domain") + "/" + os.Getenv("path_key") + "/bot",
	})
}
