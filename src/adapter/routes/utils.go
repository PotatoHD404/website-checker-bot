package routes

import (
	"fmt"
	"net/http"
)

func contains(tables []string, table string) bool {
	for _, t := range tables {
		if t == table {
			return true
		}
	}
	return false
}

func returnOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		fmt.Println(err)
		panic("can't write response")
	}
}
