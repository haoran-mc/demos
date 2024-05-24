package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	notes "github.com/haoran-mc/demos/grpc/notes_proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "the server port")
)

// 实现 the notes service (notes.NotesServer interface) ../notes_proto/notes_grpc.pb.go:55
type notesServer struct {
	notes.UnimplementedNotesServer
}

// notes.NotesServer 有下面两个方法 Save、Load

func (s *notesServer) Save(ctx context.Context, n *notes.Note) (*notes.NoteSaveReply, error) {
	log.Printf("Recieved a note to save: %v", n.Title)
	err := notes.SaveToDisk(n, "notes_data") // 保存在文件中
	if err != nil {
		return &notes.NoteSaveReply{Saved: false}, err
	}

	return &notes.NoteSaveReply{Saved: true}, nil
}

func (s *notesServer) Load(ctx context.Context, search *notes.NoteSearch) (*notes.Note, error) {
	log.Printf("Recieved a note to save: %v", search.Keyword)
	n, err := notes.LoadFromDisk(search.Keyword, "notes_data") // 从文件中加载
	if err != nil {
		return &notes.Note{}, err
	}

	return n, nil
}

func main() {
	flag.Parse() // parse arguments from the command line
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Panicf("start listen error: %v", err)
	}

	// Instanciate the server
	s := grpc.NewServer()

	// Register server method (actions the server il do)
	notes.RegisterNotesServer(s, &notesServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
