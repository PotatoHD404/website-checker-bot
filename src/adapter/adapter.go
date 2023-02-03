package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"os"
	telebot2 "website-checker-bot/telebot"
)

var adapter *httpadapter.HandlerAdapterV2

func initRouter() {
	defer main.Wg.Done()
	r := httprouter.New()
	r.GET("/init-bot", initBot)
	r.POST("/bot", handleMessage)
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Not found", r.RequestURI)
		http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	})
	adapter = httpadapter.NewV2(r)
}

func handleMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Process update
	var u telebot.Update
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println(err)
		panic("can't unmarshal")
	}
	telebot2.tgBot.ProcessUpdate(u)

	//returnOk(w)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(os.Getenv("domain") + "/" + os.Getenv("path_key") + "/bot"))
	if err != nil {
		fmt.Println(err)
		panic("can't write response")
	}
}
