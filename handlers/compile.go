package handlers

import (
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

		
		wr.Header().Add("Content-type", "application/json")
		wr.Write(outputJson)
	}
}
