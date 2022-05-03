package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	s "compiler.com/handlers/handlers_struct"
	"compiler.com/utils"
)



func HandleFmt(wr http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		req := &s.Request{}
		resp := &s.Response{}
		json.Unmarshal(body, req)
		code, err := utils.FormatFmt(req.Body)
		fmt.Println(string(code))
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = err.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			wr.Write(outputErrorJSON)
			return
		}
		err = utils.CheckCodePackageIsMain(code)
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = fmt.Sprintf("package not main")
			outputMainJson, _ := json.Marshal(resp)
			wr.Write(outputMainJson)
			return
		}
		wr.Header().Add("Content-type", "application/json")
		resp.Error = ""
		resp.Res = string(code)
		ans, err := json.Marshal(resp)
		if err != nil {
			http.Error(wr, fmt.Sprintf("%v", err), http.StatusBadRequest)
		}
		wr.Write(ans)
	default:
		wr.WriteHeader(http.StatusMethodNotAllowed)
	}

}
