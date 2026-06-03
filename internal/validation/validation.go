package validation

import (
	"fmt"

	"github.com/stripe/stripe-go/v79"
)

func ValidatePayout(payout *stripe.Payout) error {
	if payout == nil {
		return fmt.Errorf("is nil")
	}
	if payout.ID == "" {
		return fmt.Errorf("id is missing")
	}
	if payout.Created <= 0 {
		return fmt.Errorf("created is not positive")
	}
	if payout.ReconciliationStatus != "completed" {
		return fmt.Errorf("reconciliation status is not completed")
	}
	return nil
}

func ValidatePayoutTransaction(payout *stripe.BalanceTransaction) error {
	if payout == nil {
		return fmt.Errorf("is nil")
	}
	if payout.Type != "payout" {
		return fmt.Errorf("type is not payout")
	}
	if payout.ID == "" {
		return fmt.Errorf("id is missing")
	}
	if payout.Created <= 0 {
		return fmt.Errorf("created is not positive")
	}
	if payout.Amount >= 0 {
		return fmt.Errorf("amount is not negative")
	}
	if payout.Fee != 0 {
		return fmt.Errorf("fee is not 0")
	}
	if payout.Net >= 0 {
		return fmt.Errorf("net is not negative")
	}
	return nil
}

func ValidateChargeTransactions(charges []*stripe.BalanceTransaction) error {
	if len(charges) == 0 {
		return fmt.Errorf("slice is nil")
	}
	for i, charge := range charges {
		if err := validateChargeTransaction(charge); err != nil {
			return fmt.Errorf("index %d %w", i, err)
		}
	}
	return nil
}

func validateChargeTransaction(charge *stripe.BalanceTransaction) error {
	if charge == nil {
		return fmt.Errorf("is nil")
	}
	if charge.Type == "stripe_fee" {
		return fmt.Errorf("stripe_fee transactions are forbidden")
	}
	if charge.Type != "charge" && charge.Type != "payment" {
		return fmt.Errorf("type is not charge, payment, or stripe_fee")
	}
	if charge.ID == "" {
		return fmt.Errorf("id is missing")
	}
	if charge.Created <= 0 {
		return fmt.Errorf("created is not positive")
	}
	if charge.Amount <= 0 {
		return fmt.Errorf("amount is not positive")
	}
	if charge.Fee <= 0 {
		return fmt.Errorf("fee is not positive")
	}
	if charge.Net <= 0 {
		return fmt.Errorf("net is not positive")
	}
	if charge.Source == nil {
		return fmt.Errorf("source is nil")
	}
	if charge.Source.Charge == nil {
		return fmt.Errorf("charge object is nil")
	}
	if charge.Source.Charge.BillingDetails == nil {
		return fmt.Errorf("billing details is nil")
	}
	if charge.Source.Charge.BillingDetails.Email == "" {
		return fmt.Errorf("email is missing")
	}
	return nil
}

func ValidateMatchingSums(
	payout *stripe.BalanceTransaction,
	charges []*stripe.BalanceTransaction,
) (
	int,
	int,
	int,
	error,
) {
	var gross, fee, net int

	for _, charge := range charges {
		gross += int(charge.Amount)
		fee += int(charge.Fee)
	}

	net = gross - fee
	payoutAmount := int(-payout.Amount)

	if payoutAmount != net {
		return 0, 0, 0,
			fmt.Errorf(
				"payout amount mismatch. amount %v != net %v",
				payoutAmount,
				net,
			)
	}
	return gross, fee, net, nil
}
