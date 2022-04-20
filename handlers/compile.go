package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "compiler.com/sandboxproto"
	"compiler.com/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Request struct {
	Body string
}
type Client struct {
}

func NewClientGrpc() pb.RunSandboxCompileCodeClient {
	conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Can't dial localhost:8082")
	}
	client := pb.NewRunSandboxCompileCodeClient(conn)
	return client

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
		code := build.Data
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
		fmt.Println(build)
		if build.Errors != nil  {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = build.Errors.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			fmt.Println(string(outputErrorJSON))
			wr.Write(outputErrorJSON)
			return
		}
		buf, _ := utils.EncodeBinaryFile(build)

		client := NewClientGrpc()
		ctx := context.Background()
		respMsg, err := client.RunSandboxCompileCode(ctx, &pb.RequestMessage{Body: buf.String()})

		if err != nil {
			http.Error(wr, "Server error", http.StatusInternalServerError)
		}
		resp.Res = respMsg.Res
		resp.Body = string(code)
		if resp.Error != "" {
			resp.Error = respMsg.Error
		}
		outputJson, err := json.Marshal(resp)
		wr.Header().Add("Content-type", "application/json")
		wr.Write(outputJson)
	}
}
