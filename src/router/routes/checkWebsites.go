package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (env *Env) CheckWebsites(c *gin.Context) {

	websites := env.db.GetWebsites(false)
	fmt.Println(websites)
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
