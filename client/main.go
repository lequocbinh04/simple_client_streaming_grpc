package main

import (
	"bufio"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"path/filepath"
	streaming "simple_client_streaming_grpc/server/proto/steaming"
	"time"
)

const (
	remoteAddr   = "127.0.0.1:8080"
	resourceName = "test.png"
)

func main() {
	conn, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	file, err := os.Open(resourceName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := streaming.NewDataStreamingClient(conn)
	stream, err := client.UploadData(ctx)
	if err != nil {
		panic(err)
	}
	req := &streaming.UploadDataRequest{
		Data: &streaming.UploadDataRequest_Info{
			Info: &streaming.DataInfo{
				FileId:   "binh",
				FileName: resourceName,
				FileType: filepath.Ext(resourceName),
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &streaming.UploadDataRequest{
			Data: &streaming.UploadDataRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}

}
