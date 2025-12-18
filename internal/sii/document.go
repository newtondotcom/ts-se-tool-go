package sii

import (
	"strconv"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
	"github.com/robebs/ts-se-tool-go/internal/save/items"
)

// Document is a very generic representation of an SII file.
// It will be progressively specialized to mirror the original C# CustomClasses.
type Document struct {
	Blocks []Block
}

// Block represents a single SiiNunit block, such as:
// some_type : some_name { key: value }
type Block struct {
	Type       string
	Name       string
	Properties map[string][]string
	// PropertyOrder stores the order in which properties were read from the file
	PropertyOrder []string
}

// IndexByName builds a simple name -> Block map over all blocks in the
// document. It is the Go equivalent of the C# SiiNunit.SiiNitems dictionary.
func (d *Document) IndexByName() map[string]Block {
	out := make(map[string]Block, len(d.Blocks))
	for _, b := range d.Blocks {
		out[b.Name] = b
	}
	return out
}

// ToBankLoans parcourt les blocs et construit des BankLoan typés à partir des
// blocs de type "bank_loan".
func (d *Document) ToBankLoans() []items.BankLoan {
	var out []items.BankLoan

	for _, b := range d.Blocks {
		if b.Type != "bank_loan" {
			continue
		}

		var loan items.BankLoan

		if vals, ok := b.Properties["amount"]; ok && len(vals) > 0 {
			if v, err := strconv.Atoi(vals[0]); err == nil {
				loan.Amount = v
			}
		}

		if vals, ok := b.Properties["original_amount"]; ok && len(vals) > 0 {
			if v, err := strconv.Atoi(vals[0]); err == nil {
				loan.OriginalAmount = v
			}
		}

		if vals, ok := b.Properties["time_stamp"]; ok && len(vals) > 0 {
			if v, err := strconv.Atoi(vals[0]); err == nil {
				loan.TimeStamp = v
			}
		}

		if vals, ok := b.Properties["interest_rate"]; ok && len(vals) > 0 {
			if fv, err := strconv.ParseFloat(vals[0], 32); err == nil {
				loan.InterestRate = dataformat.Float(fv)
			}
		}

		if vals, ok := b.Properties["duration"]; ok && len(vals) > 0 {
			if v, err := strconv.Atoi(vals[0]); err == nil {
				loan.Duration = v
			}
		}

		out = append(out, loan)
	}

	return out
}

// ToEconomy searches for the first block of type "economy" and converts it to
// a typed items.Economy using the FromProperties helper.
func (d *Document) ToEconomy() (*items.Economy, error) {
	for _, b := range d.Blocks {
		if b.Type != "economy" {
			continue
		}
		var e items.Economy
		if err := e.FromProperties(b.Properties); err != nil {
			return nil, err
		}
		return &e, nil
	}
	return nil, nil
}

// ToPlayer searches for the first block of type "player" and converts it to a
// typed items.Player using the FromProperties helper.
func (d *Document) ToPlayer() (*items.Player, error) {
	for _, b := range d.Blocks {
		if b.Type != "player" {
			continue
		}
		var p items.Player
		if err := p.FromProperties(b.Properties); err != nil {
			return nil, err
		}
		return &p, nil
	}
	return nil, nil
}

// DebugString returns a human-friendly dump of the document.
func (d *Document) DebugString() string {
	var out string
	for _, b := range d.Blocks {
		out += b.Type + " : " + b.Name + "\n"
		for k, vals := range b.Properties {
			for _, v := range vals {
				out += "  " + k + ": " + v + "\n"
			}
		}
		out += "\n"
	}
	return out
}
