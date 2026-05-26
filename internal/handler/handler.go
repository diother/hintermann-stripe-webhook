package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/diother/hintermann-stripe-webhook/internal/repo"
	"github.com/diother/hintermann-stripe-webhook/internal/stripeapi"
	"github.com/diother/hintermann-stripe-webhook/internal/validation"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

func HandleWebhook(webhookSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			fail(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		// request materialization
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fail(w, http.StatusBadRequest, fmt.Errorf("failed to read body: %w", err))
			return
		}
		r.Body.Close()

		// webhook parsing + validation
		stripePayout, err := parseWebhookEvent(body, r.Header.Get("Stripe-Signature"), webhookSecret)
		if err != nil {
			fail(w, http.StatusBadRequest, err)
			return
		}

		// api enrichment
		payoutTxn, chargeTxns, err := stripeapi.FetchRelatedTransactions(stripePayout.ID)
		if err != nil {
			fail(w, http.StatusInternalServerError, fmt.Errorf("transaction fetch failed: %w", err))
			return
		}

		// transformation
		payout, invoices, err := processTransactions(stripePayout, payoutTxn, chargeTxns)
		if err != nil {
			fail(w, http.StatusInternalServerError, err)
			return
		}

		// persistence + dedupe
		result, err := repo.WritePayoutAndInvoices(payout, invoices)
		if err != nil {
			fail(w, http.StatusInternalServerError, fmt.Errorf("failed to persist payout %s: %w", stripePayout.ID, err))
			return
		}
		switch result {
		case repo.WriteResultCommitted:
			log.Printf("reconciliation persisted for %s\n", stripePayout.ID)

		case repo.WriteResultIdempotent:
			log.Printf("duplicate webhook ignored for %s\n", stripePayout.ID)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func fail(w http.ResponseWriter, status int, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(status), status)
}

func parseWebhookEvent(body []byte, signature string, webhookSecret string) (*stripe.Payout, error) {
	event, err := webhook.ConstructEvent(body, signature, webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	if event.Type != "payout.reconciliation_completed" {
		return nil, fmt.Errorf("unrecognized event type: %s", event.Type)
	}

	stripePayout := &stripe.Payout{}
	if err = json.Unmarshal(event.Data.Raw, stripePayout); err != nil {
		return nil, fmt.Errorf("unrecognized data object: %w", err)
	}

	if err := validation.ValidatePayout(stripePayout); err != nil {
		return nil, fmt.Errorf("stripe payout invalid: %w", err)
	}

	return stripePayout, nil
}

func processTransactions(
	stripePayout *stripe.Payout,
	payoutTxn *stripe.BalanceTransaction,
	chargeTxns []*stripe.BalanceTransaction,
) (
	*repo.Payout,
	[]repo.Invoice, error,
) {
	if err := validation.ValidatePayoutTransaction(payoutTxn); err != nil {
		return nil, nil, fmt.Errorf("payout transaction invalid: %w", err)
	}

	if err := validation.ValidateChargeTransactions(chargeTxns); err != nil {
		return nil, nil, fmt.Errorf("charge transactions invalid: %w", err)
	}

	gross, fee, net, err := validation.ValidateMatchingSums(payoutTxn, chargeTxns)
	if err != nil {
		return nil, nil, fmt.Errorf("matching sum validation failed: %w", err)
	}

	payout := repo.FromStripePayoutAndTotals(stripePayout, gross, fee, net)
	invoices := repo.FromChargeTransactionsAndPayoutId(chargeTxns, stripePayout.ID)

	return payout, invoices, nil
}
