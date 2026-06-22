// Package mock is a dependency-free Adapter that returns deterministic offers.
// It lets the whole pipeline run and be verified end-to-end before real source
// API keys (Kiwi, Travelpayouts) are wired in. The offers are route-agnostic:
// the hubs below are just illustrative, not tied to any one country.
package mock

import (
	"context"
	"fmt"
	"time"

	"flightmeta/internal/offer"
	"flightmeta/internal/sources"
)

// Adapter is the mock data source.
type Adapter struct{}

// New returns a ready-to-use mock adapter.
func New() *Adapter { return &Adapter{} }

// Name implements sources.Adapter.
func (a *Adapter) Name() string { return "mock" }

// Search returns three illustrative offers: a pricey official direct, an
// official 1-stop, and a cheaper self-transfer combo that mainstream
// aggregators do not surface as a single option.
func (a *Adapter) Search(ctx context.Context, q sources.Query) ([]offer.Offer, error) {
	if err := ctx.Err(); err != nil { // respect cancellation even for trivial work
		return nil, err
	}
	dep := q.DepartDate
	if dep.IsZero() {
		dep = time.Now().Add(14 * 24 * time.Hour)
	}
	cur := q.Currency
	if cur == "" {
		cur = "RUB"
	}

	// 1) Official direct — convenient but pricey.
	direct := offer.Offer{
		ID:         "mock-direct",
		Source:     a.Name(),
		Connection: offer.Official,
		Currency:   cur,
		PriceMinor: 4200000, // 42 000.00
		DeepLink:   "https://example-partner.test/book/mock-direct",
		Segments: []offer.Segment{
			seg(q.Origin, q.Destination, dep.Add(10*time.Hour), dep.Add(19*time.Hour), "SU", "270"),
		},
		Carriers: []string{"SU"},
	}

	// 2) Official 1-stop via a hub — single ticket, protected connection.
	hub := "IST"
	oneStop := offer.Offer{
		ID:         "mock-official-1stop",
		Source:     a.Name(),
		Connection: offer.Official,
		Currency:   cur,
		PriceMinor: 3100000,
		DeepLink:   "https://example-partner.test/book/mock-official-1stop",
		Segments: []offer.Segment{
			seg(q.Origin, hub, dep.Add(8*time.Hour), dep.Add(12*time.Hour), "TK", "414"),
			seg(hub, q.Destination, dep.Add(14*time.Hour), dep.Add(22*time.Hour), "TK", "68"),
		},
		Layovers: []offer.Layover{{Airport: hub, Duration: 2 * time.Hour}},
		Carriers: []string{"TK"},
	}

	// 3) Self-transfer combo via a different hub — two separate tickets stitched
	//    into a route aggregators do not show as one option. Cheapest here, but
	//    unprotected. This is the product's headline capability.
	hub2 := "DXB"
	selfTransfer := offer.Offer{
		ID:         "mock-self-transfer",
		Source:     a.Name(),
		Connection: offer.SelfTransfer,
		Currency:   cur,
		PriceMinor: 2450000,
		DeepLink:   "https://example-partner.test/book/mock-self-transfer",
		Unique:     true,
		Segments: []offer.Segment{
			seg(q.Origin, hub2, dep.Add(6*time.Hour), dep.Add(11*time.Hour), "FZ", "918"),
			seg(hub2, q.Destination, dep.Add(15*time.Hour), dep.Add(23*time.Hour), "EK", "418"),
		},
		Layovers: []offer.Layover{{Airport: hub2, Duration: 4 * time.Hour, SelfTransfer: true}},
		Carriers: []string{"FZ", "EK"},
	}

	// 4) Self-transfer via China — demonstrates a conditional TWOV transit-visa
	//    badge (China visa-free transit) for an RU/CIS passport.
	hub3 := "PEK"
	selfTransferCN := offer.Offer{
		ID:         "mock-self-transfer-cn",
		Source:     a.Name(),
		Connection: offer.SelfTransfer,
		Currency:   cur,
		PriceMinor: 2690000,
		DeepLink:   "https://example-partner.test/book/mock-self-transfer-cn",
		Unique:     true,
		Segments: []offer.Segment{
			seg(q.Origin, hub3, dep.Add(7*time.Hour), dep.Add(16*time.Hour), "U6", "1001"),
			seg(hub3, q.Destination, dep.Add(20*time.Hour), dep.Add(26*time.Hour), "CA", "979"),
		},
		Layovers: []offer.Layover{{Airport: hub3, Duration: 4 * time.Hour, SelfTransfer: true}},
		Carriers: []string{"U6", "CA"},
	}

	// 5) Self-transfer via India — demonstrates a "transit visa required" badge,
	//    which the "visa-free transit only" filter should drop.
	hub4 := "DEL"
	selfTransferIN := offer.Offer{
		ID:         "mock-self-transfer-in",
		Source:     a.Name(),
		Connection: offer.SelfTransfer,
		Currency:   cur,
		PriceMinor: 3350000,
		DeepLink:   "https://example-partner.test/book/mock-self-transfer-in",
		Unique:     true,
		Segments: []offer.Segment{
			seg(q.Origin, hub4, dep.Add(9*time.Hour), dep.Add(15*time.Hour), "6E", "44"),
			seg(hub4, q.Destination, dep.Add(19*time.Hour), dep.Add(24*time.Hour), "6E", "1071"),
		},
		Layovers: []offer.Layover{{Airport: hub4, Duration: 4 * time.Hour, SelfTransfer: true}},
		Carriers: []string{"6E"},
	}

	return []offer.Offer{direct, oneStop, selfTransfer, selfTransferCN, selfTransferIN}, nil
}

func seg(from, to string, dep, arr time.Time, carrier, flight string) offer.Segment {
	return offer.Segment{
		From: from, To: to,
		DepartUTC: dep.UTC(), ArriveUTC: arr.UTC(),
		MarketingCarrier: carrier,
		FlightNumber:     fmt.Sprintf("%s%s", carrier, flight),
	}
}
