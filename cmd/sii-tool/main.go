package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
	"github.com/robebs/ts-se-tool-go/internal/sii"
)

func main() {
	inPath := flag.String("in", "", "input SII file path")
	outPath := flag.String("out", "", "output SII file path (optional)")
	flag.Parse()

	if *inPath == "" {
		log.Fatal("missing -in path to SII file")
	}

	// Déchiffre/décode le fichier si nécessaire (supporte les formats encryptés ETS2),
	// puis parse le texte SII résultant.
	plain, err := siidecrypt.DecryptFile(*inPath, true)
	if err != nil {
		log.Fatalf("decrypt SII: %v", err)
	}

	doc, err := sii.ReadDocument(plain)
	if err != nil {
		log.Fatalf("parse SII: %v", err)
	}

	// Example integration with typed save-game classes: once parsing is wired,
	// this will surface structured bank loans, garages, etc.
	bankLoans := doc.ToBankLoans()
	fmt.Printf("parsed %d bank_loan blocks (typed view placeholder)\n\n", len(bankLoans))

	if *outPath == "" {
		fmt.Println(doc.DebugString())
		return
	}

	outBytes, err := sii.WriteDocument(doc)
	if err != nil {
		log.Fatalf("serialize SII: %v", err)
	}

	if err := os.WriteFile(*outPath, outBytes, 0o644); err != nil {
		log.Fatalf("write output: %v", err)
	}
}
