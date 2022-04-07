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
		build := utils.NewBuildResult()
		utils.FormatFmt(build,rsp.Res)
		wr.Header().Add("Content-type", "application/json")
		rsp.Res = string(build.Data)
		if build.Errors != nil {
			rsp.Error = fmt.Sprintf("%v", build.Errors)
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
