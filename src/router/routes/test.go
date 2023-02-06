package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) Test(c *gin.Context) {
	//url := os.Getenv("domain") + os.Getenv("path_key") + "/bot"
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
