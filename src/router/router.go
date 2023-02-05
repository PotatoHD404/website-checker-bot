package router

import (
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gin-gonic/gin"
	"website-checker-bot/router/routes"
)

func GetAdapter() *httpadapter.HandlerAdapterV2 {
	r := gin.Default()
	env := routes.GetEnv()
	r.GET("/init-bot", env.InitBot)
	r.POST("/bot", env.HandleMessage)
	//r. = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	log.Println("Not found", r.RequestURI)
	//	http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	//})
	return httpadapter.NewV2(r)
}
