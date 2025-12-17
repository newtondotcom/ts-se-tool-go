package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// Bank mirrors the C# Bank class from CustomClasses/Save/Items/Bank.cs.
type Bank struct {
	MoneyAccount            int64
	CoinsuranceFixed        int
	CoinsuranceRatio        dataformat.Float
	AccidentSeverity        dataformat.Float
	Loans                   []string
	AppEnabled              bool
	LoanLimit               int
	PaymentTimer            dataformat.Float
	Overdraft               bool
	OverdraftTimer          dataformat.Float
	OverdraftWarnCount      int
	SellPlayersTruckLater   bool
	SellPlayersTrailerLater bool
}

// FromProperties populates the Bank from SII properties.
func (b *Bank) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "money_account":
			v, _ := strconv.ParseInt(val, 10, 64)
			b.MoneyAccount = v
		case key == "coinsurance_fixed":
			b.CoinsuranceFixed = parseInt(val)
		case key == "coinsurance_ratio":
			b.CoinsuranceRatio = parseFloat(val)
		case key == "accident_severity":
			b.AccidentSeverity = parseFloat(val)
		case key == "loans":
			// capacity hint ignored
		case strings.HasPrefix(key, "loans["):
			b.Loans = append(b.Loans, val)
		case key == "app_enabled":
			b.AppEnabled = parseBool(val)
		case key == "loan_limit":
			b.LoanLimit = parseInt(val)
		case key == "payment_timer":
			b.PaymentTimer = parseFloat(val)
		case key == "overdraft":
			b.Overdraft = parseBool(val)
		case key == "overdraft_timer":
			b.OverdraftTimer = parseFloat(val)
		case key == "overdraft_warn_count":
			b.OverdraftWarnCount = parseInt(val)
		case key == "sell_players_truck_later":
			b.SellPlayersTruckLater = parseBool(val)
		case key == "sell_players_trailer_later":
			b.SellPlayersTrailerLater = parseBool(val)
		}
	}
	return nil
}

// ToProperties produces a map of SII properties equivalent to the C# PrintOut.
func (b *Bank) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["money_account"] = []string{strconv.FormatInt(b.MoneyAccount, 10)}
	props["coinsurance_fixed"] = []string{strconv.Itoa(b.CoinsuranceFixed)}
	props["coinsurance_ratio"] = []string{formatFloat(b.CoinsuranceRatio)}
	props["accident_severity"] = []string{formatFloat(b.AccidentSeverity)}

	props["loans"] = []string{strconv.Itoa(len(b.Loans))}
	for i, v := range b.Loans {
		props[fmt.Sprintf("loans[%d]", i)] = []string{v}
	}

	props["app_enabled"] = []string{formatBool(b.AppEnabled)}
	props["loan_limit"] = []string{strconv.Itoa(b.LoanLimit)}
	props["payment_timer"] = []string{formatFloat(b.PaymentTimer)}
	props["overdraft"] = []string{formatBool(b.Overdraft)}
	props["overdraft_timer"] = []string{formatFloat(b.OverdraftTimer)}
	props["overdraft_warn_count"] = []string{strconv.Itoa(b.OverdraftWarnCount)}
	props["sell_players_truck_later"] = []string{formatBool(b.SellPlayersTruckLater)}
	props["sell_players_trailer_later"] = []string{formatBool(b.SellPlayersTrailerLater)}

	return props
}
