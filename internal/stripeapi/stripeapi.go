package stripeapi

import (
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
)

func FetchRelatedTransactions(id string) (
	*stripe.BalanceTransaction,
	[]*stripe.BalanceTransaction,
	error,
) {
	params := &stripe.BalanceTransactionListParams{}
	params.Payout = &id
	params.AddExpand("data.source")

	iter := balancetransaction.List(params)

	var payout *stripe.BalanceTransaction
	if iter.Next() {
		payout = iter.BalanceTransaction()
	}

	var charges []*stripe.BalanceTransaction
	for iter.Next() {
		charges = append(charges, iter.BalanceTransaction())
	}

	if err := iter.Err(); err != nil {
		return nil, nil, err
	}
	return payout, charges, nil
}
