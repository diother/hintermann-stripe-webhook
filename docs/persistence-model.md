## Data model

### `payouts.csv` → payout ledger

```
Id, Created, Gross, Fee, Net
```

Types:

* `Id` → string (Stripe payout id)
* `Created` → string (2 Jan 2006 format)
* `Gross/Fee/Net` → string-encoded integers

### `invoices.csv` → invoice records tied to payouts

```
Id, Created, ClientName, ClientEmail, PayoutId, Gross, Fee, Net
```

Types:

* `Id` → string (Stripe balance transaction id)
* `Created` → string (2 Jan 2006 format)
* `ClientName` → string
* `ClientEmail` → string
* `PayoutId` → string (foreign key to payouts.csv)
* `Gross/Fee/Net` → string-encoded integers

## Write process

1. Acquire process-level lock
2. Read current snapshot from `data/`
3. Validate idempotency (payout ID uniqueness)
4. Build full updated snapshot in `tmp/`
5. Flush and fsync all files in `tmp/`
6. Atomically swap:

   * `data → old`
   * `tmp → data`

7. Remove `old`

## Failure model

* `CLEAN` → safe to operate
* `tmpDir exists` → interrupted write → invalid state
* `oldDir exists` → interrupted commit → invalid state

Recovery is handled manual intervention.

## Why not a database?

* allows direct inspection via CSV
* avoids external dependencies
