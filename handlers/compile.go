package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"compiler.com/utils"
)

type Request struct {
	Body string
}
type Client struct {
	
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
		build := utils.NewBuildResult()
		utils.FormatFmt(build, req.Body)
		fmt.Println(string(build.Data))
		if build.Errors != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = build.Errors.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			fmt.Println(string(outputErrorJSON))
			wr.Write(outputErrorJSON)
			return
		}
		utils.CheckCodePackageIsMain(build)
		if build.Errors != nil {
			resp.Res = ""
			resp.Error = fmt.Sprintf("package not main")
			outputMainJson, _ := json.Marshal(resp)
			fmt.Println(string(outputMainJson))
			wr.Write(outputMainJson)
			return
		}
		utils.WriteCodeFile(build)
		if build.Errors != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = build.Errors.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			fmt.Println(outputErrorJSON)
			wr.Write(outputErrorJSON)
			return
		}
		utils.CompileCode(build)
		if build.Errors != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = build.Errors.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			fmt.Println(string(outputErrorJSON))
			wr.Write(outputErrorJSON)
			return
		}
		buf, _ := utils.EncodeBinaryFile(build)
		resp.Error = ""
		resp.Res = string(buf.Bytes())
		fmt.Println(string(buf.Bytes()))
		outputJson, _ := json.Marshal(resp)
		postBody := bytes.NewBuffer(outputJson)
		client := http.Client{Timeout: 15*time.Second}
		request, err := http.NewRequest("POST","http://localhost:8081",postBody)
		responseCompile, err := client.Do(request)
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = err.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			fmt.Println(string(outputErrorJSON))
			wr.Write(outputErrorJSON)
		}
		
		wr.Header().Add("Content-type", "application/json")
		wr.Write([]byte("ok"))
	}
}
