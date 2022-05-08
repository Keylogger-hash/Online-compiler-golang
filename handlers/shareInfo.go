package handlers

import (
	"fmt"
	"net/http"

	"compiler.com/storage"
	"github.com/gorilla/mux"
)

func HandleShareInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		fmt.Println(vars)
		id := vars["id"]
		fmt.Println(id)
		client := storage.NewRedisClient()
		
		content, err := storage.GetRedisValue(client, id)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Header().Add("Content-type", "text/plain")
			return
		}
		w.Write([]byte(content))
		w.Header().Add("Content-type", "text/plain")
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
