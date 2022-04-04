package handlers

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"net/http"

	"golang.org/x/tools/imports"
)

type Response struct {
	Res  string
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
		json.Unmarshal(buf,rsp)
		formatBody, err1  := formatFmt(rsp.Res)
		wr.Header().Add("Content-type","application/json")
		rsp.Res = string(formatBody)
		if err1 != nil {
			rsp.Error = fmt.Sprintf("%v",err1)
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

func formatFmt(body string) ([]byte, error) {
	dest, err := format.Source([]byte(body))
	if err != nil {
		return nil, err
	}
	finishImports, err := imports.Process("", dest, nil)
	if err != nil {
		return nil, err
	}
	return finishImports, nil
}
