package visa

import (
	"testing"

	"flightmeta/internal/offer"
)

func TestLoadAndEnrich(t *testing.T) {
	r, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if r.Disclaimer() == "" {
		t.Fatal("expected a non-empty disclaimer")
	}

	if c, ok := r.Country("DXB"); !ok || c != "AE" {
		t.Fatalf("DXB country: got %q ok=%v", c, ok)
	}
	if _, ok := r.Country("ZZZ"); ok {
		t.Fatal("ZZZ should be an unknown airport")
	}

	offers := []offer.Offer{
		{Layovers: []offer.Layover{{Airport: "DXB", SelfTransfer: true}}},  // AE -> no_visa
		{Layovers: []offer.Layover{{Airport: "PEK", SelfTransfer: true}}},  // CN -> twov
		{Layovers: []offer.Layover{{Airport: "DEL", SelfTransfer: true}}},  // IN -> visa_required
		{Layovers: []offer.Layover{{Airport: "FRA", SelfTransfer: true}}},  // Schengen self-transfer -> visa_required
		{Layovers: []offer.Layover{{Airport: "FRA", SelfTransfer: false}}}, // Schengen airside -> no_visa
		{Layovers: []offer.Layover{{Airport: "ZZZ", SelfTransfer: true}}},  // unknown airport
	}
	r.Enrich(offers, "RU")

	want := []string{"no_visa", "twov", "visa_required", "visa_required", "no_visa", "unknown"}
	for i, w := range want {
		lay := offers[i].Layovers[0]
		if lay.VisaStatus != w {
			t.Errorf("offer %d (%s): VisaStatus = %q, want %q", i, lay.Airport, lay.VisaStatus, w)
		}
		if lay.TransitNote == "" {
			t.Errorf("offer %d: empty TransitNote", i)
		}
	}
}

func TestEnrichFallsBackToDefaultPerPassport(t *testing.T) {
	r, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	// PEK is not in the KZ ruleset -> KZ DEFAULT (unknown).
	offers := []offer.Offer{{Layovers: []offer.Layover{{Airport: "PEK"}}}}
	r.Enrich(offers, "KZ")
	if got := offers[0].Layovers[0].VisaStatus; got != "unknown" {
		t.Fatalf("KZ via PEK: got %q, want unknown (DEFAULT)", got)
	}
}
