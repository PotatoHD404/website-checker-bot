package routes

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/tucnak/telebot.v2"
	"net/http"
	"os"
)

func initBot(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Set webhook
	wg.Add(2)

	go setWebhook()
	go setupDb()

	wg.Wait()
	returnOk(w)
}

func setWebhook() {
	defer wg.Done()
	err := tgBot.SetWebhook(&telebot.Webhook{
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: os.Getenv("domain") + "/" + os.Getenv("path_key") + "/bot",
		},
	})
	if err != nil {
		fmt.Println(err)
		panic("can't set webhook")
	}
}

func setupDb() {
	defer wg.Done()
	tables, err := listTables()
	if err != nil {
		panic("can't list tables")
	}

	if !contains(tables, subscribersTable) {
		wg.Add(1)
		go createSubscribersTable()
	}

	if !contains(tables, adminsTable) {
		wg.Add(1)
		go createAdminsTable()
	}
}
