## Local development with Stripe test mode

Found in: Stripe Dashboard → Developers → API Keys

```bash
export STRIPE_SECRET=sk_test_...
````

```bash
stripe login
stripe listen --forward-to localhost:8080/webhook
```

This outputs a webhook signing secret:

```bash
export WEBHOOK_SECRET=whsec_...
```

---

## Webhook testing limitation

`payout.reconciliation_completed` is not supported by the Stripe CLI.

It cannot be triggered deterministically because it depends on Stripe’s internal payout reconciliation system.

To test this event, we use a deterministic replay flow with a fixed payout object.

---

## Replay script

Source payout:
Stripe Dashboard → Payouts → select payout → copy full payout object

Store it in:

```text
testdata/stripe/payout.reconciliation_completed.json
```

Inside the `data.object` field of a minimal event wrapper.

Then replay locally:

```bash
./scripts/stripe-replay.sh testdata/stripe/payout.reconciliation_completed.json
```

This:

* signs the payload using `WEBHOOK_SECRET`
* sends it to the local webhook endpoint
* triggers the webhook handler as if Stripe had sent the event
