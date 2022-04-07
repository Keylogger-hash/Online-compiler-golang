package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Res   string
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
		formatBody, err1 := utils.formatFmt(rsp.Res)
		wr.Header().Add("Content-type", "application/json")
		rsp.Res = string(formatBody)
		if err1 != nil {
			rsp.Error = fmt.Sprintf("%v", err1)
		} else {
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
