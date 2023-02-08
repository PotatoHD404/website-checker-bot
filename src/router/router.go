package router

import (
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gin-gonic/gin"
	"os"
	"website-checker-bot/router/routes"
)

func GetAdapter() *httpadapter.HandlerAdapterV2 {
	r := gin.Default()
	env := routes.GetEnv()
	r.GET("/", env.CheckWebsites)
	prefix := r.Group("/" + os.Getenv("path_key"))
	{
		prefix.GET("/init-bot", env.InitBot)
		prefix.POST("/bot", env.HandleMessage)
		prefix.GET("/test", env.Test)
	}
	//r. = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	log.Println("Not found", r.RequestURI)
	//	http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	//})
	return httpadapter.NewV2(r)
}
