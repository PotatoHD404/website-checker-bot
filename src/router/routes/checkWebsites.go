package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) CheckWebsites(c *gin.Context) {

	websites := env.db.GetWebsites(false)

	for _, website := range websites {
		env.pool.AddTask(func() {

			if website.CheckChanged() {

				newWebsite := env.db.GetWebsite(website.Name)

				newWebsite.Hash = website.Hash
				env.db.UpdateWebsite(newWebsite)
			}
		})
	}
	env.pool.Wait()
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
