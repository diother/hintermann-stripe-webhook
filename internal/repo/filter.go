package repo

func payoutExists(payouts []Payout, id string) bool {
	for _, p := range payouts {
		if p.Id == id {
			return true
		}
	}
	return false
}
