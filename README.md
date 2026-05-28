# Stripe payout reconciliation webhook

Small Go service that receives Stripe payout reconciliation webhooks, validates and enriches the payout data, and persists normalized payout + invoice records into append-only CSV files.

## Project structure

```text
cmd/              entrypoint
data/             CSV persistence
docs/             design notes
internal/
├── handler/      orchestration
├── repo/         CSV persistence + serialization
├── stripeapi/    Stripe API enrichment
└── validation/   validation rules
scripts/          local replay
testdata/         webhook fixtures
```

## Processing flow

```text
Stripe webhook
→ request parsing
→ validation
→ Stripe API enrichment
→ validation
→ row transformation
→ persistence
```

The handler owns orchestration directly. Effects remain visible in the main execution flow.

## Persistence model

```text
docs/persistence-model.md
```

## Environment variables

```bash
export STRIPE_SECRET=sk_live_...
export WEBHOOK_SECRET=whsec_...
```

## Running locally

```bash
go run ./cmd
```

Default webhook endpoint:

```text
POST /webhook
```

## Local development with Stripe test mode

Found in: Stripe Dashboard → Developers → API Keys

```bash
export STRIPE_SECRET=sk_test_...
```

```bash
stripe login
stripe listen --forward-to localhost:8080/webhook
```

This outputs a webhook signing secret:

```bash
export WEBHOOK_SECRET=whsec_...
```

## Testing

The system prioritizes unit testing of pure functions:

* validation rules
* transformation functions

Run all tests:

```bash
go test ./...
```

## Webhook testing limitation

`payout.reconciliation_completed` is not supported by the Stripe CLI.

It cannot be triggered deterministically because it depends on Stripe’s internal payout reconciliation system.

To test this event, we use a deterministic replay flow with a fixed payout object.

## Replay script

Source payout: Stripe Dashboard → Payouts → select payout → copy full payout object

Store it in:

```text
testdata/stripe/payout.reconciliation_completed.json
```

inside the `data.object` field of a minimal event wrapper.

Then replay locally:

```bash
./scripts/stripe-replay.sh testdata/stripe/payout.reconciliation_completed.json
```

This:

* signs the payload using `WEBHOOK_SECRET`
* sends it to the local webhook endpoint
* triggers the webhook handler as if Stripe had sent the event
