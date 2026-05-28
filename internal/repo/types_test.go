package repo

import (
	"testing"

	"github.com/stripe/stripe-go/v79"
)

func TestFromStripePayoutAndTotals(t *testing.T) {
	payout := &stripe.Payout{
		ID:      "test_id",
		Created: 1779915600,
	}

	got := FromStripePayoutAndTotals(payout, 200, 6, 194)

	want := Payout{
		Id:      "test_id",
		Created: "27 May 2026",
		Gross:   "200",
		Fee:     "6",
		Net:     "194",
	}

	if got.Id != want.Id ||
		got.Created != want.Created ||
		got.Gross != want.Gross ||
		got.Fee != want.Fee ||
		got.Net != want.Net {
		t.Errorf("mismatch:\ngot  %+v\nwant %+v", got, want)
	}
}

func TestFromChargeTransactionAndPayoutId(t *testing.T) {
	charge := &stripe.BalanceTransaction{
		ID:      "test_id",
		Created: 1779915600,
		Amount:  100,
		Fee:     10,
		Net:     90,
		Source: &stripe.BalanceTransactionSource{
			Charge: &stripe.Charge{
				BillingDetails: &stripe.ChargeBillingDetails{
					Email: "test@gmail.com",
				},
			},
		},
	}

	got := FromChargeTransactionAndPayoutId(charge, "payout_id")

	want := Invoice{
		Id:          "test_id",
		Created:     "27 May 2026",
		ClientName:  "",
		ClientEmail: "test@gmail.com",
		PayoutId:    "payout_id",
		Gross:       "100",
		Fee:         "10",
		Net:         "90",
	}

	if got.Id != want.Id ||
		got.Created != want.Created ||
		got.ClientName != want.ClientName ||
		got.ClientEmail != want.ClientEmail ||
		got.PayoutId != want.PayoutId ||
		got.Gross != want.Gross ||
		got.Fee != want.Fee ||
		got.Net != want.Net {
		t.Errorf("mismatch:\ngot  %+v\nwant %+v", got, want)
	}
}

func TestFromChargeTransactionsAndPayoutId(t *testing.T) {
	charge := &stripe.BalanceTransaction{
		ID:      "test_id",
		Created: 1779915600,
		Amount:  100,
		Fee:     10,
		Net:     90,
		Source: &stripe.BalanceTransactionSource{
			Charge: &stripe.Charge{
				BillingDetails: &stripe.ChargeBillingDetails{
					Email: "test@gmail.com",
				},
			},
		},
	}

	charges := []*stripe.BalanceTransaction{
		charge,
		charge,
	}

	got := FromChargeTransactionsAndPayoutId(charges, "payout_id")

	invoice := Invoice{
		Id:          "test_id",
		Created:     "27 May 2026",
		ClientName:  "",
		ClientEmail: "test@gmail.com",
		PayoutId:    "payout_id",
		Gross:       "100",
		Fee:         "10",
		Net:         "90",
	}

	want := []Invoice{
		invoice,
		invoice,
	}

	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if !equalInvoice(got[i], want[i]) {
			t.Errorf("mismatch at index %d:\ngot  %+v\nwant %+v", i, got[i], want[i])
		}
	}
}

func equalInvoice(a, b Invoice) bool {
	return a.Id == b.Id &&
		a.Created == b.Created &&
		a.ClientName == b.ClientName &&
		a.ClientEmail == b.ClientEmail &&
		a.PayoutId == b.PayoutId &&
		a.Gross == b.Gross &&
		a.Fee == b.Fee &&
		a.Net == b.Net
}
