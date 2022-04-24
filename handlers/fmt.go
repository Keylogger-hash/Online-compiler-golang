package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"compiler.com/utils"
)

type Response struct {
	Res   string
	Body string
	Error string
}

func HandleFmt(wr http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		rsp := &Response{}
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(wr, fmt.Sprintf("%v", err), http.StatusBadRequest)
		}
		json.Unmarshal(buf, rsp)
		code, err := utils.FormatFmt(rsp.Res)
		wr.Header().Add("Content-type", "application/json")
		if err != nil {
			rsp.Error = fmt.Sprintf("%v", err)
		} else {
			rsp.Res = string(code)
			rsp.Error = ""
		}
		ans, err := json.Marshal(rsp)
		fmt.Println(ans)
		if err != nil {
			http.Error(wr, fmt.Sprintf("%v", err), http.StatusBadRequest)
		}
		wr.Write(ans)
	default:
		wr.WriteHeader(http.StatusMethodNotAllowed)
	}

}
