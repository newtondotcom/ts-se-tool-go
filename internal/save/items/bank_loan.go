package items

import "github.com/robebs/ts-se-tool-go/internal/save/dataformat"

// BankLoan mirrors the C# Bank_Loan class and represents a single bank_loan block.
type BankLoan struct {
	Amount         int
	OriginalAmount int
	TimeStamp      int
	InterestRate   dataformat.Float
	Duration       int
}

// FromLines constructs a BankLoan from the lines of a bank_loan SII block.
// It is a partial port of the C# string[] constructor; full integration with
// a generic parser will come later.
func (b *BankLoan) FromLines(lines []string) {
	for _, line := range lines {
		// TODO: delegate to a shared SII line parser; for now this is just a stub.
		_ = line
	}
}


