package repo

import "path/filepath"

const (
	dataDir     = "data"
	tmpDir      = "tmp"
	oldDir      = "old"
	payoutsCsv  = "payouts.csv"
	invoicesCsv = "invoices.csv"
)

type WriteResult int

const (
	WriteResultCommitted WriteResult = iota
	WriteResultIdempotent
)

func payoutsPath(dir string) string {
	return filepath.Join(dir, payoutsCsv)
}

func invoicesPath(dir string) string {
	return filepath.Join(dir, invoicesCsv)
}
