package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
)

// Example usage:
//
//	go run ./examples/multi_edit path/to/save.sii
//
// This demonstrates how to:
//   - read a SII file into a generic Document
//   - print the document to the console
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s path/to/save.sii", os.Args[0])
	}

	path := os.Args[1]

	plain, err := siidecrypt.DecryptFile(path, true)
	if err != nil {
		log.Fatalf("decrypt %s: %v", path, err)
	}

	doc, err := sii.ReadDocument(plain)
	if err != nil {
		log.Fatalf("parse SII: %v", err)
	}

	fmt.Println(doc.DebugString())
}
