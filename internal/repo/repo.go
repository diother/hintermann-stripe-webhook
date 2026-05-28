package repo

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	dataDir = "data"
	tmpDir  = "tmp"
	oldDir  = "old"
)

type WriteResult int

const (
	WriteResultCommitted WriteResult = iota
	WriteResultIdempotent
)

var mu sync.Mutex

func WritePayoutAndInvoices(p Payout, ivs []Invoice) (WriteResult, error) {
	mu.Lock()
	defer mu.Unlock()

	if err := checkForInterruptedWrites(); err != nil {
		return 0, err
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return 0, err
	}

	payouts, err := readPayouts(dataDir)
	if err != nil {
		return 0, err
	}

	if payoutExists(payouts, p.Id) {
		return WriteResultIdempotent, nil
	}

	invoices, err := readInvoices(dataDir)
	if err != nil {
		return 0, err
	}

	payouts = append(payouts, p)
	invoices = append(invoices, ivs...)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return 0, err
	}

	if err := writePayouts(tmpDir, payouts); err != nil {
		return 0, err
	}

	if err := writeInvoices(tmpDir, invoices); err != nil {
		return 0, err
	}

	if err := commitSnapshot(); err != nil {
		return 0, err
	}

	return WriteResultCommitted, nil
}

func checkForInterruptedWrites() error {
	if _, err := os.Stat(tmpDir); err == nil {
		return fmt.Errorf("inconsistent state: tmpDir exists")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot verify tmpDir state: %w", err)
	}

	if _, err := os.Stat(oldDir); err == nil {
		return fmt.Errorf("inconsistent state: oldDir exists")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot verify oldDir state: %w", err)
	}

	return nil
}

func readPayouts(dir string) ([]Payout, error) {
	path := payoutsPath(dir)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Payout{}, nil
		}

		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	payouts := make([]Payout, 0, len(records))

	for _, row := range records {
		if len(row) != 5 {
			return nil, fmt.Errorf("invalid payout row length")
		}
		payouts = append(payouts, Payout{
			Id:      row[0],
			Created: row[1],
			Gross:   row[2],
			Fee:     row[3],
			Net:     row[4],
		})
	}

	return payouts, nil
}

func readInvoices(dir string) ([]Invoice, error) {
	path := invoicesPath(dir)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Invoice{}, nil
		}

		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()

	invoices := make([]Invoice, 0, len(records))

	for _, row := range records {
		if len(row) != 8 {
			return nil, fmt.Errorf("invalid invoice row length")
		}

		invoices = append(invoices, Invoice{
			Id:          row[0],
			Created:     row[1],
			ClientName:  row[2],
			ClientEmail: row[3],
			PayoutId:    row[4],
			Gross:       row[5],
			Fee:         row[6],
			Net:         row[7],
		})
	}

	return invoices, nil
}

func writePayouts(dir string, payouts []Payout) error {
	path := payoutsPath(dir)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)

	for _, p := range payouts {
		err := w.Write([]string{
			p.Id,
			p.Created,
			p.Gross,
			p.Fee,
			p.Net,
		})
		if err != nil {
			return err
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func writeInvoices(dir string, invoices []Invoice) error {
	path := invoicesPath(dir)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)

	for _, i := range invoices {
		err := w.Write([]string{
			i.Id,
			i.Created,
			i.ClientName,
			i.ClientEmail,
			i.PayoutId,
			i.Gross,
			i.Fee,
			i.Net,
		})
		if err != nil {
			return err
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func payoutExists(payouts []Payout, id string) bool {
	for _, p := range payouts {
		if p.Id == id {
			return true
		}
	}
	return false
}

func commitSnapshot() error {
	if err := os.Rename(dataDir, oldDir); err != nil {
		return err
	}

	if err := os.Rename(tmpDir, dataDir); err != nil {
		return err
	}

	return os.RemoveAll(oldDir)
}

func payoutsPath(dir string) string {
	return filepath.Join(dir, "payouts.csv")
}

func invoicesPath(dir string) string {
	return filepath.Join(dir, "invoices.csv")
}
