package repo

import "testing"

func TestPayoutExists(t *testing.T) {
	inputId := "po_1"
	inputPayouts := []Payout{
		{Id: "po_3"},
		{Id: "po_1"},
		{Id: "po_3"},
	}

	want := true
	got := payoutExists(inputPayouts, inputId)

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
