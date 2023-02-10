package routes

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (env *Env) CheckWebsites(c *gin.Context) {

	websites := env.db.GetWebsites(false)
	fmt.Println(websites)
	for _, website := range websites {
		env.pool.AddTask(func() {
			changed, err := website.CheckChanged()
			if err != nil {
				fmt.Println("Error checking website. Here is why: ", err)
				panic(err)
			}
			if !changed {
				return
			}
			newWebsite := env.db.GetWebsite(website.Name)

			newWebsite.Hash = website.Hash
			env.db.UpdateWebsite(newWebsite)
			for _, subscriber := range newWebsite.Subscribers {
				_, err := env.bot.Send(&telebot.User{ID: subscriber}, "Website "+newWebsite.Name+" changed!")
				if err != nil {
					fmt.Println("Error sending message to subscriber. Here is why: ", err)
					panic(err)
				}
			}
		})
	}
	env.pool.Wait()
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
