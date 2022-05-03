package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	s "compiler.com/handlers/handlers_struct"
	pb "compiler.com/sandboxproto"
	"compiler.com/storage"

	// s "compiler.com/handlers/handlers-struct/struct"
	"compiler.com/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
		req := &s.Request{}
		resp := &s.Response{}
		json.Unmarshal(body, req)
		code, err := utils.FormatFmt(req.Body)
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
			resp.Res = ""
			resp.Error = fmt.Sprintf("package not main")
			outputMainJson, _ := json.Marshal(resp)
			wr.Write(outputMainJson)
			return
		}
		compilePath, err := utils.WriteCodeFile(code)
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = err.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			wr.Write(outputErrorJSON)
			return
		}
		data, err := utils.CompileCode(compilePath)
		
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = err.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			wr.Write(outputErrorJSON)
			return
		}
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = err.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			wr.Write(outputErrorJSON)
			return
		}
		
		if err != nil {
			wr.Header().Add("Content-type", "application/json")
			resp.Res = ""
			resp.Error = err.Error()
			outputErrorJSON, _ := json.Marshal(resp)
			wr.Write(outputErrorJSON)
			return
		}
		key, err := utils.HashContent(code)
		sha256key := hex.EncodeToString(key)
		memcachedClient, err := storage.NewMemcachedClient()
		content, err := storage.MemcachedGetValue(memcachedClient, sha256key)
		fmt.Println(err)
		if err == nil {
			wr.Header().Add("Content-type", "application/json")
			wr.Write(content)
			return
		} else {
			buf := utils.EncodeToBase64(data)

			client := NewClientGrpc()
			ctx := context.Background()
			respMsg, _ := client.RunSandboxCompileCode(ctx, &pb.RequestMessage{Body: string(buf)})
			
			resp.Res = respMsg.Res
			resp.Body = string(code)
			resp.Error = respMsg.Error
			
			outputJson, err := json.Marshal(resp)
			storage.MemcachedAddValue(memcachedClient,sha256key,outputJson)
			if err != nil {
				wr.Header().Add("Content-type", "application/json")
				resp.Res = ""
				resp.Error = err.Error()
				outputErrorJSON, _ := json.Marshal(resp)
				wr.Write(outputErrorJSON)
			}
			wr.Header().Add("Content-type", "application/json")
			wr.Write(outputJson)
		}

	}
}
