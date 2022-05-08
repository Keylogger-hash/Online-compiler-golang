package handlers

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	h "compiler.com/handlers/handlers_struct"
	"compiler.com/storage"
	"compiler.com/utils"
	"github.com/go-redis/redis/v8"
)

func HandleShare(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req := &h.Request{}
		resp := &h.Response{}
		body, err := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, req)
		if err != nil {
			resp.Body = ""
			resp.Error = err.Error()
			resp.Res = ""
			output, _ := json.Marshal(resp)
			w.Header().Add("Content-type", "application/json")
			w.Write(output)
			return
		}
		key, err := utils.HashContent([]byte(req.Body))
		sha256key := hex.EncodeToString(key)
		// key := utils.GenerateUID()
		client := storage.NewRedisClient()
		_, err = storage.GetRedisValue(client, sha256key)
		if err != redis.Nil {
			resp.Body = ""
			resp.Error = ""
			resp.Res = sha256key
			output, _ := json.Marshal(resp)
			w.Header().Add("Content-type", "application/json")
			w.Write(output)
			return
		} 
		err = storage.AddRedisValue(client,sha256key,[]byte(req.Body))
		if err != nil {
			resp.Body = ""
			resp.Error = err.Error()
			resp.Res = ""
			output, _ := json.Marshal(resp)
			w.Header().Add("Content-type", "application/json")
			w.Write(output)
			return
		}
		resp.Body = ""
		resp.Error = ""
		resp.Res = sha256key
		output, err := json.Marshal(resp)
		w.Header().Add("Content-type", "application/json")
		w.Write(output)
		return

	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
