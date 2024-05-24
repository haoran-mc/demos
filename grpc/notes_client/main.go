package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	notes "github.com/haoran-mc/demos/grpc/notes_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := notes.NewNotesClient(conn)

	// contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// define expected flag for save
	saveCmd := flag.NewFlagSet("save", flag.ExitOnError)
	saveTitle := saveCmd.String("title", "", "Give a title to your note")
	saveBody := saveCmd.String("content", "", "Type what you like to remember")

	// define expected flags for load
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)
	loadKeyword := loadCmd.String("keyword", "", "A keyworkd you'd like to find in your notes")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'save' or 'load' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "save":
		saveCmd.Parse(os.Args[2:])
		_, err := client.Save(ctx, &notes.Note{
			Title: *saveTitle,
			Body:  []byte(*saveBody),
		})
		if err != nil {
			log.Fatalf("The note could not be saved: %v", err)
		}
		fmt.Printf("Your note was saved: %v\n", *saveTitle)

	case "load":
		loadCmd.Parse(os.Args[2:])
		note, err := client.Load(ctx, &notes.NoteSearch{
			Keyword: *loadKeyword,
		})
		if err != nil {
			log.Fatalf("The note could not be loaded: %v", err)
		}
		fmt.Printf("%v\n", note)

	default:
		fmt.Println("Expected 'save' or 'load' subcommands")
		os.Exit(1)
	}
}
