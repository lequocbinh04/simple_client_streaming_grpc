package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	streaming "simple_client_streaming_grpc/server/proto/steaming"
	"sync"
	"time"
)

const maxDataSize = 1 << 30 // 1GB

type server struct {
	streaming.UnimplementedDataStreamingServer
	store *DiskStore
}

func NewServer() *server {
	d := NewDiskStore("./upload_data")
	return &server{
		store: d,
	}
}

func (server *server) UploadData(stream streaming.DataStreaming_UploadDataServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive data info")
	}

	fmt.Println("one connection established at: ", time.Now().String())

	filename := req.GetInfo().FileName
	filetype := req.GetInfo().FileType
	fileid := req.GetInfo().FileId

	dataBuf := bytes.Buffer{}
	dataSize := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("no more data")
			break
		}
		if err != nil {
			fmt.Println("cannot receive data")
			break
		}
		chunk := req.GetChunkData()
		size := len(chunk)
		log.Printf("received a chunk with size: %d", size)

		dataSize += size
		if dataSize > maxDataSize {
			return logError(status.Errorf(codes.InvalidArgument, "data is too large: %d > %d", dataSize, maxDataSize))
		}

		_, err = dataBuf.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk to buffer: %v", err))
		}
	}

	id, err := server.store.Save(fileid, filename, filetype, dataBuf)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save: %v", err))
	}
	log.Printf("saved data with id: %s, size: %d", id, dataSize)
	return nil
}

func logError(err error) error {
	log.Println(err)
	return err
}

type DataStore interface {
	Save(dataID string, dataName string, dataType string, dataBuf bytes.Buffer) (string, error)
}

type DiskStore struct {
	mutex      sync.RWMutex
	saveFolder string
	data       map[string]*streaming.DataInfo
}

type DataInfo struct {
	DataID string
	Type   string
	Path   string
}

func NewDiskStore(saveFolder string) *DiskStore {
	if _, err := os.Stat(saveFolder); os.IsNotExist(err) {
		os.Mkdir(saveFolder, 0755)
	}
	return &DiskStore{
		saveFolder: saveFolder,
		data:       make(map[string]*streaming.DataInfo),
	}
}

func (store *DiskStore) Save(dataID string, dataName string, dataType string, dataBuf bytes.Buffer) (string, error) {
	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate data id: %w", err)
	}

	dataPath := fmt.Sprintf("%s/%s%s", store.saveFolder, fileID, dataType)

	file, err := os.Create(dataPath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	_, err = dataBuf.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write to file: %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.data[fileID.String()] = &streaming.DataInfo{
		FileId:   dataID,
		FileType: dataType,
		FileName: dataPath,
	}

	return fileID.String(), nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	s := grpc.NewServer()
	server := NewServer()
	streaming.RegisterDataStreamingServer(s, server)
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatal(s.Serve(lis))
}
