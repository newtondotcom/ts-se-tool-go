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
