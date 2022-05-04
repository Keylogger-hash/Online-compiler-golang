package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	h "compiler.com/handlers/handlers_struct"
	"compiler.com/storage"
	"compiler.com/utils"
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
		key := utils.GenerateUID()
		client, err := storage.NewMemcachedClient()
		if err != nil {
			resp.Body = ""
			resp.Error = err.Error()
			resp.Res = ""
			output, _ := json.Marshal(resp)
			w.Header().Add("Content-type", "application/json")
			w.Write(output)
			return
		}
		err = storage.MemcachedAddValue(client, key, []byte(resp.Body))
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
		resp.Res = key
		output, err := json.Marshal(resp)
		w.Header().Add("Content-type", "application/json")
		w.Write(output)
		return
	case "GET":
		resp := &h.Response{}
		client, err := storage.NewMemcachedClient()
		if err != nil {
			resp.Body = ""
			resp.Error = err.Error()
			resp.Res = ""
			output, _ := json.Marshal(resp)
			w.Header().Add("Content-type", "application/json")
			w.Write(output)
			return
		}
		storage.MemcachedGetValue()
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
