package repo

import (
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v79"
)

type Payout struct {
	Id      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func FromStripePayoutAndTotals(payout *stripe.Payout, gross, fee, net int) *Payout {
	return &Payout{
		Id:      payout.ID,
		Created: time.Unix(payout.Created, 0).UTC().Format("2 Jan 2006"),
		Gross:   strconv.Itoa(gross),
		Fee:     strconv.Itoa(fee),
		Net:     strconv.Itoa(net),
	}
}

type Invoice struct {
	Id          string
	Created     string
	ClientName  string
	ClientEmail string
	PayoutId    string
	Gross       string
	Fee         string
	Net         string
}

func FromChargeTransactionAndPayoutId(charge *stripe.BalanceTransaction, payoutId string) *Invoice {
	return &Invoice{
		Id:          charge.ID,
		Created:     time.Unix(charge.Created, 0).UTC().Format("2 Jan 2006"),
		ClientName:  charge.Source.Charge.BillingDetails.Name,
		ClientEmail: charge.Source.Charge.BillingDetails.Email,
		PayoutId:    payoutId,
		Gross:       strconv.Itoa(int(charge.Amount)),
		Fee:         strconv.Itoa(int(charge.Fee)),
		Net:         strconv.Itoa(int(charge.Net)),
	}
}

func FromChargeTransactionsAndPayoutId(charges []*stripe.BalanceTransaction, payoutId string) []Invoice {
	donations := make([]Invoice, len(charges))
	for i, d := range charges {
		donations[i] = *FromChargeTransactionAndPayoutId(d, payoutId)
	}
	return donations
}
