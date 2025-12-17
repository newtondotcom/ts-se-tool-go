package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robebs/ts-se-tool-go/internal/save/items"
	"github.com/robebs/ts-se-tool-go/internal/sii"
	"github.com/robebs/ts-se-tool-go/internal/siidecrypt"
)

// Example usage:
//
//	go run ./examples/multi_edit path/to/save.sii
//
// This demonstrates how to:
//   - read a (possibly encrypted) SII file into a generic Document
//   - load the bank block into a typed Bank struct
//   - edit the money_account and write the file back
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s path/to/save.sii", os.Args[0])
	}

	path := os.Args[1]

	// Step 1: Decrypt/decode the file (handles encrypted, BSII, or plaintext)
	fmt.Println("Step 1: Decrypting/decoding SII file...")
	plain, err := siidecrypt.DecryptFile(path, true)
	if err != nil {
		log.Fatalf("decrypt %s: %v", path, err)
	}

	// Debug: show first 500 bytes of decoded text
	if len(plain) > 500 {
		fmt.Printf("Decoded text preview (first 500 bytes):\n%s\n...\n", string(plain[:500]))
	} else {
		fmt.Printf("Decoded text:\n%s\n", string(plain))
	}

	// Step 2: Parse the text SII into a Document structure
	fmt.Println("\nStep 2: Parsing SII text...")
	doc, err := sii.ReadDocument(plain)
	if err != nil {
		log.Fatalf("parse SII: %v", err)
	}

	fmt.Printf("document has %d blocks\n", len(doc.Blocks))

	// Find the bank block.
	var bankBlock *sii.Block
	for i := range doc.Blocks {
		if doc.Blocks[i].Type == "bank" {
			bankBlock = &doc.Blocks[i]
			break
		}
	}
	if bankBlock == nil {
		log.Fatalf("no bank block found in %s", path)
	}

	// Load into typed Bank.
	var bank items.Bank
	if err := bank.FromProperties(bankBlock.Properties); err != nil {
		log.Fatalf("load bank: %v", err)
	}

	fmt.Printf("current money_account: %d\n", bank.MoneyAccount)

	// Example edit: add 100000 units of money.
	bank.MoneyAccount += 100000

	fmt.Printf("new money_account: %d\n", bank.MoneyAccount)

	// Write modified properties back into the block.
	bankBlock.Properties = bank.ToProperties()

	var economy items.Economy
	economy.FromProperties(doc.Blocks[0].Properties)
	if err := economy.FromProperties(doc.Blocks[0].Properties); err != nil {
		log.Fatalf("load economy: %v", err)
	}

	// Serialize and overwrite the file.
	/*
		out, err := sii.WriteDocument(doc)
		if err != nil {
			log.Fatalf("serialize SII: %v", err)
		}

		if err := os.WriteFile(path, out, 0o644); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}

		fmt.Println("file updated successfully")
	*/
}
