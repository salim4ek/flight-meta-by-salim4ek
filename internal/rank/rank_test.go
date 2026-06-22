package rank

import (
	"testing"

	"flightmeta/internal/offer"
	"flightmeta/internal/sources"
)

func mk(id string, price int64, conn offer.ConnectionType, carriers []string, segs int) offer.Offer {
	o := offer.Offer{ID: id, PriceMinor: price, Connection: conn, Carriers: carriers}
	for i := 0; i < segs; i++ {
		o.Segments = append(o.Segments, offer.Segment{})
	}
	return o
}

func TestApplySortsByPriceAndFilters(t *testing.T) {
	in := []offer.Offer{
		mk("a", 300, offer.Official, []string{"TK"}, 2),          // 1 stop
		mk("b", 100, offer.SelfTransfer, []string{"FZ", "EK"}, 2), // 1 stop, self-transfer
		mk("c", 200, offer.Official, []string{"SU"}, 1),          // direct
	}

	// Cheapest first, everything allowed.
	out := Apply(in, sources.Filters{AllowSelfTransfer: true})
	if len(out) != 3 || out[0].ID != "b" || out[2].ID != "a" {
		t.Fatalf("unexpected order: %v", ids(out))
	}

	// Disallowing self-transfer drops b.
	out = Apply(in, sources.Filters{AllowSelfTransfer: false})
	if len(out) != 2 || hasID(out, "b") {
		t.Fatalf("self-transfer not filtered: %v", ids(out))
	}

	// Direct-only keeps just c.
	out = Apply(in, sources.Filters{Stops: sources.StopsDirect, AllowSelfTransfer: true})
	if len(out) != 1 || out[0].ID != "c" {
		t.Fatalf("direct filter wrong: %v", ids(out))
	}

	// one_plus keeps the connecting offers, drops the direct.
	out = Apply(in, sources.Filters{Stops: sources.StopsOnePlus, AllowSelfTransfer: true})
	if len(out) != 2 || hasID(out, "c") {
		t.Fatalf("one_plus filter wrong: %v", ids(out))
	}

	// Excluding SU drops c.
	out = Apply(in, sources.Filters{ExcludeAirlines: []string{"su"}, AllowSelfTransfer: true})
	if hasID(out, "c") {
		t.Fatalf("exclude airline failed: %v", ids(out))
	}

	// Including only EK keeps just b.
	out = Apply(in, sources.Filters{IncludeAirlines: []string{"EK"}, AllowSelfTransfer: true})
	if len(out) != 1 || out[0].ID != "b" {
		t.Fatalf("include airline failed: %v", ids(out))
	}
}

func TestVisaFreeTransitFilter(t *testing.T) {
	withVisa := func(id, status string) offer.Offer {
		return offer.Offer{
			ID:       id,
			Segments: []offer.Segment{{}, {}}, // 1 stop
			Layovers: []offer.Layover{{VisaStatus: status}},
		}
	}
	in := []offer.Offer{
		withVisa("free", "no_visa"),
		withVisa("twov", "twov"),
		withVisa("need", "visa_required"),
	}
	out := Apply(in, sources.Filters{AllowSelfTransfer: true, OnlyVisaFreeTransit: true})
	if hasID(out, "need") {
		t.Fatalf("visa_required offer should be dropped: %v", ids(out))
	}
	if !hasID(out, "free") || !hasID(out, "twov") {
		t.Fatalf("no_visa/twov offers should be kept: %v", ids(out))
	}
}

func ids(os []offer.Offer) []string {
	r := make([]string, 0, len(os))
	for _, o := range os {
		r = append(r, o.ID)
	}
	return r
}

func hasID(os []offer.Offer, id string) bool {
	for _, o := range os {
		if o.ID == id {
			return true
		}
	}
	return false
}
