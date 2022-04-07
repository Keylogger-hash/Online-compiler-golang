package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"compiler.com/utils"
)

type Request struct {
	Body string
}

// Отправляет запрос на sandboxq
func HandleCompile(wr http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		req := &Request{}
		resp := &Response{}
		json.Unmarshal(body, req)
		output, err := utils.formatFmt(req.Body)
		fmt.Println(output)
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = err.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			fmt.Println(outputErrorJSON)
			wr.Write(outputErrorJSON)
			return
		}
		req.Body = string(output)
		inputJSON, _ := json.Marshal(req)
		r := bytes.NewReader(inputJSON)
		client := http.Client{}
		request, _ := http.NewRequest("POST", "http://localhost:8081/run", r)
		response, err := client.Do(request)
		fmt.Println(response)
		if err != nil {
			fmt.Println(err)
		}
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		wr.Header().Add("Content-type", "application/json")
		wr.Write(b)
	}
}
